package steps

import (
	"encoding/json"
	"time"

	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/cucumber/godog"
	"github.com/globalsign/mgo"
	"go.mongodb.org/mongo-driver/bson"
)

func (f *TopicComponent) iHaveTheseTopics(topicsWriteJson *godog.DocString) error {

	topics := []models.TopicWrite{}
	m := f.MongoClient

	err := json.Unmarshal([]byte(topicsWriteJson.Content), &topics)
	if err != nil {
		return err
	}
	s := m.Session.Copy()
	defer s.Close()

	for _, topicsDoc := range topics {
		if err := f.putTopicInDatabase(s, topicsDoc); err != nil {
			return err
		}
	}

	return nil
}

func (f *TopicComponent) putTopicInDatabase(s *mgo.Session, topicDoc models.TopicWrite) error {
	update := bson.M{
		"$set": topicDoc,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := s.DB(f.MongoClient.Database).C("topics").UpsertId(topicDoc.ID, update)
	if err != nil {
		return err
	}
	return nil
}

func (f *TopicComponent) iHaveTheseContents(contentJson *godog.DocString) error {

	collection := []models.ContentResponse{}
	m := f.MongoClient

	err := json.Unmarshal([]byte(contentJson.Content), &collection)
	if err != nil {
		return err
	}
	s := m.Session.Copy()
	defer s.Close()

	for _, topicsDoc := range collection {
		if err := f.putContentInDatabase(s, topicsDoc); err != nil {
			return err
		}
	}

	return nil
}

func (f *TopicComponent) putContentInDatabase(s *mgo.Session, contentDoc models.ContentResponse) error {
	update := bson.M{
		"$set": contentDoc,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := s.DB(f.MongoClient.Database).C("content").UpsertId(contentDoc.ID, update)
	if err != nil {
		return err
	}
	return nil
}

func (f *TopicComponent) privateEndpointsAreEnabled() error {
	f.Config.EnablePrivateEndpoints = true
	return nil
}
