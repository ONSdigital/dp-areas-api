package api_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-topic-api/api"
	"github.com/ONSdigital/dp-topic-api/api/mock"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := &mux.Router{}
		ctx := context.Background()
		cfg := &config.Config{ /*IsPublishing: true*/ } //!!! fix for bug & web differences
		api := api.Setup(ctx, cfg, r, &mock.MongoServerMock{})

		//!!! add code to test new endpoint for publishing and web, etc
		Convey("When created the following routes should have been added", func() {
			So(api, ShouldNotBeNil)
			// Replace the check below with any newly added api endpoints
			So(hasRoute(api.Router, "/hello", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
