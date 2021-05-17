package src

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func (n *Network) OperationSetValue(value int) {
	n.Coordinator.SendVoteRequestMessages("set-value", value, "", nil)
	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
}

func (n *Network) OperationRollback(steps int) {
	n.Coordinator.SendVoteRequestMessages("rollback", steps, "", nil)
	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
}

func (n *Network) OperationAdd(processName string) {

	newProcess := NewProcess(processName, false)

	n.Coordinator.SendVoteRequestMessages("add", 0, processName, newProcess)
	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()

	fmt.Println("State:", n.Coordinator.State)

	if n.Coordinator.State == "commit" {
		n.Processes[newProcess.Name] = newProcess
	}

}

func (n *Network) OperationSync() {

	n.Coordinator.SendVoteRequestMessages("synchronize", 0, "", nil)
	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
}

func (n *Network) OperationRemove(processName string) {
	n.Coordinator.SendVoteRequestMessages("remove", 0, processName, nil)
	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()

	if n.Coordinator.State == "commit" {
		delete(n.Processes, processName)

		if processName == n.Coordinator.Name {
			// choose a new coordinator
			for _, process := range n.Processes {
				n.Coordinator = process
				log.Infof("Updated cooordinator to: %v", n.Coordinator.Name)
				break
			}
		}
	}

}

func (n *Network) OperationSetTimeFailure(processName string, seconds int) {
	n.Processes[processName].TimeFailure = true

	timer := time.NewTimer(time.Duration(seconds) * time.Second)

	go func() {
		for {
			select {
			case <-timer.C:
				{
					log.Println("Time failure expired")
					n.Processes[processName].TimeFailure = false
				}
			}
		}
	}()
}

func (n *Network) OperationSetArbitraryFailure(processName string, seconds int) {
	n.Processes[processName].ArbitraryFailure = true

	timer := time.NewTimer(time.Duration(seconds) * time.Second)

	go func() {
		for {
			select {
			case <-timer.C:
				{
					log.Println("Arbitrary failure expired")
					n.Processes[processName].ArbitraryFailure = false
				}
			}
		}
	}()
}
