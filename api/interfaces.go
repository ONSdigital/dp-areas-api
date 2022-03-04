package api

import (
	"context"

	"github.com/ONSdigital/dp-areas-api/config"

	"github.com/ONSdigital/dp-areas-api/models"
)

//go:generate moq -out mock/rdsAreaStore.go -pkg mock . RDSAreaStore

// RDSAreaStore represents all the required methods from aurora DB
type RDSAreaStore interface {
	Init(ctx context.Context, cfg *config.Config) error
	Close()
	GetRelationships(areaCode, relationshipParameter string) ([]*models.AreaBasicData, error)
	ValidateArea(code string) error
	GetArea(ctx context.Context, areaId string) (*models.AreasDataResults, error)
	BuildTables(ctx context.Context, executionList []string) error
	Ping(ctx context.Context) error
	UpsertArea(ctx context.Context, area models.AreaParams) (bool, error)
	GetAncestors(areaID string) ([]models.AreasAncestors, error)
}
