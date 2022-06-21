package sdk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
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
func NewWithHealthClient(hcCli *healthcheck.Client) (*Client, error) {
	if hcCli == nil {
		return nil, errors.New("health client is nil")
	}
	return &Client{
		hcCli: healthcheck.NewClientWithClienter(service, hcCli.URL, hcCli.Client),
	}, nil
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

// callTopicAPI calls the Topic API endpoint given by path for the provided REST method, request headers, and body payload.
// It returns the response body and any error that occurred.
func (cli *Client) callTopicAPI(ctx context.Context, path, method string, headers Headers, payload []byte) ([]byte, error) {
	URL, err := url.Parse(path)
	if err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to parse path: \"%v\" error is: %v", path, err),
			Code: http.StatusInternalServerError,
		}
	}

	path = URL.String()

	var req *http.Request

	if payload != nil {
		req, err = http.NewRequest(method, path, bytes.NewReader(payload))
		req.Header.Add("Content-type", "application/json")
	} else {
		req, err = http.NewRequest(method, path, http.NoBody)
	}

	// check req, above, didn't error
	if err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to create request for call to topic api, error is: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	err = headers.Add(req)
	if err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to add headers to request, error is: %v", err),
			Code: http.StatusInternalServerError,
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

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 400 {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed as unexpected code from topic api: %v", resp.StatusCode),
			Code: resp.StatusCode,
		}
	}

	if resp.Body == nil {
		return nil, nil
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, apiError.StatusError{
			Err:  fmt.Errorf("failed to read response body from call to topic api, error is: %v", err),
			Code: http.StatusInternalServerError,
		}
	}

	return b, nil
}

// closeResponseBody closes the response body and logs an error if unsuccessful
func closeResponseBody(resp *http.Response) error {
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
