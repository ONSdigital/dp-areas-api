package mongo

import (
	"context"
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
	mongodriver.MongoConnectionConfig

	connection   *mongodriver.MongoConnection
	healthClient *mongohealth.CheckMongoClient
}

// NewMongoStore creates a new Mongo object encapsulating a connection to the mongo server/cluster with the given configuration,
// and a health client to check the health of the mongo server/cluster
func NewMongoStore(_ context.Context, cfg config.MongoConfig) (m *Mongo, err error) {
	m = &Mongo{MongoConnectionConfig: cfg}

	m.connection, err = mongodriver.Open(&m.MongoConnectionConfig)
	if err != nil {
		return nil, err
	}

	databaseCollectionBuilder := make(map[mongohealth.Database][]mongohealth.Collection)
	databaseCollectionBuilder[(mongohealth.Database)(m.Database)] = []mongohealth.Collection{(mongohealth.Collection)(m.Collection)}

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
	err := m.connection.
		GetConfiguredCollection().
		Find(bson.M{"id": id}).
		Sort(bson.D{{"version", -1}}).
		One(ctx, &area)

	if err != nil {
		if mongodriver.IsErrNoDocumentFound(err) {
			return nil, apierrors.ErrAreaNotFound
		}
		return nil, err
	}

	return &area, nil
}

// GetVersion retrieves a version document for the area
func (m *Mongo) GetVersion(ctx context.Context, id string, versionID int) (*models.Area, error) {

	selector := bson.M{
		"id":      id,
		"version": versionID,
	}

	var version models.Area
	err := m.connection.GetConfiguredCollection().FindOne(ctx, selector, &version)
	if err != nil {
		if mongodriver.IsErrNoDocumentFound(err) {
			return nil, apierrors.ErrVersionNotFound
		}
		return nil, err
	}
	return &version, nil
}

// CheckAreaExists checks that the area exists
func (m *Mongo) CheckAreaExists(ctx context.Context, id string) error {
	var query bson.M
	query = bson.M{
		"_id": id,
	}
	count, err := m.connection.
		GetConfiguredCollection().
		Find(query).
		Count(ctx)
	if err != nil {
		return err
	}

	if count == 0 {
		return apierrors.ErrAreaNotFound
	}
	return nil
}

// GetAreas retrieves all areas documents
func (m *Mongo) GetAreas(ctx context.Context, offset, limit int) (*models.AreasResults, error) {

	findQuery := m.connection.
		GetConfiguredCollection().
		Find(bson.D{})
	totalCount, err := findQuery.Count(ctx)
	if err != nil {
		log.Error(ctx, "error counting items", err)
		if mongodriver.IsErrNoDocumentFound(err) {
			return &models.AreasResults{
				Items:      &[]models.Area{},
				Count:      0,
				TotalCount: 0,
				Offset:     offset,
				Limit:      limit,
			}, nil
		}
		return nil, err
	}

	values := []models.Area{}

	if limit > 0 {
		err := findQuery.
			Skip(offset).
			Limit(limit).
			IterAll(ctx, &values)

		if err != nil {
			if mongodriver.IsErrNoDocumentFound(err) {
				return &models.AreasResults{
					Items:      &values,
					Count:      0,
					TotalCount: totalCount,
					Offset:     offset,
					Limit:      limit,
				}, nil
			}
			return nil, err
		}
	}

	return &models.AreasResults{
		Items:      &values,
		Count:      len(values),
		TotalCount: totalCount,
		Offset:     offset,
		Limit:      limit,
	}, nil
}
