package api

import (
	"net/http"

	"github.com/ONSdigital/dp-net/handlers"
	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// getTopicPublicHandler is a handler that gets a topic by its id from MongoDB
func (api *API) getTopicPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	hColID := ctx.Value(handlers.CollectionID.Context())
	logdata := log.Data{
		handlers.CollectionID.Header(): hColID,
		"request_id":                   ctx.Value(dprequest.RequestIdKey),
		"topic_id":                     id,
		"function":                     "getTopicPublicHandler",
	}

	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	// Ensure the sub document has the main document ID
	topic.Current.ID = topic.ID

	// User is not authenticated and hence has only access to current sub document
	if err := WriteJSONBody(ctx, topic.Current, w, logdata); err != nil {
		return
	}
	log.Event(ctx, "request successful", log.INFO, logdata) // NOTE: name of function is in logdata
	// NOTE 1st log.Event() in CheckIdentity() needs removing, that looks like:
	// log.Event(ctx, "checking for an identity in request context", log.HTTP(r, 0, 0, nil, nil), logData)
}

// getTopicPrivateHandler is a handler that gets a topic by its id from MongoDB
func (api *API) getTopicPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	hColID := ctx.Value(handlers.CollectionID.Context())
	logdata := log.Data{
		handlers.CollectionID.Header(): hColID,
		"request_id":                   ctx.Value(dprequest.RequestIdKey),
		"topic_id":                     id,
		"function":                     "getTopicPrivateHandler",
	}

	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	// User has valid authentication to get raw topic document
	if err := WriteJSONBody(ctx, topic, w, logdata); err != nil {
		return
	}
	log.Event(ctx, "request successful", log.INFO, logdata) // NOTE: name of function is in logdata
}
