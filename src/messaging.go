package src

type Message struct {
	From             *Process
	To               *Process
	MessageType      string // one of VOTE-REQUEST|VOTE-COMMIT|VOTE-ABORT|GLOBAL-COMMIT|GLOBAL-ABORT|
	Operation 		 string
	TransactionValue int
	History          []int
}

func (p *Process) NewMessage(to *Process, messageType string) Message {
	return Message{From: p, To: to, MessageType: messageType}
}

func (p *Process) NewFirstPhaseMessage(to *Process, messageType string, operation string, transactionValue int) Message {
	return Message{From: p, To: to, MessageType: messageType, Operation: operation, TransactionValue: transactionValue, History: nil}
}

func (p *Process) NewFirstPhaseSynchronizationMessage(to *Process, messageType string, history []int) Message {
	return Message{From: p, To: to, MessageType: messageType, Operation: "synchronize", History: history}
}

func (p *Process) NewSecondPhaseMessage(to *Process, messageType string) Message {
	return Message{From: p, To: to, MessageType: messageType}
}
/*
func (p *Process) NewRecoveryMessage(to *Process, messageType string, transactionValue int, history []int) Message {
	return Message{From: p, To: to, MessageType: messageType, TransactionValue: transactionValue, History: history}
}
*/