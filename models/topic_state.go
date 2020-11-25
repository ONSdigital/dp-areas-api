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
	state            State
	name             string
	validTransitions []State
}

var stateTransitionTable = []stateTransition{
	{
		state:            StateTopicCreated,
		name:             "created",
		validTransitions: []State{StateTopicUploaded, StateTopicDeleted}},
	{
		state:            StateTopicUploaded,
		name:             "uploaded",
		validTransitions: []State{StateTopicImporting, StateTopicDeleted}},
	{
		state:            StateTopicImporting,
		name:             "importing",
		validTransitions: []State{StateTopicImported, StateTopicFailedImport, StateTopicDeleted}},
	{
		state:            StateTopicImported,
		name:             "imported",
		validTransitions: []State{StateTopicPublished, StateTopicDeleted}},
	{
		state:            StateTopicPublished,
		name:             "published",
		validTransitions: []State{StateTopicCompleted, StateTopicFailedPublish, StateTopicDeleted}},
	{
		state:            StateTopicCompleted,
		name:             "completed",
		validTransitions: []State{StateTopicDeleted}},
	{
		state:            StateTopicDeleted,
		name:             "deleted",
		validTransitions: []State{}},
	{
		state:            StateTopicFailedImport,
		name:             "failed_import",
		validTransitions: []State{StateTopicDeleted}},
	{
		state:            StateTopicFailedPublish,
		name:             "failed_import",
		validTransitions: []State{StateTopicDeleted}},
	{
		state:            StateTopicTrue, // to be removed at some point
		name:             "true",
		validTransitions: []State{StateTopicFalse}},
	{
		state:            StateTopicFalse, // to be removed at some point
		name:             "false",
		validTransitions: []State{StateTopicTrue}},
}

// String returns the string representation of a state
func (s State) String() string {
	var name string = ""

	// search for state in table to find name
	for _, aState := range stateTransitionTable {
		if s == aState.state {
			name = aState.name
		}
	}
	return name
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
		// search for state in table
		// (to allow for states in table not being in the same order as 'StateTopicCreated' iota list)
		for _, transition := range stateTransitionTable {
			if s == transition.state {
				// see if transition allowed
				for _, allowedState := range transition.validTransitions {
					if next == allowedState {
						return true
					}
				}
			}
		}
	}

	return false
}
