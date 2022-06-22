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
	testPrivateRootTopics = models.PrivateSubtopics{
		Count:        2,
		Offset:       0,
		Limit:        100,
		TotalCount:   2,
		PrivateItems: &[]models.TopicResponse{testPrivateRootTopic1, testPrivateRootTopic2},
	}

	testPrivateRootTopic1 = models.TopicResponse{
		ID:      "1234",
		Current: &testPublicRootTopic1,
		Next:    &testPublicRootTopic1,
	}

	testPrivateRootTopic2 = models.TopicResponse{
		ID:      "5678",
		Current: &testPublicRootTopic2,
		Next:    &testPublicRootTopic2,
	}
)

func TestGetRootTopicsPrivate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	Convey("Given the private root topics is returned successfully", t, func() {
		body, err := json.Marshal(testPrivateRootTopics)
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
				So(*respRootTopics, ShouldResemble, testPrivateRootTopics)

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
