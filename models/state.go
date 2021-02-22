package models

import (
	"github.com/ONSdigital/dp-topic-api/apierrors"
)

// 7th December 2020 ... What the states are as detailed in slack message from Eleanor:
// NOTE: this comment to be eventually deleted ...
//
// possible states and transitions ...
//
// Created -> Completed
// Completed -> Published
//
// Within the next subdocument, a topic may go from published to created when it is next
// edited after being published.
//
// At the point a topic is published, both next and current should equal each other exactly,
// including having the state 'published' so when we next want to make a change to it it
// needs to go to 'created' state while being worked on.
//
// !!! (my note), so from last sentence, add a state change of: Published -> Created
//
// Deleted and Failed I dont think will be needed (we wont want to allow people to delete pages,
// and failed publishes in this area I would expect not to be reported through a state change.

// State - iota enum of possible topic states
type State int

// Possible values for a State of a topic. It can only be one of the following:
const (
	// these from dp-image-api :
	StateCreated State = iota // this is 'in_progress'
	StatePublished
	StateCompleted
)

type stateTransition struct {
	state            State
	name             string
	validTransitions []State
}

var stateTransitionTable = []stateTransition{
	{
		state:            StateCreated, // this is 'in_progress'
		name:             "created",
		validTransitions: []State{StateCompleted}},
	{
		state:            StatePublished,
		name:             "published",
		validTransitions: []State{StateCreated}},
	{
		state:            StateCompleted,
		name:             "completed",
		validTransitions: []State{StatePublished}},
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
		// (to allow for states in table not being in the same order as 'StateCreated' iota list)
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
