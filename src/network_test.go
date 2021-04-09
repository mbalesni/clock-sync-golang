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
	//verbose := true

	var processes *[]*Process

	processes, err := Parse("../processes_file.txt")

	if err != nil {

		t.Fail()

	}

	network := SpawnNetwork(processes)

	for _, process := range *processes {
		// verify election
		assert.Equal(t, 7, process.Coordinator.Id, "Coordinator is chosen correctly")

		// verify waiting
		assert.Equal(t, -1, process.WaitingCoordinator, "Not waiting for coordinator reply after election over")
		assert.Equal(t, -1, process.WaitingElection, "Not waiting for coordinator reply after election over")

		// verify names update
		assert.Equal(t, "0", string(process.Name[2:]), "Election count is correct")
	}

	network.BullyStartingFrom(4)

	for _, process := range *processes {
		// verify election
		assert.Equal(t, 7, process.Coordinator.Id, "Coordinator is chosen correctly")

		// verify waiting
		assert.Equal(t, -1, process.WaitingCoordinator, "Not waiting for coordinator reply after election over")
		assert.Equal(t, -1, process.WaitingElection, "Not waiting for coordinator reply after election over")

		// verify names update
		assert.Equal(t, "1", string(process.Name[2:]), "Election count is correct")
	}

}

func TestFreezeBullying(t *testing.T) {
	//verbose := true

	var processes *[]*Process

	processes, err := Parse("../processes_file.txt")

	if err != nil {

		t.Fail()

	}

	network := SpawnNetwork(processes)

	network.Freeze(3)

	for _, process := range *processes {
		// verify election
		assert.Equal(t, 7, process.Coordinator.Id, "Coordinator is chosen correctly")

		// verify waiting
		assert.Equal(t, -1, process.WaitingCoordinator, "Not waiting for coordinator reply after election over")
		assert.Equal(t, -1, process.WaitingElection, "Not waiting for coordinator reply after election over")

		// verify names update
		assert.Equal(t, "0", string(process.Name[2:]), "Election count is correct")
	}

}

func TestFreezeUnfreezeKillBullying(t *testing.T) {
	//verbose := true

	var processes *[]*Process

	processes, err := Parse("../processes_file.txt")

	if err != nil {
		t.Fail()
	}

	network := SpawnNetwork(processes)

	println("Froze 3")
	network.Freeze(3)

	for _, process := range *processes {

		if process.Frozen != true {
			// verify election
			assert.Equal(t, 7, process.Coordinator.Id, "Coordinator is chosen correctly")

			// verify waiting
			assert.Equal(t, -1, process.WaitingCoordinator, "Not waiting for coordinator reply after election over")
			assert.Equal(t, -1, process.WaitingElection, "Not waiting for coordinator reply after election over")

			// verify names update
			assert.Equal(t, "0", string(process.Name[2]), "Election count is correct")

		}
	}

	network.List()

	println()

	network.Clock()

	println("\nUnfroze 3")
	network.Unfreeze(3)
	println("Froze 7")
	network.Freeze(7)
	network.List()
	println()
	network.Clock()

	println("Killing 7")
	network.Kill(7)
	network.List()
	println()
	network.Clock()
	println("Berkleying")
	network.Berkley()
	network.List()
	println()
	network.Clock()

}
