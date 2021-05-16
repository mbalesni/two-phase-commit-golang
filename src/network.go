package src

import "fmt"

type Network struct {
	Processes   map[string]*Process
	Coordinator *Process
}

func SpawnNetwork(processes []*Process) *Network {

	network := Network{}
	network.Processes = make(map[string]*Process)
	for _, process := range processes {
		process.Verbose = false
		network.Processes[process.Name] = process
	}

	network.AutoDiscovery()

	return &network

}

func (n *Network) AutoDiscovery() {
	for _, currentProcess := range n.Processes {
		for _, targetProcess := range n.Processes {
			_, alreadyIn := currentProcess.OtherProcesses[targetProcess.Name]
			if !alreadyIn && currentProcess.Name != targetProcess.Name {
				currentProcess.OtherProcesses[targetProcess.Name] = targetProcess
			}

		}
	}
}

func (n *Network) Cycle() {
	for _, process := range n.Processes {
		for process.SendQueue.queue.Len() > 0 {
			message := process.SendQueue.Pop()
			fmt.Println(message)
			message.To.GetQueue.Add(message)
		}
	}

	for _, process := range n.Processes {
		process.ProcessMessages()
	}
}

func (n *Network) ListHistory() {

	fmt.Println("Listing history")

	for _, process := range n.Processes {
		fmt.Println(process.Name, process.Log)
	}
}