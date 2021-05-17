package src

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func Parse(input string) (*Network, error) {

	f, err := os.OpenFile(input, os.O_RDONLY, os.ModePerm)

	if err != nil {
		log.Fatalf("Failed to read input file %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	currentLine, header := "", ""

	var network *Network

	var processes []*Process

	var coordinator *Process

	for scanner.Scan() {

		currentLine = scanner.Text()

		if strings.ContainsAny(currentLine, "#") {

			header = strings.ReplaceAll(currentLine, " ", "")
			headerRune := []rune(header)
			headerRune = headerRune[1:len(headerRune)]
			header = string(headerRune)

		} else {

			currentLine = strings.TrimSpace(currentLine)
			currentLine = strings.ReplaceAll(currentLine, " ", "")

			if header == "System"{

				if strings.Contains(currentLine, "Coordinator") {

					lineSplit := strings.Split(currentLine, ";")

					processName := lineSplit[0]

					newProcess := NewProcess(processName, false)

					coordinator = newProcess

					processes = append(processes, coordinator)

				} else {

					if currentLine != "" {

						newProcess := NewProcess(currentLine, false)

						processes = append(processes, newProcess)

					}

				}

			} else {

				lineSplit := strings.Split(currentLine, ";")

				historyLogString := lineSplit[1]

				historyLogStringSlice := strings.Split(historyLogString, ",")

				var historyLog []int

				for _, value := range historyLogStringSlice {

					parsedValue, err := strconv.Atoi(value)

					if err != nil {
						log.Fatalf("failed to parse value %v", err)
					}

					historyLog = append(historyLog, parsedValue)

				}

				for _, value := range processes {

					value.Log = historyLog

					value.Init()

				}


			}

		}

	}

	network = SpawnNetwork(processes)
	network.Coordinator = coordinator

	return network, nil

}
