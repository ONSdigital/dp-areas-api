package api

//!!! sort the code for this, rename / adjust for /topic/<id>/content

// !!! make sure tests achieve full coverage of api/content.go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/models"
	storeMock "github.com/ONSdigital/dp-topic-api/store/datastoretest"

	. "github.com/smartystreets/goconvey/convey"
)

// Constants for testing
const ( // !!! remove not needed const's at some point / fix / rename, etc
	ctestTopicID1         = "topicTopicID1"
	ctestTopicCreatedID   = "topicCreatedID"
	ctestTopicPublishedID = "topicPublishedID"
	ctestUploadFilename   = "newimage.png"
)

const (
	chost = "http://localhost:25300"
)

//!!! build up response from following:
/*

{
    "id": "4",
    "next": {
        "spotlight": [
            {
                "Href": "/article/123",
                "Title": "Some article"
            },
            {
                "Href": "/dataset/12fasf3",
                "Title": "An interesting dataset"
            }
        ],
        "articles": [
            {
                "Href": "/article/1234",
                "Title": "Some article 2"
            },
            {
                "Href": "/article/12345",
                "Title": "Some article 3"
            }
        ],
        "bulletins": [
            {
                "Href": "/bulletins/this-month-hurray",
                "Title": "This Months Bulletin"
            }
        ],
        "timeseries": [
            {
                "Href": "/timseries/KVAC",
                "Title": "CPIH Time series"
            }
        ],
        "state" : "in_progress"
    },
    "current" : {
        "spotlight": [
            {
                "Href": "/article/123",
                "Title": "Some article"
            },
            {
                "Href": "/dataset/12fasf3",
                "Title": "An interesting dataset"
            }
        ],
        "articles": [
            {
                "Href": "/article/1234",
                "Title": "Some article 2"
            },
            {
                "Href": "/article/12345",
                "Title": "Some article 3"
            }
        ],
        "bulletins": [
            {
                "Href": "/bulletins/this-month-hurray",
                "Title": "This Months Bulletin"
            }
        ],
        "timeseries": [
            {
                "Href": "/timseries/KVAC",
                "Title": "CPIH Time series"
            }
        ],
        "state" : "published"
    }
}
*/
func dbContentWithID(state models.State, id string) *models.TopicResponse {
	return &models.TopicResponse{
		ID: id,
		Next: &models.Topic{
			ID:          id,
			Description: "next test description - 1",
			Title:       "test title - 1",
			Keywords:    []string{"keyword 1", "keyword 2", "keyword 3"},
			State:       state.String(),
			Links: &models.TopicLinks{
				Self: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s", id),
					ID:   fmt.Sprintf("%s", id),
				},
				Subtopics: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s/subtopics", id),
				},
				Content: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s/content", id),
				},
			},
		},
		Current: &models.Topic{
			ID:          id,
			Description: "current test description - 1",
			Title:       "test title - 1",
			Keywords:    []string{"keyword 1", "keyword 2", "keyword 3"},
			State:       state.String(),
			Links: &models.TopicLinks{
				Self: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s", id),
					ID:   fmt.Sprintf("%s", id),
				},
				Subtopics: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s/subtopics", id),
				},
				Content: &models.LinkObject{
					HRef: fmt.Sprintf("http://example.com/topics/%s/content", id),
				},
			},
		},
	}
}

// DB model corresponding to a topic in the provided state, without any download variant
func dbContent(state models.State) *models.TopicResponse {
	return dbContentWithID(state, ctestTopicID1)
}

// API model corresponding to dbCreatedTopic
func createdContentAll() *models.TopicResponse {
	return dbContent(models.StateTopicCreated)
}

func dbContentCurrentWithID(state models.State, id string) *models.Topic {
	return &models.Topic{
		ID:          id,
		Description: "current test description - 1",
		Title:       "test title - 1",
		Keywords:    []string{"keyword 1", "keyword 2", "keyword 3"},
		State:       state.String(),
		Links: &models.TopicLinks{
			Self: &models.LinkObject{
				HRef: fmt.Sprintf("http://example.com/topics/%s", id),
				ID:   fmt.Sprintf("%s", id),
			},
			Subtopics: &models.LinkObject{
				HRef: fmt.Sprintf("http://example.com/topics/%s/subtopics", id),
			},
			Content: &models.LinkObject{
				HRef: fmt.Sprintf("http://example.com/topics/%s/content", id),
			},
		},
	}
}

// create just the 'current' sub-document
func dbContentCurrent(state models.State) *models.Topic {
	return dbContentCurrentWithID(state, ctestTopicID1)
}

func createdContentCurrent() *models.Topic {
	return dbContentCurrent(models.StateTopicPublished)
}

// TestGetContentPublicHandler - does what the function name says
func TestGetContentPublicHandler(t *testing.T) {

	Convey("Given a topic API in web mode (private endpoints disabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = false
		Convey("And a topic API with mongoDB returning 'next' and 'current' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(id string) (*models.TopicResponse, error) {
					switch id {
					case ctestTopicID1:
						return dbContent(models.StateTopicPublished), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}

			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			Convey("When an existing 'published' topic is requested with the valid Topic-Id context value", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s", ctestTopicID1), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-document topic is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.Topic{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic, ShouldResemble, *createdContentCurrent())
				})
			})

			Convey("Requesting an nonexistent topic ID results in a NotFound response", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/inexistent", nil)
				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

func TestGetContentPrivateHandler(t *testing.T) {

	Convey("Given a topic API in publishing mode (private endpoints enabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = true
		Convey("And a topic API with mongoDB returning 'created' and 'full' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(id string) (*models.TopicResponse, error) {
					switch id {
					case ctestTopicID1:
						return dbContent(models.StateTopicCreated), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}
			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			Convey("When an existing 'created' topic is requested with the valid Topic-Id context value", func() {
				request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s", ctestTopicID1), nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected topic is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.TopicResponse{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic, ShouldResemble, *createdContentAll())
				})
			})

			Convey("Requesting an nonexistent topic ID results in a NotFound response", func() {
				request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/inexistent"), nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}
