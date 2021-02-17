package api

import (
	"encoding/json"
	"github.com/ONSdigital/dp-areas-api/utils"
	"net/http"

	errs "github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/log.go/log"
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
		log.Event(ctx, "getArea Handler: retrieving area from mongoDB returned an error", log.ERROR, log.Error(err), logdata)
		if err == errs.ErrAreaNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(area)
	if err != nil {
		log.Event(ctx, "getArea Handler: failed to marshal area resource into bytes", log.ERROR, log.Error(err), logdata)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Set headers
	setJSONContentType(w)

	if _, err := w.Write(b); err != nil {
		log.Event(ctx, "getArea Handler: error writing bytes to response", log.ERROR, log.Error(err), logdata)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Event(ctx, "getArea Handler: Successfully retrieved area", log.INFO, logdata)

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
			log.Event(ctx, "invalid query parameter: limit", log.ERROR, log.Error(err), logData)
			err = errs.ErrInvalidQueryParameter
			handleAPIErr(ctx, err, w, nil)
			return
		}
	}

	if limit > api.maxLimit {
		logData["max_limit"] = api.maxLimit
		err = errs.ErrQueryParamLimitExceedMax
		log.Event(ctx, "limit is greater than the maximum allowed", log.ERROR, logData)
		handleAPIErr(ctx, err, w, nil)
		return
	}

	if offsetParameter != "" {
		logData["offset"] = offsetParameter
		offset, err = utils.ValidatePositiveInt(offsetParameter)
		if err != nil {
			log.Event(ctx, "invalid query parameter: offset", log.ERROR, log.Error(err), logData)
			err = errs.ErrInvalidQueryParameter
			handleAPIErr(ctx, err, w, nil)
			return
		}
	}

	b, err := func() ([]byte, error) {

		logData := log.Data{}

		areasResult, err := api.areaStore.GetAreas(ctx, offset, limit)
		if err != nil {
			log.Event(ctx, "api endpoint getAreas returned an error", log.ERROR, log.Error(err))
			return nil, err
		}

		b, err := json.Marshal(areasResult)

		if err != nil {
			log.Event(ctx, "api endpoint getAreas failed to marshal resource into bytes", log.ERROR, log.Error(err), logData)
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
		log.Event(ctx, "api endpoint getAreas error writing response body", log.ERROR, log.Error(err))
		handleAPIErr(ctx, err, w, nil)
		return
	}
	log.Event(ctx, "api endpoint getAreas request successful", log.INFO)
}
