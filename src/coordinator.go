package src

import "fmt"

func (p *Process) SendVoteRequestMessages(operation string, transactionValue int) {

	p.PreCommitCoordinator(operation, transactionValue)

	// hacky solution:
	// don't send any messages if some process has TimeFailure.
	// We can only do it this way
	// because "processes" are simulated.
	for _, target := range p.OtherProcesses {
		if target.TimeFailure {
			fmt.Println("Node", target.Name, "is unreachable. Aborting.")
			return
		}
	}

	for _, target := range p.OtherProcesses {
		message := p.NewFirstPhaseMessage(target, "VOTE-REQUEST", operation, transactionValue)
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

// PreCommit when the coordinator sends VoteRequest to all, then PreCommit happens.
func (p *Process) PreCommitCoordinator(operation string, value int) {

	p.UndoLog = p.Log

	switch operation {
	case "add":
		{
			p.Log = append(p.Log, value)
		}
	case "rollback":
		{
			p.Log = p.Log[0:(len(p.Log) - value)]
		}
	}

}

func (p *Process) AddDecision(decision string) {
	p.OtherProcessesDecisions = append(p.OtherProcessesDecisions, decision)
}
