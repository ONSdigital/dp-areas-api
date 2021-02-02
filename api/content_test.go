package api

//!!! sort the code for this, rename / adjust for /topic/<id>/content

// !!! make sure tests achieve full coverage of api/content.go

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/models"
	storeMock "github.com/ONSdigital/dp-topic-api/store/datastoretest"

	. "github.com/smartystreets/goconvey/convey"
)

// Constants for testing
const ( // !!! remove not needed const's at some point / fix / rename, etc
	ctestContentID1         = "ContentID1"
	ctestContentID2         = "ContentID2"
	ctestContentCreatedID   = "ContentCreatedID"
	ctestContentPublishedID = "ContentPublishedID"
)

const (
	chost = "http://localhost:25300"
)

// build up response from following:
// Given this mongo collection document:
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
                "Href": "/article/12345",
                "Title": "Some article 2"
            },
            {
                "Href": "/article/1234",
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
                "Href": "/article/12345",
                "Title": "Some article 2"
            },
            {
                "Href": "/article/12344",
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
// NOTE: the above has to be on one line ...
// NOTE: The following HAS to be on ONE line for unmarshal to work (and all the inner double quotes need escaping)
var mongoContentJSONResponse1 string = "{\"id\": \"4\", \"next\": {\"spotlight\": [ {\"Href\": \"/article/123\", \"Title\": \"Some article\"}, {\"Href\": \"/dataset/12fasf3\", \"Title\": \"An interesting dataset\"} ], \"articles\": [ {\"Href\": \"/article/12345\", \"Title\": \"Some article 2\"}, {\"Href\": \"/article/1234\", \"Title\": \"Some article 3\"} ], \"bulletins\": [ {\"Href\": \"/bulletins/this-month-hurray\", \"Title\": \"This Months Bulletin\"} ], \"timeseries\": [ {\"Href\": \"/timseries/KVAC\", \"Title\": \"CPIH Time series\" } ], \"state\" : \"in_progress\" }, \"current\" : {\"spotlight\": [ {\"Href\": \"/article/123\", \"Title\": \"Some article\"}, { \"Href\": \"/dataset/12fasf3\", \"Title\": \"An interesting dataset\" } ], \"articles\": [ { \"Href\": \"/article/12345\", \"Title\": \"Some article 3\" }, { \"Href\": \"/article/1234\", \"Title\": \"Some article 2\" } ], \"bulletins\": [ { \"Href\": \"/bulletins/this-month-hurray\", \"Title\": \"This Months Bulletin\" } ], \"timeseries\": [ { \"Href\": \"/timseries/KVAC\", \"Title\": \"CPIH Time series\" } ], \"state\" : \"published\" } }"

// then the Get Response in Public would look like (and note article is sorted by href):
/*
{
    "offset": 0,
    "count": 6,
    "total_count": 6,
    “limit”: 0
    "items": [
        {
            "title": "Some article",
            "type": "spotlight",
            "links": {
                "self": {
                    "href": "/article/123"
                },
                "topic": {
                    "href": "/topic/4",
                    "id": "4"
                }
            }
        },
        {
            "title": "An interesting dataset",
            "type": "spotlight",
            "links": {
                "self": {
                    "href": "/dataset/12fasf3"
                },
                "topic": {
                    "href": "/topic/4",
                    "id": "4"
                }
            }
        },
        {
            "title": "Some article 2",
            "type": "article",
            "links": {
                "self": {
                    "href": "/article/1234"
                },
                "topic": {
                    "href": "/topic/4",
                    "id": "4"
                }
            }
        },
        {
            "title": "Some article 3",
            "type": "article",
            "links": {
                "self": {
                    "href": "/article/12345"
                },
                "topic": {
                    "href": "/topic/4",
                    "id": "4"
                }
            }
        },
        {
            "title": "This Months Bulletin",
            "type": "bulletin",
            "links": {
                "self": {
                    "href": "/bulletins/this-month-hurray"
                },
                "topic": {
                    "href": "/topic/4",
                    "id": "4"
                }
            }
        },
        {
            "title": "CPIH Time series",
            "type": "timeseries",
            "links": {
                "self": {
                    "href": "/timseries/KVAC"
                },
                "topic": {
                    "href": "/topic/4",
                    "id": "4"
                }
            }
        }
    ]
}
*/

// =======

// Given this mongo collection document: (with 'current' missing)
/*
{
    "id": "5",
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
        "state" : "in_progress"
    }
}
*/
// NOTE: the above has to be on one line ...
// NOTE: The following HAS to be on ONE line for unmarshal to work (and all the inner double quotes need escaping)
var mongoContentJSONResponse2 string = "{\"id\": \"5\", \"next\": {\"spotlight\": [ {\"Href\": \"/article/123\", \"Title\": \"Some article\"}, {\"Href\": \"/dataset/12fasf3\", \"Title\": \"An interesting dataset\" } ], \"state\" : \"in_progress\"} }"

// then the Get Response in Public would return a 500 error, as content.Current = nil

// =======

func dbContentWithID(state models.State, id string) *models.ContentResponse {
	var response models.ContentResponse

	switch id {
	case ctestContentID1:
		err := json.Unmarshal([]byte(mongoContentJSONResponse1), &response)
		if err != nil {
			fmt.Printf("Oops coding error in 'dbContentWithID', FIX the json 'mongoContentJSONResponse1' so that it will unmarshal correctly !")
			os.Exit(1)
		}
	case ctestContentID2:
		err := json.Unmarshal([]byte(mongoContentJSONResponse2), &response)
		if err != nil {
			fmt.Printf("Oops coding error in 'dbContentWithID', FIX the json 'mongoContentJSONResponse2' so that it will unmarshal correctly !")
			os.Exit(1)
		}
	}
	response.ID = id

	return &response
}

// DB model corresponding to content in the provided state, without any download variant
func dbContent(state models.State) *models.ContentResponse {
	return dbContentWithID(state, ctestContentID1)
}

func dbContent2(state models.State) *models.ContentResponse {
	return dbContentWithID(state, ctestContentID2)
}

// API model corresponding to ContentResponse
func createdContentAll() *models.ContentResponse {
	return dbContent(models.StateTopicCreated)
}

/*func dbContentCurrentWithID(state models.State, id string) *models.Content {
	var response models.ContentResponse

	err := json.Unmarshal([]byte(mongoContentJSONResponse1), &response)
	if err != nil {
		fmt.Printf("Oops coding error in 'dbContentWithID', FIX the json so that it will unmarshal correctly !")
		os.Exit(1)
	}
	response.ID = id
	response.Next.State = state.String()
	response.Current.State = state.String()

	return response.Current
}

// create just the 'current' sub-document
func dbContentCurrent(state models.State) *models.Content {
	return dbContentCurrentWithID(state, ctestContentID1)
}

func createdContentCurrent() *models.Content {
	return dbContentCurrent(models.StateTopicPublished) //!!! this probably needs to be generic state that is same for Topic and Content
}*/

// TestGetContentPublicHandler - does what the function name says
func TestGetContentPublicHandler(t *testing.T) {

	Convey("Given a content API in web mode (private endpoints disabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = false
		Convey("And a content API with mongoDB returning 'next' and 'current' content", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetContentFunc: func(id string) (*models.ContentResponse, error) {
					switch id {
					case ctestContentID1:
						return dbContent(models.StateTopicPublished), nil
					case ctestContentID2:
						return dbContent2(models.StateTopicPublished), nil
					default:
						return nil, apierrors.ErrContentNotFound
					}
				},
			}

			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content", ctestContentID1), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected content is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retContent := models.ContentResponseAPI{}
					err = json.Unmarshal(payload, &retContent)
					So(err, ShouldBeNil)
					So(retContent.Items, ShouldNotBeNil)
					So(retContent.Count, ShouldEqual, 6)
					So(retContent.Offset, ShouldEqual, 0)
					So(retContent.Limit, ShouldEqual, 0)
					So(retContent.TotalCount, ShouldEqual, 6)
					So(len(*retContent.Items), ShouldEqual, 6)
					// check result is sorted by Href
					So((*retContent.Items)[2].Links.Self.HRef, ShouldEqual, "/article/1234")
					So((*retContent.Items)[3].Links.Self.HRef, ShouldEqual, "/article/12345")
				})
			})

			Convey("When an existing 'published' content (with no content) is requested with the valid Topic-Id context value", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content", ctestContentID2), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no content is returned and status code 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			//!!! create test where TotalCount = 0   and check code coverage in goland

			Convey("Requesting an nonexistent content ID results in a NotFound response", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/inexistent/content", nil)
				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

func TestGetContentPrivateHandler(t *testing.T) {

	Convey("Given a content API in publishing mode (private endpoints enabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = true
		Convey("And a content API with mongoDB returning 'next' and 'current' content", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetContentFunc: func(id string) (*models.ContentResponse, error) {
					switch id {
					case ctestContentID1:
						return dbContent(models.StateTopicCreated), nil
					default:
						return nil, apierrors.ErrContentNotFound
					}
				},
			}
			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			Convey("When an existing 'created' content is requested with the valid Topic-Id context value", func() {
				request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content", ctestContentID1), nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected content is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retContentResponse := models.PrivateContentResponseAPI{}
					err = json.Unmarshal(payload, &retContentResponse)
					So(err, ShouldBeNil)
					//!!! fix following
					//		So(retContentResponse, ShouldResemble, *createdContentAll())
				})
			})

			Convey("Requesting an nonexistent content ID results in a NotFound response", func() {
				request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/inexistent/content"), nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}
