package models

import "github.com/ONSdigital/dp-topic-api/apierrors"

// Topic represents topic schema as it is stored in mongoDB
// and is used for marshaling and unmarshaling json representation for API
type Topic struct {
	ID          string   `bson:"_id,omitempty"          json:"id,omitempty"`
	Description string   `bson:"description,omitempty"  json:"description,omitempty"`
	Title       string   `bson:"title,omitempty"        json:"title,omitempty"`
	Keywords    []string `bson:"keywords,omitempty"     json:"keywords,omitempty"`
	State       string   `bson:"state,omitempty"        json:"state,omitempty"`
}

// Validate checks that a topic struct complies with the state constraints, if provided. !!! may want to add more in future
func (t *Topic) Validate() error {

	if _, err := ParseState(t.State); err != nil {
		return apierrors.ErrTopicInvalidState
	}

	// !!! add other checks, etc
	return nil
}

// ValidateTransitionFrom checks that this topic state can be validly transitioned from the existing state
func (t *Topic) ValidateTransitionFrom(existing *Topic) error {

	// check that state transition is allowed, only if state is provided
	if t.State != "" {
		if !existing.StateTransitionAllowed(t.State) {
			return apierrors.ErrTopicStateTransitionNotAllowed
		}
	}

	// if the topic is already completed, it cannot be updated
	if existing.State == StateCompleted.String() { //!!! ultimately this might not be needed
		return apierrors.ErrTopicAlreadyCompleted
	}

	return nil
}

// StateTransitionAllowed checks if the topic can transition from its current state to the provided target state
func (t *Topic) StateTransitionAllowed(target string) bool {
	currentState, err := ParseState(t.State)
	if err != nil {
		//!!! once the rest of the system is implemented, check that this logic is applicable, and adjust tests accordingly
		currentState = StateCreated // default value, if state is not present or invalid value
	}
	targetState, err := ParseState(target)
	if err != nil {
		return false
	}
	return currentState.TransitionAllowed(targetState)
}
