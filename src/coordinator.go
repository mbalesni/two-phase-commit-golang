package src

/*

func (p *Process) SendVoteRequestMessages(transactionValue int) {

	decision := "VOTE-REQUEST"

	p.State = "wait"
	for _, target := range p.OtherProcesses {
		if target.TimeFailure != true {
			message := p.NewFirstPhaseMessage(target, decision, transactionValue)
			p.SendQueue.Add(message)
		}
	}
}

func (p *Process) SendGlobalCommitMessages() {

	decision := "GLOBAL-COMMIT"

	p.State = "wait"
	for _, target := range p.OtherProcesses {
		if target.TimeFailure != true {
			message := p.NewSecondPhaseMessage(target, decision)
			p.SendQueue.Add(message)
		}
	}

}

func (p *Process) AddDecision(decision string) {
	p.OtherProcessesDecisions = append(p.OtherProcessesDecisions, decision)
}

 */