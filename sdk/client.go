package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-topic-api/models"
	apiError "github.com/ONSdigital/dp-topic-api/sdk/errors"
)

const (
	service = "dp-topic-api"
)

type Client struct {
	hcCli *healthcheck.Client
}

// New creates a new instance of Client with a given topic api url
func New(topicAPIURL string) *Client {
	return &Client{
		hcCli: healthcheck.NewClient(service, topicAPIURL),
	}
}

// NewWithHealthClient creates a new instance of topic API Client,
// reusing the URL and Clienter from the provided healthcheck client
func NewWithHealthClient(hcCli *healthcheck.Client) *Client {
	return &Client{
		hcCli: healthcheck.NewClientWithClienter(service, hcCli.URL, hcCli.Client),
	}
}

// URL returns the URL used by this client
func (cli *Client) URL() string {
	return cli.hcCli.URL
}

// Health returns the underlying Healthcheck Client for this topic API client
func (cli *Client) Health() *healthcheck.Client {
	return cli.hcCli
}

// Checker calls topic api health endpoint and returns a check object to the caller
func (cli *Client) Checker(ctx context.Context, check *health.CheckState) error {
	return cli.hcCli.Checker(ctx, check)
}

// GetRootTopicsPublic gets the public list of top level root topics for Web which returns the Current document(s) in the response
func (cli *Client) GetRootTopicsPublic(ctx context.Context, reqHeaders Headers) (*models.PublicSubtopics, apiError.Error) {
	path := fmt.Sprintf("%s/topics", cli.hcCli.URL)

	respInfo, apiErr := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	var rootTopics models.PublicSubtopics

	if err := json.Unmarshal(respInfo.Body, &rootTopics); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal rootTopics - error is: %v", err),
		}
	}

	return &rootTopics, nil
}

// GetTopicPublic gets the publicly available topics
func (cli *Client) GetTopicPublic(ctx context.Context, reqHeaders Headers, id string) (*models.Topic, apiError.Error) {
	path := fmt.Sprintf("%s/topics/%s", cli.hcCli.URL, id)

	respInfo, apiErr := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	var topic models.Topic

	if err := json.Unmarshal(respInfo.Body, &topic); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal topic - error is: %v", err),
		}
	}

	return &topic, nil
}

// GetSubtopicsPublic gets the public list of subtopics of a topic for Web which returns the Current document(s) in the response
func (cli *Client) GetSubtopicsPublic(ctx context.Context, reqHeaders Headers, id string) (*models.PublicSubtopics, apiError.Error) {
	path := fmt.Sprintf("%s/topics/%s/subtopics", cli.hcCli.URL, id)

	respInfo, apiErr := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	var subtopics models.PublicSubtopics

	if err := json.Unmarshal(respInfo.Body, &subtopics); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal subtopics - error is: %v", err),
		}
	}

	return &subtopics, nil
}

// GetNavigationPublic gets the public list of navigation items
func (cli *Client) GetNavigationPublic(ctx context.Context, reqHeaders Headers, options Options) (*models.Navigation, apiError.Error) {
	lang, err := options.Lang.String()
	if err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to get language - error is: %v", err),
		}
	}

	path := fmt.Sprintf("%s/navigation?lang=%s", cli.hcCli.URL, lang)

	respInfo, apiErr := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if apiErr != nil {
		return nil, apiErr
	}

	var navigation models.Navigation

	if err = json.Unmarshal(respInfo.Body, &navigation); err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to unmarshal navigation - error is: %v", err),
		}
	}

	return &navigation, nil
}

type ResponseInfo struct {
	Body    []byte
	Headers http.Header
	Status  int
}

// callTopicAPI calls the Topic API endpoint given by path for the provided REST method, request headers, and body payload.
// It returns the response body and any error that occurred.
func (cli *Client) callTopicAPI(ctx context.Context, path, method string, headers Headers, payload []byte) (*ResponseInfo, apiError.Error) {
	URL, err := url.Parse(path)
	if err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to parse path: \"%v\" error is: %v", path, err),
		}
	}

	path = URL.String()

	var req *http.Request

	if payload != nil {
		req, err = http.NewRequest(method, path, bytes.NewReader(payload))
	} else {
		req, err = http.NewRequest(method, path, http.NoBody)
	}

	// check req, above, didn't error
	if err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to create request for call to topic api, error is: %v", err),
		}
	}

	if payload != nil {
		req.Header.Add("Content-type", "application/json")
	}

	err = headers.Add(req)
	if err != nil {
		return nil, apiError.StatusError{
			Err: fmt.Errorf("failed to add headers to request, error is: %v", err),
		}
	}

	resp, err := cli.hcCli.Client.Do(ctx, req)
	if err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to call topic api, error is: %v", err),
			Code: http.StatusInternalServerError,
		}
	}
	defer func() {
		err = closeResponseBody(resp)
	}()

	respInfo := &ResponseInfo{
		Headers: resp.Header.Clone(),
		Status:  resp.StatusCode,
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 400 {
		return respInfo, apiError.StatusError{
			Err:  fmt.Errorf("failed as unexpected code from topic api: %v", resp.StatusCode),
			Code: resp.StatusCode,
		}
	}

	if resp.Body == nil {
		return respInfo, nil
	}

	respInfo.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		return respInfo, apiError.StatusError{
			Err:  fmt.Errorf("failed to read response body from call to topic api, error is: %v", err),
			Code: resp.StatusCode,
		}
	}

	fmt.Printf("respInfo is: %v", respInfo)

	return respInfo, nil
}

// closeResponseBody closes the response body and logs an error if unsuccessful
func closeResponseBody(resp *http.Response) apiError.Error {
	if resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			return apiError.StatusError{
				Err:  fmt.Errorf("error closing http response body from call to topic api, error is: %v", err),
				Code: http.StatusInternalServerError,
			}
		}
	}

	return nil
}
