package api

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

//API provides a struct to wrap the api around
type API struct {
	Router    *mux.Router
	areaStore AreaStore
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, areaStore AreaStore) *API {
	api := &API{
		Router:    r,
		areaStore: areaStore,
	}

	r.HandleFunc("/areas/{id}", api.getArea).Methods(http.MethodGet)
	r.HandleFunc("/areas/{id}/versions/{version}", api.getVersion).Methods(http.MethodGet)
	return api
}
