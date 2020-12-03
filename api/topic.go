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
		"function":                     "getTopicHandler",
	}

	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(req.Context(), id)
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
	log.Event(ctx, "request successful", log.INFO, logdata) //!!! is this a good log to have ? ... as in is there too much logging going on ?
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

	//logData := log.Data{}
	//authenticated := api.authenticate(req, logData) //!!! is this needed ? ... i suspect not as to get here we have already come through: api.isAuthenticated( api.isAuthorised(

	//fmt.Printf("authenticated is %v\n", authenticated)
	//fmt.Printf("logData is %+v\n", logData)

	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(req.Context(), id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	// User has valid authentication to get raw topic document
	if err := WriteJSONBody(ctx, topic, w, logdata); err != nil {
		return
	}
	log.Event(ctx, "request successful", log.INFO, logdata) //!!! is this a good log to have ? ... as in is there too much logging going on ?
}
