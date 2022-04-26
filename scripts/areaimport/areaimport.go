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
    csvFile, _ := os.Open("/Users/indra/Docs/ons/Code_History_Database_(December_2021)_UK/ChangeHistory.csv")
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

        areaInfo := models.AreaParams{
            Code:         line[0],
            AreaName:     &areaName,
            Visible:      BoolPointer(true),
            ActiveFrom:   &active_from,
            ActiveTo:     &active_to,
            ParentId:     line[7],
            AreaHectares: line[11],
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

        success = append(success, "api response for line[0]: "+resp.Status)

    }
    logs := logs{
        errors:  errors,
        success: success,
    }
    return logs
}
func main() {
    result := importChangeHistoryAreaInfo()

    fmt.Println(result)
}