package api

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-areas-api/api/stubs"
	"net/http"
	"regexp"

	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-areas-api/utils"

	errs "github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

const (
	acceptLanguageHeaderMatchString = "en|cy"
)

//GetArea is a handler that gets a area by its ID from MongoDB
func (api *API) getArea(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	areaID := vars["id"]
	logdata := log.Data{"area-id": areaID}

	//get area from mongoDB by id
	area, err := api.areaStore.GetArea(ctx, areaID)
	if err != nil {
		log.Error(ctx, "getArea Handler: retrieving area from mongoDB returned an error", err, logdata)
		if err == errs.ErrAreaNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(area)
	if err != nil {
		log.Error(ctx, "getArea Handler: failed to marshal area resource into bytes", err, logdata)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Set headers
	setJSONContentType(w)

	if _, err := w.Write(b); err != nil {
		log.Error(ctx, "getArea Handler: error writing bytes to response", err, logdata)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info(ctx, "getArea Handler: Successfully retrieved area", logdata)

}

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

	//get area from stubbed data
	area := api.GeoData[areaID]
	area.AreaType = models.AcceptLanguageMapping[req.Header.Get(models.AcceptLanguageHeaderName)]

	areaData, err := json.Marshal(area)
	if err != nil {
		responseErr := models.NewError(ctx, err, models.MarshallingAreaDataError, err.Error())
		return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
	}

	return models.NewSuccessResponse(areaData, http.StatusOK, nil), nil

}

//GetAreas is a handler that gets all the areas from MongoDB
func (api *API) getAreas(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logData := log.Data{}

	offsetParameter := req.URL.Query().Get("offset")
	limitParameter := req.URL.Query().Get("limit")

	offset := api.defaultOffset
	limit := api.defaultLimit

	var err error

	if limitParameter != "" {
		logData["limit"] = limitParameter
		limit, err = utils.ValidatePositiveInt(limitParameter)
		if err != nil {
			log.Error(ctx, "invalid query parameter: limit", err, logData)
			err = errs.ErrInvalidQueryParameter
			handleAPIErr(ctx, err, w, nil)
			return
		}
	}

	if limit > api.maxLimit {
		logData["max_limit"] = api.maxLimit
		err = errs.ErrQueryParamLimitExceedMax
		log.Error(ctx, "limit is greater than the maximum allowed", err, logData)
		handleAPIErr(ctx, err, w, nil)
		return
	}

	if offsetParameter != "" {
		logData["offset"] = offsetParameter
		offset, err = utils.ValidatePositiveInt(offsetParameter)
		if err != nil {
			log.Error(ctx, "invalid query parameter: offset", err, logData)
			err = errs.ErrInvalidQueryParameter
			handleAPIErr(ctx, err, w, nil)
			return
		}
	}

	b, err := func() ([]byte, error) {

		logData := log.Data{}

		areasResult, err := api.areaStore.GetAreas(ctx, offset, limit)
		if err != nil {
			log.Error(ctx, "api endpoint getAreas returned an error", err)
			return nil, err
		}

		b, err := json.Marshal(areasResult)

		if err != nil {
			log.Error(ctx, "api endpoint getAreas failed to marshal resource into bytes", err, logData)
			return nil, err
		}

		return b, nil
	}()

	if err != nil {
		handleAPIErr(ctx, err, w, nil)
		return
	}

	setJSONContentType(w)
	if _, err = w.Write(b); err != nil {
		log.Error(ctx, "api endpoint getAreas error writing response body", err)
		handleAPIErr(ctx, err, w, nil)
		return
	}
	log.Info(ctx, "api endpoint getAreas request successful")
}

//getAreaRelationships is a handler that gets area relationship by ID - currently from stubbed data
func (api *API) getAreaRelationships(ctx context.Context, w http.ResponseWriter, req *http.Request) (*models.SuccessResponse, *models.ErrorResponse) {
	vars := mux.Vars(req)
	areaID := vars["id"]

	//get area relationships from stubbed data
	if relationShips, ok := stubs.Relationships[areaID]; !ok {
		return nil, models.NewErrorResponse(http.StatusNotFound, nil)
	} else {
		jsonResponse, err := json.Marshal(relationShips)
		if err != nil {
			responseErr := models.NewError(ctx, err, models.MarshallingAreaRelationshipsError, err.Error())
			return nil, models.NewErrorResponse(http.StatusInternalServerError, nil, responseErr)
		}

		return models.NewSuccessResponse(jsonResponse, http.StatusOK, nil), nil

	}
}
