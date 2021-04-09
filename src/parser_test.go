package src

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {

	var processes *[]*Process

	processes, _ = Parse("../processes_file.txt")

	for _, process := range *processes {

		fmt.Println(process.Id, process.Name, (*process).Time)

	}

}
