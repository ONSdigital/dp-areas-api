package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/v2/log"
)

// GetRootTopicsPrivate gets the private list of top level root topics for Publishing which returns both Next and Current document(s) in the response
func (cli *Client) GetRootTopicsPrivate(ctx context.Context, reqHeaders Headers) (*models.PrivateSubtopics, error) {
	path := fmt.Sprintf("%s/topics", cli.hcCli.URL)

	b, err := cli.callTopicAPI(ctx, path, http.MethodGet, reqHeaders, nil)
	if err != nil {
		logData := log.Data{
			"path":        path,
			"method":      http.MethodGet,
			"req_headers": reqHeaders,
			"body":        nil,
		}
		log.Error(ctx, "failed to call topic api", err, logData)
		return nil, err
	}

	var rootTopics models.PrivateSubtopics

	if err = json.Unmarshal(b, &rootTopics); err != nil {
		logData := log.Data{
			"response_bytes": b,
		}
		log.Error(ctx, "failed to unmarshal bytes into root topics", err, logData)
		return nil, err
	}

	return &rootTopics, nil
}
