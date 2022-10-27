package models

import (
	"encoding/json"
	"io"
	"time"

	"github.com/ONSdigital/dp-topic-api/apierrors"
)

// PrivateSubtopics used for returning both Next and Current document(s) in REST API response
type PrivateSubtopics struct {
	Count        int              `bson:"count,omitempty"        json:"count"`
	Limit        int              `bson:"limit,omitempty"        json:"limit"`
	Offset       int              `bson:"offset_index,omitempty" json:"offset_index"`
	TotalCount   int              `bson:"total_count,omitempty"  json:"total_count"`
	PrivateItems *[]TopicResponse `bson:"items,omitempty"        json:"items"`
}

// PublicSubtopics used for returning just the Current document(s) in REST API response
type PublicSubtopics struct {
	Count       int      `bson:"count,omitempty"        json:"count"`
	Limit       int      `bson:"limit,omitempty"        json:"limit"`
	Offset      int      `bson:"offset_index,omitempty" json:"offset_index"`
	TotalCount  int      `bson:"total_count,omitempty"  json:"total_count"`
	PublicItems *[]Topic `bson:"items,omitempty"        json:"items"`
}

// TopicResponse represents an evolving topic with the current topic and the updated topic
// The 'Next' is what gets updated throughout the publishing journey, and then the 'publish' step copies
// the 'Next' over the 'Current' document, so that 'Current' is whats always returned in the web view.
type TopicResponse struct {
	ID      string `bson:"id,omitempty"       json:"id,omitempty"`
	Current *Topic `bson:"current,omitempty"  json:"current,omitempty"`
	Next    *Topic `bson:"next,omitempty"     json:"next,omitempty"`
}

// Topic represents topic schema as it is stored in mongoDB
// and is used for marshaling and unmarshaling json representation for API
// ID is a duplicate of ID in TopicResponse, to facilitate each subdocument being a full-formed
// response in its own right depending upon request being in publish or web and also authentication.
// Subtopics contains TopicResonse ID(s).
type Topic struct {
	ID          string      `bson:"id,omitempty"             json:"id,omitempty"`
	Description string      `bson:"description,omitempty"    json:"description,omitempty"`
	Keywords    []string    `bson:"keywords,omitempty"       json:"keywords,omitempty"`
	Links       *TopicLinks `bson:"links,omitempty"          json:"links,omitempty"`
	ReleaseDate *time.Time  `bson:"release_date,omitempty"   json:"release_date,omitempty"`
	State       string      `bson:"state,omitempty"          json:"state,omitempty"`
	SubtopicIds []string    `bson:"subtopics_ids,omitempty"  json:"-"`
	Title       string      `bson:"title,omitempty"          json:"title,omitempty"`
}

// TopicRelease represents the incoming request structure containing release content
type TopicRelease struct {
	ReleaseDate string `json:"release_date"`
}

// LinkObject represents a generic structure for all links
type LinkObject struct {
	HRef string `bson:"href,omitempty"  json:"href,omitempty"`
	ID   string `bson:"id,omitempty"    json:"id,omitempty"`
}

// TopicLinks represents a list of specific links related to the topic resource
type TopicLinks struct {
	Content   *LinkObject `bson:"content,omitempty"    json:"content,omitempty"`
	Self      *LinkObject `bson:"self,omitempty"       json:"self,omitempty"`
	Subtopics *LinkObject `bson:"subtopics,omitempty"  json:"subtopics,omitempty"`
}

// CreateReleaseDate manages the creation of a release date object from a reader
func ReadReleaseDate(r io.Reader) (*TopicRelease, error) {
	var topicRelease TopicRelease

	err := json.NewDecoder(r).Decode(&topicRelease)
	switch {
	case err == io.EOF:
		return nil, apierrors.ErrEmptyRequestBody
	case err != nil:
		return nil, apierrors.ErrUnableToReadMessage
	}

	return &topicRelease, nil
}

// Validate checks that a topic struct complies with the state constraints, if provided. TODO may want to add more in future
func (t *Topic) Validate() error {

	if _, err := ParseState(t.State); err != nil {
		return apierrors.ErrTopicInvalidState
	}

	// TODO add other checks, etc
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

	return nil
}

// StateTransitionAllowed checks if the topic can transition from its current state to the provided target state
func (t *Topic) StateTransitionAllowed(target string) bool {
	currentState, err := ParseState(t.State)
	if err != nil {
		// TODO once the rest of the system is implemented, check that this logic is applicable, and adjust tests accordingly
		currentState = StateCreated // default value, if state is not present or invalid value
		// TODO more comments needed here to state under what conditions the state may not be present or has an invalid value
	}
	targetState, err := ParseState(target)
	if err != nil {
		// TODO once the rest of the system is implemented, check that this logic is applicable, and adjust tests accordingly
		// TODO to get to here is most likely a code programming error and a panic is probably best
		//     because i believe all state changes are explicity program code specified ...
		return false
	}
	return currentState.TransitionAllowed(targetState)
}

// Validate checks the topic release object has a valid timestamp that will
// abide by standard protocol RFC3339
func (tr *TopicRelease) Validate() (*time.Time, error) {
	releaseDate, err := time.Parse(time.RFC3339, tr.ReleaseDate)
	if err != nil {
		return nil, apierrors.ErrInvalidReleaseDate
	}

	return &releaseDate, nil
}
