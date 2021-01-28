package api

import (
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// getContentPublicHandler is a handler that gets content by its id from MongoDB for Web
func (api *API) getContentPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	//!!! adjust rest of code from here for content
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getContentPublicHandler",
	}

	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(id)
	if err != nil {
		// no topic found to retrieve the subtopics from
		handleError(ctx, w, err, logdata)
		return
	}

	// User is not authenticated and hence has only access to current sub document(s)
	var result models.PublicSubtopics

	if topic.Current == nil {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	if len(topic.Current.SubtopicIds) == 0 {
		// no subtopics exist for the requested ID
		handleError(ctx, w, apierrors.ErrNotFound, logdata)
		return
	}

	for _, subTopicID := range topic.Current.SubtopicIds {
		// get sub topic from mongoDB by subTopicID
		topic, err := api.dataStore.Backend.GetTopic(subTopicID)
		if err != nil {
			logdata["missing subtopic for id"] = subTopicID
			log.Event(ctx, err.Error(), log.ERROR, logdata)
			continue
		}
		result.PublicItems = append(result.PublicItems, topic.Current)
		result.TotalCount++
	}
	if result.TotalCount == 0 {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	if err := WriteJSONBody(ctx, result, w, logdata); err != nil {
		return
	}
	log.Event(ctx, "request successful", log.INFO, logdata) // NOTE: name of function is in logdata
}

// getContentPrivateHandler is a handler that gets content by its id from MongoDB for Publishing
func (api *API) getContentPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	//!!! adjust rest of code from here for content
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getContentPrivateHandler",
	}

	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(id)
	if err != nil {
		// no topic found to retrieve the subtopics from
		handleError(ctx, w, err, logdata)
		return
	}

	// User has valid authentication to get raw full topic document(s)
	var result models.PrivateSubtopics

	if topic.Next == nil {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	if len(topic.Next.SubtopicIds) == 0 {
		// no subtopics exist for the requested ID
		handleError(ctx, w, apierrors.ErrNotFound, logdata)
		return
	}

	for _, subTopicID := range topic.Next.SubtopicIds {
		// get topic from mongoDB by subTopicID
		topic, err := api.dataStore.Backend.GetTopic(subTopicID)
		if err != nil {
			logdata["missing subtopic for id"] = subTopicID
			log.Event(ctx, err.Error(), log.ERROR, logdata)
			continue
		}
		result.PrivateItems = append(result.PrivateItems, topic)
		result.TotalCount++
	}
	if result.TotalCount == 0 {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	if err := WriteJSONBody(ctx, result, w, logdata); err != nil {
		return
	}
	log.Event(ctx, "request successful", log.INFO, logdata) // NOTE: name of function is in logdata
}
