package api

import (
	"net/http"

	"github.com/ONSdigital/dp-net/handlers"
	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// GetTopicHandler is a handler that gets a topic by its id from MongoDB
func (api *API) GetTopicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	hColID := ctx.Value(handlers.CollectionID.Context())
	logdata := log.Data{
		handlers.CollectionID.Header(): hColID,
		"request-id":                   ctx.Value(dprequest.RequestIdKey),
		"topic-id":                     id,
	}

	// get topic from mongoDB by id
	topic, err := api.mongoDB.GetTopic(req.Context(), id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	if err := WriteJSONBody(ctx, topic, w, logdata); err != nil {
		handleError(ctx, w, err, logdata)
		return
	}
	log.Event(ctx, "Successfully retrieved topic", log.INFO, logdata) //!!! is this a good log to have ? ... as in is there too much logging going on ?
}
