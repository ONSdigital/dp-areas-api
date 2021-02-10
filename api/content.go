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

func addItem(contentList *models.ContentResponseAPI, typeName string, itemLink *[]models.TypeLinkObject, id string, state string, privateResponse bool) {
	if itemLink == nil {
		return
	}

	nofItems := len(*itemLink)
	if nofItems == 0 {
		return
	}

	title := make(map[string]string)

	// Create list of sorted href's from itemLink list
	hrefs := make([]string, nofItems)
	for i, field := range *itemLink {
		hrefs[i] = field.HRef
		title[field.HRef] = field.Title
	}
	sort.Strings(hrefs)

	// Iterate through alphabeticaly sorted 'hrefs' and use each one to select corresponding title
	for _, href := range hrefs {
		// build up data items into structure 'cItem'
		//		var selfLink models.LinkObject = models.LinkObject{
		//			HRef: href}

		//		var topicLink models.LinkObject = models.LinkObject{
		//			ID:   id,
		//			HRef: "/topic/" + id}

		//		var cLinks models.ContentLinks = models.ContentLinks{
		//			Self:  &selfLink,
		//			Topic: &topicLink}

		var cItem models.ContentItem = models.ContentItem{
			Title: title[href],
			Type:  typeName,
			Links: &models.ContentLinks{
				Self: &models.LinkObject{
					HRef: href,
				},
				Topic: &models.LinkObject{
					ID:   id,
					HRef: "/topic/" + id,
				},
			},
		}

		if privateResponse {
			cItem.State = state
		}

		if contentList.Items == nil {
			contentList.Items = &[]models.ContentItem{cItem}
		} else {
			*contentList.Items = append(*contentList.Items, cItem)
		}
	}

	contentList.TotalCount = contentList.TotalCount + nofItems
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

	// check topic from mongoDB by id
	_, err := api.dataStore.Backend.GetTopic(id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
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
	addItem(&currentResult, "spotlight", content.Current.Spotlight, content.ID, content.Current.State, false)
	// then Publications (alphabetically ordered)
	addItem(&currentResult, "articles", content.Current.Articles, content.ID, content.Current.State, false)
	addItem(&currentResult, "bulletins", content.Current.Bulletins, content.ID, content.Current.State, false)
	addItem(&currentResult, "methodologies", content.Current.Methodologies, content.ID, content.Current.State, false)
	addItem(&currentResult, "methodologyArticles", content.Current.MethodologyArticles, content.ID, content.Current.State, false)
	// then Datasets (alphabetically ordered)
	addItem(&currentResult, "staticDatasets", content.Current.StaticDatasets, content.ID, content.Current.State, false)
	addItem(&currentResult, "timeseries", content.Current.Timeseries, content.ID, content.Current.State, false)

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
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"content_id": id,
		"function":   "getContentPrivateHandler",
	}

	// check topic from mongoDB by id
	_, err := api.dataStore.Backend.GetTopic(id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
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
		//TODO
		/*
			In the future: when the API becomes more than read-only
			When a document is first created, it will only have 'next' until it is published, when it gets 'current' populated.
			So current == nil is not an error.

			For now we return an error because we dont have publishing steps.
		*/
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}
	if content.Next == nil {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	var currentResult models.ContentResponseAPI

	// Add spotlight first
	addItem(&currentResult, "spotlight", content.Current.Spotlight, content.ID, content.Current.State, true)
	// then Publications (alphabetically ordered)
	addItem(&currentResult, "articles", content.Current.Articles, content.ID, content.Current.State, true)
	addItem(&currentResult, "bulletins", content.Current.Bulletins, content.ID, content.Current.State, true)
	addItem(&currentResult, "methodologies", content.Current.Methodologies, content.ID, content.Current.State, true)
	addItem(&currentResult, "methodologyArticles", content.Current.MethodologyArticles, content.ID, content.Current.State, true)
	// then Datasets (alphabetically ordered)
	addItem(&currentResult, "staticDatasets", content.Current.StaticDatasets, content.ID, content.Current.State, true)
	addItem(&currentResult, "timeseries", content.Current.Timeseries, content.ID, content.Current.State, true)

	currentResult.Count = currentResult.TotalCount

	// The 'Next' list may be a different length to the current, so we do the above again, but for Next
	var nextResult models.ContentResponseAPI

	// Add spotlight first
	addItem(&nextResult, "spotlight", content.Next.Spotlight, content.ID, content.Next.State, true)
	// then Publications (alphabetically ordered)
	addItem(&nextResult, "articles", content.Next.Articles, content.ID, content.Next.State, true)
	addItem(&nextResult, "bulletins", content.Next.Bulletins, content.ID, content.Next.State, true)
	addItem(&nextResult, "methodologies", content.Next.Methodologies, content.ID, content.Next.State, true)
	addItem(&nextResult, "methodologyArticles", content.Next.MethodologyArticles, content.ID, content.Next.State, true)
	// then Datasets (alphabetically ordered)
	addItem(&nextResult, "staticDatasets", content.Next.StaticDatasets, content.ID, content.Next.State, true)
	addItem(&nextResult, "timeseries", content.Next.Timeseries, content.ID, content.Next.State, true)

	nextResult.Count = nextResult.TotalCount

	var result models.PrivateContentResponseAPI

	result.Next = &nextResult
	result.Current = &currentResult

	if err := WriteJSONBody(ctx, result, w, logdata); err != nil {
		return
	}
	log.Event(ctx, "request successful", log.INFO, logdata) // NOTE: name of function is in logdata
}
