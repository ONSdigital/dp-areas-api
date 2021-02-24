package api

import (
	"context"
	"encoding/json"

	"github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var (
	// errors that map to a HTTP 404 response
	notFound = map[error]bool{
		apierrors.ErrAreaNotFound:    true,
		apierrors.ErrVersionNotFound: true,
	}

	// errors that should return a 400 status
	badRequest = map[error]bool{
		apierrors.ErrInvalidQueryParameter:    true,
		apierrors.ErrQueryParamLimitExceedMax: true,
	}
)

func (api *API) getVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	areaID := vars["id"]
	version := vars["version"]
	logData := log.Data{"area_id": areaID, "version": version}

	//gets version for an area from mongoDb
	b, getVersionErr := func() ([]byte, error) {
		if err := api.areaStore.CheckAreaExists(areaID); err != nil {
			log.Event(ctx, "failed to find area", log.ERROR, log.Error(err), logData)
			return nil, err
		}
		areaVersion, err := strconv.Atoi(version)
		if err != nil {
			log.Event(ctx, "failed to convert version id to areas.version int", log.ERROR, log.Error(err), logData)
			return nil, err
		}
		results, err := api.areaStore.GetVersion(areaID, areaVersion)
		if err != nil {
			log.Event(ctx, "failed to find version for areas", log.ERROR, log.Error(err), logData)
			return nil, err
		}
		b, err := json.Marshal(results)
		if err != nil {
			log.Event(ctx, "failed to marshal version resource into bytes", log.ERROR, log.Error(err), logData)
			return nil, err
		}
		return b, nil
	}()

	if getVersionErr != nil {
		handleAPIErr(ctx, getVersionErr, w, logData)
		return
	}

	setJSONContentType(w)
	_, err := w.Write(b)
	if err != nil {
		log.Event(ctx, "failed writing bytes to response", log.ERROR, log.Error(err), logData)
		handleAPIErr(ctx, err, w, logData)
	}
	log.Event(ctx, "getVersion endpoint: request successful", log.INFO, logData)
}

func handleAPIErr(ctx context.Context, err error, w http.ResponseWriter, data log.Data) {
	var status int
	switch {

	case badRequest[err]:
		status = http.StatusBadRequest
	case notFound[err]:
		status = http.StatusNotFound
	default:
		err = apierrors.ErrInternalServer
		status = http.StatusInternalServerError
	}

	if data == nil {
		data = log.Data{}
	}

	log.Event(ctx, "request unsuccessful", log.ERROR, log.Error(err), data)
	http.Error(w, err.Error(), status)
}

func setJSONContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
