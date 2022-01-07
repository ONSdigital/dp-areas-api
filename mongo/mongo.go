package mongo

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/v3/health"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/ONSdigital/log.go/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	connectTimeoutInSeconds = 5
	queryTimeoutInSeconds   = 15
)

// Mongo represents a simplistic MongoDB configuration.
type Mongo struct {
	healthClient *dpMongoHealth.CheckMongoClient
	Database     string
	Collection   string
	Connection   *dpMongoDriver.MongoConnection
	Username     string
	Password     string
	URI          string
	IsSSL        bool
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
		Collection:                    m.Collection,
		IsWriteConcernMajorityEnabled: shouldEnableWriteConcern,
		IsStrongReadConcernEnabled:    shouldEnableReadConcern,
	}
}

// Init creates a new mongoConnection with a strong consistency and a write mode of "majority".
func (m *Mongo) Init(ctx context.Context, shouldEnableReadConcern, shouldEnableWriteConcern bool) error {
	if m.Connection != nil {
		return errors.New("Datastore Connection already exists")
	}
	mongoConnection, err := dpMongoDriver.Open(m.getConnectionConfig(shouldEnableReadConcern, shouldEnableWriteConcern))
	if err != nil {
		return err
	}
	m.Connection = mongoConnection
	databaseCollectionBuilder := make(map[dpMongoHealth.Database][]dpMongoHealth.Collection)
	databaseCollectionBuilder[(dpMongoHealth.Database)(m.Database)] = []dpMongoHealth.Collection{(dpMongoHealth.Collection)(m.Collection)}

	// Create health-client from session AND collections
	m.healthClient = dpMongoHealth.NewClientWithCollections(mongoConnection, databaseCollectionBuilder)

	return nil
}

// Close closes the mongo session and returns any error
func (m *Mongo) Close(ctx context.Context) error {
	if m.Connection == nil {
		return errors.New("cannot close an empty connection")
	}
	return m.Connection.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

//GetArea retrieves a area document by its ID
func (m *Mongo) GetArea(ctx context.Context, id string) (*models.Area, error) {
	log.Info(ctx, "getting area by ID", log.Data{"id": id})

	var area models.Area
	err := m.Connection.
		GetConfiguredCollection().
		Find(bson.M{"id": id}).
		Sort(bson.D{{"version", -1}}).
		One(ctx, &area)

	if err != nil {
		if dpMongoDriver.IsErrNoDocumentFound(err) {
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
	err := m.Connection.GetConfiguredCollection().FindOne(ctx, selector, &version)
	if err != nil {
		if dpMongoDriver.IsErrNoDocumentFound(err) {
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
	count, err := m.Connection.
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

	findQuery := m.Connection.
		GetConfiguredCollection().
		Find(bson.D{})
	totalCount, err := findQuery.Count(ctx)
	if err != nil {
		log.Error(ctx, "error counting items", err)
		if dpMongoDriver.IsErrNoDocumentFound(err) {
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
			if dpMongoDriver.IsErrNoDocumentFound(err) {
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
