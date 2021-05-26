package mongo

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongodb "github.com/ONSdigital/dp-mongodb"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/v2/pkg/health"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v2/pkg/mongo-driver"
	"github.com/ONSdigital/dp-topic-api/api"
	errs "github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
)

const (
	connectTimeoutInSeconds = 5
	queryTimeoutInSeconds   = 15
)

// Mongo represents a simplistic MongoDB config, with session, health and lock clients
type Mongo struct {
	Session           *mgo.Session
	URI               string
	Database          string
	TopicsCollection  string
	ContentCollection string
	Connection        *dpMongoDriver.MongoConnection
	Username          string
	Password          string
	CAFilePath        string
	healthClient      *dpMongoHealth.CheckMongoClient
}

func (m *Mongo) getConnectionConfig() *dpMongoDriver.MongoConnectionConfig {
	return &dpMongoDriver.MongoConnectionConfig{
		CaFilePath:              m.CAFilePath,
		ConnectTimeoutInSeconds: connectTimeoutInSeconds,
		QueryTimeoutInSeconds:   queryTimeoutInSeconds,

		Username:             m.Username,
		Password:             m.Password,
		ClusterEndpoint:      m.URI,
		Database:             m.Database,
		Collection:           m.TopicsCollection,
		SkipCertVerification: true,
	}
}

// Init creates a new mongoConnection with a strong consistency and a write mode of "majority".
func (m *Mongo) Init(ctx context.Context) (err error) {
	if m.Connection != nil {
		return errors.New("Datastore Connection already exists")
	}
	mongoConnection, err := dpMongoDriver.Open(m.getConnectionConfig())
	if err != nil {
		return err
	}
	m.Connection = mongoConnection
	databaseCollectionBuilder := make(map[dpMongoHealth.Database][]dpMongoHealth.Collection)
	databaseCollectionBuilder[(dpMongoHealth.Database)(m.Database)] = []dpMongoHealth.Collection{(dpMongoHealth.Collection)(m.TopicsCollection), (dpMongoHealth.Collection)(m.ContentCollection)}

	// Create client and health-client from session AND collections
	client := dpMongoHealth.NewClientWithCollections(mongoConnection, databaseCollectionBuilder)

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

// CheckTopicExists checks that the topic exists
func (m *Mongo) CheckTopicExists(id string) error {
	s := m.Session.Copy()
	defer s.Close()

	count, err := s.DB(m.Database).C(m.TopicsCollection).Find(bson.M{"id": id}).Count()
	if err != nil {
		if err == mgo.ErrNotFound {
			return errs.ErrTopicNotFound
		}
		return err
	}

	if count == 0 {
		return errs.ErrTopicNotFound
	}

	return nil
}

// GetContent retrieves a content document by its ID
func (m *Mongo) GetContent(id string, queryTypeFlags int) (*models.ContentResponse, error) {
	s := m.Session.Copy()
	defer s.Close()

	var content models.ContentResponse
	// init default, used to minimise the mongo response to minimise go HEAP usage
	contentSelect := bson.M{
		"ID":            1,
		"next.id":       1,
		"next.state":    1,
		"current.id":    1,
		"current.state": 1,
	}

	// Add spotlight first
	if (queryTypeFlags & api.QuerySpotlightFlag) != 0 {
		contentSelect["next.spotlight"] = 1
		contentSelect["current.spotlight"] = 1
	}

	// then Publications
	if (queryTypeFlags & api.QueryArticlesFlag) != 0 {
		contentSelect["next.articles"] = 1
		contentSelect["current.articles"] = 1
	}

	if (queryTypeFlags & api.QueryBulletinsFlag) != 0 {
		contentSelect["next.bulletins"] = 1
		contentSelect["current.bulletins"] = 1
	}

	if (queryTypeFlags & api.QueryMethodologiesFlag) != 0 {
		contentSelect["next.methodologies"] = 1
		contentSelect["current.methodologies"] = 1
	}

	if (queryTypeFlags & api.QueryMethodologyArticlesFlag) != 0 {
		contentSelect["next.methodology_articles"] = 1
		contentSelect["current.methodology_articles"] = 1
	}

	// then Datasets
	if (queryTypeFlags & api.QueryStaticDatasetsFlag) != 0 {
		contentSelect["next.static_datasets"] = 1
		contentSelect["current.static_datasets"] = 1
	}

	if (queryTypeFlags & api.QueryTimeseriesFlag) != 0 {
		contentSelect["next.timeseries"] = 1
		contentSelect["current.timeseries"] = 1
	}

	err := s.DB(m.Database).C(m.ContentCollection).Find(bson.M{"id": id}).Select(contentSelect).One(&content)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, errs.ErrContentNotFound
		}
		return nil, err
	}

	return &content, nil
}
