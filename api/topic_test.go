package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/mocks"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/dp-topic-api/store"
	storeMock "github.com/ONSdigital/dp-topic-api/store/datastoretest"
	"github.com/gorilla/mux"

	. "github.com/smartystreets/goconvey/convey"
)

// Constants for testing
const ( // !!! remove not needed const's at some point
	testUserAuthToken      = "UserToken"
	testTopicID1           = "topicTopicID1"
	testTopicID2           = "topicTopicID2"
	testTopicCreatedID     = "topicCreatedID"
	testTopicUploadedID    = "topicUploadedID"
	testTopicImportingID   = "topicImportingID"
	testTopicPublishedID   = "topicPublishedID"
	testCollectionID1      = "1234"
	testVariantOriginal    = "original"
	testVariantAlternative = "bw1024"
	testUploadFilename     = "newimage.png"
	testUploadPath         = "s3://images/" + testUploadFilename
	longName               = "Llanfairpwllgwyngyllgogerychwyrndrobwllllantysiliogogogoch"
	testLockID             = "image-myID-123456789"
	testDownloadType       = "originally uploaded file"
	testPrivateHref        = "http://download.ons.gov.uk/images/imageImageID2/original/some-image-name"
	testFilename           = "some-image-name"
)

const (
	host              = "http://localhost:25300"
	authToken         = "dataset"
	healthTimeout     = 2 * time.Second
	internalServerErr = "internal server error\n"
	callerIdentity    = "someone@ons.gov.uk"
)

var errMongoDB = errors.New("MongoDB generic error")

// DB model corresponding to a topic in the provided state, without any download variant
func dbTopic(state models.State) *models.TopicUpdate {
	return dbTopicWithID(state, testTopicID1)
}

func dbTopicWithID(state models.State, id string) *models.TopicUpdate {
	return &models.TopicUpdate{
		ID: id,
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
	}
}

// API model corresponding to dbCreatedTopic
func createdTopicAll() *models.TopicUpdate {
	return dbTopic(models.StateTopicCreated)
}

// create just the 'current' sub-document
func dbTopicCurrent(state models.State) *models.Topic {
	return dbTopicCurrentWithID(state, testTopicID1)
}

func dbTopicCurrentWithID(state models.State, id string) *models.Topic {
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

func createdTopicCurrent() *models.Topic {
	return dbTopicCurrent(models.StateTopicPublished)
}

func TestGetTopicPublicHandler(t *testing.T) {

	Convey("Given a topic API in web mode (private endpoints disabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = false
		Convey("And a topic API with mongoDB returning 'next' and 'current' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(id string) (*models.TopicUpdate, error) {
					switch id {
					case testTopicID1:
						return dbTopic(models.StateTopicPublished), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}

			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			Convey("When an existing 'published' topic is requested with the valid Topic-Id context value", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s", testTopicID1), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-document topic is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.Topic{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic, ShouldResemble, *createdTopicCurrent())
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

func TestGetTopicPrivateHandler(t *testing.T) {

	Convey("Given a topic API in publishing mode (private endpoints enabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = true
		Convey("And a topic API with mongoDB returning 'created' and 'full' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(id string) (*models.TopicUpdate, error) {
					switch id {
					case testTopicID1:
						return dbTopic(models.StateTopicCreated), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}
			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			Convey("When an existing 'created' topic is requested with the valid Topic-Id context value", func() {
				request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s", testTopicID1), nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected topic is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.TopicUpdate{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic, ShouldResemble, *createdTopicAll())
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

// GetAPIWithMocks also used in other tests, so exported
func GetAPIWithMocks(cfg *config.Configuration, mockedDataStore store.Storer) *API {
	mu.Lock()
	defer mu.Unlock()
	//	urlBuilder := url.NewBuilder("http://example.com")

	topicPermissions := mocks.NewAuthHandlerMock()
	permissions := mocks.NewAuthHandlerMock()

	return Setup(testContext, cfg, mux.NewRouter(), store.DataStore{Backend: mockedDataStore}, topicPermissions, permissions)
}

func BenchmarkGetTopicPrivateHandler(b *testing.B) {

	// Given a topic API in publishing mode (private endpoints enabled)
	cfg, err := config.Get()
	if err != nil {
		fmt.Printf("config fail\n")
		return
	}
	cfg.EnablePrivateEndpoints = true
	// And a topic API with mongoDB returning 'created' and 'full' topics

	mongoDBMock := &storeMock.MongoDBMock{
		GetTopicFunc: func(id string) (*models.TopicUpdate, error) {
			switch id {
			case testTopicID1:
				return dbTopic(models.StateTopicCreated), nil
			default:
				return nil, apierrors.ErrTopicNotFound
			}
		},
	}
	topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

	b.ReportAllocs()

	// test all event types
	for i := 0; i < b.N; i++ {

		// When an existing 'created' topic is requested with the valid Topic-Id context value
		request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s", testTopicID1), nil)
		if err != nil {
			fmt.Printf("request fail\n")
			return
		}

		w := httptest.NewRecorder()
		topicAPI.Router.ServeHTTP(w, request)
		// Then the expected topic is returned with status code 200
		payload, err := ioutil.ReadAll(w.Body)
		if err != nil {
			fmt.Printf("readall fail\n")
			return
		}
		retTopic := models.TopicUpdate{}
		err = json.Unmarshal(payload, &retTopic)
		// NOTE: to check that the correct structure is unmarshaled and the contents of the structure
		// are as expected, run this benchmark in the debugger with a breakpoint on the next line ...
		if err != nil {
			fmt.Printf("unmarshal fail\n")
			return
		}
	}
}

func BenchmarkGetDatasetPrivate(b *testing.B) {

	// Given a topic API in publishing mode (private endpoints enabled)
	cfg, err := config.Get()
	if err != nil {
		fmt.Printf("config fail\n")
		return
	}
	cfg.EnablePrivateEndpoints = true
	//	And a topic API with mongoDB returning 'created' and 'full' topics

	mongoDBMock := &storeMock.MongoDBMock{
		GetTopicFunc: func(id string) (*models.TopicUpdate, error) {
			switch id {
			case testTopicID1:
				return dbTopic(models.StateTopicCreated), nil
			default:
				return nil, apierrors.ErrTopicNotFound
			}
		},
	}
	topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

	b.ReportAllocs()

	// test all event types
	for i := 0; i < b.N; i++ {

		//	When an existing 'created' topic is requested with the valid Topic-Id context value
		request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/datasets/%s", testTopicID1), nil)
		if err != nil {
			fmt.Printf("request fail\n")
			return
		}

		w := httptest.NewRecorder()
		topicAPI.Router.ServeHTTP(w, request)
		// Then the expected topic is returned with status code 200
		payload, err := ioutil.ReadAll(w.Body)
		if err != nil {
			fmt.Printf("readall fail\n")
			return
		}
		retTopic := models.TopicUpdate{}
		err = json.Unmarshal(payload, &retTopic)
		// NOTE: to check that the correct structure is unmarshaled and the contents of the structure
		// are as expected, run this benchmark in the debugger with a breakpoint on the next line ...
		if err != nil {
			fmt.Printf("unmarshal fail\n")
			return
		}
	}
}

func BenchmarkGetTopicPublicHandler(b *testing.B) {

	// Given a topic API in web mode (private endpoints disabled)
	cfg, err := config.Get()
	if err != nil {
		fmt.Printf("config fail\n")
		return
	}

	cfg.EnablePrivateEndpoints = false
	// And a topic API with mongoDB returning 'next' and 'current' topics

	mongoDBMock := &storeMock.MongoDBMock{
		GetTopicFunc: func(id string) (*models.TopicUpdate, error) {
			switch id {
			case testTopicID1:
				return dbTopic(models.StateTopicCreated), nil
			default:
				return nil, apierrors.ErrTopicNotFound
			}
		},
	}

	topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

	b.ReportAllocs()

	// test all event types
	for i := 0; i < b.N; i++ {

		// When an existing 'current' topic is requested with the valid Topic-Id context value
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s", testTopicID1), nil)

		w := httptest.NewRecorder()
		topicAPI.Router.ServeHTTP(w, request)
		// Then the expected topic is returned with status code 200
		payload, err := ioutil.ReadAll(w.Body)
		if err != nil {
			fmt.Printf("readall fail\n")
			return
		}
		retTopic := models.Topic{}
		// NOTE: to check that the correct structure is unmarshaled and the contents of the structure
		// are as expected, run this benchmark in the debugger with a breakpoint on the next line ...
		err = json.Unmarshal(payload, &retTopic)
		if err != nil {
			fmt.Printf("unmarshal fail\n")
			return
		}
	}
}

func BenchmarkGetDatasetPublic(b *testing.B) {

	// Given a topic API in publishing mode (private endpoints enabled)
	cfg, err := config.Get()
	if err != nil {
		fmt.Printf("config fail\n")
		return
	}
	cfg.EnablePrivateEndpoints = false
	//	And a topic API with mongoDB returning 'created' and 'full' topics

	mongoDBMock := &storeMock.MongoDBMock{
		GetTopicFunc: func(id string) (*models.TopicUpdate, error) {
			switch id {
			case testTopicID1:
				return dbTopic(models.StateTopicCreated), nil
			default:
				return nil, apierrors.ErrTopicNotFound
			}
		},
	}
	topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

	b.ReportAllocs()

	// test all event types
	for i := 0; i < b.N; i++ {

		// When an existing 'created' topic is requested with the valid Topic-Id context value
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/datasets/%s", testTopicID1), nil)

		w := httptest.NewRecorder()
		topicAPI.Router.ServeHTTP(w, request)
		// Then the expected topic is returned with status code 200
		payload, err := ioutil.ReadAll(w.Body)
		if err != nil {
			fmt.Printf("readall fail\n")
			return
		}
		retTopic := models.Topic{}
		// NOTE: to check that the correct structure is unmarshaled and the contents of the structure
		// are as expected, run this benchmark in the debugger with a breakpoint on the next line ...
		err = json.Unmarshal(payload, &retTopic)
		if err != nil {
			fmt.Printf("unmarshal fail\n")
			return
		}
	}
}
