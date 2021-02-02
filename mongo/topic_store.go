package mongo

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongodb "github.com/ONSdigital/dp-mongodb"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/health"
	errs "github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
)

// Mongo represents a simplistic MongoDB config, with session, health and lock clients
type Mongo struct {
	Session           *mgo.Session
	URI               string
	Database          string
	TopicsCollection  string
	ContentCollection string
	healthClient      *dpMongoHealth.CheckMongoClient
}

// Init creates a new mgo.Session with a strong consistency and a write mode of "majority".
func (m *Mongo) Init(ctx context.Context) (err error) {
	if m.Session != nil {
		return errors.New("session already exists")
	}

	// Create session
	if m.Session, err = mgo.Dial(m.URI); err != nil {
		return err
	}
	m.Session.EnsureSafe(&mgo.Safe{WMode: "majority"})
	m.Session.SetMode(mgo.Strong, true)

	databaseCollectionBuilder := make(map[dpMongoHealth.Database][]dpMongoHealth.Collection)
	databaseCollectionBuilder[(dpMongoHealth.Database)(m.Database)] = []dpMongoHealth.Collection{(dpMongoHealth.Collection)(m.TopicsCollection), (dpMongoHealth.Collection)(m.ContentCollection)}

	// Create client and healthclient from session AND collections
	client := dpMongoHealth.NewClientWithCollections(m.Session, databaseCollectionBuilder)

	m.healthClient = &dpMongoHealth.CheckMongoClient{
		Client:      *client,
		Healthcheck: client.Healthcheck,
	}

	return nil
}

// Close closes the mongo session and returns any error
func (m *Mongo) Close(ctx context.Context) error {
	if m.Session == nil {
		return errors.New("cannot close a mongoDB connection without a valid session")
	}
	return dpMongodb.Close(ctx, m.Session)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

// GetTopic retrieves a topic document by its ID
func (m *Mongo) GetTopic(id string) (*models.TopicResponse, error) {
	s := m.Session.Copy()
	defer s.Close()

	var topic models.TopicResponse

	err := s.DB(m.Database).C(m.TopicsCollection).Find(bson.M{"id": id}).One(&topic)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errs.ErrTopicNotFound
		}
		return nil, err
	}

	return &topic, nil
}

// GetContent retrieves a content document by its ID
func (m *Mongo) GetContent(id string) (*models.ContentResponse, error) {
	s := m.Session.Copy()
	defer s.Close()

	var content models.ContentResponse

	err := s.DB(m.Database).C(m.ContentCollection).Find(bson.M{"id": id}).One(&content)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errs.ErrContentNotFound
		}
		return nil, err
	}

	return &content, nil
}
