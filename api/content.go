package api

import (
	"fmt"
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

func addPublicItem(contentList *models.PublicContent, typeName string, itemLink *[]models.TypeLinkObject, id string, state string) (count int) {
	count = 0

	if itemLink == nil {
		return
	}

	for _, item := range *itemLink {
		var topicLink models.LinkObject
		var selfLink models.LinkObject

		selfLink.HRef = item.HRef
		topicLink.ID = id
		topicLink.HRef = "/topic/" + id

		var cLinks models.ContentLinks
		cLinks.Self = &selfLink
		cLinks.Topic = &topicLink

		var cItem models.ContentItem
		cItem.Title = item.Title
		cItem.Type = typeName
		//		*cItem.State = state  // true for published
		cItem.Links = &cLinks

		if contentList.PublicItems == nil {
			contentList.PublicItems = &[]models.ContentItem{cItem}
		} else {
			*contentList.PublicItems = append(*contentList.PublicItems, cItem)
		}

		count++
		fmt.Printf("item: %+v\n", item)
	}

	return count
}

// getContentPublicHandler is a handler that gets content by its id from MongoDB for Web
func (api *API) getContentPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	//!!! adjust rest of code from here for content
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"content_id": id,
		"function":   "getContentPublicHandler",
	}

	// get content from mongoDB by id
	content, err := api.dataStore.Backend.GetContent(id)
	if err != nil {
		// no content found
		handleError(ctx, w, err, logdata)
		return
	}

	// User is not authenticated and hence has only access to current sub document(s)
	var result models.PublicContent

	if content.Current == nil {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	// Add spotlight first
	result.TotalCount = addPublicItem(&result, "spotlight", content.Current.Spotlight, content.ID, content.Current.State)
	// !!! add the following in the order as per doc's elsewhere
	result.TotalCount += addPublicItem(&result, "articles", content.Current.Articles, content.ID, content.Current.State)
	result.TotalCount += addPublicItem(&result, "bulletins", content.Current.Bulletins, content.ID, content.Current.State)
	result.TotalCount += addPublicItem(&result, "methodologies", content.Current.Methodologies, content.ID, content.Current.State)
	result.TotalCount += addPublicItem(&result, "methodologyArticles", content.Current.MethodologyArticles, content.ID, content.Current.State)
	result.TotalCount += addPublicItem(&result, "staticDatasets", content.Current.StaticDatasets, content.ID, content.Current.State)
	result.TotalCount += addPublicItem(&result, "timeseries", content.Current.Timeseries, content.ID, content.Current.State)

	if result.TotalCount == 0 {
		// no content exist for the requested ID
		handleError(ctx, w, apierrors.ErrNotFound, logdata)
		// !!! OR should this be, go over with Eleanor
		// 		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
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
		"content_id": id,
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
