package api

import (
	"context"
	"net/http"

	dpauth "github.com/ONSdigital/dp-authorisation/auth"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-topic-api/models"
)

//go:generate moq -out mock/mongo.go -pkg mock . MongoServer
//go:generate moq -out mock/auth.go -pkg mock . AuthHandler

// MongoServer defines the required methods from MongoDB
type MongoServer interface {
	Close(ctx context.Context) error
	Checker(ctx context.Context, state *healthcheck.CheckState) (err error)
	GetTopic(ctx context.Context, id string) (topic *models.Topic, err error)
}

// AuthHandler interface for adding auth to endpoints
type AuthHandler interface {
	Require(required dpauth.Permissions, handler http.HandlerFunc) http.HandlerFunc
}
