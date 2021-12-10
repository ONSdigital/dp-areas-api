package mongo

import (
	"context"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-topic-api/api"
	errs "github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/models"

	mongohealth "github.com/ONSdigital/dp-mongodb/v3/health"
	mongodriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

type Mongo struct {
	mongodriver.MongoConnectionConfig

	Connection   *mongodriver.MongoConnection
	healthClient *mongohealth.CheckMongoClient

	ContentCollection string
}

func getConnectionConfig(cfg config.MongoConfig) mongodriver.MongoConnectionConfig {
	return mongodriver.MongoConnectionConfig{
		ClusterEndpoint:         cfg.BindAddr,
		Username:                cfg.Username,
		Password:                cfg.Password,
		Database:                cfg.Database,
		Collection:              cfg.TopicsCollection,
		ConnectTimeoutInSeconds: cfg.ConnectTimeoutInSeconds,
		QueryTimeoutInSeconds:   cfg.QueryTimeoutInSeconds,

		IsWriteConcernMajorityEnabled: cfg.IsWriteConcernMajorityEnabled,
		IsStrongReadConcernEnabled:    cfg.IsStrongReadConcernEnabled,

		TLSConnectionConfig: cfg.TLSConnectionConfig,
	}
}

// NewDBConnection creates a new mongodb.MongoConnection with the given configuration
func NewDBConnection(_ context.Context, cfg config.MongoConfig) (m *Mongo, err error) {
	m = &Mongo{MongoConnectionConfig: getConnectionConfig(cfg), ContentCollection: cfg.ContentCollection}

	m.Connection, err = mongodriver.Open(&m.MongoConnectionConfig)
	if err != nil {
		return nil, err
	}

	databaseCollectionBuilder := make(map[mongohealth.Database][]mongohealth.Collection)
	databaseCollectionBuilder[(mongohealth.Database)(m.Database)] = []mongohealth.Collection{(mongohealth.Collection)(m.Collection), (mongohealth.Collection)(m.ContentCollection)}

	m.healthClient = mongohealth.NewClientWithCollections(m.Connection, databaseCollectionBuilder)

	return m, nil
}

// Close closes the mongo session and returns any error
// It is an error to call m.Close if m.Init() returned an error, and there is no open connection
func (m *Mongo) Close(ctx context.Context) error {
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
		if mongodriver.IsErrNoDocumentFound(err) {
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
		if mongodriver.IsErrNoDocumentFound(err) {
			return nil, errs.ErrContentNotFound
		}
		return nil, err
	}

	return &content, nil
}
