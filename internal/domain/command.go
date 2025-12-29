package domain

type CommandState string

const (
	StateCreated  CommandState = "CREATED"
	StateSent     CommandState = "SENT"
	StateAcked    CommandState = "ACKED"
	StateExecuted CommandState = "EXECUTED"
	StateFailed   CommandState = "FAILED"
)

type Command struct {
	CommandID string
	RobotID   string
	State     CommandState
}
