package api

import (
	"net/http"
	"sort"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

func addPublicItem(contentList *models.ContentResponseAPI, typeName string, itemLink *[]models.TypeLinkObject, id string, state string, privateResponse bool) {
	var count int

	if itemLink == nil {
		return
	}

	nofItems := len(*itemLink)
	if nofItems == 0 {
		return
	}

	// Create list of sorted href's from itemLink list
	hrefs := make([]string, 0, nofItems)
	for _, field := range *itemLink {
		hrefs = append(hrefs, field.HRef)
	}
	sort.Strings(hrefs)

	// Iterate through sorted hrefs and use each one to select item from
	// itemLink in alphabetical order
	for _, href := range hrefs {
		for _, item := range *itemLink {
			if href == item.HRef {
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
				if privateResponse {
					cItem.State = state
				}
				cItem.Links = &cLinks

				if contentList.Items == nil {
					contentList.Items = &[]models.ContentItem{cItem}
				} else {
					*contentList.Items = append(*contentList.Items, cItem)
				}

				count++
			}
		}
	}

	contentList.TotalCount = contentList.TotalCount + count
}

// getContentPublicHandler is a handler that gets content by its id from MongoDB for Web
func (api *API) getContentPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
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

	if content.Current == nil {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	var currentResult models.ContentResponseAPI

	// Add spotlight first
	addPublicItem(&currentResult, "spotlight", content.Current.Spotlight, content.ID, content.Current.State, false)
	// then Publications (alphabetically ordered)
	addPublicItem(&currentResult, "articles", content.Current.Articles, content.ID, content.Current.State, false)
	addPublicItem(&currentResult, "bulletins", content.Current.Bulletins, content.ID, content.Current.State, false)
	addPublicItem(&currentResult, "methodologies", content.Current.Methodologies, content.ID, content.Current.State, false)
	addPublicItem(&currentResult, "methodologyArticles", content.Current.MethodologyArticles, content.ID, content.Current.State, false)
	// then Datasets (alphabetically ordered)
	addPublicItem(&currentResult, "staticDatasets", content.Current.StaticDatasets, content.ID, content.Current.State, false)
	addPublicItem(&currentResult, "timeseries", content.Current.Timeseries, content.ID, content.Current.State, false)

	currentResult.Count = currentResult.TotalCount // This may be '0' which is the case for some existing ONS pages (like: bankruptcyinsolvency as of 3.feb.2021)

	if err := WriteJSONBody(ctx, currentResult, w, logdata); err != nil {
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

	// get content from mongoDB by id
	content, err := api.dataStore.Backend.GetContent(id)
	if err != nil {
		// no content found
		handleError(ctx, w, err, logdata)
		return
	}

	// User has valid authentication to get raw full content document(s)

	if content.Current == nil {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}
	if content.Next == nil {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	var currentResult models.ContentResponseAPI

	// Add spotlight first
	addPublicItem(&currentResult, "spotlight", content.Current.Spotlight, content.ID, content.Current.State, true)
	// then Publications (alphabetically ordered)
	addPublicItem(&currentResult, "articles", content.Current.Articles, content.ID, content.Current.State, true)
	addPublicItem(&currentResult, "bulletins", content.Current.Bulletins, content.ID, content.Current.State, true)
	addPublicItem(&currentResult, "methodologies", content.Current.Methodologies, content.ID, content.Current.State, true)
	addPublicItem(&currentResult, "methodologyArticles", content.Current.MethodologyArticles, content.ID, content.Current.State, true)
	// then Datasets (alphabetically ordered)
	addPublicItem(&currentResult, "staticDatasets", content.Current.StaticDatasets, content.ID, content.Current.State, true)
	addPublicItem(&currentResult, "timeseries", content.Current.Timeseries, content.ID, content.Current.State, true)

	if currentResult.TotalCount == 0 {
		// no content exist for the requested ID
		handleError(ctx, w, apierrors.ErrNotFound, logdata)
		// !!! OR should this be, go over with Eleanor
		// 		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}
	currentResult.Count = currentResult.TotalCount

	// The 'Next' list may be a different length to the current, so we do the above again, but for Next
	var nextResult models.ContentResponseAPI

	// Add spotlight first
	addPublicItem(&nextResult, "spotlight", content.Next.Spotlight, content.ID, content.Next.State, true)
	// then Publications (alphabetically ordered)
	addPublicItem(&nextResult, "articles", content.Next.Articles, content.ID, content.Next.State, true)
	addPublicItem(&nextResult, "bulletins", content.Next.Bulletins, content.ID, content.Next.State, true)
	addPublicItem(&nextResult, "methodologies", content.Next.Methodologies, content.ID, content.Next.State, true)
	addPublicItem(&nextResult, "methodologyArticles", content.Next.MethodologyArticles, content.ID, content.Next.State, true)
	// then Datasets (alphabetically ordered)
	addPublicItem(&nextResult, "staticDatasets", content.Next.StaticDatasets, content.ID, content.Next.State, true)
	addPublicItem(&nextResult, "timeseries", content.Next.Timeseries, content.ID, content.Next.State, true)

	if nextResult.TotalCount == 0 {
		// no content exist for the requested ID
		handleError(ctx, w, apierrors.ErrNotFound, logdata)
		// !!! OR should this be, go over with Eleanor
		// 		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}
	nextResult.Count = nextResult.TotalCount

	var result models.PrivateContentResponseAPI

	result.Next = &nextResult
	result.Current = &currentResult

	if err := WriteJSONBody(ctx, result, w, logdata); err != nil {
		return
	}
	log.Event(ctx, "request successful", log.INFO, logdata) // NOTE: name of function is in logdata
}
