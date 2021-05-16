package src

func (n *Network) OperationSetValue(value int) {
	n.Coordinator.SendVoteRequestMessages("add", value)
	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
}

func (n *Network) OperationRollback(steps int) {
	n.Coordinator.SendVoteRequestMessages("rollback", steps)
	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
}

func (n *Network) OperationAdd(processName string) {

	newProcess := NewProcess(processName, false)

	n.Processes[newProcess.Name] = newProcess

	n.AutoDiscovery()

	n.Coordinator.SendVoteRequestMessages("synchronize", 0)

	// Send VOTE REQUEST
	n.Cycle()
	// Reason about whether anyone sent VOTE-COMMIT or VOTE-ABORT and send GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()
	// Consume GLOBAL-COMMIT or GLOBAL-ABORT
	n.Cycle()

}

