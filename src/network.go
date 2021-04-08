package src

import "fmt"

type Network struct {
	Processes map[int]*Process
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

	return n.Processes[processId].Coordinator

}

func (n *Network) List() {

	for _, process := range n.Processes {
		coordinatorString := ""
		if process.Coordinator.Id == process.Id {
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
