package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dp-areas-api/api"
	"github.com/ONSdigital/dp-areas-api/api/mock"
	"github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testAreaId = "W06000015"
)

func dbArea(id string) *models.Area {
	return &models.Area{
		ID:      id,
		Name:    "Cardiff",
		Version: 1,
	}
}

func TestGetAreaReturnsOk(t *testing.T) {

	Convey("Given a successful request to get specific area", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25500/areas/%s", testAreaId), nil)
		w := httptest.NewRecorder()
		mockedAreaStore := &mock.AreaStoreMock{
			GetAreaFunc: func(ctx context.Context, id string) (*models.Area, error) {
				return &models.Area{ID: "W06000015", Name: "Cardiff", Version: 1}, nil
			},
		}
		Convey("When the request is served", func() {
			areaApi := api.Setup(context.Background(), mux.NewRouter(), mockedAreaStore)
			areaApi.Router.ServeHTTP(w, r)
			Convey("Then an OK response is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedArea := models.Area{}
				err = json.Unmarshal(payload, &returnedArea)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(err, ShouldBeNil)
				So(returnedArea, ShouldResemble, *dbArea(testAreaId))
			})
		})
	})
}

func TestGetAreaReturnsError(t *testing.T) {

	Convey("Given an api request that cannot find the area", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25500/areas/%s", testAreaId), nil)
		w := httptest.NewRecorder()
		mockedAreaStore := &mock.AreaStoreMock{
			GetAreaFunc: func(ctx context.Context, id string) (*models.Area, error) {
				return nil, apierrors.ErrAreaNotFound
			},
		}
		Convey("When the request is served", func() {
			areaApi := api.Setup(context.Background(), mux.NewRouter(), mockedAreaStore)
			areaApi.Router.ServeHTTP(w, r)
			Convey("Then status not found,404 response is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
				So(mockedAreaStore.GetAreaCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given the api cannot connect to datastore", t, func() {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:25500/areas/%s", testAreaId), nil)
		w := httptest.NewRecorder()
		mockedAreaStore := &mock.AreaStoreMock{
			GetAreaFunc: func(ctx context.Context, id string) (*models.Area, error) {
				return nil, errors.Errorf("Internal server error")
			},
		}
		Convey("When the request is served", func() {
			areaApi := api.Setup(context.Background(), mux.NewRouter(), mockedAreaStore)
			areaApi.Router.ServeHTTP(w, r)
			Convey("Then an internal server error is returned", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
				So(mockedAreaStore.GetAreaCalls(), ShouldHaveLength, 1)
			})
		})
	})
}
