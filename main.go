package main

import (
	prompt "github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"two-phase-program/src"
)

var network *src.Network

func executor(in string) {

	whitespaceSplit := strings.Fields(in)

	if len(whitespaceSplit) == 0 {
		return
	}

	if whitespaceSplit[0] != "Set-value" &&
		whitespaceSplit[0] != "Rollback" &&
		whitespaceSplit[0] != "Add" &&
		whitespaceSplit[0] != "Remove" &&
		whitespaceSplit[0] != "Time-failure" &&
		whitespaceSplit[0] != "Arbitrary-failure" &&
		whitespaceSplit[0] != "List" {

		log.Error("Invalid command")

	} else {
		switch command := whitespaceSplit[0]; command {
		case "Set-value":
			{
				if len(whitespaceSplit) != 2 {
					log.Error("Set-value takes the next value as an argument")
				} else {
					value, err := strconv.Atoi(whitespaceSplit[1])
					if err != nil {
						log.Errorf("%v", err)
					} else {
						log.Println("Wants to set value:", value)
						network.OperationSetValue(value)
						network.ListHistory()
					}
				}
			}
		case "Rollback":
			{
				if len(whitespaceSplit) != 2 {
					log.Error("Rollback takes the number of steps to roll back by as an argument")
				} else {
					value, err := strconv.Atoi(whitespaceSplit[1])
					if err != nil {
						log.Errorf("%v", err)
					} else {
						if len(network.Coordinator.Log) < value {
							log.Error("the system cannot reverse to that long state")
						} else {
							log.Println("Wants to roll back by:", value)
							network.OperationRollback(value)
							network.ListHistory()
						}
					}
				}
			}
		case "Add":
			{
				if len(whitespaceSplit) != 2 {
					log.Error("Add takes a process name as an argument")
				} else {
					if network.Processes[whitespaceSplit[1]] != nil {
						log.Error("Can't add process with the same name")
					} else {
						network.OperationAdd(whitespaceSplit[1])
						network.ListHistory()
					}
				}
			}
		case "Remove":
			{
				if len(whitespaceSplit) != 2 {
					log.Error("Remove takes a process name as an argument")
				} else {
					if network.Processes[whitespaceSplit[1]] == nil {
						log.Error("Can't remove what's not there :|")
					} else if network.Coordinator.Name == whitespaceSplit[1] {
						log.Error("Can't remove the coordinator!")
					} else {
						network.OperationRemove(whitespaceSplit[1])
						network.ListHistory()
					}
				}
			}
		case "Time-failure":
			{
				if len(whitespaceSplit) != 3 {
					log.Error("Time-failure takes a process name FIRST and duration in seconds as an argument")
				} else {
					value, err := strconv.Atoi(whitespaceSplit[2])
					if err != nil {
						log.Errorf("%v", err)
					} else if network.Processes[whitespaceSplit[1]] == nil {
						log.Error("Can't cause what's not there to fail :|")
					} else if network.Coordinator.Name == whitespaceSplit[1] {
						log.Error("Can't cause the coordinator to fail!")
					} else {
						network.OperationSetTimeFailure(whitespaceSplit[1], value)
						network.ListHistory()
					}
				}
			}
		case "Arbitrary-failure":
			{
				if len(whitespaceSplit) != 3 {
					log.Error("Arbitrary-failure takes a process name FIRST and duration in seconds as an argument")
				} else {
					value, err := strconv.Atoi(whitespaceSplit[2])
					if err != nil {
						log.Errorf("%v", err)
					} else if network.Processes[whitespaceSplit[1]] == nil {
						log.Error("Can't cause what's not there to fail :|")
					} else if network.Coordinator.Name == whitespaceSplit[1] {
						log.Error("Can't cause the coordinator to fail!")
					} else {
						network.OperationSetArbitraryFailure(whitespaceSplit[1], value)
						network.ListHistory()
					}
				}
			}
		case "List":
			{
				network.ListHistory()
			}
		}
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "Set-value", Description: "Commit a new value to the system"},
		{Text: "Rollback", Description: "Roll back to previous values"},
		{Text: "Add", Description: "Add a new process participant to the system"},
		{Text: "Remove", Description: "Remove a process from the system"},
		{Text: "Time-failure", Description: "Kills the node (for S seconds)"},
		{Text: "Arbitrary-failure", Description: "Swaps node's response to commit requests (for S seconds)"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func main() {

	if len(os.Args) < 2 {

		log.Println("No argument given")
		os.Exit(1)

	} else {

		parsedArgs := strings.TrimSpace(os.Args[1])

		if strings.HasSuffix(parsedArgs, ".txt") {

			var err error
			network, err = src.Parse(parsedArgs)

			if err != nil {
				log.Fatalf("failed to read file %v", err)
			}

			p := prompt.New(
				executor,
				completer,
				prompt.OptionPrefix("λ "),
				prompt.OptionTitle("prompt for Huber's take on 2PC"),
			)
			p.Run()
		} else {

			log.Println("Please provide a .txt file")
			os.Exit(1)

		}
	}
}
