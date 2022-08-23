package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-areas-api/api/geodata"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/models"

	"github.com/gorilla/mux"
)

var (
	fls = []string{
		geodata.E92000001PropertyData,
		geodata.W92000004PropertyData,
		geodata.E34002743PropertyData,
		geodata.W37000454PropertyData,
	}
)

// API provides a struct to wrap the api around
type API struct {
	Router       *mux.Router
	GeoData      map[string]models.AreasDataResults
	rdsAreaStore RDSAreaStore
}

type baseHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) (*models.SuccessResponse, *models.ErrorResponse)

func contextAndErrors(h baseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		response, err := h(ctx, w, req)
		if err != nil {
			writeErrorResponse(ctx, w, err)
			return
		}
		writeSuccessResponse(ctx, w, response)
	}
}

// Setup function sets up the api and returns an api
func Setup(ctx context.Context, cfg *config.Config, r *mux.Router, rdsStore RDSAreaStore) (*API, error) {
	// initialised stubbed geo data
	geoData, err := initialiseStubbedAreaData(ctx)
	if err != nil {
		return nil, err
	}

	api := &API{
		Router:       r,
		GeoData:      geoData,
		rdsAreaStore: rdsStore,
	}

	r.HandleFunc("/v1/areas/{id}", contextAndErrors(api.getAreaData)).Methods(http.MethodGet)
	r.HandleFunc("/v1/areas/{id}/relations", contextAndErrors(api.getAreaRelationships)).Methods(http.MethodGet)

	if cfg.EnablePrivateEndpoints {
		r.HandleFunc("/v1/areas/{id}", contextAndErrors(api.updateArea)).Methods(http.MethodPut)
	}

	r.HandleFunc("/v1/boundaries/{id}", contextAndErrors(api.getBoundary)).Methods(http.MethodGet)

	return api, nil
}

func initialiseStubbedAreaData(_ context.Context) (map[string]models.AreasDataResults, error) {
	geoData := make(map[string]models.AreasDataResults, 2)
	for _, geoDataFile := range fls {
		var data models.AreasDataResults
		if err := json.Unmarshal([]byte(geoDataFile), &data); err != nil {
			return nil, err
		}
		geoData[data.Code] = data
	}
	return geoData, nil
}

func writeErrorResponse(ctx context.Context, w http.ResponseWriter, errorResponse *models.ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	// process custom headers
	if errorResponse.Headers != nil {
		for key := range errorResponse.Headers {
			w.Header().Set(key, errorResponse.Headers[key])
		}
	}
	w.WriteHeader(errorResponse.Status)

	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		responseErr := models.NewError(ctx, err, "JSONMarshalError", "failed to write http response")
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		responseErr := models.NewError(ctx, err, "WriteResponseError", "failed to write http response")
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}
}

func writeSuccessResponse(ctx context.Context, w http.ResponseWriter, successResponse *models.SuccessResponse) {
	w.Header().Set("Content-Type", "application/json")
	// process custom headers
	if successResponse.Headers != nil {
		for key := range successResponse.Headers {
			w.Header().Set(key, successResponse.Headers[key])
		}
	}
	w.WriteHeader(successResponse.Status)

	_, err := w.Write(successResponse.Body)
	if err != nil {
		responseErr := models.NewError(ctx, err, "WriteResponseError", "failed to write http response")
		http.Error(w, responseErr.Description, http.StatusInternalServerError)
		return
	}
}
