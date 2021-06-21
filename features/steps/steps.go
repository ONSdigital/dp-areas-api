package steps

import (
	"context"
	"encoding/json"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v2/pkg/mongodb"
	"time"

	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/cucumber/godog"
	"go.mongodb.org/mongo-driver/bson"
)

func (f *TopicComponent) iHaveTheseTopics(topicsWriteJson *godog.DocString) error {
	ctx := context.Background()
	topics := []models.TopicWrite{}
	m := f.MongoClient

	err := json.Unmarshal([]byte(topicsWriteJson.Content), &topics)
	if err != nil {
		return err
	}

	for _, topicsDoc := range topics {
		if err := f.putTopicInDatabase(ctx, m.Connection, topicsDoc); err != nil {
			return err
		}
	}

	return nil
}

func (f *TopicComponent) putTopicInDatabase(ctx context.Context, mongoConnection *dpMongoDriver.MongoConnection, topicDoc models.TopicWrite) error {
	update := bson.M{
		"$set": topicDoc,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := mongoConnection.GetConfiguredCollection().UpsertId(ctx, topicDoc.ID, update)
	if err != nil {
		return err
	}
	return nil
}

func (f *TopicComponent) iHaveTheseContents(contentJson *godog.DocString) error {
	ctx := context.Background()
	collection := []models.ContentResponse{}
	m := f.MongoClient

	err := json.Unmarshal([]byte(contentJson.Content), &collection)
	if err != nil {
		return err
	}

	for _, topicsDoc := range collection {
		if err := f.putContentInDatabase(ctx, m.Connection, topicsDoc); err != nil {
			return err
		}
	}

	return nil
}

func (f *TopicComponent) putContentInDatabase(ctx context.Context, mongoConnection *dpMongoDriver.MongoConnection, contentDoc models.ContentResponse) error {
	update := bson.M{
		"$set": contentDoc,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := mongoConnection.C(f.MongoClient.ContentCollection).UpsertId(ctx, contentDoc.ID, update)
	if err != nil {
		return err
	}
	return nil
}

func (f *TopicComponent) privateEndpointsAreEnabled() error {
	f.Config.EnablePrivateEndpoints = true
	return nil
}
