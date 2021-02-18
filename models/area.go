package models

import (
	"time"
)

// AreaType defines possible area types
type AreaType int

// possible area types
const (
	Country AreaType = iota
	Region
	UnitaryAuthorities
	CombinedAuthorities
	MetropolitanCounties
	Counties
	LondonBoroughs
	MetropolitanDistricts
	NonMetropolitanDistricts
	ElectoralWards
	Invalid
)

//Area represents the structure for an area
type Area struct {
	ID                string        `bson:"id,omitempty"                  json:"id,omitempty"`
	Version           int           `bson:"version,omitempty"             json:"version,omitempty"`
	Name              string        `bson:"name,omitempty"                json:"name,omitempty"`
	ReleaseDate       string        `bson:"release_date,omitempty"        json:"release_date,omitempty"`
	LastUpdated       time.Time     `bson:"last_updated,omitempty"        json:"last_updated,omitempty"`
	Type              string        `bson:"type,omitempty"                json:"type,omitempty"`
	ParentAreas       []LinkedAreas `bson:"parent_areas,omitempty"        json:"parent_areas,omitempty"`
	ChildAreas        []LinkedAreas `bson:"child_areas,omitempty"         json:"child_areas,omitempty"`
	NeighbouringAreas []LinkedAreas `bson:"neighbouring_areas,omitempty"  json:"neighbouring_areas,omitempty"`
}

type LinkedAreas struct {
	ID   string `bson:"id,omitempty"         json:"id,omitempty"`
	Type string `bson:"type,omitempty"       json:"type,omitempty"`
	Name string `bson:"name,omitempty"       json:"name,omitempty"`
}
