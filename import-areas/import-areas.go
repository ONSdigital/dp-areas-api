package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/log.go/log"
	"github.com/globalsign/mgo"
	"gopkg.in/mgo.v2/bson"
)

type AreaStruct struct {
	Results BindingsStruct `json:"results,omitempty"`
}
type BindingsStruct struct {
	Bindings []BindingsData `json:"bindings,omitempty"`
}
type BindingsData struct {
	AreaName   GenStruct `json:"areaname,omitempty"`
	AreaCode   GenStruct `json:"areacode,omitempty"`
	Code       GenStruct `json:"code,omitempty"`
	Name       GenStruct `json:"name,omitempty"`
	ParentCode GenStruct `json:"parentcode,omitempty"`
	ParentName GenStruct `json:"parentname,omitempty"`
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

var countryIDs = map[string]string{"England": "E92000001", "Wales": "W92000004"}
var countryStatIDs = map[string]string{"England": "E92", "Wales": "W92"}

var queryRegionEngland = `PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
PREFIX statdef: <http://statistics.data.gov.uk/def/statistical-entity#>
PREFIX geodef: <http://statistics.data.gov.uk/def/statistical-geography#>
PREFIX statid: <http://statistics.data.gov.uk/id/statistical-entity/>
PREFIX pmdfoi: <http://publishmydata.com/def/ontology/foi/>
SELECT DISTINCT ?areacode ?areaname ?parentcode ?parentname
WHERE {
	VALUES ?types { statid:E12 }
	?area statdef:code ?types ;
		geodef:status "live" ;
		geodef:officialname ?areaname ;
		rdfs:label ?areacode ;
		pmdfoi:parent ?parent .
	?parent geodef:status "live" ;
		geodef:officialname ?parentname ;
		rdfs:label ?parentcode .
}`

func buildCountryQuery(statid, geoid string) string {

	queryCountry := fmt.Sprintf(
		`PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX statdef: <http://statistics.data.gov.uk/def/statistical-entity#>
		PREFIX geodef: <http://statistics.data.gov.uk/def/statistical-geography#>
		PREFIX statid: <http://statistics.data.gov.uk/id/statistical-entity/>
		PREFIX pmdfoi: <http://publishmydata.com/def/ontology/foi/>
		PREFIX geoid: <http://statistics.data.gov.uk/id/statistical-geography/>
		SELECT DISTINCT ?areacode ?areaname ?code ?name
		WHERE {
			VALUES ?types { statid:%v }
			?area statdef:code ?types ;
				geodef:status "live" ;
				geodef:officialname ?areaname ;
				rdfs:label ?areacode .
			?child pmdfoi:parent geoid:%v ;
				pmdfoi:code ?code ;
				geodef:officialname ?name ;
		}`, statid, geoid)
	return queryCountry
}

func buildRegionQuery(geoid string) string {

	queryRegionChildEngland := fmt.Sprintf(
		`PREFIX geodef: <http://statistics.data.gov.uk/def/statistical-geography#>
		PREFIX pmdfoi: <http://publishmydata.com/def/ontology/foi/>
		PREFIX geoid: <http://statistics.data.gov.uk/id/statistical-geography/>
		SELECT DISTINCT ?code ?name
		WHERE {
			?child pmdfoi:parent geoid:%v;
				pmdfoi:code ?code ;
				geodef:officialname ?name.
		}`, geoid)

	return queryRegionChildEngland
}

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
	req, err := http.NewRequest("POST", "http://statistics.data.gov.uk/sparql", strings.NewReader(v.Encode()))
	if err != nil {
		log.Event(ctx, "error returned from building the country request", log.ERROR, log.Error(err))
		return nil, err
	}

	req.Header.Add("Accept", "application/sparql-results+json")
	resp, err := client.Do(req)
	if err != nil {
		log.Event(ctx, "error returned from country POST request", log.ERROR, log.Error(err))
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("Unexpected status code")
		log.Event(ctx, "Response from POST request is not a 200", log.ERROR, log.Error(err))
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Event(ctx, "error returned from reading the response body", log.ERROR, log.Error(err))
		return nil, err
	}

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

func buildAreas(areaStruct AreaStruct) ([]models.Area, error) {
	ctx := context.Background()
	areas := make([]models.Area, 0)

	for _, binding := range areaStruct.Results.Bindings {
		logData := log.Data{"AreaCode": binding.AreaCode.Value}
		areaType, err := lookupAreaType(binding.AreaCode.Value)
		if err != nil {
			log.Event(ctx, "error returned from areaType lookup", log.ERROR, log.Error(err), logData)
			return nil, err
		}
		area := models.Area{
			ID:   binding.AreaCode.Value,
			Name: binding.AreaName.Value,
			Type: areaType,
		}

		logData = log.Data{"ParentAreaCode": binding.ParentCode.Value}
		parentType, err := lookupAreaType(binding.ParentCode.Value)
		if err != nil {
			log.Event(ctx, "error returned from areaType lookup", log.ERROR, log.Error(err), logData)
			return nil, err
		}
		area.ParentAreas = append(area.ParentAreas, models.LinkedAreas{
			ID:   binding.ParentCode.Value,
			Name: binding.ParentName.Value,
			Type: parentType,
		})

		areas = append(areas, area)
	}

	return areas, nil
}

func postAreaQuery(query string) ([]models.Area, error) {
	ctx := context.Background()

	v := url.Values{}
	v.Set("query", query)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://statistics.data.gov.uk/sparql", strings.NewReader(v.Encode()))
	if err != nil {
		log.Event(ctx, "error returned from building an area request", log.ERROR, log.Error(err))
		return nil, err
	}

	req.Header.Add("Accept", "application/sparql-results+json")
	resp, err := client.Do(req)
	if err != nil {
		log.Event(ctx, "error returned from the POST request", log.ERROR, log.Error(err))
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("Unexpected status code")
		log.Event(ctx, "Response from POST request is not a 200", log.ERROR, log.Error(err))
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Event(ctx, "error returned from reading the response body", log.ERROR, log.Error(err))
		return nil, err
	}

	var returnedAreaData AreaStruct

	if err := json.Unmarshal(data, &returnedAreaData); err != nil {
		logData := log.Data{"json file": returnedAreaData}
		log.Event(ctx, "failed to unmarshal json", log.ERROR, log.Error(err), logData)
		return nil, err
	}

	areaData, err := buildAreas(returnedAreaData)
	if err != nil {
		log.Event(ctx, "error returned from building the area", log.ERROR, log.Error(err))
		return nil, err
	}

	areasWithChild := make([]models.Area, 0)

	for _, region := range areaData {

		areaID := region.ID
		query := buildRegionQuery(areaID)

		v := url.Values{}
		v.Set("query", query)

		client := &http.Client{}
		req, err := http.NewRequest("POST", "http://statistics.data.gov.uk/sparql", strings.NewReader(v.Encode()))
		if err != nil {
			log.Event(ctx, "error returned from building a request for the children of an area", log.ERROR, log.Error(err))
			return nil, err
		}

		req.Header.Add("Accept", "application/sparql-results+json")
		resp, err := client.Do(req)
		if err != nil {
			log.Event(ctx, "error returned from the POST request", log.ERROR, log.Error(err))
			return nil, err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			err = errors.New("Unexpected status code")
			log.Event(ctx, "Response from POST request is not a 200", log.ERROR, log.Error(err))
			return nil, err
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Event(ctx, "error returned from reading the response body", log.ERROR, log.Error(err))
			return nil, err
		}

		var returnedChildData AreaStruct

		if err := json.Unmarshal(data, &returnedChildData); err != nil {
			logData := log.Data{"json file": returnedChildData}
			log.Event(ctx, "failed to unmarshal json", log.ERROR, log.Error(err), logData)
			return nil, err
		}

		regionDataWithChild := buildChildAreas(returnedChildData, region)

		areasWithChild = append(areasWithChild, regionDataWithChild)
	}
	return areasWithChild, nil
}

func buildChildAreas(childStruct AreaStruct, region models.Area) models.Area {

	ctx := context.Background()

	for _, binding := range childStruct.Results.Bindings {

		childType, err := lookupAreaType(binding.Code.Value)
		logData := log.Data{"ChildAreaCode": binding.Code.Value}
		if err != nil {
			log.Event(ctx, "error returned from areaType lookup", log.ERROR, log.Error(err), logData)
			continue
		}

		child := &models.LinkedAreas{
			ID:   binding.Code.Value,
			Name: binding.Name.Value,
			Type: childType,
		}
		region.ChildAreas = append(region.ChildAreas, *child)

	}
	return region

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

	englandQuery := buildCountryQuery(countryStatIDs["England"], countryIDs["England"])
	englandData, err := postCountryQuery(englandQuery)
	if err != nil {
		log.Event(ctx, "error returned from the country query POST", log.ERROR, log.Error(err))
		os.Exit(1)
	}

	logData := log.Data{"AreaData": englandData}

	if _, err = session.DB("areas").C("areas").Upsert(bson.M{"id": countryIDs["England"]}, englandData); err != nil {
		log.Event(ctx, "failed to insert England area data document, data lost in mongo but exists in this log", log.ERROR, log.Error(err), logData)
		os.Exit(1)
	}

	log.Event(ctx, "successfully put England area data into Mongo", log.INFO, logData)

	walesQuery := buildCountryQuery(countryStatIDs["Wales"], countryIDs["Wales"])
	walesData, err := postCountryQuery(walesQuery)
	if err != nil {
		log.Event(ctx, "error returned from the country query POST", log.ERROR, log.Error(err))
		os.Exit(1)
	}

	logData = log.Data{"AreaData": walesData}

	if _, err = session.DB("areas").C("areas").Upsert(bson.M{"id": countryIDs["Wales"]}, walesData); err != nil {
		log.Event(ctx, "failed to insert Wales area data document, data lost in mongo but exists in this log", log.ERROR, log.Error(err), logData)
		os.Exit(1)
	}

	log.Event(ctx, "successfully put Wales area data into Mongo", log.INFO, logData)

	regionData, err := postAreaQuery(queryRegionEngland)
	if err != nil {
		log.Event(ctx, "error returned from the region query POST", log.ERROR, log.Error(err))
		os.Exit(1)
	}

	for _, region := range regionData {
		logData = log.Data{"AreaData": region}
		if _, err = session.DB("areas").C("areas").Upsert(bson.M{"id": region.ID}, region); err != nil {
			log.Event(ctx, "failed to insert region area data document, data lost in mongo but exists in this log", log.ERROR, log.Error(err), logData)
			os.Exit(1)
		}
	}

	logData = log.Data{"AreaData": regionData}
	log.Event(ctx, "successfully put region area data into Mongo", log.INFO, logData)
}
