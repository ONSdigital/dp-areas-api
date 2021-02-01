package models

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
