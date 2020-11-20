package models

//!!! this will all need a re-work for topic-api

import "github.com/ONSdigital/dp-topic-api/apierrors"

// State - iota enum of possible topic states
type State int

// Possible values for a State of a topic. It can only be one of the following:
const (
	StateCreated State = iota
	StateUploaded
	StateImporting
	StateImported
	StatePublished
	StateCompleted
	StateDeleted
	StateFailedImport
	StateFailedPublish
	StateTrue  // to be removed at some point
	StateFalse // to be removed at some point
)

var stateValues = []string{"created", "uploaded", "importing", "imported", "published", "completed", "deleted", "failed_import", "failed_publish", "true", "false"}

// String returns the string representation of a state
func (s State) String() string {
	return stateValues[s]
}

// ParseState returns a state from its string representation
func ParseState(stateStr string) (State, error) {
	for s, validState := range stateValues {
		if stateStr == validState {
			return State(s), nil
		}
	}
	return -1, apierrors.ErrTopicInvalidState
}

// TransitionAllowed returns true only if the transition from the current state and the provided targetState is allowed
func (s State) TransitionAllowed(target State) bool {
	switch s {
	case StateCreated:
		switch target {
		case StateUploaded, StateDeleted:
			return true
		default:
			return false
		}
	case StateUploaded:
		switch target {
		case StateImporting, StateDeleted:
			return true
		default:
			return false
		}
	case StateImporting:
		switch target {
		case StateImported, StateFailedImport, StateDeleted:
			return true
		default:
			return false
		}
	case StateImported:
		switch target {
		case StatePublished, StateDeleted:
			return true
		default:
			return false
		}
	case StatePublished:
		switch target {
		case StateCompleted, StateFailedPublish, StateDeleted:
			return true
		default:
			return false
		}
	case StateCompleted:
		switch target {
		case StateDeleted:
			return true
		default:
			return false
		}
	case StateFailedImport:
		switch target {
		case StateDeleted:
			return true
		default:
			return false
		}
	case StateFailedPublish:
		switch target {
		case StateDeleted:
			return true
		default:
			return false
		}
	default:
		return false
	}
}
