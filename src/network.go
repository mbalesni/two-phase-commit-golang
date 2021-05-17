package src

import "log"

type Network struct {
	Processes   map[string]*Process
	Coordinator *Process
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

func (n *Network) Cycle() {
	for _, process := range n.Processes {
		for process.SendQueue.queue.Len() > 0 {
			message := process.SendQueue.Pop()
			message.To.GetQueue.Add(message)
		}
	}

	for _, process := range n.Processes {
		process.ProcessMessages()
	}
}

func (n *Network) ListHistory() {

	log.Println("Listing history")
	log.Printf("Cooordinator:%v", n.Coordinator.Name)

	for _, process := range n.Processes {
		log.Println(process.Name, process.Log)
	}
}
