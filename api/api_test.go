package api_test

import (
	"context"
	"github.com/ONSdigital/dp-areas-api/config"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ONSdigital/dp-areas-api/api"
	"github.com/ONSdigital/dp-areas-api/api/mock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		os.Clearenv()
		mongoMock := &mock.AreaStoreMock{}
		r := mux.NewRouter()
		ctx := context.Background()
		cfg, err := config.Get()
		So(err, ShouldBeNil)
		api := api.Setup(ctx, cfg, r, mongoMock)
		Convey("When created the following routes should have been added", func() {
			So(hasRoute(api.Router, "/areas", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/areas/{id}", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/areas/{id}/versions/{version}", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
