package src

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func (n *Network) OperationSetValue(value int) {
	n.Coordinator.SendVoteRequestMessages("add", value, "")
	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
}

func (n *Network) OperationRollback(steps int) {
	n.Coordinator.SendVoteRequestMessages("rollback", steps, "")
	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
}

func (n *Network) OperationAdd(processName string) {

	newProcess := NewProcess(processName, false)

	n.Coordinator.SendVoteRequestMessages("synchronize", 0, processName)
	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()

	n.Processes[newProcess.Name] = newProcess

	n.AutoDiscovery()
}

func (n *Network) OperationRemove(processName string) {
	n.Coordinator.SendVoteRequestMessages("remove", 0, processName)
	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()

	if n.Coordinator.UndoOtherProcesses[processName] == nil {
		delete(n.Processes, processName)
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

