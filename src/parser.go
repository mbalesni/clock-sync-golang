package src

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func Parse(input string) (*[]*Process, error) {

	f, err := os.OpenFile(input, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to read input file %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	currentLine := ""

	var processes []*Process

	for scanner.Scan() {

		currentLine = scanner.Text()
		currentLine = strings.TrimSpace(currentLine)
		currentLine = strings.ReplaceAll(currentLine, " ", "")

		lineSplit := strings.Split(currentLine, ",")

		if len(lineSplit) == 3 {

			id := lineSplit[0]
			name := lineSplit[1]
			time := lineSplit[2]

			timeList := strings.Split(time, ":")
			idInt, _ := strconv.ParseInt(id, 10, 64)
			hoursParsed, _ := strconv.ParseInt(timeList[0], 10, 64)
			minutesParsed, _ := strconv.ParseInt(timeList[1], 10, 64)

			process := NewProcess(int(idInt), name, &Time{Hours: int(hoursParsed), Minutes: int(minutesParsed)}, name[2:], false)

			processes = append(processes, process)

		} else {

			return &processes, fmt.Errorf("Incomplete line")

		}

	}

	return &processes, nil

}
