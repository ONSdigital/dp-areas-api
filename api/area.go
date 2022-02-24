package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/ONSdigital/log.go/v2/log"

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

func (api *API) updateArea(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	defer func() {
		if err := req.Body.Close(); err != nil {
			_ = models.NewError(ctx, err, models.BodyCloseError, models.BodyClosedFailedDescription)
		}
	}()

	vars := mux.Vars(req)
	areaCode := vars["id"]
	logData := log.Data{"area code": areaCode}
	log.Info(ctx, "received request to upsert area", logData)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, models.NewBodyReadError(ctx, err)
	}

	area := models.AreaParams{}

	err = json.Unmarshal(body, &area)
	if err != nil {
		return nil, models.NewBodyUnmarshalError(ctx, err)
	}
	area.Code = areaCode
	area.SetAreaType(ctx)
	validationErrors := area.ValidateAreaRequest(ctx)

	if len(validationErrors) != 0 {
		return nil, models.NewErrorResponse(http.StatusNotFound, nil, validationErrors...)
	}

	isInserted, err := api.rdsAreaStore.UpsertArea(ctx, area)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.AreaDataIdUpsertError, err.Error())
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	if isInserted {
		return models.NewSuccessResponse(nil, http.StatusCreated, nil), nil
	} else {
		return models.NewSuccessResponse(nil, http.StatusOK, nil), nil
	}

}
