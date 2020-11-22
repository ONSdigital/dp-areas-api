package models

//!!! this will all need a re-work for topic-api

import (
	"github.com/ONSdigital/dp-topic-api/apierrors"
)

// State - iota enum of possible topic states
type State int

// Possible values for a State of a topic. It can only be one of the following:
const (
	StateTopicCreated State = iota
	StateTopicUploaded
	StateTopicImporting
	StateTopicImported
	StateTopicPublished
	StateTopicCompleted
	StateTopicDeleted
	StateTopicFailedImport
	StateTopicFailedPublish
	StateTopicTrue  // to be removed at some point
	StateTopicFalse // to be removed at some point
)

type stateTransition struct {
	name             string
	validTransitions []State
}

// This List MUST be in the same order as the 'iota' list above that starts with 'StateTopicCreated' for
// the the following state transition table and code to work and be efficient
var stateTransitionTable = []stateTransition{
	{ // StateTopicCreated
		name:             "created",
		validTransitions: []State{StateTopicUploaded, StateTopicDeleted}},
	{ // StateTopicUploaded
		name:             "uploaded",
		validTransitions: []State{StateTopicImporting, StateTopicDeleted}},
	{ // StateTopicImporting
		name:             "importing",
		validTransitions: []State{StateTopicImported, StateTopicFailedImport, StateTopicDeleted}},
	{ // StateTopicImported
		name:             "imported",
		validTransitions: []State{StateTopicPublished, StateTopicDeleted}},
	{ // StateTopicPublished
		name:             "published",
		validTransitions: []State{StateTopicCompleted, StateTopicFailedPublish, StateTopicDeleted}},
	{ // StateTopicCompleted
		name:             "completed",
		validTransitions: []State{StateTopicDeleted}},
	{ // StateTopicDeleted
		name:             "deleted",
		validTransitions: []State{}},
	{ // StateTopicFailedImport
		name:             "failed_import",
		validTransitions: []State{StateTopicDeleted}},
	{ // StateTopicFailedPublish
		name:             "failed_import",
		validTransitions: []State{StateTopicDeleted}},
	{ // StateTopicTrue // to be removed at some point
		name:             "true",
		validTransitions: []State{StateTopicFalse}},
	{ // StateTopicFalse // to be removed at some point
		name:             "false",
		validTransitions: []State{StateTopicTrue}},
}

// String returns the string representation of a state
func (s State) String() string {
	return stateTransitionTable[s].name
}

// ParseState returns a state from its string representation
func ParseState(stateStr string) (State, error) {
	for i, aState := range stateTransitionTable {
		if stateStr == aState.name {
			return State(i), nil
		}
	}
	return -1, apierrors.ErrTopicInvalidState
}

// TransitionAllowed returns true only if the transition from the current state and the provided next is allowed
func (s State) TransitionAllowed(next State) bool {
	if (s >= 0) && (int(s) < len(stateTransitionTable)) {
		for _, allowedState := range stateTransitionTable[s].validTransitions {
			if next == allowedState {
				return true
			}
		}
	}

	return false
}
