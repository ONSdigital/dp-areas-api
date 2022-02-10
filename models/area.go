package models

// AreaType defines possible area types
type AreaType int

var (
	AcceptLanguageHeaderName = "Accept-Language"
	AcceptLanguageMapping    = map[string]string{
		"en": "English",
		"cy": "Cymraeg",
	}
)

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

// Struct for seeding area_type table with test data
type AreaTypeSeeding struct {
	AreaTypes map[string]map[string]interface{}
}

// Struct for seeding area table with test data
type AreaSeeding struct {
	Areas map[string]interface{}
}
