package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	errs "github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/gorilla/mux"
)

const (
	acceptLanguageHeaderMatchString = "en|cy"
)

var (
	queryStr = "select id, code, active from areas_basic where id=$1"
)

//getBoundaryAreaData is a handler that gets a boundary data by ID - currently stubbed
func (api *API) getAreaData(ctx context.Context, _ http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	areaID := vars["id"]

	// error if accept language header not found
	var validationErrs []error
	if req.Header.Get(models.AcceptLanguageHeaderName) == "" {
		validationErrs = append(validationErrs, models.NewValidationError(ctx, models.AcceptLanguageHeaderError, models.AcceptLanguageHeaderNotFoundDescription))
	} else if m, _ := regexp.MatchString(acceptLanguageHeaderMatchString, req.Header.Get(models.AcceptLanguageHeaderName)); !m {
		validationErrs = append(validationErrs, models.NewValidationError(ctx, models.AcceptLanguageHeaderError, models.AcceptLanguageHeaderInvalidDescription))
	}
	if api.GeoData[areaID].Code == "" {
		validationErrs = append(validationErrs, models.NewValidationError(ctx, models.AreaDataIdGetError, models.AreaDataGetErrorDescription))
	}
	//handle errors
	if len(validationErrs) != 0 {
		return nil, models.NewErrorResponse(http.StatusNotFound, nil, validationErrs...)
	}

	// get ancestry data
	ancestryData, err := api.ancestorStore.GetAncestors(areaID)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.AncestryDataGetError, err.Error())
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	//get area from stubbed data
	area := api.GeoData[areaID]

	// update area data with ancestry data
	area.Ancestors = ancestryData

	area.AreaType = models.AcceptLanguageMapping[req.Header.Get(models.AcceptLanguageHeaderName)]

	areaData, err := json.Marshal(area)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.MarshallingAreaDataError, err.Error())
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(areaData, http.StatusOK, nil), nil
}

//getAreaRDSData test endpoint to demo rds database interaction
//TODO: remove this handler once rds transaction endpoints get added to the service - this is just an example
//Note: See TestGetAreaDataRFromRDS(t *testing.T) for mocking example
func (api *API) getAreaRDSData(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	areaID := vars["id"]
	var validationErrs []error

	area, err := api.rdsAreaStore.GetArea(areaID)
	if err != nil {
		validationErrs = append(validationErrs, models.NewValidationError(ctx, models.AreaDataIdGetError, err.Error()))
		return nil, models.NewErrorResponse(http.StatusNotFound, nil, validationErrs...)
	}

	areaData, err := json.Marshal(area)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.MarshallingAreaDataError, err.Error())
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(areaData, http.StatusOK, nil), nil
}

//getAreaRelationships is a handler that gets area relationship by ID - currently from stubbed data
func (api *API) getAreaRelationships(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	areaID := vars["id"]

	err := api.rdsAreaStore.ValidateArea(areaID)

	if err != nil {
		if err.Error() == errs.ErrNoRows.Error() {
			responseErr := models.NewError(ctx, err, models.InvalidAreaCodeError, err.Error())
			return nil, models.NewErrorResponse(http.StatusNotFound, nil, responseErr)
		}
		responseErr := models.NewError(ctx, err, models.AreaDataIdGetError, err.Error())
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	relatedAreaDetails, err := api.rdsAreaStore.GetRelationships(areaID)
	if err != nil {
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, err)
	}

	relationShips := make([]*models.AreaRelationShips, 0)
	for _, area := range relatedAreaDetails {
		relationShips = append(relationShips, &models.AreaRelationShips{
			AreaCode: area.Code,
			AreaName: area.Name,
			Href:     fmt.Sprintf("/v1/area/%s", area.Code),
		})
	}

	jsonResponse, err := json.Marshal(relationShips)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.MarshallingAreaRelationshipsError, err.Error())
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil

}
