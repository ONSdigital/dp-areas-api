package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/globalsign/mgo"
)

type AreaStruct struct {
	Results BindingsStruct `json:"results,omitempty"`
}
type BindingsStruct struct {
	Bindings []BindingsData `json:"bindings,omitempty"`
}
type BindingsData struct {
	AreaName GenStruct `json:"areaname,omitempty"`
	AreaCode GenStruct `json:"areacode,omitempty"`
	Code     GenStruct `json:"code,omitempty"`
	Name     GenStruct `json:"name,omitempty"`
}
type GenStruct struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

var areaTypeAndCode = map[string]string{
	"E92": "Country",
	"W92": "Country",
	"E12": "Region",
	"E06": "Unitary Authorities",
	"W06": "Unitary Authorities",
	"E47": "Combined Authorities",
	"E11": "Metropolitan Counties",
	"E10": "Counties",
	"E09": "London Boroughs",
	"E08": "Metropolitan Districts",
	"E07": "Non-metropolitan Districts",
	"E05": "Electoral Wards",
	"W05": "Electoral Wards",
}

var queryCountryEngland = `PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX statdef: <http://statistics.data.gov.uk/def/statistical-entity#>
PREFIX geodef: <http://statistics.data.gov.uk/def/statistical-geography#>
PREFIX statid: <http://statistics.data.gov.uk/id/statistical-entity/>
PREFIX pmdfoi: <http://publishmydata.com/def/ontology/foi/>
PREFIX geoid: <http://statistics.data.gov.uk/id/statistical-geography/>
SELECT DISTINCT ?areacode ?areaname ?code ?name
WHERE {
	VALUES ?types { statid:E92 }
	?area statdef:code ?types ;
		geodef:status "live" ;
		geodef:officialname ?areaname ;
		rdfs:label ?areacode .
	?child pmdfoi:parent geoid:E92000001 ;
		pmdfoi:code ?code ;
		geodef:officialname ?name ;
}`

var queryCountryWales = `PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX statdef: <http://statistics.data.gov.uk/def/statistical-entity#>
PREFIX geodef: <http://statistics.data.gov.uk/def/statistical-geography#>
PREFIX statid: <http://statistics.data.gov.uk/id/statistical-entity/>
PREFIX pmdfoi: <http://publishmydata.com/def/ontology/foi/>
PREFIX geoid: <http://statistics.data.gov.uk/id/statistical-geography/>
SELECT DISTINCT ?areacode ?areaname ?code ?name
WHERE {
	VALUES ?types { statid:W92 }
	?area statdef:code ?types ;
		geodef:status "live" ;
		geodef:officialname ?areaname ;
		rdfs:label ?areacode .
	?child pmdfoi:parent geoid:W92000004 ;
		pmdfoi:code ?code ;
	  	geodef:officialname ?name ;
}`

func lookupAreaType(code string) (string, error) {
	prefix := code[:3]

	value, found := areaTypeAndCode[prefix]
	if !found {
		return "", errors.New("unable to find an area type for this code")
	}
	return value, nil
}

func buildCountryArea(areaStruct AreaStruct) (*models.Area, error) {
	ctx := context.Background()

	area := &models.Area{}
	for _, binding := range areaStruct.Results.Bindings {
		area.ID = binding.AreaCode.Value
		area.Name = binding.AreaName.Value
		areaType, err := lookupAreaType(binding.AreaCode.Value)
		logData := log.Data{"AreaCode": binding.AreaCode.Value}
		if err != nil {
			log.Event(ctx, "error returned from areaType lookup", log.ERROR, log.Error(err), logData)
			return nil, err
		}
		area.Type = areaType
		childType, err := lookupAreaType(binding.Code.Value)
		logData = log.Data{"ChildAreaCode": binding.Code.Value}
		if err != nil {
			log.Event(ctx, "error returned from areaType lookup", log.ERROR, log.Error(err), logData)
			return nil, err
		}
		child := &models.LinkedAreas{
			ID:   binding.Code.Value,
			Name: binding.Name.Value,
			Type: childType,
		}
		area.ChildAreas = append(area.ChildAreas, *child)
	}
	return area, nil
}

func postCountryQuery(query string) (*models.Area, error) {
	ctx := context.Background()

	v := url.Values{}
	v.Set("query", query)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://statistics.data.gov.uk/sparql", strings.NewReader(v.Encode()))
	req.Header.Add("Accept", "application/sparql-results+json")
	resp, _ := client.Do(req)

	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var returnedCountryData AreaStruct

	if err := json.Unmarshal(data, &returnedCountryData); err != nil {
		logData := log.Data{"json file": returnedCountryData}
		log.Event(ctx, "failed to unmarshal json", log.ERROR, log.Error(err), logData)
		return nil, err
	}

	countryData, err := buildCountryArea(returnedCountryData)
	if err != nil {
		log.Event(ctx, "error returned from building the country area", log.ERROR, log.Error(err))
		return nil, err
	}

	return countryData, nil
}

func main() {

	var (
		mongoURL string
	)

	flag.StringVar(&mongoURL, "mongo-url", "localhost:27017", "mongoDB URL")
	flag.Parse()

	ctx := context.Background()

	session, err := mgo.Dial(mongoURL)
	if err != nil {
		log.Event(ctx, "unable to create mongo session", log.ERROR, log.Error(err))
		os.Exit(1)
	}
	defer session.Close()

	englandData, err := postCountryQuery(queryCountryEngland)
	if err != nil {
		log.Event(ctx, "error returned from the country query POST", log.ERROR, log.Error(err))
	}

	logData := log.Data{"AreaData": englandData}

	if err = session.DB("areas").C("areas").Insert(englandData); err != nil {
		log.Event(ctx, "failed to insert England area data document, data lost in mongo but exists in this log", log.ERROR, log.Error(err), logData)
		os.Exit(1)
	}

	walesData, err := postCountryQuery(queryCountryWales)
	if err != nil {
		log.Event(ctx, "error returned from the country query POST", log.ERROR, log.Error(err))
	}

	log.Event(ctx, "successfully put England area data into Mongo", log.INFO, logData)

	logData = log.Data{"AreaData": walesData}

	if err = session.DB("areas").C("areas").Insert(walesData); err != nil {
		log.Event(ctx, "failed to insert Wales area data document, data lost in mongo but exists in this log", log.ERROR, log.Error(err), logData)
		os.Exit(1)
	}

	log.Event(ctx, "successfully put Wales area data into Mongo", log.INFO, logData)
}
