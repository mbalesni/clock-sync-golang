package main

import (
	"clock_sync_golang/src"
	"fmt"
	prompt "github.com/c-bata/go-prompt"
	"os"
	"strings"
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
				if len(whitespaceSplit) < 2 || len(whitespaceSplit) > 2 {

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
		}
	}

}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		//{Text: "Materialize", Description: "Loads the default config, no input needed"},
		//{Text: "Read", Description: "Reads a .txt file in the specified format"},
		{Text: "List", Description: "Lists all current active nodes in the ring, no input"},
		{Text: "Lookup", Description: "Lookups up a node, key:start_node"},
		{Text: "Join", Description: "Joins the given node Id with the ring"},
		{Text: "Leave", Description: "Shuts down the specified Node"},
		{Text: "Shortcut", Description: "Adds a shortcut to the specified node"},
		//{Text: "Shutdown", Description: "Shuts down the whole cluster"},
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

			executor(fmt.Sprintf("Read %s", parsedArgs))

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
