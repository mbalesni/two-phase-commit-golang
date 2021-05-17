package src

import log "github.com/sirupsen/logrus"

func (p *Process) SendVoteRequestMessages(operation string, transactionValue int, processName string) {

	for _, target := range p.OtherProcesses {
		if target.TimeFailure {
			log.Errorf("Node %v is unreachable. Aborting.", target.Name)
			return
		}
	}

	p.PreCommitCoordinator(operation, transactionValue, processName)

	for _, target := range p.OtherProcesses {
		var message Message

		if operation == "synchronize" {
			message = p.NewFirstPhaseSynchronizationMessage(target, "VOTE-REQUEST", p.Log)
		} else if operation == "remove" && target.Name != processName {
			message = p.NewFirstPhaseRemoveMessage(target, "VOTE-REQUEST", processName)
		} else {
			message = p.NewFirstPhaseMessage(target, "VOTE-REQUEST", operation, transactionValue)
		}
		p.SendQueue.Add(message)
	}

	p.State = "wait"

}

func (p *Process) SendGlobalCommitMessages() {

	p.State = "init"
	p.Commit()

	for _, target := range p.OtherProcesses {
		message := p.NewSecondPhaseMessage(target, "GLOBAL-COMMIT")
		p.SendQueue.Add(message)
	}

}

func (p *Process) SendGlobalAbortMessages() {

	p.State = "init"
	p.Abort()

	for _, target := range p.OtherProcesses {
		message := p.NewSecondPhaseMessage(target, "GLOBAL-ABORT")
		p.SendQueue.Add(message)
	}

}

// PreCommit when the coordinator sends VoteRequest to all, then PreCommit happens.
func (p *Process) PreCommitCoordinator(operation string, value int, processName string) {

	p.UndoLog = p.Log
	p.UndoOtherProcesses = p.OtherProcesses

	switch operation {
	case "add":
		{
			p.Log = append(p.Log, value)
		}
	case "rollback":
		{
			p.Log = p.Log[0:(len(p.Log) - value)]
		}
	case "synchronize":
		{
			p.Log = p.Log
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
