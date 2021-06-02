package mongo

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/v2/pkg/health"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v2/pkg/mongo-driver"
	"github.com/ONSdigital/log.go/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

const (
	connectTimeoutInSeconds = 5
	queryTimeoutInSeconds   = 15
)

// Mongo represents a simplistic MongoDB configuration.
type Mongo struct {
	Session      *mgo.Session
	healthClient *dpMongoHealth.CheckMongoClient
	Database     string
	Collection   string
	Connection   *dpMongoDriver.MongoConnection
	Username     string
	Password     string
	CAFilePath   string
	URI          string
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
		Collection:           m.Collection,
		SkipCertVerification: true,
	}
}

// Init creates a new mongoConnection with a strong consistency and a write mode of "majority".
func (m *Mongo) Init(ctx context.Context) error {
	if m.Connection != nil {
		return errors.New("Datastore Connection already exists")
	}
	mongoConnection, err := dpMongoDriver.Open(m.getConnectionConfig())
	if err != nil {
		return err
	}
	m.Connection = mongoConnection
	databaseCollectionBuilder := make(map[dpMongoHealth.Database][]dpMongoHealth.Collection)
	databaseCollectionBuilder[(dpMongoHealth.Database)(m.Database)] = []dpMongoHealth.Collection{(dpMongoHealth.Collection)(m.Collection)}

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
	if m.Connection == nil {
		return errors.New("cannot close a empty connection")
	}
	return m.Connection.Close(ctx)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

//GetArea retrieves a area document by its ID
func (m *Mongo) GetArea(ctx context.Context, id string) (*models.Area, error) {
	s := m.Session.Copy()
	defer s.Close()
	log.Info(ctx, "getting area by ID", log.Data{"id": id})

	var area models.Area
	err := s.DB(m.Database).C(m.Collection).Find(bson.M{"id": id}).Sort("-version").One(&area)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, apierrors.ErrAreaNotFound
		}
		return nil, err
	}

	return &area, nil
}

// GetVersion retrieves a version document for the area
func (m *Mongo) GetVersion(id string, versionID int) (*models.Area, error) {
	s := m.Session.Copy()
	defer s.Close()

	selector := bson.M{
		"id":      id,
		"version": versionID,
	}

	var version models.Area
	err := s.DB(m.Database).C("areas").Find(selector).One(&version)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, apierrors.ErrVersionNotFound
		}
		return nil, err
	}
	return &version, nil
}

// CheckAreaExists checks that the area exists
func (m *Mongo) CheckAreaExists(id string) error {
	s := m.Session.Copy()
	defer s.Close()

	var query bson.M
	query = bson.M{
		"_id": id,
	}
	count, err := s.DB(m.Database).C("areas").Find(query).Count()
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
	s := m.Session.Copy()
	defer s.Close()

	var q *mgo.Query

	q = s.DB(m.Database).C("areas").Find(nil)

	totalCount, err := q.Count()
	if err != nil {
		log.Error(ctx, "error counting items", err)
		if err == mgo.ErrNotFound {
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
		iter := q.Skip(offset).Limit(limit).Iter()

		defer func() {
			err := iter.Close()
			if err != nil {
				log.Event(ctx, "error closing iterator", log.ERROR, log.Error(err))
			}
		}()

		if err := iter.All(&values); err != nil {
			if err == mgo.ErrNotFound {
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
