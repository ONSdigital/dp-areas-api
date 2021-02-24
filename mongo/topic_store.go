package mongo

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongodb "github.com/ONSdigital/dp-mongodb"
	dpMongoHealth "github.com/ONSdigital/dp-mongodb/health"
	"github.com/ONSdigital/dp-topic-api/api"
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
func (m *Mongo) GetContent(id string, queryType int) (*models.ContentResponse, error) {
	s := m.Session.Copy()
	defer s.Close()

	var content models.ContentResponse
	contentSelect := bson.M{"ID": 1} // int default, used to minimise the mongo response to minimise go HEAP usage

	/* NOTE: if we have:

			contentSelect = bson.M{
			"ID":                           1,
			"next.id":                      1,
			"next.state":                   1,
			"next.spotlight":               0,
			"next.articles":                1,
			"current.id":                   1,
			"current.state":                1,
			"current.spotlight":            0,
			"current.articles":             1,
		}

		the line where spotlight is flaged as `0` causes the search to return no fields,

		so .. due to not being able to `de-select` items with a `0` parameter the
	         below has to select just the required type(s) for the desired query.
	*/
	switch queryType {
	case api.QuerySpotlight:
		contentSelect = bson.M{
			"ID":                1,
			"next.id":           1,
			"next.state":        1,
			"next.spotlight":    1,
			"current.id":        1,
			"current.state":     1,
			"current.spotlight": 1,
		}

	case api.QueryArticles:
		contentSelect = bson.M{
			"ID":               1,
			"next.id":          1,
			"next.state":       1,
			"next.articles":    1,
			"current.id":       1,
			"current.state":    1,
			"current.articles": 1,
		}

	case api.QueryBulletins:
		contentSelect = bson.M{
			"ID":                1,
			"next.id":           1,
			"next.state":        1,
			"next.bulletins":    1,
			"current.id":        1,
			"current.state":     1,
			"current.bulletins": 1,
		}

	case api.QueryMethodologies:
		contentSelect = bson.M{
			"ID":                    1,
			"next.id":               1,
			"next.state":            1,
			"next.methodologies":    1,
			"current.id":            1,
			"current.state":         1,
			"current.methodologies": 1,
		}
	case api.QueryMethodologyArticles:
		contentSelect = bson.M{
			"ID":                           1,
			"next.id":                      1,
			"next.state":                   1,
			"next.methodology_articles":    1,
			"current.id":                   1,
			"current.state":                1,
			"current.methodology_articles": 1,
		}

	case api.QueryStaticDatasets:
		contentSelect = bson.M{
			"ID":                      1,
			"next.id":                 1,
			"next.state":              1,
			"next.static_datasets":    1,
			"current.id":              1,
			"current.state":           1,
			"current.static_datasets": 1,
		}

	case api.QueryTimeseries:
		contentSelect = bson.M{
			"ID":                 1,
			"next.id":            1,
			"next.state":         1,
			"next.timeseries":    1,
			"current.id":         1,
			"current.state":      1,
			"current.timeseries": 1,
		}

	// Publications:
	case api.QueryArticles | api.QueryBulletins | api.QueryMethodologies | api.QueryMethodologyArticles:
		contentSelect = bson.M{
			"ID":                           1,
			"next.id":                      1,
			"next.state":                   1,
			"next.articles":                1,
			"next.bulletins":               1,
			"next.methodologies":           1,
			"next.methodology_articles":    1,
			"current.id":                   1,
			"current.state":                1,
			"current.articles":             1,
			"current.bulletins":            1,
			"current.methodologies":        1,
			"current.methodology_articles": 1,
		}

	// Datasets:
	case api.QueryStaticDatasets | api.QueryTimeseries:
		contentSelect = bson.M{
			"ID":                      1,
			"next.id":                 1,
			"next.state":              1,
			"next.static_datasets":    1,
			"next.timeseries":         1,
			"current.id":              1,
			"current.state":           1,
			"current.static_datasets": 1,
			"current.timeseries":      1,
		}

	// All types, that is a request for the content with no query parameter
	case api.QuerySpotlight | api.QueryArticles | api.QueryBulletins | api.QueryMethodologies |
		api.QueryMethodologyArticles | api.QueryStaticDatasets | api.QueryTimeseries:
		contentSelect = bson.M{
			"ID":                           1,
			"next.id":                      1,
			"next.state":                   1,
			"next.spotlight":               1,
			"next.articles":                1,
			"next.bulletins":               1,
			"next.methodologies":           1,
			"next.methodology_articles":    1,
			"next.static_datasets":         1,
			"next.timeseries":              1,
			"current.id":                   1,
			"current.state":                1,
			"current.spotlight":            1,
			"current.articles":             1,
			"current.bulletins":            1,
			"current.methodologies":        1,
			"current.methodology_articles": 1,
			"current.static_datasets":      1,
			"current.timeseries":           1,
		}
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
