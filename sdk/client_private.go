package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-topic-api/models"
	apiError "github.com/ONSdigital/dp-topic-api/sdk/errors"
)

// GetRootTopicsPrivate gets the private list of top level root topics for Publishing which returns both Next and Current document(s) in the response
func (cli *Client) GetRootTopicsPrivate(ctx context.Context, reqHeaders Headers) (*models.PrivateSubtopics, apiError.Error) {
	path := fmt.Sprintf("%s/topics", cli.hcCli.URL)

	respInfo, apiErr := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	var rootTopics models.PrivateSubtopics

	if err := json.Unmarshal(respInfo.Body, &rootTopics); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal rootTopics - error is: %v", err),
		}
	}

	return &rootTopics, nil
}

// GetTopicPrivate gets the full topic resource including the Next and Current nested objects
func (cli *Client) GetTopicPrivate(ctx context.Context, reqHeaders Headers, id string) (*models.TopicResponse, apiError.Error) {
	path := fmt.Sprintf("%s/topics/%s", cli.hcCli.URL, id)

	respInfo, apiErr := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	var topic models.TopicResponse

	if err := json.Unmarshal(respInfo.Body, &topic); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal topic - error is: %v", err),
		}
	}

	return &topic, nil
}

// GetSubtopicsPrivate gets the private list of subtopics of a topic for Publishing which returns both Next and Current document(s) in the response
func (cli *Client) GetSubtopicsPrivate(ctx context.Context, reqHeaders Headers, id string) (*models.PrivateSubtopics, apiError.Error) {
	path := fmt.Sprintf("%s/topics/%s/subtopics", cli.hcCli.URL, id)

	respInfo, apiErr := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	var subtopics models.PrivateSubtopics

	if err := json.Unmarshal(respInfo.Body, &subtopics); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal subtopics - error is: %v", err),
		}
	}

	return &subtopics, nil
}

type Result struct {
}

// PutTopicReleasePrivate inserts the release date into the topic next object ready for publishing
func (cli *Client) PutTopicReleasePrivate(ctx context.Context, reqHeaders Headers, id string, payload []byte) (*ResponseInfo, apiError.Error) {
	path := fmt.Sprintf("%s/topics/%s/release-date", cli.hcCli.URL, id)

	respInfo, apiErr := cli.callTopicAPI(ctx, path, http.MethodPut, reqHeaders, payload)
	if apiErr != nil {
		return respInfo, apiErr
	}

	fmt.Printf("got here last with: %v", respInfo)

	return respInfo, nil
}
