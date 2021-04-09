package parser

import (
	"bufio"
	"errors"
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

	currentLine := 0

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

			time := strings.Split(val, ":")
			hoursParsed, _ := strconv.ParseInt(time[0], 10, 64)
			minutesParsed, _ := strconv.ParseInt(time[1], 10, 64)

			process := NewProcess(id, name, &Time{Hours: hoursParsed, Minutes: minutesParsed}, false)

			processes = append(processes, &process)

		} else {

			return processes, fmt.Errorf("Incomplete line")

		}

	}

	return &processes, nil

}
