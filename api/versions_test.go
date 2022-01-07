package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-areas-api/api/mock"
	"github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/dp-areas-api/models"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetVersionReturnsOk(t *testing.T) {
	Convey("Given a successful request to get version", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:25500/areas/W06000015/versions/1", nil)
		w := httptest.NewRecorder()
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
		Convey("When the request is served", func() {
			apiMock, _ := GetAPIWithMocks(mockedAreaStore)
			apiMock.Router.ServeHTTP(w, r)
			Convey("Then an OK response 200 is returned", func() {
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				returnedArea := models.Area{}
				err = json.Unmarshal(payload, &returnedArea)
				So(err, ShouldBeNil)
				So(w.Code, ShouldEqual, http.StatusOK)
				So(mockedAreaStore.CheckAreaExistsCalls(), ShouldHaveLength, 1)
				So(mockedAreaStore.GetVersionCalls(), ShouldHaveLength, 1)
				So(returnedArea, ShouldResemble, *dbArea(testAreaId1,testAreaName1))
			})
		})
	})
}

func TestGetVersionReturnsError(t *testing.T) {
	Convey("Given a request for a version for an area that does not exist ", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:25500/areas/W0600001/versions/1", nil)
		w := httptest.NewRecorder()
		mockedAreaStore := &mock.AreaStoreMock{
			CheckAreaExistsFunc: func(id string) error {
				return apierrors.ErrAreaNotFound
			},
			GetVersionFunc: func(id string, versionID int) (*models.Area, error) {
				return nil, nil
			},
		}
		Convey("When the request is served", func() {
			apiMock, _ := GetAPIWithMocks(mockedAreaStore)
			apiMock.Router.ServeHTTP(w, r)
			Convey("Then an error area not found is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
				So(w.Body.String(), ShouldContainSubstring, apierrors.ErrAreaNotFound.Error())
				So(mockedAreaStore.CheckAreaExistsCalls(), ShouldHaveLength, 1)
				So(mockedAreaStore.GetVersionCalls(), ShouldHaveLength, 0)
			})
		})
	})

	Convey("Given aq request for a version that does not exist", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:25500/areas/W06000015/versions/2", nil)
		w := httptest.NewRecorder()
		mockedAreaStore := &mock.AreaStoreMock{
			CheckAreaExistsFunc: func(id string) error {
				return nil
			},
			GetVersionFunc: func(id string, versionID int) (*models.Area, error) {
				return nil, apierrors.ErrVersionNotFound
			},
		}
		Convey("When the request is served", func() {
			apiMock, _ := GetAPIWithMocks(mockedAreaStore)
			apiMock.Router.ServeHTTP(w, r)
			Convey("Then an error version not found is returned", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
				So(w.Body.String(), ShouldContainSubstring, apierrors.ErrVersionNotFound.Error())
				So(mockedAreaStore.CheckAreaExistsCalls(), ShouldHaveLength, 1)
				So(mockedAreaStore.GetVersionCalls(), ShouldHaveLength, 1)
			})
		})
	})

	Convey("Given the api cannot connect to areastore return an internal server error", t, func() {
		r := httptest.NewRequest("GET", "http://localhost:25500/areas/W06000015/versions/1", nil)
		w := httptest.NewRecorder()
		mockedAreaStore := &mock.AreaStoreMock{
			CheckAreaExistsFunc: func(id string) error {
				return apierrors.ErrInternalServer
			},
			GetVersionFunc: func(id string, versionID int) (*models.Area, error) {
				return nil, nil
			},
		}
		Convey("When the request is served", func() {
			apiMock, _ := GetAPIWithMocks(mockedAreaStore)
			apiMock.Router.ServeHTTP(w, r)
			Convey("Then an internal server error is returned", func() {
				assertInternalServerErr(w)

				So(mockedAreaStore.CheckAreaExistsCalls(), ShouldHaveLength, 1)
				So(mockedAreaStore.GetVersionCalls(), ShouldHaveLength, 0)
			})
		})
	})

}

func assertInternalServerErr(w *httptest.ResponseRecorder) {
	So(w.Code, ShouldEqual, http.StatusInternalServerError)
	So(strings.TrimSpace(w.Body.String()), ShouldEqual, apierrors.ErrInternalServer.Error())
}
