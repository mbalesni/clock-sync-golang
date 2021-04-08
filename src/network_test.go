package src

import (
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

	var processes []*Process

	processes = append(processes, NewProcess(1, "A_0", &Time{Hours: 11, Minutes: 0}, verbose))
	processes = append(processes, NewProcess(3, "B_0", &Time{Hours: 13, Minutes: 33}, verbose))
	processes = append(processes, NewProcess(4, "D_0", &Time{Hours: 17, Minutes: 30}, verbose))
	processes = append(processes, NewProcess(7, "E_0", &Time{Hours: 23, Minutes: 0}, verbose))
	processes = append(processes, NewProcess(5, "F_0", &Time{Hours: 3, Minutes: 0}, verbose))

	network := SpawnNetwork(&processes)

	network.BullyStartingFrom(1)

	for _, process := range processes {
		// verify election
		assert.Equal(t, 7, process.Coordinator.Id, "Coordinator is chosen correctly")

		// verify waiting
		assert.Equal(t, -1, process.WaitingCoordinator, "Not waiting for coordinator reply after election over")
		assert.Equal(t, -1, process.WaitingElection, "Not waiting for coordinator reply after election over")

		// verify names update
		assert.Equal(t, "1", string(process.Name[2:]), "Election count is correct")
	}

	// ELECTION from Id=4
	network.BullyStartingFrom(4)

	for _, process := range processes {
		// verify election
		assert.Equal(t, 7, process.Coordinator.Id, "Coordinator is chosen correctly")

		// verify waiting
		assert.Equal(t, -1, process.WaitingCoordinator, "Not waiting for coordinator reply after election over")
		assert.Equal(t, -1, process.WaitingElection, "Not waiting for coordinator reply after election over")

		// verify names update
		assert.Equal(t, "2", string(process.Name[2:]), "Election count is correct")
	}

	network.List()

	println()

	network.Clock()

}
