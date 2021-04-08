package src

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 1, A_0, 11:00
// 3, B_0, 13:33
// 4, D_0, 17:30
// 7, E_0, 23:00
// 5, F_0, 3:00

func GetProcessIdxById(id int, processes []*Process) int {
	for i, process := range processes {
		if process.Id == id {
			return i
		}
	}
	return -1
}

func TestBullying(t *testing.T) {
	verbose := true
	var processes map[int]*Process
	processes = make(map[int]*Process)

	// Inserting the processes

	processes[1] = NewProcess(1, "A_0", &Time{Hours: 11, Minutes: 0}, verbose)
	processes[3] = NewProcess(3, "B_0", &Time{Hours: 13, Minutes: 33}, verbose)
	processes[4] = NewProcess(4, "D_0", &Time{Hours: 17, Minutes: 30}, verbose)
	processes[7] = NewProcess(7, "E_0", &Time{Hours: 23, Minutes: 0}, verbose)
	processes[5] = NewProcess(5, "F_0", &Time{Hours: 3, Minutes: 0}, verbose)

	// add higher & lower processes
	for _, currentProcess := range processes {
		for targetProcessId, targetProcess := range processes {

			if targetProcessId > currentProcess.Id {

				currentProcess.HigherProcesses[targetProcessId] = targetProcess
				currentProcess.MaxCoordinatorWait += 1 // will need to wait for 1 cycle longer for each new higher process

			} else if targetProcess.Id < currentProcess.Id {

				currentProcess.LowerProcesses[targetProcessId] = targetProcess

			}

		}
	}

	/*
		for i, process := range processes {

			assert.Equal(t, len(processes)-i, len(process.HigherProcesses))
			assert.Equal(t, i, len(process.LowerProcesses))

		}
	*/

	// Asserting that the process distribution was ok
	assert.Equal(t, 4, len(processes[1].HigherProcesses))
	assert.Equal(t, 4, len(processes[7].LowerProcesses))

	// start election from lowest process
	// What's the point of adding an ID to the election??
	processes[1].RunElection(-1)

	// I didn't get this at all
	for i := 0; i < 10; i++ {
		// sync network (move messages from Send to Get queues)
		for _, process := range processes {
			for process.SendQueue.queue.Len() > 0 {
				message := process.SendQueue.Pop()
				message.To.GetQueue.Add(message)
			}
		}

		for _, process := range processes {
			process.Cycle()
		}
	}

	for _, process := range processes {
		// verify election
		assert.Equal(t, 7, process.Coordinator.Id, "Coordinator is chosen correctly")

		// verify waiting
		assert.Equal(t, -1, process.WaitingCoordinator, "Not waiting for coordinator reply after election over")
		assert.Equal(t, -1, process.WaitingElection, "Not waiting for coordinator reply after election over")

		// verify names update
		assert.Equal(t, "1", string(process.Name[2:]), "Election count is correct")
	}

	// list
	for _, process := range processes {
		coordinatorString := ""
		if process.Coordinator.Id == process.Id {
			coordinatorString = "(Coordinator)"
		}
		fmt.Println(process.Id, process.Name, coordinatorString)
	}
	println()

	// ELECTION from Id=4
	processes[4].RunElection(-1)

	for i := 0; i < 10; i++ {
		// sync network (move messages from Send to Get queues)
		for _, process := range processes {
			for process.SendQueue.queue.Len() > 0 {
				message := process.SendQueue.Pop()
				message.To.GetQueue.Add(message)
			}
		}

		for _, process := range processes {
			process.Cycle()
		}
	}

	for _, process := range processes {
		// verify election
		assert.Equal(t, 7, process.Coordinator.Id, "Coordinator is chosen correctly")

		// verify waiting
		assert.Equal(t, -1, process.WaitingCoordinator, "Not waiting for coordinator reply after election over")
		assert.Equal(t, -1, process.WaitingElection, "Not waiting for coordinator reply after election over")

		// verify names update
		assert.Equal(t, "2", string(process.Name[2:]), "Election count is correct")
	}
	println()

	// list
	for _, process := range processes {
		coordinatorString := ""
		if process.Coordinator.Id == process.Id {
			coordinatorString = "(Coordinator)"
		}
		fmt.Println(process.Id, process.Name, coordinatorString)
	}

	println()

	// time
	for _, process := range processes {
		fmt.Println(process.Name, process.Time)
	}

}
