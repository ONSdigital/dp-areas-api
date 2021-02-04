package api

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

//GetAreaHandler is a handler that gets a area by its ID from MongoDB
func (api *API) getArea(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	areaID := vars["id"]
	logdata := log.Data{"area-id": areaID}

	//get area from mongoDB by id
	area, err := api.areaStore.GetArea(ctx, areaID)
	if err != nil {
		log.Event(ctx, "getArea Handler: retrieving area from mongoDB returned an error", log.ERROR, log.Error(err), logdata)
		if err == apierrors.ErrAreaNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var b []byte
	b, err = json.Marshal(area)
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
