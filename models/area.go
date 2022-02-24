package models

import (
	"context"
	"time"
)

var areaTypeAndCode = map[string]string{
	"E92": "Country",
	"E06": "Unitary Authorities",
	"W92": "Country",
	"E12": "Region",
	"W06": "Unitary Authorities",
	"E47": "Combined Authorities",
	"E11": "Metropolitan Counties",
	"E10": "Counties",
	"E09": "London Boroughs",
	"E08": "Metropolitan Districts",
	"E07": "Non-metropolitan Districts",
	"E05": "Electoral Wards",
	"W05": "Electoral Wards",
}

// AreaType defines possible area types
type AreaType int

var (
	AcceptLanguageHeaderName = "Accept-Language"
	AcceptLanguageMapping    = map[string]string{
		"en": "English",
		"cy": "Cymraeg",
	}
)

// AreaName represents the structure of the area name details used for update request
type AreaName struct {
	Name       string     `json:"name"`
	ActiveFrom *time.Time `json:"active_from"`
	ActiveTo   *time.Time `json:"active_to"`
}

// AreaParams represents the structure of the area used for create or update
type AreaParams struct {
	Code          string     `json:"code"`
	AreaName      *AreaName  `json:"area_name"`
	GeometricData string     `json:"geometry"`
	ActiveFrom    *time.Time `json:"active_from"`
	ActiveTo      *time.Time `json:"active_to"`
	Visible       *bool      `json:"visible"`
	AreaType      string
}

func (a *AreaParams) ValidateAreaRequest(ctx context.Context) []error {
	var validationErrs []error
	if a.Code == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidAreaCodeError, InvalidAreaCodeErrorDescription))
	}

	if a.AreaType == "" {
		validationErrs = append(validationErrs, NewValidationError(ctx, InvalidAreaTypeError, InvalidAreaTypeErrorDescription))
	}

	if a.AreaName == nil {
		validationErrs = append(validationErrs, NewValidationError(ctx, AreaNameDetailsNotProvidedError, AreaNameDetailsNotProvidedErrorDescription))
	} else {

		if a.AreaName.Name == "" {
			validationErrs = append(validationErrs, NewValidationError(ctx, AreaNameNotProvidedError, AreaNameNotProvidedErrorDescription))
		}

		if a.AreaName.ActiveFrom == nil {
			validationErrs = append(validationErrs, NewValidationError(ctx, AreaNameActiveFromNotProvidedError, AreaNameActiveFromNotProvidedErrorDescription))
		}

		if a.AreaName.ActiveTo == nil {
			validationErrs = append(validationErrs, NewValidationError(ctx, AreaNameActiveToNotProvidedError, AreaNameActiveToNotProvidedErrorDescription))
		}
	}

	return validationErrs
}

func (a *AreaParams) SetAreaType(ctx context.Context) {
	if len(a.Code) > 3 {
		areaCodePrefix := a.Code[:3]
		if areaType, ok := areaTypeAndCode[areaCodePrefix]; ok {
			a.AreaType = areaType
		}
	}
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

// Struct for seeding area_type table with test data
type AreaTypeSeeding struct {
	AreaTypes map[string]map[string]interface{}
}

// Struct for seeding area table with test data
type AreaSeeding struct {
	Areas map[string]interface{}
}
