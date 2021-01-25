package api_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-areas-api/api"
	"github.com/ONSdigital/dp-areas-api/api/mock"
	"github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/dp-areas-api/models"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetVersionReturnsOk(t *testing.T) {
	Convey("Given a successful request to get version, it returns 200 OK response", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:25500/areas/W06000015/versions/1", nil)
		w := httptest.NewRecorder()
		route := mux.NewRouter()
		mockedAreaStore := &mock.AreaStoreMock{
			CheckAreaExistsFunc: func(id string) error {
				return nil
			},
			GetVersionFunc: func(id string, versionID int) (*models.Area, error) {
				return &models.Area{
					ID:      id,
					Name:    "Cardiff",
					Version: 1,
				}, nil
			},
		}
		var testContext = context.Background()
		apiMock := api.Setup(testContext, route, mockedAreaStore)
		apiMock.Router.ServeHTTP(w, r)

		payload, err := ioutil.ReadAll(w.Body)
		So(err, ShouldBeNil)
		returnedArea := models.Area{}
		err = json.Unmarshal(payload, &returnedArea)
		So(err, ShouldBeNil)
		So(w.Code, ShouldEqual, http.StatusOK)
		So(mockedAreaStore.CheckAreaExistsCalls(), ShouldHaveLength, 1)
		So(mockedAreaStore.GetVersionCalls(), ShouldHaveLength, 1)
		So(returnedArea, ShouldResemble, *dbArea(testAreaId))
	})
}

func TestGetVersionReturnsError(t *testing.T) {
	Convey("Given area does not exist returns area not found", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:25500/areas/W06000015/versions/1", nil)
		w := httptest.NewRecorder()
		route := mux.NewRouter()
		mockedAreaStore := &mock.AreaStoreMock{
			CheckAreaExistsFunc: func(id string) error {
				return apierrors.ErrAreaNotFound
			},
			GetVersionFunc: func(id string, versionID int) (*models.Area, error) {
				return &models.Area{}, nil
			},
		}
		var testContext = context.Background()
		apiMock := api.Setup(testContext, route, mockedAreaStore)
		apiMock.Router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusNotFound)
		So(w.Body.String(), ShouldContainSubstring, apierrors.ErrAreaNotFound.Error())
		So(mockedAreaStore.CheckAreaExistsCalls(), ShouldHaveLength, 1)
		So(mockedAreaStore.GetVersionCalls(), ShouldHaveLength, 0)
	})

	Convey("Given version does not exist for an area returns version not found", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:25500/areas/W06000015/versions/1", nil)
		w := httptest.NewRecorder()
		route := mux.NewRouter()
		mockedAreaStore := &mock.AreaStoreMock{
			CheckAreaExistsFunc: func(id string) error {
				return nil
			},
			GetVersionFunc: func(id string, versionID int) (*models.Area, error) {
				return nil, apierrors.ErrVersionNotFound
			},
		}
		var testContext = context.Background()
		apiMock := api.Setup(testContext, route, mockedAreaStore)
		apiMock.Router.ServeHTTP(w, r)
		So(w.Code, ShouldEqual, http.StatusNotFound)
		So(w.Body.String(), ShouldContainSubstring, apierrors.ErrVersionNotFound.Error())
		So(mockedAreaStore.CheckAreaExistsCalls(), ShouldHaveLength, 1)
		So(mockedAreaStore.GetVersionCalls(), ShouldHaveLength, 1)
	})

	Convey("Given the api cannot connect to areastore return an internal server error", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:25500/areas/W06000015/versions/1", nil)
		w := httptest.NewRecorder()
		route := mux.NewRouter()
		mockedAreaStore := &mock.AreaStoreMock{
			CheckAreaExistsFunc: func(id string) error {
				return apierrors.ErrInternalServer
			},
			GetVersionFunc: func(id string, versionID int) (*models.Area, error) {
				return nil, nil
			},
		}
		var testContext = context.Background()
		apiMock := api.Setup(testContext, route, mockedAreaStore)
		apiMock.Router.ServeHTTP(w, r)
		assertInternalServerErr(w)

		So(mockedAreaStore.CheckAreaExistsCalls(), ShouldHaveLength, 1)
		So(mockedAreaStore.GetVersionCalls(), ShouldHaveLength, 0)
	})
}

func assertInternalServerErr(w *httptest.ResponseRecorder) {
	So(w.Code, ShouldEqual, http.StatusInternalServerError)
	So(strings.TrimSpace(w.Body.String()), ShouldEqual, apierrors.ErrInternalServer.Error())
}
