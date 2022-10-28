package store

import (
	"context"
	"time"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-topic-api/models"
)

// DataStore provides a datastore.Storer interface used to store, retrieve, remove or update topics
type DataStore struct {
	Backend Storer
}

//go:generate moq -out mock/mongo.go -pkg mock . MongoDB
//go:generate moq -out mock/datastore.go -pkg mock . Storer

// dataMongoDB represents the required methods to access data from mongoDB
type dataMongoDB interface {
	GetTopic(ctx context.Context, id string) (*models.TopicResponse, error)
	CheckTopicExists(ctx context.Context, id string) error
	GetContent(ctx context.Context, id string, queryTypeFlags int) (*models.ContentResponse, error)
	UpdateReleaseDate(ctx context.Context, id string, releaseDate time.Time) error
	UpdateState(ctx context.Context, id, state string) error
	UpdateTopic(ctx context.Context, id string, topic *models.TopicResponse) error
}

// MongoDB represents all the required methods from mongo DB
type MongoDB interface {
	dataMongoDB
	Close(context.Context) error
	Checker(context.Context, *healthcheck.CheckState) error
}

// Storer represents basic data access via Get, Remove and Upsert methods, abstracting it from mongoDB
type Storer interface {
	dataMongoDB
}
