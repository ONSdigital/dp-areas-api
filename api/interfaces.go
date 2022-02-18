package api

import (
	"context"

	"github.com/ONSdigital/dp-areas-api/config"

	"github.com/ONSdigital/dp-areas-api/models"
)

//go:generate moq -out mock/rdsAreaStore.go -pkg mock . RDSAreaStore
//go:generate moq -out mock/ancestorStore.go -pkg mock . AncestorStore

// RDSAreaStore represents all the required methods from aurora DB
type RDSAreaStore interface {
	Init(ctx context.Context, cfg *config.Config) error
	Close()
	GetRelationships(areaCode string) ([]*models.AreaBasicData, error)
	ValidateArea(code string) error
	GetArea(areaId string) (*models.AreaDataRDS, error)
	BuildTables(ctx context.Context, executionList []string) error
	Ping(ctx context.Context) error
}

// Ancestors defines a method to lookup ancestors from an areaID
type AncestorStore interface {
	GetAncestors(areaID string) ([]models.AreasAncestors, error)
}
