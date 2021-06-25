package api

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-areas-api/utils"

	errs "github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
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
