package store

import (
	"context"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-topic-api/models"
)

// DataStore provides a datastore.Storer interface used to store, retrieve, remove or update topics
type DataStore struct {
	Backend Storer
}

//go:generate moq -out datastoretest/mongo.go -pkg storetest . MongoDB
//go:generate moq -out datastoretest/datastore.go -pkg storetest . Storer

// dataMongoDB represents the required methods to access data from mongoDB
type dataMongoDB interface {
	GetTopic(id string) (*models.TopicResponse, error)
	CheckTopicExists(id string) error
	GetContent(id string, queryTypeFlags int) (*models.ContentResponse, error)
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
