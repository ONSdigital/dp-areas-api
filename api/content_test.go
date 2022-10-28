package api

import (
	"context"
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
	storeMock "github.com/ONSdigital/dp-topic-api/store/mock"

	. "github.com/smartystreets/goconvey/convey"
)

// Constants for testing
const (
	ctestContentID1 = "ContentID1"
	ctestContentID2 = "ContentID2"
	ctestContentID3 = "ContentID3"
	ctestContentID4 = "ContentID4"
	ctestContentID5 = "ContentID5"
	ctestContentID7 = "ContentID7"
	ctestContentID8 = "ContentID8"
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
        "state" : "published"
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
var mongoContentJSONResponse1 string = "{\"id\": \"4\", \"next\": {\"spotlight\": [ {\"Href\": \"/article/123\", \"Title\": \"Some article\"}, {\"Href\": \"/dataset/12fasf3\", \"Title\": \"An interesting dataset\"} ], \"articles\": [ {\"Href\": \"/article/12345\", \"Title\": \"Some article 2\"}, {\"Href\": \"/article/1234\", \"Title\": \"Some article 3\"} ], \"bulletins\": [ {\"Href\": \"/bulletins/this-month-hurray\", \"Title\": \"This Months Bulletin\"} ], \"timeseries\": [ {\"Href\": \"/timseries/KVAC\", \"Title\": \"CPIH Time series\" } ], \"state\" : \"published\" }, \"current\" : {\"spotlight\": [ {\"Href\": \"/article/123\", \"Title\": \"Some article\"}, { \"Href\": \"/dataset/12fasf3\", \"Title\": \"An interesting dataset\" } ], \"articles\": [ { \"Href\": \"/article/12345\", \"Title\": \"Some article 3\" }, { \"Href\": \"/article/1234\", \"Title\": \"Some article 2\" } ], \"bulletins\": [ { \"Href\": \"/bulletins/this-month-hurray\", \"Title\": \"This Months Bulletin\" } ], \"timeseries\": [ { \"Href\": \"/timseries/KVAC\", \"Title\": \"CPIH Time series\" } ], \"state\" : \"published\" } }"

// then the Get Response in Public would look like (and note article is sorted by href):
// (in Private mode, Next & Current contain the following)
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
        "state" : "published"
    }
}
*/
// NOTE: the above has to be on one line ...
// NOTE: The following HAS to be on ONE line for unmarshal to work (and all the inner double quotes need escaping)
var mongoContentJSONResponse2 string = "{\"id\": \"5\", \"next\": {\"spotlight\": [ {\"Href\": \"/article/123\", \"Title\": \"Some article\"}, {\"Href\": \"/dataset/12fasf3\", \"Title\": \"An interesting dataset\" } ], \"state\" : \"published\"} }"

// then the Get Response in Public would return a 500 error, as content.Current = nil
// (Private also returns 500)

// =======

// Given this mongo collection document: (with no items)
/*
{
    "id": "6",
    "next": {
        "state" : "published"
    },
    "current" : {
        "state" : "published"
    }
}
*/
// NOTE: the above has to be on one line ...
// NOTE: The following HAS to be on ONE line for unmarshal to work (and all the inner double quotes need escaping)
var mongoContentJSONResponse3 string = "{\"id\": \"4\", \"next\": {\"state\" : \"published\"}, \"current\" : {\"state\" : \"published\"} }"

// then the Get Response in Public would this, where TotalCount = 0
// (in Private mode, Next & Current contain the following)
/*
{
    "offset": 0,
    "count": 0,
    "total_count": 0,
    “limit”: 0
    "items": [
    ]
}
*/

// =======

// Given this mongo collection document: (with 'next' missing)
/*
{
    "id": "7",
    "current": {
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
        "state" : "published"
    }
}
*/
// NOTE: the above has to be on one line ...
// NOTE: The following HAS to be on ONE line for unmarshal to work (and all the inner double quotes need escaping)
var mongoContentJSONResponse4 string = "{\"id\": \"5\", \"current\": {\"spotlight\": [ {\"Href\": \"/article/123\", \"Title\": \"Some article\"}, {\"Href\": \"/dataset/12fasf3\", \"Title\": \"An interesting dataset\" } ], \"state\" : \"published\"} }"

// then the Get Response in Private would return a 500 error, as content.Next = nil

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
	case ctestContentID3:
		err := json.Unmarshal([]byte(mongoContentJSONResponse3), &response)
		if err != nil {
			fmt.Printf("Oops coding error in 'dbContentWithID', FIX the json 'mongoContentJSONResponse3' so that it will unmarshal correctly !")
			os.Exit(1)
		}
	case ctestContentID4:
		err := json.Unmarshal([]byte(mongoContentJSONResponse4), &response)
		if err != nil {
			fmt.Printf("Oops coding error in 'dbContentWithID', FIX the json 'mongoContentJSONResponse4' so that it will unmarshal correctly !")
			os.Exit(1)
		}
	case ctestContentID7:
		err := json.Unmarshal([]byte(mongoContentJSONResponse7), &response)
		if err != nil {
			fmt.Printf("Oops coding error in 'dbContentWithID', FIX the json 'mongoContentJSONResponse7' so that it will unmarshal correctly !")
			os.Exit(1)
		}
	case ctestContentID8:
		err := json.Unmarshal([]byte(mongoContentJSONResponse8), &response)
		if err != nil {
			fmt.Printf("Oops coding error in 'dbContentWithID', FIX the json 'mongoContentJSONResponse8' so that it will unmarshal correctly !")
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

func dbContent3(state models.State) *models.ContentResponse {
	return dbContentWithID(state, ctestContentID3)
}

func dbContent4(state models.State) *models.ContentResponse {
	return dbContentWithID(state, ctestContentID4)
}

func dbContent7(state models.State) *models.ContentResponse {
	return dbContentWithID(state, ctestContentID7)
}

func dbContent8(state models.State) *models.ContentResponse {
	return dbContentWithID(state, ctestContentID8)
}

// TestGetContentPublicHandler - does what the function name says
func TestGetContentPublicHandler(t *testing.T) {

	Convey("Given a content API in web mode (private endpoints disabled)", t, func() {
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		cfg.EnablePrivateEndpoints = false
		Convey("And a content API with mongoDB returning 'next' and 'current' content", func() {

			mongoDBMock := &storeMock.MongoDBMock{
				GetContentFunc: func(ctx context.Context, id string, queryTypeFlags int) (*models.ContentResponse, error) {
					switch id {
					case ctestContentID1:
						return dbContent(models.StatePublished), nil
					case ctestContentID2:
						return dbContent2(models.StatePublished), nil
					case ctestContentID3:
						return dbContent3(models.StatePublished), nil
					case ctestContentID7:
						return dbContent7(models.StatePublished), nil
					case ctestContentID8:
						return dbContent8(models.StatePublished), nil
					default:
						return nil, apierrors.ErrContentNotFound
					}
				},
				CheckTopicExistsFunc: func(ctx context.Context, id string) error {
					switch id {
					case ctestContentID1,
						ctestContentID2,
						ctestContentID3,
						ctestContentID5,
						ctestContentID7,
						ctestContentID8:
						return nil
					default:
						return apierrors.ErrTopicNotFound
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

			Convey("When an existing 'published' content (with no current) is requested with the valid Topic-Id context value", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content", ctestContentID2), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no content is returned and status code 500", func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("content not found\n"))
				})
			})

			Convey("When an existing 'published' content (with no items in current) is requested with the valid Topic-Id context value", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content", ctestContentID3), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no content is returned with status code 404", func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("content not found\n"))
				})
			})

			// the following two tests cover different failure modes and code coverage in 'getContentPublicHandler'
			Convey("Requesting an nonexistent content & topic ID results in a NotFound response (topic read fails)", func() {
				request := httptest.NewRequest(http.MethodGet, "http://localhost:25300/topics/inexistent/content", nil)
				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				So(payload, ShouldResemble, []byte("topic not found\n"))
			})

			Convey("Requesting an nonexistent content ID results in a NotFound response (content read fails)", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content", ctestContentID5), nil)
				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				So(payload, ShouldResemble, []byte("content not found\n"))
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query type: spotlight", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=spotlight", ctestContentID7), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected query type is returned with status code 200, and result is sorted", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retContent := models.ContentResponseAPI{}
					err = json.Unmarshal(payload, &retContent)
					So(err, ShouldBeNil)
					So(retContent.Items, ShouldNotBeNil)
					So(retContent.Count, ShouldEqual, 2)
					So(retContent.Offset, ShouldEqual, 0)
					So(retContent.Limit, ShouldEqual, 0)
					So(retContent.TotalCount, ShouldEqual, 2)
					So(len(*retContent.Items), ShouldEqual, 2)
					// check result is sorted by unique Href
					So((*retContent.Items)[0].Links.Self.HRef, ShouldEqual, "/h1")
					So((*retContent.Items)[1].Links.Self.HRef, ShouldEqual, "/h2")
					So((*retContent.Items)[1].Links.Self.ID, ShouldEqual, "")
					So((*retContent.Items)[1].Links.Topic.HRef, ShouldEqual, "/topic/"+ctestContentID7)
					So((*retContent.Items)[1].Links.Topic.ID, ShouldEqual, ctestContentID7)
					So((*retContent.Items)[1].Title, ShouldEqual, "Labour disputes")
					So((*retContent.Items)[1].Type, ShouldEqual, spotlightStr)
					So((*retContent.Items)[1].State, ShouldEqual, "published")
				})
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query type: articles", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=articles", ctestContentID7), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected query type is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retContent := models.ContentResponseAPI{}
					err = json.Unmarshal(payload, &retContent)
					So(err, ShouldBeNil)
					// check result is sorted by unique Href
					So((*retContent.Items)[0].Links.Self.HRef, ShouldEqual, "/a1")
				})
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query type: bulletins", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=bulletins", ctestContentID7), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected query type is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retContent := models.ContentResponseAPI{}
					err = json.Unmarshal(payload, &retContent)
					So(err, ShouldBeNil)
					// check result is sorted by unique Href
					So((*retContent.Items)[0].Links.Self.HRef, ShouldEqual, "/b1")
				})
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query type: methodologies", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=methodologies", ctestContentID7), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected query type is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retContent := models.ContentResponseAPI{}
					err = json.Unmarshal(payload, &retContent)
					So(err, ShouldBeNil)
					// check result is sorted by unique Href
					So((*retContent.Items)[0].Links.Self.HRef, ShouldEqual, "/m1")
				})
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query type: methodologyarticles", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=methodologyarticles", ctestContentID7), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected query type is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retContent := models.ContentResponseAPI{}
					err = json.Unmarshal(payload, &retContent)
					So(err, ShouldBeNil)
					// check result is sorted by unique Href
					So((*retContent.Items)[0].Links.Self.HRef, ShouldEqual, "/ma1")
				})
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query type: staticdatasets", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=staticdatasets", ctestContentID7), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected query type is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retContent := models.ContentResponseAPI{}
					err = json.Unmarshal(payload, &retContent)
					So(err, ShouldBeNil)
					// check result is sorted by unique Href
					So((*retContent.Items)[0].Links.Self.HRef, ShouldEqual, "/s1")
				})
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query with wrong type: fred", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=fred", ctestContentID7), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected empty response is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusBadRequest)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("content query not recognised\n"))
				})
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query type: publications", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=publications", ctestContentID7), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected query type is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retContent := models.ContentResponseAPI{}
					err = json.Unmarshal(payload, &retContent)
					So(err, ShouldBeNil)

					So(retContent.TotalCount, ShouldEqual, 4)

					// check result is sorted by unique Href
					So((*retContent.Items)[0].Links.Self.HRef, ShouldEqual, "/a1")
					So((*retContent.Items)[0].Type, ShouldEqual, articlesStr)
					So((*retContent.Items)[1].Links.Self.HRef, ShouldEqual, "/b1")
					So((*retContent.Items)[1].Type, ShouldEqual, bulletinsStr)
					So((*retContent.Items)[2].Links.Self.HRef, ShouldEqual, "/m1")
					So((*retContent.Items)[2].Type, ShouldEqual, methodologiesStr)
					So((*retContent.Items)[3].Links.Self.HRef, ShouldEqual, "/ma1")
					So((*retContent.Items)[3].Type, ShouldEqual, methodologyarticlesStr)
				})
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query type: datasets", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=datasets", ctestContentID7), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected query type is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retContent := models.ContentResponseAPI{}
					err = json.Unmarshal(payload, &retContent)
					So(err, ShouldBeNil)

					So(retContent.TotalCount, ShouldEqual, 2)

					// check result is sorted by unique Href
					So((*retContent.Items)[0].Links.Self.HRef, ShouldEqual, "/s1")
					So((*retContent.Items)[0].Type, ShouldEqual, staticdatasetsStr)
					So((*retContent.Items)[1].Links.Self.HRef, ShouldEqual, "/t1")
					So((*retContent.Items)[1].Type, ShouldEqual, timeseriesStr)
				})
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query type: articles AND prefix and postfix spaces in query", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=%%20articles%%20", ctestContentID7), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then the expected query type is returned with status code 200", func() {
					So(w.Code, ShouldEqual, http.StatusOK)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					retContent := models.ContentResponseAPI{}
					err = json.Unmarshal(payload, &retContent)
					So(err, ShouldBeNil)
					// check result is sorted by unique Href
					So((*retContent.Items)[0].Links.Self.HRef, ShouldEqual, "/a1")
				})
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query type: spotlight AND page has not content", func() {
				request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=spotlight", ctestContentID8), nil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				So(payload, ShouldResemble, []byte("content not found\n"))
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
				GetContentFunc: func(ctx context.Context, id string, queryTypeFlags int) (*models.ContentResponse, error) {
					switch id {
					case ctestContentID1:
						return dbContent(models.StatePublished), nil
					case ctestContentID2:
						return dbContent2(models.StatePublished), nil
					case ctestContentID3:
						return dbContent3(models.StatePublished), nil
					case ctestContentID4:
						return dbContent4(models.StatePublished), nil
					case ctestContentID7:
						return dbContent7(models.StatePublished), nil
					case ctestContentID8:
						return dbContent8(models.StatePublished), nil
					default:
						return nil, apierrors.ErrContentNotFound
					}
				},
				CheckTopicExistsFunc: func(ctx context.Context, id string) error {
					switch id {
					case ctestContentID1,
						ctestContentID2,
						ctestContentID3,
						ctestContentID4,
						ctestContentID5,
						ctestContentID7,
						ctestContentID8:
						return nil
					default:
						return apierrors.ErrTopicNotFound
					}
				},
			}
			topicAPI := GetAPIWithMocks(cfg, mongoDBMock)

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value", func() {
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

					So(retContentResponse.Next.Items, ShouldNotBeNil)
					So(retContentResponse.Next.Count, ShouldEqual, 6)
					So(retContentResponse.Next.Offset, ShouldEqual, 0)
					So(retContentResponse.Next.Limit, ShouldEqual, 0)
					So(retContentResponse.Next.TotalCount, ShouldEqual, 6)
					So(len(*retContentResponse.Next.Items), ShouldEqual, 6)
					// check result is sorted by Href
					So((*retContentResponse.Next.Items)[2].Links.Self.HRef, ShouldEqual, "/article/1234")
					So((*retContentResponse.Next.Items)[3].Links.Self.HRef, ShouldEqual, "/article/12345")

					So(retContentResponse.Current.Items, ShouldNotBeNil)
					So(retContentResponse.Current.Count, ShouldEqual, 6)
					So(retContentResponse.Current.Offset, ShouldEqual, 0)
					So(retContentResponse.Current.Limit, ShouldEqual, 0)
					So(retContentResponse.Current.TotalCount, ShouldEqual, 6)
					So(len(*retContentResponse.Current.Items), ShouldEqual, 6)
					// check result is sorted by Href
					So((*retContentResponse.Current.Items)[2].Links.Self.HRef, ShouldEqual, "/article/1234")
					So((*retContentResponse.Current.Items)[3].Links.Self.HRef, ShouldEqual, "/article/12345")
				})
			})

			Convey("When an existing 'published' content (with no current) is requested with the valid Topic-Id context value", func() {
				request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content", ctestContentID2), nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no content is returned and status code 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When an existing 'published' content (with no next) is requested with the valid Topic-Id context value", func() {
				request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content", ctestContentID4), nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no content is returned and status code 500", func() {
					So(w.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})

			Convey("When an existing 'published' content (with no items in next and current) is requested with the valid Topic-Id context value", func() {
				request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content", ctestContentID3), nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				Convey("Then no content is returned with status code 404", func() {
					So(w.Code, ShouldEqual, http.StatusNotFound)
					payload, err := ioutil.ReadAll(w.Body)
					So(err, ShouldBeNil)
					So(payload, ShouldResemble, []byte("content not found\n"))
				})
			})

			// the following two tests cover different failure modes and code coverage in 'getContentPublicHandler'
			Convey("Requesting an nonexistent content & topic ID results in a NotFound response (topic read fails)", func() {
				request, err := createRequestWithAuth(http.MethodGet, "http://localhost:25300/topics/inexistent/content", nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				So(payload, ShouldResemble, []byte("topic not found\n"))
			})

			Convey("Requesting an nonexistent content ID results in a NotFound response (content read fails)", func() {
				request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content", ctestContentID5), nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				So(payload, ShouldResemble, []byte("content not found\n"))
			})

			Convey("When an existing 'published' content is requested with the valid Topic-Id context value for a query type: spotlight AND page has not content for next and current", func() {
				request, err := createRequestWithAuth(http.MethodGet, fmt.Sprintf("http://localhost:25300/topics/%s/content?type=spotlight", ctestContentID8), nil)
				So(err, ShouldBeNil)

				w := httptest.NewRecorder()
				topicAPI.Router.ServeHTTP(w, request)
				So(w.Code, ShouldEqual, http.StatusNotFound)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				So(payload, ShouldResemble, []byte("content not found\n"))
			})
		})
	})
}

// Given this mongo collection document that contains examples of all types and two of 'spotlight'
// (this is used for query tests)
/*
{
    "id": "7",
    "next": {
        "state": "published",
        "spotlight": [
            {
                "href": "/h2",
                "title": "Labour disputes"
            },
            {
                "href": "/h1",
                "title": "Labour disputes in the UK"
            }
        ],
        "articles": [
            {
                "href": "/a1",
                "title": "Labour disputes in the UK1"
            }
        ],
        "bulletins": [
            {
                "href": "/b1",
                "title": "Labour market overview, UK2"
            }
        ],
        "methodologies": [
            {
                "href": "/m1",
                "title": "** broken **"
            }
        ],
        "methodology_articles": [
            {
                "href": "/ma1",
                "title": "Labour Disputes Inquiry QMI"
            }
        ],
        "static_datasets": [
            {
                "href": "/s1",
                "title": "LABD01: Labour disputes"
            }
        ],
        "timeseries": [
            {
                "href": "/t1",
                "title": "Labour disputes;UK"
            }
        ]
    },
    "current": {
        "state": "published",
        "spotlight": [
            {
                "href": "/h2",
                "title": "Labour disputes"
            },
            {
                "href": "/h1",
                "title": "Labour disputes in the UK"
            }
        ],
        "articles": [
            {
                "href": "/a1",
                "title": "Labour disputes in the UK1"
            }
        ],
        "bulletins": [
            {
                "href": "/b1",
                "title": "Labour market overview, UK2"
            }
        ],
        "methodologies": [
            {
                "href": "/m1",
                "title": "** broken **"
            }
        ],
        "methodology_articles": [
            {
                "href": "/ma1",
                "title": "Labour Disputes Inquiry QMI"
            }
        ],
        "static_datasets": [
            {
                "href": "/s1",
                "title": "LABD01: Labour disputes"
            }
        ],
        "timeseries": [
            {
                "href": "/t1",
                "title": "Labour disputes;UK"
            }
        ]
    }
}
*/

// NOTE: the above has to be on one line ...
// NOTE: The following HAS to be on ONE line for unmarshal to work (and all the inner double quotes need escaping)
var mongoContentJSONResponse7 string = "{\"id\": \"workplacedisputesandworkingconditions\",\"next\": {\"state\": \"published\",\"spotlight\": [{\"href\": \"/h2\",\"title\": \"Labour disputes\"},{\"href\": \"/h1\",\"title\": \"Labour disputes in the UK\"}],\"articles\": [{\"href\": \"/a1\",\"title\": \"Labour disputes in the UK1\"}],\"bulletins\": [{\"href\": \"/b1\",\"title\": \"Labour market overview, UK2\"}],\"methodologies\": [{\"href\": \"/m1\",\"title\": \"** broken **\"}],\"methodology_articles\": [{\"href\": \"/ma1\",\"title\": \"Labour Disputes Inquiry QMI\"}],\"static_datasets\": [{\"href\": \"/s1\",\"title\": \"LABD01: Labour disputes\"}],\"timeseries\": [{\"href\": \"/t1\",\"title\": \"Labour disputes;UK\"}]},\"current\": {\"state\": \"published\",\"state\": \"published\",\"spotlight\": [{\"href\": \"/h2\",\"title\": \"Labour disputes\"},{\"href\": \"/h1\",\"title\": \"Labour disputes in the UK\"}],\"articles\": [{\"href\": \"/a1\",\"title\": \"Labour disputes in the UK1\"}],\"bulletins\": [{\"href\": \"/b1\",\"title\": \"Labour market overview, UK2\"}],\"methodologies\": [{\"href\": \"/m1\",\"title\": \"** broken **\"}],\"methodology_articles\": [{\"href\": \"/ma1\",\"title\": \"Labour Disputes Inquiry QMI\"}],\"static_datasets\": [{\"href\": \"/s1\",\"title\": \"LABD01: Labour disputes\"}],\"timeseries\": [{\"href\": \"/t1\",\"title\": \"Labour disputes;UK\"}]}}"

// then the Get Response in Public would look like (and note spotlight is sorted by href):
// (in Private mode, Next & Current contain the following)
/*
{
    "next": {
        "count": 8,
        "offset_index": 0,
        "limit": 0,
        "total_count": 8,
        "items": [
            {
                "title": "Labour disputes",
                "type": "spotlight",
                "links": {
                    "self": {
                        "href": "/h1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                },
                "state": "published"
            },
            {
                "title": "Labour disputes in the UK",
                "type": "spotlight",
                "links": {
                    "self": {
                        "href": "/h2"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                },
                "state": "published"
            },
            {
                "title": "Labour disputes in the UK1",
                "type": "articles",
                "links": {
                    "self": {
                        "href": "/a1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            },
            {
                "title": "Labour disputes in the UK2",
                "type": "bulletins",
                "links": {
                    "self": {
                        "href": "/b1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            },
            {
                "title": "** broken **",
                "type": "methodologies",
                "links": {
                    "self": {
                        "href": "/m1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            },
            {
                "title": "Labour Disputes Inquiry QMI",
                "type": "methodologyArticles",
                "links": {
                    "self": {
                        "href": "/ma1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            },
            {
                "title": "LABD01: Labour disputes",
                "type": "staticDatasets",
                "links": {
                    "self": {
                        "href": "/s1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            },
            {
                "title": "Labour disputes;UK",
                "type": "timeseries",
                "links": {
                    "self": {
                        "href": "/t1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            }
        ]
    },
    "current": {
        "count": 8,
        "offset_index": 0,
        "limit": 0,
        "total_count": 8,
        "items": [
            {
                "title": "Labour disputes",
                "type": "spotlight",
                "links": {
                    "self": {
                        "href": "/h1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                },
                "state": "published"
            },
            {
                "title": "Labour disputes in the UK",
                "type": "spotlight",
                "links": {
                    "self": {
                        "href": "/h2"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                },
                "state": "published"
            },
            {
                "title": "Labour disputes in the UK1",
                "type": "articles",
                "links": {
                    "self": {
                        "href": "/a1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            },
            {
                "title": "Labour disputes in the UK2",
                "type": "bulletins",
                "links": {
                    "self": {
                        "href": "/b1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            },
            {
                "title": "** broken **",
                "type": "methodologies",
                "links": {
                    "self": {
                        "href": "/m1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            },
            {
                "title": "Labour Disputes Inquiry QMI",
                "type": "methodologyArticles",
                "links": {
                    "self": {
                        "href": "/ma1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            },
            {
                "title": "LABD01: Labour disputes",
                "type": "staticDatasets",
                "links": {
                    "self": {
                        "href": "/s1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            },
            {
                "title": "Labour disputes;UK",
                "type": "timeseries",
                "links": {
                    "self": {
                        "href": "/t1"
                    },
                    "topic": {
						"href": "/topic/7",
						"id": "7"
                    }
                }
            }
        ]
    }
}
*/

// -=-=-

// Given this mongo collection document that contains no content
// (this is used for query tests)
/*
{
    "id": "8",
    "next": {
        "state": "published"
    },
    "current": {
        "state": "published"
    }
}
*/

// NOTE: the above has to be on one line ...
// NOTE: The following HAS to be on ONE line for unmarshal to work (and all the inner double quotes need escaping)
var mongoContentJSONResponse8 string = "{\"id\": \"8\", \"next\": { \"state\": \"published\" }, \"current\": { \"state\": \"published\" }}"

// then the Get Response in Public would look like (and note spotlight is sorted by href):
// (in Private mode, Next & Current contain the following)
/*
{
    "next": {
        "count": 0,
        "offset_index": 0,
        "limit": 0,
        "total_count": 0,
        "items": [
        ]
    },
    "current": {
        "count": 0,
        "offset_index": 0,
        "limit": 0,
        "total_count": 0,
        "items": [
        ]
    }
}
*/
