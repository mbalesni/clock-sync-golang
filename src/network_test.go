package src

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 1, A_0, 11:00am
// 3, B_0, 13:33pm
// 4, D_0, 17:30pm
// 7, E_0, 23:00pm
// 5, F_0, 3:00am

func GetProcessIdxById(id int, processes []*Process) int {
	for i, process := range processes {
		if process.Id == id {
			return i
		}
	}
	return -1
}

func TestBullying(t *testing.T) {
	verbose := false
	var processes []*Process

	processes = append(processes, NewProcess(1, "A_0", "11:00", verbose))
	processes = append(processes, NewProcess(3, "B_0", "13:33", verbose))
	processes = append(processes, NewProcess(4, "D_0", "17:30", verbose))
	processes = append(processes, NewProcess(7, "E_0", "23:00", verbose))
	processes = append(processes, NewProcess(5, "F_0", "3:00", verbose))

	// add higher & lower processes
	for i, process := range processes {
		for _, target_process := range processes {

			if target_process.Id > process.Id {

				processes[i].HigherProcesses = append(processes[i].HigherProcesses, &ProcessShallow{Id: target_process.Id})
				processes[i].MaxCoordinatorWait += 1 // will need to wait for 1 cycle longer for each new higher process

			} else if target_process.Id < process.Id {

				processes[i].LowerProcesses = append(processes[i].LowerProcesses, &ProcessShallow{Id: target_process.Id})

			}

		}
	}

	assert.Equal(t, 4, len(processes[0].HigherProcesses))
	assert.Equal(t, 4, len(processes[3].LowerProcesses))

	// start election from lowest process
	processes[0].RunElection(-1)

	for i := 0; i < 10; i++ {
		// sync network (move messages from Send to Get queues)
		for _, process := range processes {
			for process.SendQueue.queue.Len() > 0 {
				message := process.SendQueue.Pop()
				sendToIdx := GetProcessIdxById(message.To, processes)
				processes[sendToIdx].GetQueue.Add(message)
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
		fmt.Println(process.Id, process.Name, process.Time)
	}

	// ELECTION from Id=4
	processes[2].RunElection(-1)

	for i := 0; i < 10; i++ {
		// sync network (move messages from Send to Get queues)
		for _, process := range processes {
			for process.SendQueue.queue.Len() > 0 {
				message := process.SendQueue.Pop()
				sendToIdx := GetProcessIdxById(message.To, processes)
				processes[sendToIdx].GetQueue.Add(message)
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

	// list
	for _, process := range processes {
		fmt.Println(process.Id, process.Name, process.Time)
	}

}
