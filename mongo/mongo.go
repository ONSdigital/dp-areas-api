package mongo

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"

	mongohealth "github.com/ONSdigital/dp-mongodb/v3/health"
	mongodriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

type Mongo struct {
	mongodriver.MongoDriverConfig

	connection   *mongodriver.MongoConnection
	healthClient *mongohealth.CheckMongoClient
}

// NewMongoStore creates a new Mongo object encapsulating a connection to the mongo server/cluster with the given configuration,
// and a health client to check the health of the mongo server/cluster
func NewMongoStore(_ context.Context, cfg config.MongoConfig) (m *Mongo, err error) {
	m = &Mongo{MongoDriverConfig: cfg}

	m.connection, err = mongodriver.Open(&m.MongoDriverConfig)
	if err != nil {
		return nil, err
	}

	databaseCollectionBuilder := map[mongohealth.Database][]mongohealth.Collection{
		mongohealth.Database(m.Database): {mongohealth.Collection(m.ActualCollectionName(config.AreasCollection))}}
	m.healthClient = mongohealth.NewClientWithCollections(m.connection, databaseCollectionBuilder)

	return m, nil
}

// Close the mongo session and returns any error
// It is an error to call m.Close if m.Init() returned an error, and there is no open connection
func (m *Mongo) Close(ctx context.Context) error {
	return m.connection.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

//GetArea retrieves a area document by its ID
func (m *Mongo) GetArea(ctx context.Context, id string) (*models.Area, error) {
	log.Info(ctx, "getting area by ID", log.Data{"id": id})

	var area models.Area
	err := m.connection.Collection(m.ActualCollectionName(config.AreasCollection)).
		FindOne(ctx, bson.M{"id": id}, &area, mongodriver.Sort(bson.D{{"version", -1}}))
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocumentFound) {
			return nil, apierrors.ErrAreaNotFound
		}
		return nil, err
	}

	return &area, nil
}

// GetVersion retrieves a version document for the area
func (m *Mongo) GetVersion(ctx context.Context, id string, versionID int) (*models.Area, error) {

	var version models.Area
	err := m.connection.Collection(m.ActualCollectionName(config.AreasCollection)).
		FindOne(ctx, bson.M{"id": id, "version": versionID}, &version)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocumentFound) {
			return nil, apierrors.ErrVersionNotFound
		}
		return nil, err
	}

	return &version, nil
}

// GetAreas retrieves all areas documents
func (m *Mongo) GetAreas(ctx context.Context, offset, limit int) (*models.AreasResults, error) {

	var result = []models.Area{}
	totalCount, err := m.connection.Collection(m.ActualCollectionName(config.AreasCollection)).
		Find(ctx, bson.D{}, &result, mongodriver.Sort(bson.M{"_id": 1}), mongodriver.Offset(offset), mongodriver.Limit(limit))
	if err != nil {
		log.Error(ctx, "error finding areas", err)
		return nil, err
	}

	return &models.AreasResults{
		Items:      &result,
		Count:      len(result),
		TotalCount: totalCount,
		Offset:     offset,
		Limit:      limit,
	}, nil
}
