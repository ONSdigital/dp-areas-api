package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-topic-api/models"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testPrivateTopics = models.PrivateSubtopics{
		Count:        2,
		Offset:       0,
		Limit:        100,
		TotalCount:   2,
		PrivateItems: &[]models.TopicResponse{testPrivateTopic1, testPrivateTopic2},
	}

	testPrivateTopic1 = models.TopicResponse{
		ID:      "1234",
		Current: &testPublicTopic1,
		Next:    &testPublicTopic1,
	}

	testPrivateTopic2 = models.TopicResponse{
		ID:      "5678",
		Current: &testPublicTopic2,
		Next:    &testPublicTopic2,
	}

	topicRelease = models.TopicRelease{
		ReleaseDate: "2022-11-11T09:30:00Z",
	}
)

func TestGetRootTopicsPrivate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	Convey("Given the private root topics is returned successfully", t, func() {
		body, err := json.Marshal(testPrivateTopics)
		if err != nil {
			t.Errorf("failed to setup test data, error: %v", err)
		}

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			},
			nil)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetRootTopicsPrivate is called", func() {
			respRootTopics, err := topicAPIClient.GetRootTopicsPrivate(ctx, Headers{
				ServiceAuthToken: "valid-service-token",
			})

			Convey("Then the expected private root topics is returned", func() {
				So(*respRootTopics, ShouldResemble, testPrivateTopics)

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics")
					})
				})
			})
		})
	})

	Convey("Given a 500 response from topic api", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetRootTopicsPrivate is called", func() {
			respRootTopics, err := topicAPIClient.GetRootTopicsPrivate(ctx, Headers{})

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)

				Convey("And the expected private root topics should be nil", func() {
					So(respRootTopics, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics")
					})
				})
			})
		})
	})

	Convey("Given the client returns an unexpected error", t, func() {
		clientError := errors.New("unexpected error")
		httpClient := newMockHTTPClient(&http.Response{}, clientError)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetRootTopicsPrivate is called", func() {
			respRootTopics, err := topicAPIClient.GetRootTopicsPrivate(ctx, Headers{})

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)

				Convey("And the expected private root topics should be nil", func() {
					So(respRootTopics, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics")
					})
				})
			})
		})
	})
}

func TestGetTopicPrivate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	Convey("Given the private topic is returned successfully", t, func() {
		body, err := json.Marshal(testPrivateTopic1)
		if err != nil {
			t.Errorf("failed to setup test data, error: %v", err)
		}

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			},
			nil)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetTopicPrivate is called", func() {
			respTopic, err := topicAPIClient.GetTopicPrivate(ctx, Headers{
				ServiceAuthToken: "valid-service-token",
			}, "1234")

			Convey("Then the expected private topic is returned", func() {
				So(*respTopic, ShouldResemble, testPrivateTopic1)

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1234")
					})
				})
			})
		})
	})

	Convey("Given a 500 response from topic api", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetTopicPrivate is called", func() {
			respTopic, err := topicAPIClient.GetTopicPrivate(ctx, Headers{}, "1234")

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)

				Convey("And the expected private topic should be nil", func() {
					So(respTopic, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1234")
					})
				})
			})
		})
	})

	Convey("Given the client returns an unexpected error", t, func() {
		clientError := errors.New("unexpected error")
		httpClient := newMockHTTPClient(&http.Response{}, clientError)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetTopicPrivate is called", func() {
			respTopic, err := topicAPIClient.GetTopicPrivate(ctx, Headers{}, "1234")

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)

				Convey("And the expected private topic should be nil", func() {
					So(respTopic, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1234")
					})
				})
			})
		})
	})
}

func TestGetSubtopicsPrivate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	Convey("Given private subtopics is returned successfully", t, func() {
		body, err := json.Marshal(testPrivateTopics)
		if err != nil {
			t.Errorf("failed to setup test data, error: %v", err)
		}

		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(body)),
			},
			nil)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetSubtopicsPrivate is called", func() {
			respSubtopics, err := topicAPIClient.GetSubtopicsPrivate(ctx, Headers{
				ServiceAuthToken: "valid-service-token",
			}, "1357")

			Convey("Then the expected private subtopics is returned", func() {
				So(*respSubtopics, ShouldResemble, testPrivateTopics)

				Convey("And no error is returned", func() {
					So(err, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1357/subtopics")
					})
				})
			})
		})
	})

	Convey("Given a 500 response from topic api", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetSubtopicsPrivate is called", func() {
			respSubtopics, err := topicAPIClient.GetSubtopicsPrivate(ctx, Headers{}, "1357")

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)

				Convey("And the expected private subtopics should be nil", func() {
					So(respSubtopics, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1357/subtopics")
					})
				})
			})
		})
	})

	Convey("Given the client returns an unexpected error", t, func() {
		clientError := errors.New("unexpected error")
		httpClient := newMockHTTPClient(&http.Response{}, clientError)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When GetSubtopicsPrivate is called", func() {
			respSubtopics, err := topicAPIClient.GetSubtopicsPrivate(ctx, Headers{}, "1357")

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)

				Convey("And the expected private subtopics should be nil", func() {
					So(respSubtopics, ShouldBeNil)

					Convey("And client.Do should be called once with the expected parameters", func() {
						doCalls := httpClient.DoCalls()
						So(doCalls, ShouldHaveLength, 1)
						So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1357/subtopics")
					})
				})
			})
		})
	})
}

func TestPutTopicReleasePrivate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	body, err := json.Marshal(topicRelease)
	if err != nil {
		t.Errorf("failed to setup test data, error: %v", err)
	}

	Convey("Given private put topic release is successful", t, func() {
		httpClient := newMockHTTPClient(
			&http.Response{
				StatusCode: http.StatusOK,
				Body:       nil,
				Header:     nil,
			},
			nil)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When PutTopicReleasePrivate is called", func() {
			respInfo, err := topicAPIClient.PutTopicReleasePrivate(ctx, Headers{
				ServiceAuthToken: "valid-service-token",
			}, "1357", body)

			Convey("Then it succeeds with no errors returned", func() {
				So(err, ShouldBeNil)
				So(respInfo, ShouldNotBeNil)
				So(respInfo.Status, ShouldEqual, http.StatusOK)
				So(respInfo.Body, ShouldBeNil)
				So(respInfo.Headers, ShouldBeNil)
			})

			Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				So(doCalls, ShouldHaveLength, 1)
				So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1357/release-date")
			})
		})
	})

	Convey("Given a 500 response from topic api", t, func() {
		httpClient := newMockHTTPClient(&http.Response{StatusCode: http.StatusInternalServerError}, nil)
		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When PutTopicReleasePrivate is called", func() {
			respInfo, err := topicAPIClient.PutTopicReleasePrivate(ctx, Headers{}, "1357", body)

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)
				So(respInfo, ShouldNotBeNil)
				So(respInfo.Body, ShouldBeNil)
				So(respInfo.Headers, ShouldBeNil)
				So(respInfo.Status, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				So(doCalls, ShouldHaveLength, 1)
				So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1357/release-date")
			})
		})
	})

	Convey("Given the client returns an unexpected error", t, func() {
		clientError := errors.New("unexpected error")
		httpClient := newMockHTTPClient(&http.Response{}, clientError)

		topicAPIClient := newTopicAPIClient(t, httpClient)

		Convey("When PutTopicReleasePrivate is called", func() {
			respInfo, err := topicAPIClient.PutTopicReleasePrivate(ctx, Headers{}, "1357", body)

			Convey("Then an error should be returned ", func() {
				So(err, ShouldNotBeNil)
				So(err.Status(), ShouldEqual, http.StatusInternalServerError)
				So(respInfo, ShouldBeNil)
			})

			Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClient.DoCalls()
				So(doCalls, ShouldHaveLength, 1)
				So(doCalls[0].Req.URL.Path, ShouldEqual, "/topics/1357/release-date")

			})
		})
	})
}
