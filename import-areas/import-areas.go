package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ONSdigital/log.go/log"
)

type returnedData struct {
	head    map[string][]string
	results map[string][]*countryData
}
type countryData struct {
	areaname []map[string]string
	areacode []map[string]string
	code     []map[string]string
	name     []map[string]string
}

func main() {

	/*
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
	*/

	queryCountryEngland := `PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
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

	v := url.Values{}
	v.Set("query", queryCountryEngland)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://statistics.data.gov.uk/sparql", strings.NewReader(v.Encode()))
	req.Header.Add("Accept", "application/sparql-results+json")
	resp, _ := client.Do(req)

	data, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Printf("%s\n", data)

	//var returnedEnglandData map[string]interface{}
	var returnedEnglandData returnedData
	ctx := context.Background()

	if err := json.Unmarshal(data, &returnedEnglandData); err != nil {
		logData := log.Data{"json file": returnedEnglandData}
		log.Event(ctx, "failed to unmarshal json", log.ERROR, log.Error(err), logData)
		os.Exit(1)
	}

	fmt.Printf("%#v\n", returnedEnglandData)

	queryCountryWales := `PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
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

	v = url.Values{}
	v.Set("query", queryCountryWales)

	client = &http.Client{}
	req, _ = http.NewRequest("POST", "http://statistics.data.gov.uk/sparql", strings.NewReader(v.Encode()))
	req.Header.Add("Accept", "application/sparql-results+json")
	resp, _ = client.Do(req)

	data, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	//fmt.Printf("%s\n", data)

	/*
			queryCountryWales := `PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		    PREFIX statdef: <http://statistics.data.gov.uk/def/statistical-entity#>
		    PREFIX geodef: <http://statistics.data.gov.uk/def/statistical-geography#>
		    PREFIX statid: <http://statistics.data.gov.uk/id/statistical-entity/>
		    PREFIX pmdfoi: <http://publishmydata.com/def/ontology/foi/>
		    SELECT DISTINCT ?areacode ?areaname ?parentcode ?parentname ?sibling
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
				?sibling geodef:parentcode ?parent;
			}
			LIMIT 5
			`
	*/
	//        VALUES ?types { statid:E12 statid:E06 statid:E47 statid:E11 statid:E10 statid:E09 statid:E08 statid:E07 statid:E05 }

	/*
		v2 := url.Values{}
		v2.Set("query", query2)

		client = &http.Client{}
		req, _ = http.NewRequest("POST", "http://statistics.data.gov.uk/sparql", strings.NewReader(v2.Encode()))
		req.Header.Add("Accept", "application/sparql-results+json")
		resp, _ = client.Do(req)

		data, _ = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		fmt.Printf("%s\n", data)
	*/
}
