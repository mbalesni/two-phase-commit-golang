package src

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
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
	UndoOtherProcesses      map[string]*Process
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
func (p *Process) PreCommit(operation string, value int, history []int, processName string, process *Process) string {

	// log.Infof("Pre")
	// log.Println(p.Name, "precommits. Op:", operation, "Value", value, "P name", processName)
	if p.ArbitraryFailure {
		p.State = "abort"
		return "VOTE-ABORT"
	}

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
			p.Log = history
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

	p.State = "ready"
	return "VOTE-COMMIT"
}

// Commit when commiting we do not rollback the changes and instead only delete the undo log.
func (p *Process) Commit() {
	p.State = "commit"
	p.UndoLog = nil
	p.UndoOtherProcesses = nil
}

// Abort when aborting we rollback the changes and go to the "abort" stage
func (p *Process) Abort() {
	p.State = "abort"
	p.Log = p.UndoLog
	p.UndoLog = nil
	if p.UndoOtherProcesses != nil {
		p.OtherProcesses = p.UndoOtherProcesses
		p.UndoOtherProcesses = nil
	}
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

				decision := p.PreCommit(message.Operation, message.TransactionValue, message.History, message.ProcessName, message.Process)
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
					// keep this weird if/else structure. it's magic.
					// otherwise everything breaks
				} else {
					p.AddDecision(message.MessageType)
					var commitDecisionCount int

					for _, val := range p.OtherProcessesDecisions {
						if val == "VOTE-COMMIT" {
							commitDecisionCount = commitDecisionCount + 1
						}
					}

					if commitDecisionCount == len(p.UndoOtherProcesses) {
						log.Infof("Operation commited successfully!")
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
					// fmt.Println(p.Name, "did not request an abort. Something is wrong!")
				} else {
					log.Errorln("Operation aborted by:", message.From.Name)
					p.OtherProcessesDecisions = []string{}
					p.SendGlobalAbortMessages()
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
					// fmt.Println(p.Name, " did not wait for a global commit. Something is wrong!")
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
					// fmt.Println(p.Name, " did not wait for a global abort. Something is wrong!")
				} else {
					p.Abort()
				}
			}
		}
	}
}
