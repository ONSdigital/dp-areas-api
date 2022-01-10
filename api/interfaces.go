package api

import (
	"context"

	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out mock/areaStore.go -pkg mock . AreaStore

// AreaStore represents all the required methods from mongo DB
type AreaStore interface {
	Close(ctx context.Context) error
	Checker(context.Context, *healthcheck.CheckState) error
	GetArea(ctx context.Context, id string) (*models.Area, error)
	GetVersion(ctx context.Context, id string, versionID int) (*models.Area, error)
	GetAreas(ctx context.Context, offset, limit int) (*models.AreasResults, error)
}
