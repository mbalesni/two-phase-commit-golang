package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"two-phase-program/src"

	prompt "github.com/c-bata/go-prompt"
)

var network *src.Network
var processes []*src.Process

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
		whitespaceSplit[0] != "Arbitrary-failure" {

		fmt.Println("Invalid command")

	} else {

		switch command := whitespaceSplit[0]; command {
		case "Set-value":
			{
				if len(whitespaceSplit) != 2 {
					fmt.Println("Set-value takes the next value as an argument")
					return
				} else {
					value, err := strconv.Atoi(whitespaceSplit[1])
					fmt.Println("Wants to set value:", value)
					if err == nil {
						network.SetValue(value)
					} else {
						return
					}
				}
			}
		case "Rollback":
			{
				if len(whitespaceSplit) != 2 {
					fmt.Println("Rollback takes the number of steps to roll back by as an argument")
				} else {
					value, err := strconv.Atoi(whitespaceSplit[1])
					fmt.Println("Wants to roll back by:", value)
					if err == nil {
						for _, process := range network.Processes {
							if len(process.History) < value {
								fmt.Println(fmt.Errorf("the system cannot reverse to that long state"))
								return
							}
						}
						network.Rollback(value)
					}
				}
			}
			// case "List":
			// 	{
			// if len(whitespaceSplit) != 1 {
			// 	fmt.Println("List takes no argument")
			// } else {
			// 	network.List()
			// }
			// 	}
			// case "Clock":
			// 	{
			// 		if len(whitespaceSplit) != 1 {
			// 			fmt.Println("Clock takes no argument")
			// 		} else {
			// 			network.Clock()
			// 		}
			// 	}
			// case "Set-time":
			// 	{
			// 		if len(whitespaceSplit) != 3 {
			// 			fmt.Println("Set-time takes 2 arguments, the process id and the time hour:minute")
			// 		} else {
			// 			processId, err := strconv.ParseInt(whitespaceSplit[1], 10, 64)

			// 			if err != nil {

			// 				fmt.Println("Failed to parse process id")
			// 			} else {
			// 				_, exists := network.Processes[int(processId)]
			// 				if !exists {
			// 					fmt.Println("Process does not exist.")
			// 					return
			// 				}

			// 				hourMinutes := strings.Split(whitespaceSplit[2], ":")

			// 				if len(hourMinutes) != 2 {

			// 					fmt.Println("There cannot be more, or less, than 2 values delimited by :")

			// 				} else {

			// 					hours, err := strconv.ParseInt(hourMinutes[0], 10, 64)

			// 					if err != nil {

			// 						fmt.Println("Failed to parse the hours")

			// 					} else {

			// 						minutes, err := strconv.ParseInt(hourMinutes[1], 10, 64)

			// 						if err != nil {

			// 							fmt.Println("Failed to parse the minutes")

			// 						} else {

			// 							network.SetTime(int(processId), src.Time{Hours: int(hours), Minutes: int(minutes)})

			// 						}
			// 					}
			// 				}
			// 			}
			// 		}
			// 	}
			// case "Reload":
			// 	{
			// 		if len(whitespaceSplit) != 2 {

			// 			fmt.Println("Read takes one argument, a text file")

			// 		} else {

			// 			file_location := whitespaceSplit[1]

			// 			processes, err := src.Parse(file_location)

			// 			if err != nil {

			// 				fmt.Printf("Something bad happened whilst reloading", err)

			// 			} else {

			// 				network.Reload(processes)

			// 			}
			// 		}
			// 	}
			// case "Freeze":
			// 	{
			// 		if len(whitespaceSplit) != 2 {

			// 			fmt.Println("Freeze takes one argument, a processId")

			// 		} else {

			// 			processId, err := strconv.ParseInt(whitespaceSplit[1], 10, 64)

			// 			if err != nil {

			// 				fmt.Println("Failed to parse process id")
			// 			} else {
			// 				_, exists := network.Processes[int(processId)]
			// 				if !exists {
			// 					fmt.Println("Process does not exist.")
			// 					return
			// 				}

			// 				network.Freeze(int(processId))

			// 			}

			// 		}
			// 	}
			// case "Unfreeze":
			// 	{
			// 		if len(whitespaceSplit) != 2 {

			// 			fmt.Println("Unfreeze takes one argument, a processId")

			// 		} else {

			// 			processId, err := strconv.ParseInt(whitespaceSplit[1], 10, 64)

			// 			if err != nil {

			// 				fmt.Println("Failed to parse process id")
			// 			} else {
			// 				_, exists := network.Processes[int(processId)]
			// 				if !exists {
			// 					fmt.Println("Process does not exist.")
			// 					return
			// 				}
			// 				network.Unfreeze(int(processId))

			// 			}

			// 		}
			// 	}

			// case "Kill":
			// 	{
			// 		if len(whitespaceSplit) != 2 {

			// 			fmt.Println("Kill takes one argument, a processId")

			// 		} else {

			// 			processId, err := strconv.ParseInt(whitespaceSplit[1], 10, 64)

			// 			if err != nil {

			// 				fmt.Println("Failed to parse process id")
			// 			} else {
			// 				_, exists := network.Processes[int(processId)]
			// 				if !exists {
			// 					fmt.Println("Process does not exist.")
			// 					return
			// 				}
			// 				network.Kill(int(processId))

			// 			}
			// 		}
			// 	}
		}
		network.ListHistory()
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

		fmt.Println("No argument given")
		os.Exit(1)

	} else {

		parsedArgs := strings.TrimSpace(os.Args[1])

		if strings.HasSuffix(parsedArgs, ".txt") {

			// TODO: parse from file

			processes = append(processes, src.NewProcess("P1", true))
			processes = append(processes, src.NewProcess("P2", true))
			processes = append(processes, src.NewProcess("P3", true))
			processes = append(processes, src.NewProcess("P4", true))

			network = src.SpawnNetwork(&processes)
			network.Coordinator = (processes)[3]

			// file_location := whitespaceSplit[1]

			// file, err := src.Parse(file_location)

			// if err != nil {

			// 	fmt.Printf("Something bad happened whilst parsing", err)

			// } else {

			// 	network = src.SpawnNetwork(file)

			// }

			// // Synchronizing
			// go func() {
			// 	for {
			// 		select {
			// 		case <-timer1.C:
			// 			{
			// 				//fmt.Println("Berkleying")
			// 				network.Berkley()
			// 			}
			// 		}
			// 	}
			// }()

			// go func() {
			// 	for {
			// 		select {
			// 		case <-timer2.C:
			// 			//fmt.Println("Time-ing")
			// 			tempCurrentTime := src.CurrentTime()
			// 			diff := currentTime.Distance(tempCurrentTime)
			// 			// Clock ticking
			// 			for _, process := range network.Processes {

			// 				process.SyncTime(diff)

			// 			}
			// 			currentTime = tempCurrentTime

			// 		}
			// 	}
			// }()

			p := prompt.New(
				executor,
				completer,
				prompt.OptionPrefix("λ "),
				prompt.OptionTitle("prompt for Huber's take on 2PC"),
			)
			p.Run()
		} else {

			fmt.Println("Please provide a .txt file")
			os.Exit(1)

		}
	}
}
