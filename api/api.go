package api

import (
	"context"

	"github.com/gorilla/mux"
)

//API provides a struct to wrap the api around
type API struct {
	Router  *mux.Router
	mongoDB MongoServer
}

//Setup function sets up the api and returns an api
func Setup(ctx context.Context, r *mux.Router, mongoDB MongoServer) *API {
	api := &API{
		Router:  r,
		mongoDB: mongoDB,
	}

	//!!! see dp-image for possible best code ...

	r.HandleFunc("/hello", HelloHandler()).Methods("GET")
	return api
}
