package models

import (
	"sort"

	"github.com/ONSdigital/dp-topic-api/apierrors"
)

// ContentResponse represents an evolving content with the current content and the updated content.
// This is for mongo storage / retrieval.
// The 'Next' is what gets updated throughout the publishing journey, and then the 'publish' step copies
// the 'Next' over the 'Current' document, so that 'Current' is whats always returned in the web view.
// ID is a duplicate of ID in TopicResponseStore.
type ContentResponse struct {
	ID      string   `bson:"id,omitempty"       json:"id,omitempty"`
	Next    *Content `bson:"next,omitempty"     json:"next,omitempty"`
	Current *Content `bson:"current,omitempty"  json:"current,omitempty"`
}

// Content represents content schema as it is stored in mongoDB
// and is used for marshaling and unmarshaling json representation for API
// ID is a duplicate of ID in TopicResponse, to facilitate each subdocument being a full-formed
// response in its own right depending upon request being in publish or web and also authentication.
type Content struct {
	State               string            `bson:"state,omitempty"                 json:"state,omitempty"`
	Spotlight           *[]TypeLinkObject `bson:"spotlight,omitempty"             json:"spotlight,omitempty"`
	Articles            *[]TypeLinkObject `bson:"articles,omitempty"              json:"articles,omitempty"`
	Bulletins           *[]TypeLinkObject `bson:"bulletins,omitempty"             json:"bulletins,omitempty"`
	Methodologies       *[]TypeLinkObject `bson:"methodologies,omitempty"         json:"methodologies,omitempty"`
	MethodologyArticles *[]TypeLinkObject `bson:"methodology_articles,omitempty"  json:"methodology_articles,omitempty"`
	StaticDatasets      *[]TypeLinkObject `bson:"static_datasets,omitempty"       json:"static_datasets,omitempty"`
	Timeseries          *[]TypeLinkObject `bson:"timeseries,omitempty"            json:"timeseries,omitempty"`
}

// TypeLinkObject represents a generic structure for all type links
type TypeLinkObject struct {
	HRef  string `bson:"href,omitempty"   json:"href,omitempty"`
	Title string `bson:"title,omitempty"  json:"title,omitempty"`
}

// PrivateContentResponseAPI represents an evolving content with the current content and the updated content.
// This is for the REST API response.
// The 'Next' is what gets updated throughout the publishing journey, and then the 'publish' step copies
// the 'Next' over the 'Current' document, so that 'Current' is whats always returned in the web view.
type PrivateContentResponseAPI struct {
	Next    *ContentResponseAPI `json:"next,omitempty"`
	Current *ContentResponseAPI `json:"current,omitempty"`
}

// ContentResponseAPI used for returning the Current OR Next & Current document(s) in REST API response
type ContentResponseAPI struct {
	Count      int            `json:"count"`
	Offset     int            `json:"offset_index"`
	Limit      int            `json:"limit"`
	TotalCount int            `json:"total_count"`
	Items      *[]ContentItem `json:"items"`
}

// ContentItem is an individual content item
type ContentItem struct {
	Title string        `json:"title,omitempty"`
	Type  string        `json:"type,omitempty"`
	Links *ContentLinks `json:"links,omitempty"`
	State string        `json:"state,omitempty"`
}

// ContentLinks are content links
type ContentLinks struct {
	Self  *LinkObject `json:"self,omitempty"`
	Topic *LinkObject `json:"topic,omitempty"`
}

// !!! add code to validate state transitions as per topic.go
//!!! fix the following, and sort test code elsewhere, as per topic

// Validate checks that a content struct complies with the state constraints, if provided. !!! may want to add more in future
func (t *Content) Validate() error {

	if _, err := ParseState(t.State); err != nil {
		return apierrors.ErrTopicInvalidState
	}

	// !!! add other checks, etc
	return nil
}

// ValidateTransitionFrom checks that this content state can be validly transitioned from the existing state
func (t *Content) ValidateTransitionFrom(existing *Content) error {

	// check that state transition is allowed, only if state is provided
	if t.State != "" {
		if !existing.StateTransitionAllowed(t.State) {
			return apierrors.ErrTopicStateTransitionNotAllowed
		}
	}

	// if the topic is already completed, it cannot be updated

	return nil
}

// StateTransitionAllowed checks if the content can transition from its current state to the provided target state
func (t *Content) StateTransitionAllowed(target string) bool {
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

// AppendLinkInfo appends to list more links sorted by HRef
func (contentList *ContentResponseAPI) AppendLinkInfo(typeName string, itemLink *[]TypeLinkObject, id string, state string) {
	if itemLink == nil {
		return
	}

	nofItems := len(*itemLink)
	if nofItems == 0 {
		return
	}

	title := make(map[string]string)

	// Create list of sorted href's from itemLink list
	hrefs := make([]string, nofItems)
	for i, field := range *itemLink {
		hrefs[i] = field.HRef
		title[field.HRef] = field.Title
	}
	sort.Strings(hrefs)

	// Iterate through alphabeticaly sorted 'hrefs' and use each one to select corresponding title
	for _, href := range hrefs {
		// build up data items into structure
		var cItem ContentItem = ContentItem{
			Title: title[href],
			Type:  typeName,
			Links: &ContentLinks{
				Self: &LinkObject{
					HRef: href,
				},
				Topic: &LinkObject{
					ID:   id,
					HRef: "/topic/" + id,
				},
			},
			State: state,
		}

		if contentList.Items == nil {
			contentList.Items = &[]ContentItem{cItem}
		} else {
			*contentList.Items = append(*contentList.Items, cItem)
		}
	}

	contentList.TotalCount += nofItems
}
