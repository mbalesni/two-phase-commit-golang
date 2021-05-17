package src

import (
	log "github.com/sirupsen/logrus"
)

func (p *Process) SendVoteRequestMessages(operation string, transactionValue int, processName string, process *Process) {

	for _, target := range p.OtherProcesses {
		if target.TimeFailure {
			log.Errorf("Node %v is unreachable. Aborting.", target.Name)
			p.State = "abort"
			return
		}
	}

	for _, target := range p.OtherProcesses {
		var message Message

		if operation == "synchronize" {
			message = p.NewFirstPhaseSynchronizationMessage(target, "VOTE-REQUEST", p.Log)
		} else if operation == "remove" {
			message = p.NewFirstPhaseRemoveMessage(target, "VOTE-REQUEST", processName)
		} else {
			message = p.NewFirstPhaseMessage(target, "VOTE-REQUEST", operation, transactionValue)
		}
		p.SendQueue.Add(message)
	}

	p.PreCommitCoordinator(operation, transactionValue, processName, process)

}

func (p *Process) SendGlobalCommitMessages() {

	p.Commit()

	for _, target := range p.OtherProcesses {
		message := p.NewSecondPhaseMessage(target, "GLOBAL-COMMIT")
		p.SendQueue.Add(message)
	}

}

func (p *Process) SendGlobalAbortMessages() {

	p.Abort()

	for _, target := range p.OtherProcesses {
		message := p.NewSecondPhaseMessage(target, "GLOBAL-ABORT")
		p.SendQueue.Add(message)
	}

}

// PreCommit when the coordinator sends VoteRequest to all, then PreCommit happens.
func (p *Process) PreCommitCoordinator(operation string, value int, processName string, process *Process) {

	p.State = "wait"
	p.UndoLog = p.Log
	p.UndoOtherProcesses = MapCopy(p.OtherProcesses)

	switch operation {
	case "set-value":
		{
			p.Log = append(p.Log, value)
		}
	case "rollback":
		{
			p.Log = p.Log[0:(len(p.Log) - value)]
		}
	case "synchronize":
		{
			// synchronization is done through the coordinator's log
			p.Log = p.Log
		}
	case "add":
		{
			p.OtherProcesses[processName] = process
		}
	case "remove":
		{
			delete(p.OtherProcesses, processName)
		}
	}

}

func (p *Process) AddDecision(decision string) {
	p.OtherProcessesDecisions = append(p.OtherProcessesDecisions, decision)
}
