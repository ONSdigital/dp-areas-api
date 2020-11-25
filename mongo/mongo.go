package mongo

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongodb "github.com/ONSdigital/dp-mongodb"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/health"
	errs "github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
)

// Mongo represents a simplistic MongoDB configuration, with session, health and lock clients
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

	// Create client and healthclient from session
	client := dpMongoHealth.NewClient(m.Session)
	m.healthClient = &dpMongoHealth.CheckMongoClient{
		Client:      *client,
		Healthcheck: client.Healthcheck,
	}

	return nil
}

// Close closes the mongo session and returns any error
func (m *Mongo) Close(ctx context.Context) error {
	return dpMongodb.Close(ctx, m.Session)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

// GetTopic retrieves a topic document by its ID
func (m *Mongo) GetTopic(ctx context.Context, id string) (*models.TopicUpdate, error) {
	s := m.Session.Copy()
	defer s.Close()
	log.Event(ctx, "getting topic by ID", log.Data{"id": id})

	var topic models.TopicUpdate

	err := s.DB(m.Database).C(m.TopicsCollection).Find(bson.M{"_id": id}).One(&topic)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errs.ErrTopicNotFound
		}
		return nil, err
	}

	return &topic, nil
}
