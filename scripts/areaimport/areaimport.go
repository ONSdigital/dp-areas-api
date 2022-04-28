package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ONSdigital/dp-areas-api/models"
)

var areaTypeAndCode = map[string]string{
	"E92": "Country",
	"E06": "Unitary Authorities",
	"W92": "Country",
	"E12": "Region",
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

// http://127.0.0.1:25500/v1/areas/
// /Users/indra/Docs/ons/Code_History_Database_(December_2021)_UK/ChangeHistory.csv
type Config struct {
	CSVFilePath   string `envconfig:"CSV_FILE_PATH" required:"true"`
	AreaUpdateUrl string `envconfig:"AREA_UPDATE_URL" required:"true"`
}
type logs struct {
	errors  []string
	success []string
}

func BoolPointer(b bool) *bool {
	return &b
}

func importChangeHistoryAreaInfo() logs {
	var errors []string
	var success []string
	csvFile, err := os.Open(config.CSVFilePath)
	if err != nil {
		log.Fatalf("Failed to open the CSV on path %+v:", err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if line[0] == "GEOGCD" {
			continue
		}
		if line[10] != "live" {
			continue
		}
		if line[1] == "" {
			errors = append(errors, "found area info with empty name: "+line[0])
			continue
		}

		_, areaTypeFound := areaTypeAndCode[line[8]]
		if !areaTypeFound {
			errors = append(errors, "area type not drive for: "+line[0])
			continue
		}

		const layout = "02/01/2006 15:04:05"
		active_from, err := time.Parse(layout, line[5])
		if err != nil {
			log.Fatal("Could not parse time:", err)
		}
		active_to, _ := time.Parse(layout, line[6])

		areaName := models.AreaName{
			Name:       line[1],
			ActiveFrom: &active_from,
			ActiveTo:   &active_to,
		}

		hectares, _ := strconv.ParseFloat(line[11], 64)
		areaInfo := models.AreaParams{
			Code:         line[0],
			AreaName:     &areaName,
			Visible:      BoolPointer(true),
			ActiveFrom:   &active_from,
			ActiveTo:     &active_to,
			ParentCode:   line[7],
			AreaHectares: hectares,
		}

		json, err := json.Marshal(areaInfo)

		if err != nil {
			log.Fatalf("json marshalling failed: %+v", err)
		}

		client := &http.Client{}

		req, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:25500/v1/areas/"+line[0], bytes.NewBuffer(json))
		if err != nil {
			log.Fatalf("request failed %+v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("error response from server: %+v", err)
		}
		defer resp.Body.Close()

		success = append(success, "api response for line[0]: "+resp.Status)

	}
	logs := logs{
		errors:  errors,
		success: success,
	}
	return logs
}
func main() {
	config := getConfig()
	result := importChangeHistoryAreaInfo(config)

	fmt.Println(result)
}
