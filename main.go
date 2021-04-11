package main

import (
	"clock_sync_golang/src"
	"time"

	"fmt"
	"os"
	"strconv"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

var network *src.Network

func executor(in string) {

	whitespaceSplit := strings.Fields(in)

	if len(whitespaceSplit) == 0 {
		return
	}

	if whitespaceSplit[0] != "Read" &&
		whitespaceSplit[0] != "List" &&
		whitespaceSplit[0] != "Clock" &&
		whitespaceSplit[0] != "Kill" &&
		whitespaceSplit[0] != "Set-time" &&
		whitespaceSplit[0] != "Freeze" &&
		whitespaceSplit[0] != "Unfreeze" &&
		whitespaceSplit[0] != "Reload" {

		fmt.Println("Invalid command")

	} else {

		switch command := whitespaceSplit[0]; command {
		case "Read":
			{
				if len(whitespaceSplit) != 2 {

					fmt.Println("Read takes one argument, a text file")

				} else {

					file_location := whitespaceSplit[1]

					file, err := src.Parse(file_location)

					if err != nil {

						fmt.Printf("Something bad happened whilst parsing", err)

					} else {

						network = src.SpawnNetwork(file)

					}

				}
			}
		case "List":
			{
				if len(whitespaceSplit) != 1 {
					fmt.Println("List takes no argument")
				} else {
					network.List()
				}
			}
		case "Clock":
			{
				if len(whitespaceSplit) != 1 {
					fmt.Println("Clock takes no argument")
				} else {
					network.Clock()
				}
			}
		case "Set-time":
			{
				if len(whitespaceSplit) != 3 {
					fmt.Println("Set-time takes 2 arguments, the process id and the time hour:minute")
				} else {
					processId, err := strconv.ParseInt(whitespaceSplit[1], 10, 64)

					if err != nil {

						fmt.Println("Failed to parse process id")
					} else {
						_, exists := network.Processes[int(processId)]
						if !exists {
							fmt.Println("Process does not exist.")
							return
						}

						hourMinutes := strings.Split(whitespaceSplit[2], ":")

						if len(hourMinutes) != 2 {

							fmt.Println("There cannot be more, or less, than 2 values delimited by :")

						} else {

							hours, err := strconv.ParseInt(hourMinutes[0], 10, 64)

							if err != nil {

								fmt.Println("Failed to parse the hours")

							} else {

								minutes, err := strconv.ParseInt(hourMinutes[1], 10, 64)

								if err != nil {

									fmt.Println("Failed to parse the minutes")

								} else {

									network.SetTime(int(processId), src.Time{Hours: int(hours), Minutes: int(minutes)})

								}
							}
						}
					}
				}
			}
		case "Reload":
			{
				if len(whitespaceSplit) != 2 {

					fmt.Println("Read takes one argument, a text file")

				} else {

					file_location := whitespaceSplit[1]

					processes, err := src.Parse(file_location)

					if err != nil {

						fmt.Printf("Something bad happened whilst reloading", err)

					} else {

						network.Reload(processes)

					}
				}
			}
		case "Freeze":
			{
				if len(whitespaceSplit) != 2 {

					fmt.Println("Freeze takes one argument, a processId")

				} else {

					processId, err := strconv.ParseInt(whitespaceSplit[1], 10, 64)

					if err != nil {

						fmt.Println("Failed to parse process id")
					} else {
						_, exists := network.Processes[int(processId)]
						if !exists {
							fmt.Println("Process does not exist.")
							return
						}

						network.Freeze(int(processId))

					}

				}
			}
		case "Unfreeze":
			{
				if len(whitespaceSplit) != 2 {

					fmt.Println("Unfreeze takes one argument, a processId")

				} else {

					processId, err := strconv.ParseInt(whitespaceSplit[1], 10, 64)

					if err != nil {

						fmt.Println("Failed to parse process id")
					} else {
						_, exists := network.Processes[int(processId)]
						if !exists {
							fmt.Println("Process does not exist.")
							return
						}
						network.Unfreeze(int(processId))

					}

				}
			}

		case "Kill":
			{
				if len(whitespaceSplit) != 2 {

					fmt.Println("Kill takes one argument, a processId")

				} else {

					processId, err := strconv.ParseInt(whitespaceSplit[1], 10, 64)

					if err != nil {

						fmt.Println("Failed to parse process id")
					} else {
						_, exists := network.Processes[int(processId)]
						if !exists {
							fmt.Println("Process does not exist.")
							return
						}
						network.Kill(int(processId))

					}
				}
			}
		}
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "List", Description: "List all processes"},
		{Text: "Clock", Description: "List all processes' clocks"},
		{Text: "Set-time", Description: "Sets time of a process"},
		{Text: "Reload", Description: "Reloads file, needs file location"},
		{Text: "Freeze", Description: "Freeze a process"},
		{Text: "Unfreeze", Description: "Unfreezes a process"},
		{Text: "Kill", Description: "Murders a process"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {

	if len(os.Args) < 2 {

		fmt.Println("No argument given")
		os.Exit(1)

	} else {

		parsedArgs := strings.TrimSpace(os.Args[1])

		if strings.HasSuffix(parsedArgs, ".txt") {

			currentTime := src.CurrentTime()
			executor(fmt.Sprintf("Read %s", parsedArgs))

			timer1 := time.NewTicker(5 * time.Second)
			timer2 := time.NewTicker(1 * time.Minute)

			// Synchronizing
			go func() {
				for {
					select {
					case <-timer1.C:
						{
							//fmt.Println("Berkleying")
							network.Berkley()
						}
					}
				}
			}()

			go func() {
				for {
					select {
					case <-timer2.C:
						//fmt.Println("Time-ing")
						tempCurrentTime := src.CurrentTime()
						diff := currentTime.Distance(tempCurrentTime)
						// Clock ticking
						for _, process := range network.Processes {

							process.SyncTime(diff)

						}
						currentTime = tempCurrentTime

					}
				}
			}()

			p := prompt.New(
				executor,
				completer,
				prompt.OptionPrefix("λ "),
				prompt.OptionTitle("prompt for huber's take on bully + berkley"),
			)
			p.Run()
		} else {

			fmt.Println("Please provide a .txt file")
			os.Exit(1)

		}
	}
}
