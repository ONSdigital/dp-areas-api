package api

import (
	"errors"

	"github.com/ONSdigital/dp-areas-api/api/stubs"
	"github.com/ONSdigital/dp-areas-api/models"
)

// Ancestry implements the AcestorStore interface
type Ancestry struct{}

// GetAncestors retrieves AreaAncestors by ID from stubbed data for now
func (a Ancestry) GetAncestors(areaID string) ([]models.AreasAncestors, error) {
	if _, ok := stubs.Ancestors[areaID]; !ok {
		return nil, errors.New(models.AncestryDataGetError)
	}
	return stubs.Ancestors[areaID], nil
}
