package src

import "fmt"

type Network struct {
	Processes   map[string]*Process
	Coordinator *Process
}

func SpawnNetwork(processes *[]*Process) *Network {

	network := Network{}
	network.Processes = make(map[string]*Process)
	for _, process := range *processes {
		process.Verbose = false
		process.Init()
		network.Processes[process.Name] = process
	}

	network.AutoDiscovery()

	return &network

}

func (n *Network) AutoDiscovery() {
	// populate LowerProcesses and HigherProcesses for each process
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
		process.Cycle()
	}
}

func (n *Network) ListHistory() {

	fmt.Println("Listing history")

	for _, process := range n.Processes {
		fmt.Println(process.Name, process.History)
	}
}

func (n *Network) SetValue(value int) {

	n.Coordinator.InitCommit(value)
	for i := 0; i < 3; i++ {
		n.Cycle()
	}
}

func (n *Network) Rollback(value int) {

	fmt.Println("Rolling back")

	for _, process := range n.Processes {
		process.History = process.History[:len(process.History)-value]
	}
}
