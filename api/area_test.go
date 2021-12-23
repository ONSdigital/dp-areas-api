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
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testAreaId1   = "W06000015"
	testAreaId2   = "W02000004"
	testAreaName1 = "Cardiff"
	testAreaName2 = "Penylan"
	EnglandAreaData = "E92000001"
	WalesAreaData = "W92000004"
)

var mu sync.Mutex

func dbArea(id string, name string) *models.Area {
	return &models.Area{
		ID:      id,
		Name:    name,
		Version: 1,
	}
}

// GetAPIWithMocks also used in other tests, so exported
func GetAPIWithMocks(mockedAreaStore api.AreaStore) (*api.API, error) {
	mu.Lock()
	defer mu.Unlock()
	cfg, err := config.Get()
	So(err, ShouldBeNil)
	cfg.DefaultLimit = 0
	cfg.DefaultOffset = 0
	return api.Setup(context.Background(), cfg, mux.NewRouter(), mockedAreaStore)
}

func TestGetAreaReturnsOk(t *testing.T) {

	Convey("Given a successful request to get specific area", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25500/areas/%s", testAreaId1), nil)
		w := httptest.NewRecorder()
		//cfg,err:=config.Get()
		//So(err, ShouldBeNil)

		mockedAreaStore := &mock.AreaStoreMock{
			GetAreaFunc: func(ctx context.Context, id string) (*models.Area, error) {
				return &models.Area{ID: "W06000015", Name: "Cardiff", Version: 1}, nil
			},
		}
		Convey("When the request is served", func() {
			areaApi, _ := GetAPIWithMocks(mockedAreaStore)
			areaApi.Router.ServeHTTP(w, r)
			Convey("Then an OK response is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedArea := models.Area{}
				err = json.Unmarshal(payload, &returnedArea)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
				So(returnedArea, ShouldResemble, *dbArea(testAreaId1, testAreaName1))
			})
		})
	})
}

func TestGetAreaReturnsError(t *testing.T) {

	Convey("Given an api request that cannot find the area", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25500/areas/%s", testAreaId1), nil)
		w := httptest.NewRecorder()

		mockedAreaStore := &mock.AreaStoreMock{
			GetAreaFunc: func(ctx context.Context, id string) (*models.Area, error) {
				return nil, apierrors.ErrAreaNotFound
			},
		}
		Convey("When the request is served", func() {
			areaApi, _ := GetAPIWithMocks(mockedAreaStore)
			areaApi.Router.ServeHTTP(w, r)
			Convey("Then status not found,404 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
				So(mockedAreaStore.GetAreaCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given the api cannot connect to datastore", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25500/areas/%s", testAreaId1), nil)
		w := httptest.NewRecorder()
		mockedAreaStore := &mock.AreaStoreMock{
			GetAreaFunc: func(ctx context.Context, id string) (*models.Area, error) {
				return nil, errors.Errorf("Internal server error")
			},
		}
		Convey("When the request is served", func() {
			areaApi, _ := GetAPIWithMocks(mockedAreaStore)
			areaApi.Router.ServeHTTP(w, r)
			Convey("Then an internal server error is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
				So(mockedAreaStore.GetAreaCalls(), ShouldHaveLength, 1)
			})
		})
	})
}

func TestGetAreasReturnsOK(t *testing.T) {
	t.Parallel()
	Convey("Given a successful request to getAreas", t, func() {
		Convey("When the request uses the default offset and limit values", func() {
			r := httptest.NewRequest("GET", "http://localhost:22000/areas", nil)
			w := httptest.NewRecorder()
			mockedAreaStore := &mock.AreaStoreMock{
				GetAreasFunc: func(ctx context.Context, offset int, limit int) (*models.AreasResults, error) {
					return &models.AreasResults{
						Items: &[]models.Area{
							{ID: "W06000015", Name: "Cardiff", Version: 1}, {ID: "W02000004", Name: "Penylan", Version: 1},
						},
						Count:      2,
						Offset:     0,
						Limit:      20,
						TotalCount: 2,
					}, nil
				},
			}

			api, _ := GetAPIWithMocks(mockedAreaStore)
			api.Router.ServeHTTP(w, r)
			Convey("Then a list of areas with an ok response is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedAreas := models.AreasResults{}
				err = json.Unmarshal(payload, &returnedAreas)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(len(mockedAreaStore.GetAreasCalls()), ShouldEqual, 1)
				So(returnedAreas, ShouldResemble, models.AreasResults{
					Items:      &[]models.Area{*dbArea(testAreaId1, testAreaName1), *dbArea(testAreaId2, testAreaName2)},
					Count:      2,
					Offset:     0,
					Limit:      20,
					TotalCount: 2,
				})
			})
		})
	})

	// func to unmarshal and validate body bytes
	validateBody := func(bytes []byte, expected models.AreasResults) {
		var response models.AreasResults
		err := json.Unmarshal(bytes, &response)
		So(err, ShouldBeNil)
		So(response, ShouldResemble, expected)
	}

	Convey("Given a successful request to get areas endpoint", t, func() {
		Convey("When valid limit and offset query parameters are provided", func() {
			r := httptest.NewRequest("GET", "http://localhost:22000/areas?offset=2&limit=2", nil)
			w := httptest.NewRecorder()

			mockedAreaStore := &mock.AreaStoreMock{
				GetAreasFunc: func(ctx context.Context, offset int, limit int) (*models.AreasResults, error) {
					return &models.AreasResults{
						Items:      &[]models.Area{},
						Count:      2,
						Offset:     offset,
						Limit:      limit,
						TotalCount: 5,
					}, nil
				},
			}
			api, _ := GetAPIWithMocks(mockedAreaStore)
			api.Router.ServeHTTP(w, r)

			Convey("Then the call succeeds with 200 OK code with expected body and calls", func() {
				expectedResponse := models.AreasResults{
					Items:      &[]models.Area{},
					Count:      2,
					Offset:     2,
					Limit:      2,
					TotalCount: 5,
				}
				So(w.Code, ShouldEqual, http.StatusOK)
				validateBody(w.Body.Bytes(), expectedResponse)
			})
		})
	})

	Convey("Given a request to the areas endpoint with non existing areas", t, func() {
		Convey("When valid limit  and offset query parameters are provided", func() {

			r := httptest.NewRequest("GET", "http://localhost:22000/areas?offset=2&limit=7", nil)
			w := httptest.NewRecorder()
			mockedAreaStore := &mock.AreaStoreMock{
				GetAreasFunc: func(ctx context.Context, offset int, limit int) (*models.AreasResults, error) {
					return &models.AreasResults{
						Items:      &[]models.Area{},
						Count:      0,
						Offset:     offset,
						Limit:      limit,
						TotalCount: 0,
					}, nil
				},
			}

			api, _ := GetAPIWithMocks(mockedAreaStore)
			api.Router.ServeHTTP(w, r)

			Convey("Then the call succeeds with 200 OK code, expected body and calls", func() {
				expectedResponse := models.AreasResults{
					Items:      &[]models.Area{},
					Count:      0,
					Offset:     2,
					Limit:      7,
					TotalCount: 0,
				}
				So(w.Code, ShouldEqual, http.StatusOK)
				validateBody(w.Body.Bytes(), expectedResponse)
			})
		})
	})

	Convey("Given a request to the areas api", t, func() {
		Convey("When a given offset is greater than the total count", func() {
			r := httptest.NewRequest("GET", "http://localhost:22000/areas?offset=4&limit=3", nil)
			w := httptest.NewRecorder()
			mockedAreaStore := &mock.AreaStoreMock{
				GetAreasFunc: func(ctx context.Context, offset, limit int) (*models.AreasResults, error) {
					return &models.AreasResults{
						Items:      &[]models.Area{},
						Count:      0,
						Offset:     offset,
						Limit:      limit,
						TotalCount: 3,
					}, nil
				},
			}
			api, _ := GetAPIWithMocks(mockedAreaStore)
			api.Router.ServeHTTP(w, r)
			Convey("Then an empty array returned", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
}

func TestGetAreasReturnsError(t *testing.T) {
	t.Parallel()
	Convey("Given a request to the areas api", t, func() {
		Convey("When the api cannot connect to datastore return an internal server error", func() {
			r := httptest.NewRequest("GET", "http://localhost:22000/areas", nil)
			w := httptest.NewRecorder()
			mockedAreaStore := &mock.AreaStoreMock{
				GetAreasFunc: func(ctx context.Context, offset, limit int) (*models.AreasResults, error) {
					return nil, apierrors.ErrInternalServer
				},
			}

			apiMock, _ := GetAPIWithMocks(mockedAreaStore)
			Convey("Then an internal server error is returned", func() {
				apiMock.Router.ServeHTTP(w, r)
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
				So(len(mockedAreaStore.GetAreasCalls()), ShouldEqual, 1)
			})
		})
	})

	Convey("Given a request to the areas api", t, func() {
		Convey("When a negative limit and offset query parameters are provided, then return a 400 error", func() {

			r := httptest.NewRequest("GET", "http://localhost:22000/areas?offset=-2&limit=-7", nil)
			w := httptest.NewRecorder()
			mockedAreaStore := &mock.AreaStoreMock{
				GetAreasFunc: func(ctx context.Context, offset, limit int) (*models.AreasResults, error) {
					return nil, apierrors.ErrInvalidQueryParameter
				},
			}
			api, _ := GetAPIWithMocks(mockedAreaStore)
			api.Router.ServeHTTP(w, r)
			Convey("Then a status code of 400, StatusBadRequest is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
				So(strings.TrimSpace(w.Body.String()), ShouldEqual, apierrors.ErrInvalidQueryParameter.Error())

			})
		})
	})

	Convey("Given a request to the areas api", t, func() {
		Convey("When a limit higher than the maximum default limit value is provided ", func() {

			r := httptest.NewRequest("GET", "http://localhost:22000/areas?offset=2&limit=1500", nil)
			w := httptest.NewRecorder()
			mockedAreaStore := &mock.AreaStoreMock{
				GetAreasFunc: func(ctx context.Context, offset, limit int) (*models.AreasResults, error) {
					return nil, apierrors.ErrQueryParamLimitExceedMax
				},
			}
			api, _ := GetAPIWithMocks(mockedAreaStore)
			api.Router.ServeHTTP(w, r)
			Convey("Then a status code of 400, StatusBadRequest is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
				So(strings.TrimSpace(w.Body.String()), ShouldEqual, apierrors.ErrQueryParamLimitExceedMax.Error())
			})
		})
	})

	Convey("Given a request to the areas api", t, func() {
		Convey("When an invalid user defined offset value is provided", func() {

			r := httptest.NewRequest("GET", "http://localhost:22000/areas?offset=h&limit=1", nil)
			w := httptest.NewRecorder()
			mockedAreaStore := &mock.AreaStoreMock{
				GetAreasFunc: func(ctx context.Context, offset, limit int) (*models.AreasResults, error) {
					return nil, apierrors.ErrInvalidQueryParameter
				},
			}
			api, _ := GetAPIWithMocks(mockedAreaStore)
			api.Router.ServeHTTP(w, r)
			Convey("Then a status code of 400, StatusBadRequest is returned", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
				So(strings.TrimSpace(w.Body.String()), ShouldEqual, apierrors.ErrInvalidQueryParameter.Error())
			})
		})
	})
}

func TestGetAreaDataReturnsOkForEngland(t *testing.T) {
	Convey("Given a successful request to stubbed area data - E92000001", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s", EnglandAreaData), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "en")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithMocks(&mock.AreaStoreMock{
			GetAreasFunc: func(ctx context.Context, offset, limit int) (*models.AreasResults, error) {
				return nil, apierrors.ErrInvalidQueryParameter
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
				So(returnedArea.WelshName, ShouldEqual, "Lloegr")
				So(returnedArea.Name, ShouldEqual, "England")
				So(returnedArea.AreaType, ShouldEqual, "English")
			})
		})
	})
}

func TestGetAreaDataReturnsOkForWales(t *testing.T) {
	Convey("Given a successful request to stubbed area data - W92000004", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s", WalesAreaData), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "cy")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithMocks(&mock.AreaStoreMock{
			GetAreasFunc: func(ctx context.Context, offset, limit int) (*models.AreasResults, error) {
				return nil, apierrors.ErrInvalidQueryParameter
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
				So(returnedArea.WelshName, ShouldEqual, "Cymru")
				So(returnedArea.Name, ShouldEqual, "Wales")
				So(returnedArea.AreaType, ShouldEqual, "Cymraeg")
			})
		})
	})
}

func TestGetAreaDataReturnsValidationErrors(t *testing.T) {
	Convey("Given a request to stubbed area data - validation errors returned", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s", WalesAreaData+"z"), nil)
		r.Header.Set(models.AcceptLanguageHeaderName, "cz")
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithMocks(&mock.AreaStoreMock{
			GetAreasFunc: func(ctx context.Context, offset, limit int) (*models.AreasResults, error) {
				return nil, apierrors.ErrInvalidQueryParameter
			},
		})
		areaApi.Router.ServeHTTP(w, r)

		Convey("When request area data is served", func() {
			
			Convey("Then validation failure details returned - 2 validation errors generated", func() {
				payload, _ := ioutil.ReadAll(w.Body)
				var responseBody map[string]interface{}
				_ = json.Unmarshal(payload, &responseBody)
				So(len(responseBody["errors"].([]interface{})), ShouldEqual, 2)

				// error #1 - accept language header problem
				error1 := responseBody["errors"].([]interface{})[0].(map[string]interface{})
				So(error1["code"], ShouldEqual, models.AcceptLanguageHeaderError)
				So(error1["description"], ShouldEqual, models.AcceptLanguageHeaderInvalidDescription)

				// error #2 - are code not found
				error2 := responseBody["errors"].([]interface{})[1].(map[string]interface{})
				So(error2["code"], ShouldEqual, models.AreaDataIdGetError)
				So(error2["description"], ShouldEqual, models.AreaDataGetErrorDescription)
			})
			
		})
	})
}

func TestGetAreaDataReturnsValidationError(t *testing.T) {
	Convey("Given a request to stubbed area data - validation errors returned", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:2200/v1/areas/%s", WalesAreaData), nil)
		w := httptest.NewRecorder()

		areaApi, _ := GetAPIWithMocks(&mock.AreaStoreMock{
			GetAreasFunc: func(ctx context.Context, offset, limit int) (*models.AreasResults, error) {
				return nil, apierrors.ErrInvalidQueryParameter
			},
		})
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
