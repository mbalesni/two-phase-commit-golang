package src

import (
	"fmt"
	"sync"
)

type Message struct {
	From             *Process
	To               *Process
	MessageType      string // one of VOTE-REQUEST|VOTE-COMMIT|VOTE-ABORT|GLOBAL-COMMIT|GLOBAL-ABORT|
	TransactionValue int
	History          []int
}

type Process struct {
	ArbitraryFailure        bool
	GetQueue                MessageQueue
	History                 []int
	Name                    string
	NextHistory             []int
	OtherProcesses          map[string]*Process
	OtherProcessesDecisions []string
	SendQueue               MessageQueue
	State                   string // one of init|wait|ready
	TimeFailure             bool
	Verbose                 bool
	mu                      sync.Mutex
	// WaitingElection    int
	// WaitingCoordinator int
	// MaxElectionWait    int
	// MaxCoordinatorWait int
}

func NewProcess(name string, verbose bool) *Process {

	process := Process{Name: name, Verbose: verbose}
	process.Init()
	return &process
}

func (p *Process) InitCommit(transactionValue int) {
	p.State = "wait"
	p.NextCommitValue = transactionValue // TODO: replace "nextcommitvalue" logic with "nexthistory" logic

	// TODO: instead of sending messages immediately
	// let's check if any process has TimeFailure.
	// If so, just don't initiate the commit and fail immediately.
	for _, target := range p.OtherProcesses {
		if target.TimeFailure != true {
			message := p.NewVoteMessage(target, "VOTE-REQUEST", transactionValue)
			p.SendQueue.Add(message)
		}
	}
}

func (p *Process) Commit() {
	p.History = append(p.History, p.NextCommitValue)
	p.State = "init"
}

func (p *Process) Abort() {
	p.NextCommitValue = -1
	p.State = "init"
}

func (p *Process) NewMessage(to *Process, messageType string) Message {
	return Message{From: p, To: to, MessageType: messageType}
}

func (p *Process) NewVoteMessage(to *Process, messageType string, transactionValue int) Message {
	return Message{From: p, To: to, MessageType: messageType, TransactionValue: transactionValue}
}

func (p *Process) NewRecoveryMessage(to *Process, messageType string, transactionValue int, history []int) Message {
	return Message{From: p, To: to, MessageType: messageType, TransactionValue: transactionValue, History: history}
}

func (p *Process) AddDecision(decision string) {
	p.OtherProcessesDecisions = append(p.OtherProcessesDecisions, decision)
}

func (p *Process) RunGlobalCommit() {
	decision := "GLOBAL-COMMIT"

	for _, target := range p.OtherProcesses {
		if target.TimeFailure != true {
			message := p.NewMessage(target, decision)
			p.SendQueue.Add(message)
		}
	}

	p.Commit()
}

func (p *Process) RunGlobalAbort() {
	decision := "GLOBAL-ABORT"

	for _, target := range p.OtherProcesses {
		if target.TimeFailure != true {
			message := p.NewMessage(target, decision)
			p.SendQueue.Add(message)
		}
	}

	p.Abort()
}

func (p *Process) RunRecovery() {
	messageType := "RECOVERY-REQUEST"

	// only sends to 1 other process
	for _, target := range p.OtherProcesses {
		if target.TimeFailure != true {
			message := p.NewMessage(target, messageType)
			p.SendQueue.Add(message)
			break // breaks after first message
		}
	}
}

func (p *Process) PreCommit(transactionValue int) string {
	if p.ArbitraryFailure {
		p.State = "init"
		return "VOTE-ABORT"
	}
	// in the real world, this would run the operation
	// but will not commit it.
	// if it succeeds, it answers OK
	// else it answers ABORT.
	// If it later gets a global commit
	// it finishes the operation and commits it.
	p.State = "ready"
	p.NextCommitValue = transactionValue
	return "VOTE-COMMIT"
}

func (p *Process) ProcessMessages() {
	for p.GetQueue.queue.Len() > 0 {
		message := p.GetQueue.Pop()
		switch messageType := message.MessageType; messageType {

		case "VOTE-REQUEST":
			{
				if p.Verbose {
					fmt.Println("P=", p.Name, "got a VOTE-REQUEST from", message.From.Name)
				}
				decision := p.PreCommit(message.TransactionValue)
				responseMessage := p.NewMessage(message.From, decision)
				p.SendQueue.Add(responseMessage)
			}
		case "VOTE-COMMIT":
			{
				if p.Verbose {
					fmt.Println(p.Name, "got a VOTE-COMMIT from", message.From.Name)
				}
				if p.State != "wait" {
					fmt.Println(p.Name, " did not request a commit. Something is wrong!")
				} else {
					p.AddDecision(message.MessageType)
					// TODO: fix by comparing # responses to actual # of requests instead of total # of processes
					if len(p.OtherProcessesDecisions) == len(p.OtherProcesses) {
						p.OtherProcessesDecisions = []string{}
						p.RunGlobalCommit()
					}
				}

			}
		case "VOTE-ABORT":
			{
				if p.Verbose {
					fmt.Println(p.Name, "got a VOTE-ABORT from", message.From.Name)
				}
				if p.State != "wait" {
					fmt.Println(p.Name, " did not request a commit. Something is wrong!")
				} else {
					p.OtherProcessesDecisions = []string{}
					p.RunGlobalAbort()
				}

			}
		case "GLOBAL-COMMIT":
			{
				if p.Verbose {
					fmt.Println(p.Name, "got a GLOBAL-COMMIT from", message.From.Name)
				}
				if p.State != "ready" {
					fmt.Println(p.Name, " did not wait for a global commit. Something is wrong!")
				} else {
					p.Commit()
				}
			}
		case "GLOBAL-ABORT":
			{
				if p.Verbose {
					fmt.Println(p.Name, "got a GLOBAL-COMMIT from", message.From.Name)
				}
				if p.State != "ready" {
					fmt.Println(p.Name, " did not wait for a global abort. Something is wrong!")
				} else {
					p.Abort()
				}
			}
		case "RECOVERY-REQUEST":
			{
				if p.Verbose {
					fmt.Println(p.Name, "got a RECOVERY-REQUEST from", message.From.Name)
				}
				responseMessage := p.NewRecoveryMessage(message.From, "RECOVER-RESPONSE", p.NextCommitValue, p.History)
				p.SendQueue.Add(responseMessage)
			}
		case "RECOVERY-RESPONSE":
			{
				if p.Verbose {
					fmt.Println(p.Name, "got a RECOVERY-RESPONSE from", message.From.Name)
				}
				p.NextCommitValue = message.TransactionValue
				p.History = message.History
			}
		}

	}
}

func (p *Process) Cycle() {
	p.ProcessMessages()
}

func (p *Process) Init() {
	p.SendQueue = MessageQueue{}
	p.GetQueue = MessageQueue{}
	p.OtherProcesses = make(map[string]*Process)
	p.OtherProcessesDecisions = []string{}
	p.State = "init"
	p.SendQueue.Init()
	p.GetQueue.Init()
}
