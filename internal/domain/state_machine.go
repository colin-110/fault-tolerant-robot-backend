package domain

import "fmt"

var allowedTransitions = map[CommandState][]CommandState{
	StateCreated:  {StateSent},
	StateSent:     {StateAcked, StateFailed},
	StateAcked:    {StateExecuted, StateFailed},
	StateExecuted: {},
	StateFailed:   {},
}

func CanTransition(from, to CommandState) bool {
	for _, allowed := range allowedTransitions[from] {
		if allowed == to {
			return true
		}
	}
	return false
}

func Transition(cmd *Command, next CommandState) error {
	if !CanTransition(cmd.State, next) {
		return fmt.Errorf("invalid state transition: %s -> %s", cmd.State, next)
	}
	cmd.State = next
	return nil
}
