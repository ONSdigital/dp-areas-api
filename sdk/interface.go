package sdk

import (
	"context"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-topic-api/models"
	apiError "github.com/ONSdigital/dp-topic-api/sdk/errors"
)

//go:generate moq -out ./mocks/client.go -pkg mocks . Clienter

type Clienter interface {
	Checker(ctx context.Context, check *health.CheckState) error
	GetNavigationPublic(ctx context.Context, reqHeaders Headers, options Options) (*models.Navigation, apiError.Error)
	GetRootTopicsPrivate(ctx context.Context, reqHeaders Headers) (*models.PrivateSubtopics, apiError.Error)
	GetRootTopicsPublic(ctx context.Context, reqHeaders Headers) (*models.PublicSubtopics, apiError.Error)
	GetTopicPrivate(ctx context.Context, reqHeaders Headers, id string) (*models.TopicResponse, apiError.Error)
	GetTopicPublic(ctx context.Context, reqHeaders Headers, id string) (*models.Topic, apiError.Error)
	GetSubtopicsPrivate(ctx context.Context, reqHeaders Headers, id string) (*models.PrivateSubtopics, apiError.Error)
	GetSubtopicsPublic(ctx context.Context, reqHeaders Headers, id string) (*models.PublicSubtopics, apiError.Error)
	Health() *healthcheck.Client
	URL() string
}
