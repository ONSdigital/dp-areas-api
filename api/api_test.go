package api_test

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/jackc/pgx/v4/pgxpool"

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
		api, _ := api.Setup(ctx, cfg, r, mongoMock, &pgxpool.Pool{})
		Convey("When created the following routes should have been added", func() {
			So(hasRoute(api.Router, "/areas", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/areas/{id}", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/areas/{id}/versions/{version}", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/areas/{id}", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/areas/{id}/relations", "GET"), ShouldBeTrue)
			So(hasRoute(api.Router, "/v1/rds/areas/{id}", "GET"), ShouldBeTrue)
		})
	})
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
