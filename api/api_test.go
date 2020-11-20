package api_test

import (
	"context"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ONSdigital/dp-topic-api/api"
	"github.com/ONSdigital/dp-topic-api/api/mock"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	mu          sync.Mutex
	testContext = context.Background()
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := &mux.Router{}
		ctx := context.Background()

		Convey("When created in Publishing mode", func() {
			cfg := &config.Config{EnablePrivateEndpoints: true}
			api := api.Setup(ctx, cfg, r, &mock.MongoServerMock{})

			Convey("When created the following routes should have been added", func() {
				So(api, ShouldNotBeNil)
				So(hasRoute(api.Router, "/hello", "GET"), ShouldBeTrue) //!!! remove at some point
				So(hasRoute(api.Router, "/topics/{id}", "GET"), ShouldBeTrue)
			})
		})

		Convey("When created in Web mode", func() {
			cfg := &config.Config{EnablePrivateEndpoints: false}
			api := api.Setup(ctx, cfg, r, &mock.MongoServerMock{})

			Convey("When created the following routes should have been added", func() {
				So(api, ShouldNotBeNil)
				So(hasRoute(api.Router, "/hello", "GET"), ShouldBeTrue) // !!! remove at some point
				So(hasRoute(api.Router, "/topics/{id}", "GET"), ShouldBeTrue)
			})
		})

	})
}

// GetAPIWithMocks also used in other tests
func GetAPIWithMocks(cfg *config.Config, mongoDbMock *mock.MongoServerMock /*, authHandlerMock *mock.AuthHandlerMock*/) *api.API {
	mu.Lock()
	defer mu.Unlock()
	//	urlBuilder := url.NewBuilder("http://example.com")
	return api.Setup(testContext, cfg, mux.NewRouter() /*authHandlerMock,*/, mongoDbMock /*, urlBuilder*/)
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
