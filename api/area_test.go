package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/ONSdigital/dp-areas-api/api"
	"github.com/ONSdigital/dp-areas-api/api/mock"
	"github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	EnglandAreaData       = "E92000001"
	WalesAreaData         = "W92000004"
	SwanseaAirportBuaData = "W37000454"
	YorkshireAreaData     = "E12000003"
	SheffieldAreaData     = "E08000019"
)

var (
	EnglandName             = "England"
	WalesName               = "Wales"
	SheffieldName           = "Sheffield"
	isVisible               = true
	countryAreaType         = "Country"
	latitude        float64 = 53.65955162358695
	longitude       float64 = -1.434126224128561
)

var mu sync.Mutex

var ancestors = map[string][]models.AreasAncestors{
	EnglandAreaData: {},
	WalesAreaData:   {},
	SheffieldAreaData: {
		{YorkshireAreaData, "Yorkshire and the Humber"},
		{EnglandAreaData, "England"},
	},
}

func GetAPIWithRDSMocks(mockedRDSStore api.RDSAreaStore) (*api.API, error) {
	mu.Lock()
	defer mu.Unlock()
	cfg, err := config.Get()
	So(err, ShouldBeNil)
	return api.Setup(context.Background(), cfg, mux.NewRouter(), mockedRDSStore)
}

func TestGetBoundaryDataReturnsOk(t *testing.T) {
	Convey("Given a successful request - E92000001", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/boundaries/%s", EnglandAreaData), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "en")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			GetAreaFunc: func(ctx context.Context, areaId string) (*models.AreasDataResults, error) {
				return &models.AreasDataResults{Code: "E92000001", Name: &EnglandName, GeometricData: testGeometricData(), Visible: &isVisible, AreaType: &countryAreaType}, nil
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request boundary data is served", func() {
			Convey("Then an OK response is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedArea := models.BoundaryDataResults{}
				err = json.Unmarshal(payload, &returnedArea)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
				So(returnedArea.Columns, ShouldEqual, "area_id, centroid_bng, centroid, boundary")
				So(returnedArea.Values.AreaID, ShouldEqual, EnglandAreaData)
				So(returnedArea.Values.Boundary, ShouldEqual, "[[[-4.333344310304969,51.7249345567844],[-4.333714799280278,51.72459176590723],[-4.333228471051473,51.72485054948815],[-4.332484708106176,51.72356674989395],[-4.331536255737817,51.72317184002521],[-4.331998050964103,51.72328323509102],[-4.331985242138722,51.72303436581642],[-4.331926419513477,51.72324150346475],[-4.331169197082952,51.72279244850181],[-4.332242159236069,51.72127380479937],[-4.332572468607179,51.72125016436441],[-4.332702576920521,51.7210758648653],[-4.332779021783574,51.720677780977],[-4.334653645452485,51.72266301535556],[-4.333985539454364,51.72284809808241],[-4.334603522966874,51.7227584105357],[-4.334424154980614,51.72464599982296],[-4.334172914205174,51.72488393285414],[-4.334206192564066,51.72446153276184],[-4.333894953495085,51.7251385780543],[-4.333344310304969,51.7249345567844]]]")
				So(returnedArea.Values.Centroid, ShouldEqual, "[-4.333344310304969,51.7249345567844]")
				So(returnedArea.Values.CentroidBng, ShouldEqual, "[-4.333344310304969,51.7249345567844]")
			})
		})
	})
}

func TestGetBoundaryDataWithUnknownIDReturnsNotFound(t *testing.T) {
	Convey("Given a unsuccessful request - E92000002", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/boundaries/%s", "E92000002"), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "en")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			GetAreaFunc: func(ctx context.Context, areaId string) (*models.AreasDataResults, error) {
				return &models.AreasDataResults{Code: "E92000002", Name: &EnglandName, GeometricData: testGeometricData(), Visible: &isVisible, AreaType: &countryAreaType}, nil
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request boundary data is served", func() {
			Convey("Then an NOT OK response is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedArea := models.BoundaryDataResults{}
				err = json.Unmarshal(payload, &returnedArea)
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

func TestGetAreaDataReturnsOkForEngland(t *testing.T) {
	Convey("Given a successful request - E92000001", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s", EnglandAreaData), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "en")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			GetAreaFunc: func(ctx context.Context, areaId string) (*models.AreasDataResults, error) {
				return &models.AreasDataResults{Code: "E92000001", Name: &EnglandName, GeometricData: testGeometricData(), Visible: &isVisible, AreaType: &countryAreaType}, nil
			},
			GetAncestorsFunc: func(areaCode string) ([]models.AreasAncestors, error) {
				return ancestors[WalesAreaData], nil
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area data is served", func() {

			Convey("Then an OK response is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedArea := models.AreasDataResults{}
				err = json.Unmarshal(payload, &returnedArea)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
				So(returnedArea.Code, ShouldEqual, EnglandAreaData)
				So(returnedArea.GeometricData[0][0], ShouldResemble, [2]float64{longitude, latitude})
				So(*returnedArea.Name, ShouldEqual, "England")
				So(*returnedArea.AreaType, ShouldEqual, "Country")
				So(*returnedArea.Visible, ShouldEqual, true)
				So(returnedArea.Ancestors, ShouldResemble, ancestors[EnglandAreaData])
			})
		})
	})
}

func TestGetAreaDataReturnsOkForWales(t *testing.T) {
	Convey("Given a successful request for - W92000004", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s", WalesAreaData), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "cy")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			GetAreaFunc: func(ctx context.Context, areaId string) (*models.AreasDataResults, error) {
				return &models.AreasDataResults{Code: "W92000004", Name: &WalesName, Visible: &isVisible, AreaType: &countryAreaType}, nil
			},
			GetAncestorsFunc: func(areaCode string) ([]models.AreasAncestors, error) {
				return ancestors[WalesAreaData], nil
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area data is served", func() {

			Convey("Then an OK response is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedArea := models.AreasDataResults{}
				err = json.Unmarshal(payload, &returnedArea)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
				So(returnedArea.Code, ShouldEqual, WalesAreaData)
				So(*returnedArea.Name, ShouldEqual, "Wales")
				So(*returnedArea.AreaType, ShouldEqual, "Country")
				So(returnedArea.Ancestors, ShouldResemble, ancestors[WalesAreaData])
			})
		})
	})
}

func TestGetAreaDataReturnsValidationErrors(t *testing.T) {
	Convey("Given a request with invalid area code", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s", WalesAreaData+"z"), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "cz")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			GetAreaFunc: func(ctx context.Context, areaId string) (*models.AreasDataResults, error) {
				return nil, apierrors.ErrNoRows
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area data is served", func() {

			Convey("Then 404 http status should be returned", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})

		})
	})
}

func TestGetAreaDataReturnsValidationError(t *testing.T) {
	Convey("Given a request to stubbed area data - validation errors returned", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s", WalesAreaData), nil)
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area data is served", func() {

			Convey("Then validation failure details returned - 1 validation error generated", func() {
				payload, _ := ioutil.ReadAll(w.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(payload, &responseBody)
				So(len(responseBody["errors"].([]interface{})), ShouldEqual, 1)

				// error - accept language header not found
				error := responseBody["errors"].([]interface{})[0].(map[string]interface{})
				So(error["code"], ShouldEqual, models.AcceptLanguageHeaderError)
				So(error["description"], ShouldEqual, models.AcceptLanguageHeaderNotFoundDescription)
			})

		})
	})
}

func TestGetAreaRelationshipsReturnsOk(t *testing.T) {
	Convey("Given a successful request to area relationship data - E92000001", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s/relations", EnglandAreaData), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "en")
		w := httptest.NewRecorder()
		relatedAreas := []*models.AreaBasicData{
			{"E12000001", "North East"},
			{"E12000002", "North West"},
			{"E12000003", "Yorkshire and The Humber"},
		}

		expectedRelationShips := []*models.AreaRelationShips{
			{"E12000001", "North East", "/v1/area/E12000001"},
			{"E12000002", "North West", "/v1/area/E12000002"},
			{"E12000003", "Yorkshire and The Humber", "/v1/area/E12000003"},
		}

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			ValidateAreaFunc: func(areaCode string) error {
				return nil
			},
			GetRelationshipsFunc: func(areaCode, relationshipParameter string) ([]*models.AreaBasicData, error) {
				return relatedAreas, nil
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area relationship data is served", func() {

			Convey("Then an OK response is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				var relationsShips []*models.AreaRelationShips
				err = json.Unmarshal(payload, &relationsShips)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
				So(relationsShips, ShouldResemble, expectedRelationShips)
			})
		})
	})
}

func TestGetAreaRelationshipsFailsForInvalidIds(t *testing.T) {
	Convey("Given a successful request to area relationship data - invalid", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s/relations", "InvalidAreaCode"), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "en")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			ValidateAreaFunc: func(areaCode string) error {
				return apierrors.ErrNoRows
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area data is served", func() {

			Convey("Then an 404 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
	})
}

func TestGetAreaRelationshipsWithParameterReturnsOk(t *testing.T) {
	Convey("Given a successful request for child area relationship data - E92000001", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s/relations?relationship=child", YorkshireAreaData), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "en")
		w := httptest.NewRecorder()
		childRelatedAreas := []*models.AreaBasicData{
			{SheffieldAreaData, SheffieldName},
		}
		relatedAreas := []*models.AreaBasicData{
			{"E92000001", "North East"},
			{"E12000002", "North West"},
			{SheffieldAreaData, SheffieldName},
		}

		expectedChildRelationShips := []*models.AreaRelationShips{
			{SheffieldAreaData, SheffieldName, "/v1/area/E08000019"},
		}

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			ValidateAreaFunc: func(areaCode string) error {
				return nil
			},
			GetRelationshipsFunc: func(areaCode, relationshipParameter string) ([]*models.AreaBasicData, error) {
				if relationshipParameter == "child" {
					return childRelatedAreas, nil
				}
				return relatedAreas, nil
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request child area relationship data is served", func() {

			Convey("Then an OK response is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				var relationsShips []*models.AreaRelationShips
				err = json.Unmarshal(payload, &relationsShips)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
				So(relationsShips, ShouldResemble, expectedChildRelationShips)
			})
		})
	})
}

func TestGetAreaRelationshipsWithParameterFailsForInvalidIds(t *testing.T) {
	Convey("Given a request for child area relationship data - invalid", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s/relations?relationship=child", "InvalidAreaCode"), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "en")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			ValidateAreaFunc: func(areaCode string) error {
				return apierrors.ErrNoRows
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area data is served", func() {

			Convey("Then an 404 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
				payload, _ := ioutil.ReadAll(w.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(payload, &responseBody)
				So(len(responseBody["errors"].([]interface{})), ShouldEqual, 1)
				error := responseBody["errors"].([]interface{})[0].(map[string]interface{})
				So(error["code"], ShouldEqual, models.InvalidAreaCodeError)
			})
		})
	})
}

func TestGetAreaDataReturnsOKWithOrderedAncestryList(t *testing.T) {
	Convey("Given a successful request to area data - Ancestry hierarchy ordered response", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s", SheffieldAreaData), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "en")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			GetAreaFunc: func(ctx context.Context, areaId string) (*models.AreasDataResults, error) {
				return &models.AreasDataResults{Code: SheffieldAreaData, Name: &SheffieldName, AreaType: &countryAreaType}, nil
			},
			GetAncestorsFunc: func(areaCode string) ([]models.AreasAncestors, error) {
				return ancestors[SheffieldAreaData], nil
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area data is served", func() {

			Convey("Then an OK response is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedArea := models.AreasDataResults{}
				err = json.Unmarshal(payload, &returnedArea)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
				So(returnedArea.Code, ShouldEqual, SheffieldAreaData)
				So(*returnedArea.Name, ShouldEqual, SheffieldName)
				So(*returnedArea.AreaType, ShouldEqual, "Country")
				So(len(returnedArea.Ancestors), ShouldEqual, 2)
				So(returnedArea.Ancestors[0], ShouldResemble, ancestors[SheffieldAreaData][0])
				So(returnedArea.Ancestors[1], ShouldResemble, ancestors[SheffieldAreaData][1])
			})
		})
	})
}

func TestGetAreaDataReturnsOkForEmptyAncestorsList(t *testing.T) {
	Convey("Given a successful request to stubbed area data - empty Ancestors list", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s", WalesAreaData), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "cy")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			GetAreaFunc: func(ctx context.Context, areaId string) (*models.AreasDataResults, error) {
				return &models.AreasDataResults{Code: WalesAreaData, Name: &WalesName, AreaType: &countryAreaType}, nil
			},
			GetAncestorsFunc: func(areaCode string) ([]models.AreasAncestors, error) {
				return ancestors[WalesAreaData], nil
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area data is served", func() {

			Convey("Then an OK response is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedArea := models.AreasDataResults{}
				err = json.Unmarshal(payload, &returnedArea)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
				So(returnedArea.Code, ShouldEqual, WalesAreaData)
				So(*returnedArea.Name, ShouldEqual, "Wales")
				So(*returnedArea.AreaType, ShouldEqual, "Country")
				So(len(returnedArea.Ancestors), ShouldEqual, 0)
			})
		})
	})
}

func TestGetAreaDataFailsForAncestorsDataError(t *testing.T) {
	Convey("Given a successful request to stubbed area data - Ancestors Data Error", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s", SwanseaAirportBuaData), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "cy")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			ValidateAreaFunc: func(areaCode string) error {
				return nil
			},
			GetAncestorsFunc: func(areaCode string) ([]models.AreasAncestors, error) {
				return ancestors[SwanseaAirportBuaData], apierrors.ErrInternalServer
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area data is served", func() {

			Convey("Then an internal server error is returned", func() {
				payload, _ := ioutil.ReadAll(w.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(payload, &responseBody)
				So(len(responseBody["errors"].([]interface{})), ShouldEqual, 1)
				error := responseBody["errors"].([]interface{})[0].(map[string]interface{})
				So(error["code"], ShouldEqual, models.AncestryDataGetError)
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestUpdateAreaData(t *testing.T) {
	Convey("Given a request to update a new area data - W92000004", t, func() {
		reader := strings.NewReader(`{"area_name": {"name": "Wales", "active_from": "2022-01-01T00:00:00Z", "active_to": "2022-02-01T00:00:00Z"}}`)
		r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:2200/v1/areas/%s", "W92000004"), reader)
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{
			UpsertAreaFunc: func(ctx context.Context, area models.AreaParams) (bool, error) { return true, nil },
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area data from rds instance is served", func() {

			Convey("Then an 201 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusCreated)
			})
		})
	})
}

func TestUpdateAreaDataReturnsValidationError(t *testing.T) {
	Convey("Given a request without area details area name details", t, func() {
		reader := strings.NewReader(`{}`)
		r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:2200/v1/areas/%s", WalesAreaData), reader)
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When update area is served", func() {

			Convey("Then validation errors are returned", func() {
				payload, _ := ioutil.ReadAll(w.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(payload, &responseBody)
				So(len(responseBody["errors"].([]interface{})), ShouldEqual, 1)

				error := responseBody["errors"].([]interface{})[0].(map[string]interface{})
				So(error["code"], ShouldEqual, models.AreaNameDetailsNotProvidedError)
				So(error["description"], ShouldEqual, models.AreaNameDetailsNotProvidedErrorDescription)
			})

		})
	})

	Convey("Given invalid area name details", t, func() {
		reader := strings.NewReader(`{"area_name":{}}`)
		r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:2200/v1/areas/%s", WalesAreaData), reader)
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithRDSMocks(&mock.RDSAreaStoreMock{})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When update area is served", func() {

			Convey("Then validation errors are returned", func() {
				payload, _ := ioutil.ReadAll(w.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(payload, &responseBody)
				So(len(responseBody["errors"].([]interface{})), ShouldEqual, 3)

				nameError := responseBody["errors"].([]interface{})[0].(map[string]interface{})
				So(nameError["code"], ShouldEqual, models.AreaNameNotProvidedError)
				So(nameError["description"], ShouldEqual, models.AreaNameNotProvidedErrorDescription)

				validFromError := responseBody["errors"].([]interface{})[1].(map[string]interface{})
				So(validFromError["code"], ShouldEqual, models.AreaNameActiveFromNotProvidedError)
				So(validFromError["description"], ShouldEqual, models.AreaNameActiveFromNotProvidedErrorDescription)

				validToError := responseBody["errors"].([]interface{})[2].(map[string]interface{})
				So(validToError["code"], ShouldEqual, models.AreaNameActiveToNotProvidedError)
				So(validToError["description"], ShouldEqual, models.AreaNameActiveToNotProvidedErrorDescription)
			})

		})
	})
}

func testGeometricData() [][][2]float64 {
	var gd [][][2]float64
	gd = make([][][2]float64, 1)

	for i := range gd {
		gd[i] = make([][2]float64, 1)
		for _ = range gd[i] {
			gd[0][0] = [2]float64{longitude, latitude}
		}
	}
	return gd
}
