package api

import (
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-topic-api/models"
)

// Flag values for a query type:
const (
	QuerySpotlight int = 1 << iota // powers of 2, for combining bit flags

	// Publications:
	QueryArticles
	QueryBulletins
	QueryMethodologies
	QueryMethodologyArticles

	// Datasets:
	QueryStaticDatasets
	QueryTimeseries
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
	spotlightStr:           QuerySpotlight,
	articlesStr:            QueryArticles,
	bulletinsStr:           QueryBulletins,
	methodologiesStr:       QueryMethodologies,
	methodologyarticlesStr: QueryMethodologyArticles,
	staticdatasetsStr:      QueryStaticDatasets,
	timeseriesStr:          QueryTimeseries,

	publicationsStr: QueryArticles | QueryBulletins | QueryMethodologies | QueryMethodologyArticles,

	datasetsStr: QueryStaticDatasets | QueryTimeseries,
}

// getContentTypeParameter obtains a filter that defines a set of possible types
func getContentTypeParameter(queryVars url.Values) int {
	valArray, found := queryVars["type"]
	if !found {
		// no type specified, so return flags for all types
		return QuerySpotlight | QueryArticles | QueryBulletins | QueryMethodologies |
			QueryMethodologyArticles | QueryStaticDatasets | QueryTimeseries
	}

	// make query type lower case for following comparison to cope with wrong case of letter(s)
	lowerVal := strings.ToLower(valArray[0])

	// also remove leading and trailing whitespace as it casuses the check to fail
	trimmedVal := strings.TrimSpace(lowerVal)

	set, ok := querySets[trimmedVal]
	if ok {
		return set // return bit flag or flags for requested query
	}

	return 0 // query not recognised, so bad request
}

// getRequiredItems builds up a list of required links info in specifc order as commented within function
func getRequiredItems(queryType int, content *models.Content, id string) *models.ContentResponseAPI {
	var result models.ContentResponseAPI

	// Add spotlight first
	if (queryType & QuerySpotlight) != 0 {
		result.AppendLinkInfo(spotlightStr, content.Spotlight, id, content.State)
	}

	// then Publications (alphabetically ordered)
	if (queryType & QueryArticles) != 0 {
		result.AppendLinkInfo(articlesStr, content.Articles, id, content.State)
	}
	if (queryType & QueryBulletins) != 0 {
		result.AppendLinkInfo(bulletinsStr, content.Bulletins, id, content.State)
	}
	if (queryType & QueryMethodologies) != 0 {
		result.AppendLinkInfo(methodologiesStr, content.Methodologies, id, content.State)
	}
	if (queryType & QueryMethodologyArticles) != 0 {
		result.AppendLinkInfo(methodologyarticlesStr, content.MethodologyArticles, id, content.State)
	}

	// then Datasets (alphabetically ordered)
	if (queryType & QueryStaticDatasets) != 0 {
		result.AppendLinkInfo(staticdatasetsStr, content.StaticDatasets, id, content.State)
	}
	if (queryType & QueryTimeseries) != 0 {
		result.AppendLinkInfo(timeseriesStr, content.Timeseries, id, content.State)
	}

	result.Count = result.TotalCount

	return &result
}
