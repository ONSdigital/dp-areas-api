package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	brokenrecorder "github.com/ONSdigital/dp-topic-api/broken"
	"github.com/ONSdigital/log.go/v2/log"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWriteJSONBody(t *testing.T) {

	Convey("Exercise WriteJSONBody", t, func() {

		Convey("Given a GOOD recorder to write to, with valid data", func() {

			ctx := context.Background()
			logdata := log.Data{
				"content_id": "ok",
			}
			v := "ok value"
			w := httptest.NewRecorder()

			err := WriteJSONBody(ctx, v, w, logdata)
			Convey("Then the expected write success is seen", func() {
				So(err, ShouldBeNil)
				So(w.Code, ShouldEqual, http.StatusOK)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				var retType string
				err = json.Unmarshal(payload, &retType)
				So(err, ShouldBeNil)
				So(retType, ShouldResemble, v)
			})
		})

		Convey("Given a GOOD recorder to write to, with invalid data", func() {

			ctx := context.Background()
			logdata := log.Data{
				"content_id": "broken JSON",
			}
			v := math.Inf(1)
			w := httptest.NewRecorder()

			err := WriteJSONBody(ctx, v, w, logdata)
			Convey("Then the expected JSON error is seen", func() {
				res := fmt.Sprintf("%+v\n", err)
				So(res, ShouldResemble, "json: unsupported value: +Inf\n")
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				So(string(payload), ShouldResemble, "internal error\n")
			})
		})

		Convey("Given a Broken Recorder to write to, with valid data", func() {

			ctx := context.Background()
			logdata := log.Data{
				"content_id": "broken write",
			}
			v := "fail"

			w := brokenrecorder.NewRecorder()

			err := WriteJSONBody(ctx, v, w, logdata)
			Convey("Then the expected write failure is seen", func() {
				errorResult := errors.New("broken write")
				So(err, ShouldResemble, errorResult)
				payload, err := ioutil.ReadAll(w.Body)
				So(err, ShouldBeNil)
				So(string(payload), ShouldContainSubstring, "This is the broken write")
			})
		})
	})
}

//
//func TestNavigationCacheMaxAge(t *testing.T) {
//	Convey("Excercise setUp", t, func() {
//		ctx := context.Background()
//		cfg, err := config.Get()
//		So(err, ShouldBeNil)
//		mockedDataStore := &storetest.StorerMock{}
//		permissions := mocks.NewAuthHandlerMock()
//		api := Setup(ctx, cfg, mux.NewRouter(), store.DataStore{Backend: mockedDataStore}, permissions)
//
//		So(api.navigationCacheMaxAge, ShouldResemble, fmt.Sprintf("%f", cfg.NavigationCacheMaxAge.Seconds()))
//
//		w := httptest.NewRecorder()
//
//		err := api.ReadJSONBody(ctx, b, v, w, logdata)
//
//	})
