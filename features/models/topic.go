package models

import (
	"time"

	"github.com/ONSdigital/dp-topic-api/models"
)

// TopicWrite is used for component testing
type TopicWrite struct {
	ID      string  `bson:"id,omitempty"       json:"id,omitempty"`
	Current *TopicW `bson:"current,omitempty"  json:"current,omitempty"`
	Next    *TopicW `bson:"next,omitempty"     json:"next,omitempty"`
}

// TopicW is used for component testing
type TopicW struct {
	ID          string             `bson:"id,omitempty"             json:"id,omitempty"`
	Description string             `bson:"description,omitempty"    json:"description,omitempty"`
	Keywords    []string           `bson:"keywords,omitempty"       json:"keywords,omitempty"`
	Links       *models.TopicLinks `bson:"links,omitempty"          json:"links,omitempty"`
	ReleaseDate *time.Time         `bson:"release_date,omitempty"          json:"release_date,omitempty"`
	State       string             `bson:"state,omitempty"          json:"state,omitempty"`
	SubtopicIds []string           `bson:"subtopics_ids,omitempty"  json:"subtopics_ids,omitempty"`
	Title       string             `bson:"title,omitempty"          json:"title,omitempty"`
}
