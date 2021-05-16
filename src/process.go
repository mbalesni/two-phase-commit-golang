package src

import (
	"fmt"
	"sync"
)

type Process struct {
	Name                    string
	State                   string // one of init|wait|ready
	Verbose                 bool
	mu                      sync.Mutex
	ArbitraryFailure        bool
	TimeFailure             bool
	GetQueue                MessageQueue
	SendQueue               MessageQueue
	Log                     []int
	UndoLog                 []int
	OtherProcesses          map[string]*Process
	OtherProcessesDecisions []string
}

func NewProcess(name string, verbose bool) *Process {
	process := Process{Name: name, Verbose: verbose}
	process.Init()
	return &process
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

// PreCommit when the coordinator sends VoteRequest to all, then PreCommit happens.
func (p *Process) PreCommit(operation string, value int) string {
	if p.ArbitraryFailure {
		p.State = "init"
		return "VOTE-ABORT"
	}

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

	p.State = "ready"
	return "VOTE-COMMIT"
}

// Commit when commiting we do not rollback the changes and instead only delete the undo log.
func (p *Process) Commit() {
	p.UndoLog = nil
	p.State = "init"
}

// Abort when aborting we rollback the changes and go back to the "init" stage
func (p *Process) Abort() {
	p.State = "init"
	p.Log = p.UndoLog // TODO: maybe also set Undo Log to nil
}

func (p *Process) ProcessMessages() {
	for p.GetQueue.queue.Len() > 0 {
		message := p.GetQueue.Pop()
		switch messageType := message.MessageType; messageType {

		// First phase messages
		// Only participants receive
		case "VOTE-REQUEST":
			{
				if p.Verbose {
					fmt.Println("P=", p.Name, "got a VOTE-REQUEST from", message.From.Name)
				}
				decision := p.PreCommit(message.Operation, message.TransactionValue)
				responseMessage := p.NewMessage(message.From, decision)
				p.SendQueue.Add(responseMessage)
			}
		// Only coordinator receives
		case "VOTE-COMMIT":
			{
				if p.Verbose {
					fmt.Println(p.Name, "got a VOTE-COMMIT from", message.From.Name)
				}
				if p.State != "wait" {
					fmt.Println(p.Name, "did not request a commit. Something is wrong! State:", p.State)
				} else {
					p.AddDecision(message.MessageType)
					// TODO: fix by comparing # responses to actual # of requests instead of total # of processes
					// or maybe actually keep it this way. Specification doesn't mention this.
					if len(p.OtherProcessesDecisions) == len(p.OtherProcesses) {
						p.State = "commit"
						p.OtherProcessesDecisions = []string{}
						p.SendGlobalCommitMessages()
					}
				}
			}
		// Only coordinator receives
		case "VOTE-ABORT":
			{
				if p.Verbose {
					fmt.Println(p.Name, "got a VOTE-ABORT from", message.From.Name)
				}
				if p.State != "wait" {
					fmt.Println(p.Name, "did not request an abort. Something is wrong!")
				} else {
					p.AddDecision(message.MessageType)
				}
			}
		// Second phase messages
		// Only participants receive
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
		// Only participants receive
		case "GLOBAL-ABORT":
			{
				if p.Verbose {
					fmt.Println(p.Name, "got a GLOBAL-ABORT from", message.From.Name)
				}
				if p.State != "ready" {
					fmt.Println(p.Name, " did not wait for a global abort. Something is wrong!")
				} else {
					p.Abort()
				}
			}

			// // From Coordinator to Participants
			// case "SYNCHRONIZE-REQUEST":
			// 	{
			// 		if p.Verbose {
			// 			fmt.Println(p.Name, "got a SYNCHRONIZE-REQUEST from", message.From.Name)
			// 		}
			// 		responseMessage := p.NewRecoveryMessage(message.From, "SYNCHRONIZE-RESPONSE", p.NextCommitValue, p.History)
			// 		p.SendQueue.Add(responseMessage)
			// 	}
			// // From Participants to Coordinator
			// case "SYNCHRONIZE-RESPONSE":
			// 	{
			// 		if p.Verbose {
			// 			fmt.Println(p.Name, "got a RECOVERY-RESPONSE from", message.From.Name)
			// 		}
			// 		p.NextCommitValue = message.TransactionValue
			// 		p.History = message.History
			// 	}
			// // From Coordinator to Participants
			// case "GLOBAL-SYNCHRONIZE":
			// 	{

			// 	}

		}
	}
}
