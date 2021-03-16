package steps

import (
	"encoding/json"
	"time"

	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/cucumber/godog"
	"github.com/globalsign/mgo"
	"go.mongodb.org/mongo-driver/bson"
)

func (f *TopicComponent) iHaveThisRootTopic(topicsWriteJson *godog.DocString) error {

	topics := []models.TopicWrite{}
	m := f.MongoClient

	err := json.Unmarshal([]byte(topicsWriteJson.Content), &topics)
	if err != nil {
		return err
	}
	s := m.Session.Copy()
	defer s.Close()

	for _, topicsDoc := range topics {
		if err := f.putDatasetInDatabase(s, topicsDoc); err != nil {
			return err
		}
	}

	return nil
}

func (f *TopicComponent) putDatasetInDatabase(s *mgo.Session, topicDoc models.TopicWrite) error {
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
