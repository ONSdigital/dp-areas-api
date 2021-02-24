package api

import (
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-topic-api/models"
)

// Flag values for a query type:
const (
	QuerySpotlightFlag int = 1 << iota // powers of 2, for combining bit flags

	// Publications:
	QueryArticlesFlag
	QueryBulletinsFlag
	QueryMethodologiesFlag
	QueryMethodologyArticlesFlag

	// Datasets:
	QueryStaticDatasetsFlag
	QueryTimeseriesFlag
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
	spotlightStr:           QuerySpotlightFlag,
	articlesStr:            QueryArticlesFlag,
	bulletinsStr:           QueryBulletinsFlag,
	methodologiesStr:       QueryMethodologiesFlag,
	methodologyarticlesStr: QueryMethodologyArticlesFlag,
	staticdatasetsStr:      QueryStaticDatasetsFlag,
	timeseriesStr:          QueryTimeseriesFlag,

	publicationsStr: QueryArticlesFlag | QueryBulletinsFlag | QueryMethodologiesFlag | QueryMethodologyArticlesFlag,

	datasetsStr: QueryStaticDatasetsFlag | QueryTimeseriesFlag,
}

// getContentTypeParameter obtains a filter that defines a set of possible types
func getContentTypeParameter(queryVars url.Values) int {
	valArray, found := queryVars["type"]
	if !found {
		// no type specified, so return flags for all types
		return QuerySpotlightFlag | QueryArticlesFlag | QueryBulletinsFlag | QueryMethodologiesFlag |
			QueryMethodologyArticlesFlag | QueryStaticDatasetsFlag | QueryTimeseriesFlag
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
func getRequiredItems(queryTypeFlags int, content *models.Content, id string) *models.ContentResponseAPI {
	var result models.ContentResponseAPI

	// Add spotlight first
	if (queryTypeFlags & QuerySpotlightFlag) != 0 {
		result.AppendLinkInfo(spotlightStr, content.Spotlight, id, content.State)
	}

	// then Publications (alphabetically ordered)
	if (queryTypeFlags & QueryArticlesFlag) != 0 {
		result.AppendLinkInfo(articlesStr, content.Articles, id, content.State)
	}
	if (queryTypeFlags & QueryBulletinsFlag) != 0 {
		result.AppendLinkInfo(bulletinsStr, content.Bulletins, id, content.State)
	}
	if (queryTypeFlags & QueryMethodologiesFlag) != 0 {
		result.AppendLinkInfo(methodologiesStr, content.Methodologies, id, content.State)
	}
	if (queryTypeFlags & QueryMethodologyArticlesFlag) != 0 {
		result.AppendLinkInfo(methodologyarticlesStr, content.MethodologyArticles, id, content.State)
	}

	// then Datasets (alphabetically ordered)
	if (queryTypeFlags & QueryStaticDatasetsFlag) != 0 {
		result.AppendLinkInfo(staticdatasetsStr, content.StaticDatasets, id, content.State)
	}
	if (queryTypeFlags & QueryTimeseriesFlag) != 0 {
		result.AppendLinkInfo(timeseriesStr, content.Timeseries, id, content.State)
	}

	result.Count = result.TotalCount

	return &result
}
