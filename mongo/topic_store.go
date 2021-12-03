package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/v3/health"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/ONSdigital/dp-topic-api/api"
	errs "github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
)

const (
	connectTimeoutInSeconds = 5
	queryTimeoutInSeconds   = 15
)

// Mongo represents a simplistic MongoDB config, with session, health and lock clients
type Mongo struct {
	URI               string
	Database          string
	TopicsCollection  string
	ContentCollection string
	Connection        *dpMongoDriver.MongoConnection
	Username          string
	Password          string
	healthClient      *dpMongoHealth.CheckMongoClient
	IsSSL             bool
}

func (m *Mongo) getConnectionConfig(shouldEnableReadConcern, shouldEnableWriteConcern bool) *dpMongoDriver.MongoConnectionConfig {
	return &dpMongoDriver.MongoConnectionConfig{
		TLSConnectionConfig: dpMongoDriver.TLSConnectionConfig{
			IsSSL: m.IsSSL,
		},
		ConnectTimeoutInSeconds: connectTimeoutInSeconds,
		QueryTimeoutInSeconds:   queryTimeoutInSeconds,

		Username:                      m.Username,
		Password:                      m.Password,
		ClusterEndpoint:               m.URI,
		Database:                      m.Database,
		Collection:                    m.TopicsCollection,
		IsWriteConcernMajorityEnabled: shouldEnableWriteConcern,
		IsStrongReadConcernEnabled:    shouldEnableReadConcern,
	}
}

// Init creates a new mongoConnection with a strong consistency and a write mode of "majority".
func (m *Mongo) Init(ctx context.Context, shouldEnableReadConcern, shouldEnableWriteConcern bool) (err error) {
	if m.Connection != nil {
		return errors.New("Datastore Connection already exists")
	}
	mongoConnection, err := dpMongoDriver.Open(m.getConnectionConfig(shouldEnableReadConcern, shouldEnableWriteConcern))
	if err != nil {
		return err
	}
	m.Connection = mongoConnection
	databaseCollectionBuilder := make(map[dpMongoHealth.Database][]dpMongoHealth.Collection)
	databaseCollectionBuilder[(dpMongoHealth.Database)(m.Database)] = []dpMongoHealth.Collection{(dpMongoHealth.Collection)(m.TopicsCollection), (dpMongoHealth.Collection)(m.ContentCollection)}

	// Create health-client from session AND collections
	m.healthClient = dpMongoHealth.NewClientWithCollections(mongoConnection, databaseCollectionBuilder)

	return nil
}

// Close closes the mongo session and returns any error
func (m *Mongo) Close(ctx context.Context) error {
	if m.Connection == nil {
		return errors.New("cannot close a empty connection")
	}
	return m.Connection.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

// GetTopic retrieves a topic document by its ID
func (m *Mongo) GetTopic(ctx context.Context, id string) (*models.TopicResponse, error) {
	var topic models.TopicResponse

	err := m.Connection.GetConfiguredCollection().FindOne(ctx, bson.M{"id": id}, &topic)
	if err != nil {
		if dpMongoDriver.IsErrNoDocumentFound(err) {
			return nil, errs.ErrTopicNotFound
		}
		return nil, err
	}

	return &topic, nil
}

// CheckTopicExists checks that the topic exists
func (m *Mongo) CheckTopicExists(ctx context.Context, id string) error {

	count, err := m.Connection.
		GetConfiguredCollection().
		Find(bson.M{"id": id}).
		Count(ctx)
	if err != nil {
		if dpMongoDriver.IsErrNoDocumentFound(err) {
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
func (m *Mongo) GetContent(ctx context.Context, id string, queryTypeFlags int) (*models.ContentResponse, error) {
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

	err := m.Connection.
		C(m.ContentCollection).
		Find(bson.M{"id": id}).
		Select(contentSelect).
		One(ctx, &content)
	if err != nil {
		if dpMongoDriver.IsErrNoDocumentFound(err) {
			return nil, errs.ErrContentNotFound
		}
		return nil, err
	}

	return &content, nil
}
