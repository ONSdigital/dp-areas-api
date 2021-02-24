package mongo

import (
	"context"
	"errors"
	"github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpmongo "github.com/ONSdigital/dp-mongodb"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/health"
	"github.com/ONSdigital/log.go/log"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Mongo represents a simplistic MongoDB configuration.
type Mongo struct {
	Session      *mgo.Session
	healthClient *dpMongoHealth.CheckMongoClient
	Database     string
	Collection   string
}

// Init creates a new mgo.Session with a strong consistency and a write mode of "majority".
func (m *Mongo) Init(mongoConf config.MongoConfig) error {
	if m.Session != nil {
		return errors.New("session already exists")
	}

	var err error
	if m.Session, err = mgo.Dial(mongoConf.BindAddr); err != nil {
		return err
	}

	m.Session.EnsureSafe(&mgo.Safe{WMode: "majority"})
	m.Session.SetMode(mgo.Strong, true)

	m.Database = mongoConf.Database
	m.Collection = mongoConf.Collection

	databaseCollectionBuilder := make(map[dpMongoHealth.Database][]dpMongoHealth.Collection)
	databaseCollectionBuilder[(dpMongoHealth.Database)(mongoConf.Database)] = []dpMongoHealth.Collection{(dpMongoHealth.Collection)(mongoConf.Collection)}

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
	return dpmongo.Close(ctx, m.Session)
}

// Checker is called by the healthcheck library to check the health state of this mongoDB instance
func (m *Mongo) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return m.healthClient.Checker(ctx, state)
}

//GetArea retrieves a area document by its ID
func (m *Mongo) GetArea(ctx context.Context, id string) (*models.Area, error) {
	s := m.Session.Copy()
	defer s.Close()
	log.Event(ctx, "getting area by ID", log.INFO, log.Data{"id": id})

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
		"id":     id,
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
		log.Event(ctx, "error counting items", log.ERROR, log.Error(err))
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

