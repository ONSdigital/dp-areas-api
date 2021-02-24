package api

import (
	"context"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/gorilla/mux"
	"net/http"
)

//API provides a struct to wrap the api around
type API struct {
	Router        *mux.Router
	areaStore     AreaStore
	defaultLimit  int
	defaultOffset int
	maxLimit      int
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, cfg *config.Config, r *mux.Router, areaStore AreaStore) *API {
	api := &API{
		Router:        r,
		areaStore:     areaStore,
		defaultLimit:  cfg.DefaultLimit,
		defaultOffset: cfg.DefaultOffset,
		maxLimit:      cfg.DefaultMaxLimit,
	}
	r.HandleFunc("/areas", api.getAreas).Methods(http.MethodGet)
	r.HandleFunc("/areas/{id}", api.getArea).Methods(http.MethodGet)
	r.HandleFunc("/areas/{id}/versions/{version}", api.getVersion).Methods(http.MethodGet)
	return api
}
