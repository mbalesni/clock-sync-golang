package src

import (
	"fmt"
)

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

	network.AutoDiscovery()

	for _, process := range network.Processes {

		if !process.Frozen {

			network.BullyStartingFrom(process.Id)

			break

		}

	}

	for _, process := range network.Processes {

		if !process.Frozen {

			process.Name = fmt.Sprintf("%s%s", string(process.Name[:2]), process.InitialCount)

		}

	}

	fmt.Println("Coordinator chosen: ", network.Coordinator.Id, network.Coordinator.Name)

	return network

}

func (n *Network) AutoDiscovery() {
	// populate LowerProcesses and HigherProcesses for each process
	for _, currentProcess := range n.Processes {
		for targetProcessId, targetProcess := range n.Processes {

			if targetProcessId > currentProcess.Id {

				_, alreadyIn := currentProcess.HigherProcesses[targetProcessId]

				if !alreadyIn {
					currentProcess.HigherProcesses[targetProcessId] = targetProcess
					currentProcess.MaxCoordinatorWait += 1 // will need to wait for 1 cycle longer for each new higher process
				}

			} else if targetProcess.Id < currentProcess.Id {

				_, alreadyIn := currentProcess.LowerProcesses[targetProcessId]

				if !alreadyIn {
					currentProcess.LowerProcesses[targetProcessId] = targetProcess
				}

			}

		}
	}
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
		//fmt.Println(receiver.Time, receiver.InitialTime)

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
		fmt.Println(process.Id, process.Name, FormatTime(process.Time.Hours, process.Time.Minutes))
	}

}

func (n *Network) Freeze(processId int) {

	n.Processes[processId].Frozen = true

	if n.Coordinator != nil && n.Coordinator.Id == processId {

		for _, process := range n.Processes {

			process.Coordinator = nil

			// from https://piazza.com/class/klsd1d3jd5q6qi?cid=20#
			// “If coordinator is removed in any condition, then default time should be used”
			// we assume that "removed in any condition = Killed or Frozen"
			if !process.Frozen && process.Id != processId {

				newTime := *(process.InitialTime)
				process.Time = &newTime

			}

		}

		n.Coordinator = nil

		for _, process := range n.Processes {

			if !process.Frozen {
				fmt.Println("Starting election from:", process.Id)
				n.BullyStartingFrom(process.Id)
				//n.Berkley()

				break

			}

		}
		fmt.Println("The new coordinator is:", n.Coordinator.Id)
	}

}

func (n *Network) Unfreeze(processId int) {

	n.Processes[processId].Frozen = false

	if processId > n.Coordinator.Id {

		newTime := *(n.Processes[processId].InitialTime)
		n.Processes[processId].Time = &newTime

		n.BullyStartingFrom(processId)
		//n.Berkley()

	} else {

		n.Coordinator = nil

		for _, process := range n.Processes {

			if !process.Frozen {

				n.BullyStartingFrom(process.Id)
				//n.Berkley()

				return

			}

		}

	}
}

func (n *Network) Kill(processId int) {

	for _, currentProcess := range n.Processes {

		delete(currentProcess.HigherProcesses, processId)
		delete(currentProcess.LowerProcesses, processId)

	}

	delete(n.Processes, processId)

	processesLeft := len(n.Processes)
	if processesLeft == 0 {
		fmt.Println("No processes left.")
		return
	}

	if processId == n.Coordinator.Id {

		n.Coordinator = nil

		for _, process := range n.Processes {

			if !process.Frozen {

				newTime := *(process.InitialTime)
				process.Time = &newTime

			}

		}

		for _, process := range n.Processes {

			if !process.Frozen {

				fmt.Println("Starting election from:", process.Id)

				n.BullyStartingFrom(process.Id)
				//n.Berkley()
				break

			}

		}

		fmt.Println("The new coordinator is:", n.Coordinator.Id)

	}

}

func (n *Network) Reload(processes *[]*Process) {

	for _, process := range *processes {

		_, exists := n.Processes[process.Id]

		if !exists {

			process.Verbose = false
			process.Init()
			n.Processes[process.Id] = process

		}
	}

	n.AutoDiscovery()

	for _, process := range n.Processes {

		if !process.Frozen {

			n.BullyStartingFrom(process.Id)
			//n.Berkley()
			break

		}

	}

}

func (n *Network) SetTime(processId int, time Time) *Process {

	n.Processes[processId].Time = &time

	return n.Processes[processId]

}
