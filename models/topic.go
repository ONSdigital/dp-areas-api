package models

import "github.com/ONSdigital/dp-topic-api/apierrors"

// TopicResponse represents an evolving topic with the current topic and the updated topic
// The 'Next' is what gets updated throughout the publishing journey, and then the 'publish' step copies
// the 'Next' over the 'Current' document, so that 'Current' is whats always returned in the web view.
type TopicResponse struct {
	ID      string `bson:"id,omitempty"       json:"id,omitempty"`
	Next    *Topic `bson:"next,omitempty"     json:"next,omitempty"`
	Current *Topic `bson:"current,omitempty"  json:"current,omitempty"`
}

// Topic represents topic schema as it is stored in mongoDB
// and is used for marshaling and unmarshaling json representation for API
// ID is a duplicate of ID in TopicResponse, to facilitate each subdocument being a full-formed
// response in its own right depending upon request being in publish or web and also authentication.
// Subtopics contains TopicResonse ID(s).
type Topic struct {
	ID          string      `bson:"_id,omitempty"            json:"id,omitempty"`
	Description string      `bson:"description,omitempty"    json:"description,omitempty"`
	Title       string      `bson:"title,omitempty"          json:"title,omitempty"`
	Keywords    []string    `bson:"keywords,omitempty"       json:"keywords,omitempty"`
	State       string      `bson:"state,omitempty"          json:"state,omitempty"`
	Links       *TopicLinks `bson:"links,omitempty"          json:"links,omitempty"`
	SubtopicIds []string    `bson:"subtopics_ids,omitempty"  json:"-"`
}

// LinkObject represents a generic structure for all links
type LinkObject struct {
	HRef string `bson:"href,omitempty"  json:"href,omitempty"`
	ID   string `bson:"id,omitempty"    json:"id,omitempty"`
}

// TopicLinks represents a list of specific links related to the topic resource
type TopicLinks struct {
	Self      *LinkObject `bson:"self,omitempty"       json:"self,omitempty"`
	Subtopics *LinkObject `bson:"subtopics,omitempty"  json:"subtopics,omitempty"`
	Content   *LinkObject `bson:"content,omitempty"    json:"content,omitempty"`
}

// PublicSubtopics used for returning just the Current document(s) in REST API response
type PublicSubtopics struct {
	Count       int      `bson:"count,omitempty"        json:"count"`
	Offset      int      `bson:"offset_index,omitempty" json:"offset_index"`
	Limit       int      `bson:"limit,omitempty"        json:"limit"`
	TotalCount  int      `bson:"total_count,omitempty"  json:"total_count"`
	PublicItems *[]Topic `bson:"items,omitempty"        json:"items"`
}

// PrivateSubtopics used for returning both Next and Current document(s) in REST API response
type PrivateSubtopics struct {
	Count        int              `bson:"count,omitempty"        json:"count"`
	Offset       int              `bson:"offset_index,omitempty" json:"offset_index"`
	Limit        int              `bson:"limit,omitempty"        json:"limit"`
	TotalCount   int              `bson:"total_count,omitempty"  json:"total_count"`
	PrivateItems *[]TopicResponse `bson:"items,omitempty"        json:"items"`
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
	//	if existing.State == StateCompleted.String() { //!!! ultimately this might not be needed
	//		return apierrors.ErrTopicAlreadyCompleted
	//	}

	return nil
}

// StateTransitionAllowed checks if the topic can transition from its current state to the provided target state
func (t *Topic) StateTransitionAllowed(target string) bool {
	currentState, err := ParseState(t.State)
	if err != nil {
		//!!! once the rest of the system is implemented, check that this logic is applicable, and adjust tests accordingly
		currentState = StateCreated // default value, if state is not present or invalid value
		// !!! more comments needed here to state under what conditions the state may not be present or has an invalid value
	}
	targetState, err := ParseState(target)
	if err != nil {
		//!!! once the rest of the system is implemented, check that this logic is applicable, and adjust tests accordingly
		// !!! to get to here is most likely a code programming error and a panic is probably best
		//     because i believe all state changes are explicity program code specified ...
		return false
	}
	return currentState.TransitionAllowed(targetState)
}
