package api

import (
	"net/http"
	"net/url"
	"sort"
	"strings"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

func addItem(contentList *models.ContentResponseAPI, typeName string, itemLink *[]models.TypeLinkObject, id string, state string) {
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
		// build up data items into structure
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
			State: state,
		}

		if contentList.Items == nil {
			contentList.Items = &[]models.ContentItem{cItem}
		} else {
			*contentList.Items = append(*contentList.Items, cItem)
		}
	}

	contentList.TotalCount = contentList.TotalCount + nofItems
}

func getRequiredItems(queryType int, content *models.Content, id string) models.ContentResponseAPI {
	var result models.ContentResponseAPI

	// Add spotlight first
	if (queryType & querySpotlight) != 0 {
		addItem(&result, spotlightStr, content.Spotlight, id, content.State)
	}

	// then Publications (alphabetically ordered)
	if (queryType & queryAarticles) != 0 {
		addItem(&result, articlesStr, content.Articles, id, content.State)
	}
	if (queryType & queryBulletins) != 0 {
		addItem(&result, bulletinsStr, content.Bulletins, id, content.State)
	}
	if (queryType & queryMethodologies) != 0 {
		addItem(&result, methodologiesStr, content.Methodologies, id, content.State)
	}
	if (queryType & queryMethodologyArticles) != 0 {
		addItem(&result, methodologyarticlesStr, content.MethodologyArticles, id, content.State)
	}

	// then Datasets (alphabetically ordered)
	if (queryType & queryStaticDatasets) != 0 {
		addItem(&result, staticdatasetsStr, content.StaticDatasets, id, content.State)
	}
	if (queryType & queryTimeseries) != 0 {
		addItem(&result, timeseriesStr, content.Timeseries, id, content.State)
	}

	result.Count = result.TotalCount

	return result
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

	// get type from query parameters, or default value
	queryType := getContentTypeParameter(req.URL.Query())

	// check topic from mongoDB by id
	err := api.dataStore.Backend.CheckTopicExists(id)
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
		handleError(ctx, w, apierrors.ErrContentNotFound, logdata)
		return
	}

	currentResult := getRequiredItems(queryType, content.Current, content.ID)
	// currentResult.TotalCount // This may be '0' which is the case for some existing ONS pages (like: bankruptcyinsolvency as of 3.feb.2021)

	if queryType != 0 && currentResult.TotalCount == 0 {
		handleError(ctx, w, apierrors.ErrContentNotFound, logdata)
		return
	}

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

	// get type from query parameters, or default value
	queryType := getContentTypeParameter(req.URL.Query())

	// check topic from mongoDB by id
	err := api.dataStore.Backend.CheckTopicExists(id)
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

	currentResult := getRequiredItems(queryType, content.Current, content.ID)

	// The 'Next' type items may have a different length to the current, so we do the above again, but for Next
	nextResult := getRequiredItems(queryType, content.Next, content.ID)

	if queryType != 0 && currentResult.TotalCount == 0 && nextResult.TotalCount == 0 {
		handleError(ctx, w, apierrors.ErrContentNotFound, logdata)
		return
	}

	var result models.PrivateContentResponseAPI
	result.Next = &nextResult
	result.Current = &currentResult

	if err := WriteJSONBody(ctx, result, w, logdata); err != nil {
		return
	}
	log.Event(ctx, "request successful", log.INFO, logdata) // NOTE: name of function is in logdata
}

// Flag values for a query type:
const (
	querySpotlight int = 1 << iota // powers of 2, for bit flags

	// Publications:
	queryAarticles
	queryBulletins
	queryMethodologies
	queryMethodologyArticles

	// Datasets:
	queryStaticDatasets
	queryTimeseries
)

const (
	spotlightStr           = "spotlight"
	articlesStr            = "articles"
	bulletinsStr           = "bulletins"
	methodologiesStr       = "methodologies"
	methodologyarticlesStr = "methodologyarticles"
	staticdatasetsStr      = "staticdatasets"
	timeseriesStr          = "timeseries"
	publicationsStr        = "publications"
	datasetsStr            = "datasets"
)

var querySets map[string]int = map[string]int{
	// search keys are done as lower case to make searches work regardless of case
	spotlightStr:           querySpotlight,
	articlesStr:            queryAarticles,
	bulletinsStr:           queryBulletins,
	methodologiesStr:       queryMethodologies,
	methodologyarticlesStr: queryMethodologyArticles,
	staticdatasetsStr:      queryStaticDatasets,
	timeseriesStr:          queryTimeseries,

	publicationsStr: queryAarticles | queryBulletins | queryMethodologies | queryMethodologyArticles,

	datasetsStr: queryStaticDatasets | queryTimeseries,
}

// getContentTypeParameter obtains a filter that defines a set of possible types
func getContentTypeParameter(queryVars url.Values) int {
	valArray, found := queryVars["type"]
	if !found {
		// no type specified, so return flags for all types
		return querySpotlight | queryAarticles | queryBulletins | queryMethodologies | queryMethodologyArticles | queryStaticDatasets | queryTimeseries
	}

	// make query type lower case for following comparison to cope with wrong case of letter(s)
	lowerVal := strings.ToLower(valArray[0])

	// also remove leading and trailing whitespace as it casuses the check to fail
	trimmedVal := strings.TrimSpace(lowerVal)

	set, ok := querySets[trimmedVal]
	if ok {
		return set
	}

	return 0
}
