package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/mocks"
	"github.com/ONSdigital/dp-topic-api/store"
	storetest "github.com/ONSdigital/dp-topic-api/store/datastoretest"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	mu          sync.Mutex
	testContext = context.Background()
)

func TestPublishedSubnetEndpointsAreDisabled(t *testing.T) {

	type testEndpoint struct {
		Method string
		URL    string
	}

	publishSubnetEndpoints := map[testEndpoint]int{
		// Dataset Endpoints
		/*{Method: "POST", URL: "http://localhost:22000/datasets/1234"}:                            http.StatusMethodNotAllowed,
		{Method: "PUT", URL: "http://localhost:22000/datasets/1234"}:                             http.StatusMethodNotAllowed,
		{Method: "PUT", URL: "http://localhost:22000/datasets/1234/editions/1234/versions/2123"}: http.StatusMethodNotAllowed,*/

		// Instance endpoints
		{Method: "GET", URL: "http://localhost:22000/instances"}: http.StatusNotFound,
		/*{Method: "POST", URL: "http://localhost:22000/instances"}:                           http.StatusNotFound,
		{Method: "GET", URL: "http://localhost:22000/instances/1234"}:                       http.StatusNotFound,
		{Method: "PUT", URL: "http://localhost:22000/instances/123"}:                        http.StatusNotFound,
		{Method: "PUT", URL: "http://localhost:22000/instances/123/dimensions/test"}:        http.StatusNotFound,
		{Method: "POST", URL: "http://localhost:22000/instances/1/events"}:                  http.StatusNotFound,
		{Method: "PUT", URL: "http://localhost:22000/instances/1/inserted_observations/11"}: http.StatusNotFound,
		{Method: "PUT", URL: "http://localhost:22000/instances/1/import_tasks"}:             http.StatusNotFound,*/

		// Dimension endpoints
		/*{Method: "GET", URL: "http://localhost:22000/instances/1/dimensions"}:                       http.StatusNotFound,
		{Method: "POST", URL: "http://localhost:22000/instances/1/dimensions"}:                      http.StatusNotFound,
		{Method: "GET", URL: "http://localhost:22000/instances/1/dimensions/1/options"}:             http.StatusNotFound,
		{Method: "PUT", URL: "http://localhost:22000/instances/1/dimensions/1/options/1/node_id/1"}: http.StatusNotFound,*/
	}

	Convey("When the API is started with private endpoints disabled", t, func() {

		for endpoint, expectedStatusCode := range publishSubnetEndpoints {
			Convey("The following endpoint "+endpoint.URL+"(Method:"+endpoint.Method+") should return 404", func() {
				request, err := createRequestWithAuth(endpoint.Method, endpoint.URL, nil)
				So(err, ShouldBeNil)

				cfg, err := config.Get()
				So(err, ShouldBeNil)
				cfg.EnablePrivateEndpoints = false

				w := httptest.NewRecorder()
				mockedDataStore := &storetest.StorerMock{}
				api := GetWebAPIWithMocks(testContext, cfg, mockedDataStore, nil, nil)

				api.Router.ServeHTTP(w, request)

				So(w.Code, ShouldEqual, expectedStatusCode)
			})
		}
	})
}

func TestSetup(t *testing.T) {
	Convey("Given an API insetAuttance", t, func() {
		topicPermissions := mocks.NewAuthHandlerMock()
		permissions := mocks.NewAuthHandlerMock()

		Convey("When created in Publishing mode", func() {
			cfg, err := config.Get()
			So(err, ShouldBeNil)
			cfg.EnablePrivateEndpoints = true
			mockedDataStore := &storetest.StorerMock{}

			api := GetWebAPIWithMocks(testContext, cfg, mockedDataStore, topicPermissions, permissions)

			Convey("When created the following routes should have been added", func() {
				So(api, ShouldNotBeNil)
				So(hasRoute(api.Router, "/hello", "GET"), ShouldBeTrue) //!!! remove at some point
				So(hasRoute(api.Router, "/topics/{id}", "GET"), ShouldBeTrue)
			})
		})

		Convey("When created in Web mode", func() {
			cfg, err := config.Get()
			So(err, ShouldBeNil)
			cfg.EnablePrivateEndpoints = false
			mockedDataStore := &storetest.StorerMock{}

			api := GetWebAPIWithMocks(testContext, cfg, mockedDataStore, topicPermissions, permissions)

			Convey("When created the following routes should have been added", func() {
				So(api, ShouldNotBeNil)
				So(hasRoute(api.Router, "/hello", "GET"), ShouldBeTrue) // !!! remove at some point
				So(hasRoute(api.Router, "/topics/{id}", "GET"), ShouldBeTrue)
			})
		})

	})
}

// GetWebAPIWithMocks also used in other tests, so exported
func GetWebAPIWithMocks(ctx context.Context, cfg *config.Configuration, mockedDataStore store.Storer, topicPermissions AuthHandler, permissions AuthHandler) *API {
	mu.Lock()
	defer mu.Unlock()
	//	urlBuilder := url.NewBuilder("http://example.com")

	//	cfg.ServiceAuthToken = authToken
	//	cfg.DatasetAPIURL = host

	return Setup(ctx, cfg, mux.NewRouter(), store.DataStore{Backend: mockedDataStore}, topicPermissions, permissions)
}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}

func createRequestWithAuth(method, URL string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, URL, body)
	ctx := request.Context()
	ctx = dprequest.SetCaller(ctx, "someone@ons.gov.uk")
	request = request.WithContext(ctx)
	return request, err
}
