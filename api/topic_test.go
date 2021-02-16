package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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
const (
	testTopicID1 = "topicTopicID1"
)

func dbTopicWithID(state models.State, id string) *models.TopicResponse {
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
					ID:   id,
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
					ID:   id,
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
func dbTopic(state models.State) *models.TopicResponse {
	return dbTopicWithID(state, testTopicID1)
}

// API model corresponding to TopicResponse
func createdTopicAll() *models.TopicResponse {
	return dbTopic(models.StateCreated)
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
				ID:   id,
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
func dbTopicCurrent(state models.State) *models.Topic {
	return dbTopicCurrentWithID(state, testTopicID1)
}

func createdTopicCurrent() *models.Topic {
	return dbTopicCurrent(models.StatePublished)
}

func TestGetTopicPublicHandler(t *testing.T) {

	Convey("Given a topic API in web mode (private endpoints disabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = false
		Convey("And a topic API with mongoDB returning 'next' and 'current' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(id string) (*models.TopicResponse, error) {
					switch id {
					case testTopicID1:
						return dbTopic(models.StatePublished), nil
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
				GetTopicFunc: func(id string) (*models.TopicResponse, error) {
					switch id {
					case testTopicID1:
						return dbTopic(models.StateCreated), nil
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
					retTopic := models.TopicResponse{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic, ShouldResemble, *createdTopicAll())
				})
			})

			Convey("Requesting an nonexistent topic ID results in a NotFound response", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/inexistent", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

// NOTE: The data within the following four sets of data returning functions
//       are interlinked with one another by the SubtopicsIds

// ================= - 1 has subtopics & points to 2 & 3
// DB model corresponding to a topic in the provided state, without any download variant
func dbTopic1(state models.State) *models.TopicResponse {
	return &models.TopicResponse{
		ID: "1",
		Next: &models.Topic{
			ID:    "1",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/1/subtopics",
				},
			},
			SubtopicIds: []string{"2", "3"},
		},
		Current: &models.Topic{
			ID:    "1",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/1/subtopics",
				},
			},
			SubtopicIds: []string{"2", "3"},
		},
	}
}

// ================= - 2 has subtopics & points to 4, 6 (but ID 6 does not exist)
// DB model corresponding to a topic in the provided state, without any download variant
func dbTopic2(state models.State) *models.TopicResponse {
	return &models.TopicResponse{
		ID: "2",
		Next: &models.Topic{
			ID:    "2",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/2/subtopics",
				},
			},
			SubtopicIds: []string{"4", "6"},
		},
		Current: &models.Topic{
			ID:    "2",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/2/subtopics",
				},
			},
			SubtopicIds: []string{"4", "6"},
		},
	}
}

// ================= - 3 has subtopics, but the ID 5 in the list does not exist
// DB model corresponding to a topic in the provided state, without any download variant
func dbTopic3(state models.State) *models.TopicResponse {
	return &models.TopicResponse{
		ID: "3",
		Next: &models.Topic{
			ID:    "3",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/3/subtopics",
				},
			},
			SubtopicIds: []string{"5"},
		},
		Current: &models.Topic{
			ID:    "3",
			State: state.String(),
			Links: &models.TopicLinks{
				Subtopics: &models.LinkObject{
					HRef: "http://example.com/topics/3/subtopics",
				},
			},
			SubtopicIds: []string{"5"},
		},
	}
}

// ================= - 4 has NO subtopics, so is an end node that has a content link
// DB model corresponding to a topic in the provided state, without any download variant
func dbTopic4(state models.State) *models.TopicResponse {
	return &models.TopicResponse{
		ID: "4",
		Next: &models.Topic{
			ID:    "4",
			State: state.String(),
			Links: &models.TopicLinks{
				Content: &models.LinkObject{
					HRef: "http://example.com/topics/4/content",
				},
			},
		},
		Current: &models.Topic{
			ID:    "4",
			State: state.String(),
			Links: &models.TopicLinks{
				Content: &models.LinkObject{
					HRef: "http://example.com/topics/4/content",
				},
			},
		},
	}
}

func TestGetSubtopicsPublicHandler(t *testing.T) {

	Convey("Given a topic API in web mode (private endpoints disabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = false
		Convey("And a topic API with mongoDB returning 'next' and 'current' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(id string) (*models.TopicResponse, error) {
					switch id {
					case "1":
						return dbTopic1(models.StatePublished), nil
					case "2":
						return dbTopic2(models.StatePublished), nil
					case "3":
						return dbTopic3(models.StatePublished), nil
					case "4":
						return dbTopic4(models.StatePublished), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}

			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			// 1 has subtopics & points to 2 & 3
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 1", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/1/subtopics", nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-documents is returned with status code 200, and documents with ID's 2 & 3 returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PublicSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 2)
					So((*retTopic.PublicItems)[0].ID, ShouldEqual, "2")
					So((*retTopic.PublicItems)[1].ID, ShouldEqual, "3")
				})
			})

			// 2 has subtopics & points to 4, 6 (but ID 6 does not exist)
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 2", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/2/subtopics", nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-document is returned with status code 200, and document with ID 4 is returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PublicSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 1)
					So((*retTopic.PublicItems)[0].ID, ShouldEqual, "4")
				})
			})

			// 3 has subtopics, but the ID 5 in the list does not exist
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 3", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/3/subtopics", nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no sub-documents are returned and we get status code 500", func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("content not found\n"))
				})
			})

			// 4 has NO subtopics, so is an end node that has a content link
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 4", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/4/subtopics", nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no sub-documents are returned and we get status code 404", func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("not found\n"))
				})
			})

			Convey("Requesting an nonexistent topic ID results in a NotFound response", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/inexistent/subtopics", nil)
				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

func TestGetSubtopicsPrivateHandler(t *testing.T) {

	Convey("Given a topic API in web mode (private endpoints enabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = true
		Convey("And a topic API with mongoDB returning 'next' and 'current' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(id string) (*models.TopicResponse, error) {
					switch id {
					case "1":
						return dbTopic1(models.StatePublished), nil
					case "2":
						return dbTopic2(models.StatePublished), nil
					case "3":
						return dbTopic3(models.StatePublished), nil
					case "4":
						return dbTopic4(models.StatePublished), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}

			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			// 1 has subtopics & points to 2 & 3
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 1", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/1/subtopics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-documents is returned with status code 200, and documents with ID's 2 & 3 returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PrivateSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 2)
					So((*retTopic.PrivateItems)[0].Current.ID, ShouldEqual, "2")
					So((*retTopic.PrivateItems)[1].Current.ID, ShouldEqual, "3")
				})
			})

			// 2 has subtopics & points to 4, 6 (but ID 6 does not exist)
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 2", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/2/subtopics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-document is returned with status code 200, and document with ID 4 is returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PrivateSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 1)
					So((*retTopic.PrivateItems)[0].Current.ID, ShouldEqual, "4")
				})
			})

			// 3 has subtopics, but the ID 5 in the list does not exist
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 3", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/3/subtopics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no sub-documents are returned and we get status code 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("internal error\n"))
				})
			})

			// 4 has NO subtopics, so is an end node that has a content link
			Convey("When an existing 'published' subtopic is requested with the valid Topic-Id value 4", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/4/subtopics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no sub-documents are returned and we get status code 404", func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("not found\n"))
				})
			})

			Convey("Requesting an nonexistent topic ID results in a NotFound response", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/inexistent/subtopics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

func TestGetTopicsListPublicHandler(t *testing.T) {
	Convey("Given a topic API in web mode (private endpoints disabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = false
		Convey("And a topic API with mongoDB returning 'next' and 'current' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(id string) (*models.TopicResponse, error) {
					switch id {
					case "2":
						return dbTopic2(models.StatePublished), nil
					case "3":
						return dbTopic3(models.StatePublished), nil
					case "topic_root":
						return dbTopic1(models.StatePublished), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}

			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			// topic_root for test uses 1 which has subtopics & points to 2 & 3
			Convey("When an existing 'published' /topics list is requested", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics", nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-documents is returned with status code 200, and documents with ID's 2 & 3 returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PublicSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 2)
					So((*retTopic.PublicItems)[0].ID, ShouldEqual, "2")
					So((*retTopic.PublicItems)[1].ID, ShouldEqual, "3")
				})
			})

			// No more tests needed because getting the 'topic_root' makes use of
			// function getSubtopicsPublicByID() which is fully tested in
			// TestGetSubtopicsPublicHandler() above, preventing duplication of tests.
		})
	})
}

func TestGetTopicsListPrivateHandler(t *testing.T) {

	Convey("Given a topic API in web mode (private endpoints enabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = true
		Convey("And a topic API with mongoDB returning 'next' and 'current' topics", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetTopicFunc: func(id string) (*models.TopicResponse, error) {
					switch id {
					case "2":
						return dbTopic2(models.StatePublished), nil
					case "3":
						return dbTopic3(models.StatePublished), nil
					case "topic_root":
						return dbTopic1(models.StatePublished), nil
					default:
						return nil, apierrors.ErrTopicNotFound
					}
				},
			}

			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			// topic_root for test uses 1 which has subtopics & points to 2 & 3
			Convey("When an existing 'published' /topics list is requested", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected sub-documents is returned with status code 200, and documents with ID's 2 & 3 returned", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retTopic := models.PrivateSubtopics{}
					err = json.Unmarshal(payload, &retTopic)
					So(err, ShouldBeNil)
					So(retTopic.TotalCount, ShouldEqual, 2)
					So((*retTopic.PrivateItems)[0].Current.ID, ShouldEqual, "2")
					So((*retTopic.PrivateItems)[1].Current.ID, ShouldEqual, "3")
				})
			})

			// No more tests needed because getting the 'topic_root' makes use of
			// function getSubtopicsPublicByID() which is fully tested in
			// TestGetSubtopicsPublicHandler() above, preventing duplication of tests.
		})
	})
}

// GetAPIWithMocks also used in other tests, so exported
func GetAPIWithMocks(cfg *config.Config, mockedDataStore store.Storer) *API {
	mu.Lock()
	defer mu.Unlock()
	//	urlBuilder := url.NewBuilder("http://example.com")

	permissions := mocks.NewAuthHandlerMock()

	return Setup(testContext, cfg, mux.NewRouter(), store.DataStore{Backend: mockedDataStore}, permissions)
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
		GetTopicFunc: func(id string) (*models.TopicResponse, error) {
			switch id {
			case testTopicID1:
				return dbTopic(models.StateCreated), nil
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
		retTopic := models.TopicResponse{}
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
		GetTopicFunc: func(id string) (*models.TopicResponse, error) {
			switch id {
			case testTopicID1:
				return dbTopic(models.StateCreated), nil
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
		retTopic := models.TopicResponse{}
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
		GetTopicFunc: func(id string) (*models.TopicResponse, error) {
			switch id {
			case testTopicID1:
				return dbTopic(models.StateCreated), nil
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
		GetTopicFunc: func(id string) (*models.TopicResponse, error) {
			switch id {
			case testTopicID1:
				return dbTopic(models.StateCreated), nil
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
