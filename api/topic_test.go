package api_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-topic-api/api/mock"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/models"

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

var errMongoDB = errors.New("MongoDB generic error")

// DB model corresponding to a topic in the provided state, without any download variant !!! fix this comment for topic-api
func dbTopic(state models.State) *models.Topic {
	return dbTopicWithId(state, testTopicID1)
}

func dbTopicWithId(state models.State, id string) *models.Topic {
	return &models.Topic{
		ID:          id,
		Description: "test description - 1",
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

// API model corresponding to dbCreatedTopic !!! ?
func createdTopic() *models.Topic {
	return dbTopic(models.StateTopicCreated)
}

func TestGetTopicHandler(t *testing.T) {

	Convey("Given a topic API in publishing mode", t, func() {
		cfg, err := config.Get()
		cfg.EnablePrivateEndpoints = true
		So(err, ShouldBeNil)
		doTestGetTopicHandler(cfg)
	})

	Convey("Given a topic API in web mode", t, func() {
		cfg, err := config.Get()
		cfg.EnablePrivateEndpoints = false
		So(err, ShouldBeNil)
		doTestGetTopicHandler(cfg)
	})
}

func doTestGetTopicHandler(cfg *config.Config) {

	Convey("And a topic API with mongoDB returning 'created' and 'published' topics", func() {

		mongoDBMock := &mock.MongoServerMock{
			GetTopicFunc: func(ctx context.Context, id string) (*models.Topic, error) {
				switch id {
				case testTopicID1:
					return dbTopic(models.StateTopicCreated), nil //!!! might want to change this to StateTopicTrue
					//				case testImageID2:
					//					return dbFullImageWithDownloads(models.StateTopicPublished, dbDownload(models.StateDownloadPublished)), nil
				default:
					return nil, apierrors.ErrTopicNotFound
				}
			},
		}
		/*authHandlerMock := &mock.AuthHandlerMock{
			RequireFunc: func(required dpauth.Permissions, handler http.HandlerFunc) http.HandlerFunc {
				return handler
			},
		}*/
		topicApi := GetAPIWithMocks(cfg, mongoDBMock /*, authHandlerMock*/)

		// !!! this message probably needs changing once the system implements the ret of the spec more fully.
		Convey("When an existing 'created' topic is requested with the valid Collection-Id context value", func() {
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s", testTopicID1), nil)
			//			r = r.WithContext(context.WithValue(r.Context(), dphttp.FlorenceIdentityKey, testUserAuthToken))
			w := httptest.NewRecorder()
			topicApi.Router.ServeHTTP(w, r)
			Convey("Then the expected topic is returned with status code 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				retTopic := models.Topic{}
				err = json.Unmarshal(payload, &retTopic)
				So(err, ShouldBeNil)
				So(retTopic, ShouldResemble, *createdTopic())
			})
		})

		// !!! this message probably needs changing once the system implements the rest of the spec more fully.

		/*		Convey("When an existing 'published' topic is requested without a Collection-Id context value", func() {
				r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/images/%s", testImageID2), nil)
				r = r.WithContext(context.WithValue(r.Context(), dphttp.FlorenceIdentityKey, testUserAuthToken))
				w := httptest.NewRecorder()
				topicApi.Router.ServeHTTP(w, r)
				Convey("Then the published topic is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retImage := models.Image{}
					err = json.Unmarshal(payload, &retImage)
					So(err, ShouldBeNil)
					So(retImage, ShouldResemble, *apiFullImage(models.StateTopicPublished))
				})
			})*/

		Convey("Requesting an nonexistent topic ID results in a NotFound response", func() {
			r := httptest.NewRequest(http.MethodGet, "http://localhost:24700/topics/inexistent", nil)
			//			r = r.WithContext(context.WithValue(r.Context(), dphttp.FlorenceIdentityKey, testUserAuthToken))
			w := httptest.NewRecorder()
			topicApi.Router.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, http.StatusNotFound)
		})
	})
}
