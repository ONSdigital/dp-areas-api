package models

import (
	"time"
)

// AreaType defines possible area types
type AreaType int

var (
	AcceptLanguageHeaderName = "Accept-Language"
	AcceptLanguageMapping    = map[string]string{
		"en": "English",
		"cy": "Cymraeg",
	}
)

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

// AreasResults represents a structure for a list of areas
type AreasResults struct {
	Items      *[]Area `json:"items"`
	Count      int     `json:"count"`
	Offset     int     `json:"offset"`
	Limit      int     `json:"limit"`
	TotalCount int     `json:"total_count"`
}

// AreasDataResults represents the structure for an area in api v1.
type AreasDataResults struct {
	Code          string                 `json:"code"`
	Name          string                 `json:"name"`
	ValidFrom     string                 `json:"date_start"`
	ValidTo       string                 `json:"date_end"`
	WelshName     string                 `json:"name_welsh"`
	GeometricData map[string]interface{} `json:"geometry"`
	Visible       bool                   `json:"visible"`
	AreaType      string                 `json:"area_type"`
	Ancestors     []AreasAncestors       `json:"ancestors"`
}

// AreasAncestors represents the Ancestry structure.
type AreasAncestors struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// AreaRelationShips represents the related areas with self ref
type AreaRelationShips struct {
	AreaCode string `json:"area_code"`
	AreaName string `json:"area_name"`
	Href     string `json:"href"`
}

// AreaRelationShips represents the related areas with self ref
type AreaBasicData struct {
	Code string `json:"code"`
	Name string `json:"name"`
}


// basic area data
type AreaDataRDS struct {
	Id     int64  `json:"id"`
	Code   string `json:"code"`
	Active bool   `json:"active"`
}
