package src

import "fmt"

type Network struct {
	Processes   map[int]*Process
	Coordinator *Process
}

func SpawnNetwork(processes *[]*Process) *Network {

	network := &Network{}
	network.Processes = make(map[int]*Process)
	for _, process := range *processes {
		process.Verbose = false
		process.Init()
		network.Processes[process.Id] = process
	}
	// Auto discovery
	for _, currentProcess := range network.Processes {
		for targetProcessId, targetProcess := range network.Processes {

			if targetProcessId > currentProcess.Id {

				currentProcess.HigherProcesses[targetProcessId] = targetProcess
				currentProcess.MaxCoordinatorWait += 1 // will need to wait for 1 cycle longer for each new higher process

			} else if targetProcess.Id < currentProcess.Id {

				currentProcess.LowerProcesses[targetProcessId] = targetProcess

			}

		}
	}

	return network

}

func (n *Network) BullyStartingFrom(processId int) *Process {

	n.Processes[processId].RunElection(-1)

	for i := 0; i < len(n.Processes)+2; i++ {
		for _, process := range n.Processes {
			for process.SendQueue.queue.Len() > 0 {
				message := process.SendQueue.Pop()
				message.To.GetQueue.Add(message)
			}
		}

		for _, process := range n.Processes {
			process.Cycle()
		}
	}

	n.Coordinator = n.Processes[processId].Coordinator

	return n.Coordinator

}

func (n *Network) Berkley() {

	// Create messages to send to all other processes

	for _, target := range n.Coordinator.LowerProcesses {
		// Send how much they have to adjust their clock
		message := n.Coordinator.NewSyncMessage(target, "CLOCK_SYNC", target.Time.Distance(n.Coordinator.Time))
		n.Coordinator.SendQueue.Add(message)
	}

	// Send all messages
	for n.Coordinator.SendQueue.queue.Len() > 0 {
		message := n.Coordinator.SendQueue.Pop()
		message.To.GetQueue.Add(message)
	}
	// Process the messages
	for _, receiver := range n.Processes {

		receiver.ProcessMessages()

	}

}

func (n *Network) List() {

	for _, process := range n.Processes {
		coordinatorString := ""
		if process.Coordinator != nil && process.Coordinator.Id == process.Id {
			coordinatorString = "(Coordinator)"
		}
		fmt.Println(process.Id, process.Name, coordinatorString)
	}

}

func (n *Network) Clock() {

	for _, process := range n.Processes {
		fmt.Println(process.Id, process.Name, process.Time)
	}

}

func (n *Network) Freeze(processId int) *Process {

	n.Processes[processId].Frozen = true

	if n.Coordinator != nil && n.Coordinator.Id == processId {

		for _, process := range n.Processes {

			process.Coordinator = nil

		}

		n.Coordinator = nil

		for _, process := range n.Processes {

			if process.Frozen != true {

				n.BullyStartingFrom(process.Id)

				return n.Processes[processId]

			}

		}

	}

	return n.Processes[processId]

}

func (n *Network) Unfreeze(processId int) *Process {

	n.Processes[processId].Frozen = false

	if processId > n.Coordinator.Id {

		n.BullyStartingFrom(processId)

	} else {

		n.Coordinator = nil

		for _, process := range n.Processes {

			if process.Frozen != true {

				n.BullyStartingFrom(process.Id)

				return n.Processes[processId]

			}

		}

	}

	return n.Processes[processId]

}

func (n *Network) Kill(processId int) {

	for _, currentProcess := range n.Processes {

		delete(currentProcess.HigherProcesses, processId)
		delete(currentProcess.LowerProcesses, processId)

	}

	delete(n.Processes, processId)

	for _, process := range n.Processes {

		if process.Frozen != true {

			n.BullyStartingFrom(process.Id)

			break

		}

	}

}
