package steps

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-areas-api/api"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/ONSdigital/dp-areas-api/service"
	componenttest "github.com/ONSdigital/dp-component-test"

	"github.com/cucumber/godog"
)

var (
	BuildTime = strconv.Itoa(time.Now().Nanosecond())
	GitCommit = "component test commit"
	Version   = "component test version"
)

type Component struct {
	componenttest.ErrorFeature
	AuthServiceInjector *componenttest.AuthorizationFeature
	APIInjector         *componenttest.APIFeature
	RDSInjector         *RDSFeature

	svc      *service.Service
	api      *api.API
	response *httptest.ResponseRecorder
	payload  []byte
}

func NewComponent(t *testing.T) *Component {
	cfg, err := config.Get()
	if err != nil {
		return nil
	}

	component := &Component{
		ErrorFeature:        componenttest.ErrorFeature{TB: t},
		AuthServiceInjector: componenttest.NewAuthorizationFeature(),
		RDSInjector:         NewRDSFeature(componenttest.ErrorFeature{TB: t}, cfg),
	}

	component.RDSInjector = NewRDSFeature(component.ErrorFeature, cfg)
	return component
}

// Reset re-initialises the service under test and the api mocks.
func (c *Component) Reset() {
	c.AuthServiceInjector.Reset()
	c.APIInjector.Reset()
	c.RDSInjector.Reset()
}

func (c *Component) Close() {
	c.AuthServiceInjector.Close()
	c.RDSInjector.Close()
}

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I GET "([^"]*)"$`, c.iGETFor)
	ctx.Step(`^I should receive the following JSON response:$`, c.iShouldReceiveTheFollowingJSONResponse)
	ctx.Step(`^the HTTP status code should be "([^"]*)"$`, c.theHTTPStatusCodeShouldBe)
	ctx.Step(`^the response header "([^"]*)" should be "([^"]*)"$`, c.theResponseHeaderShouldBe)
}

func (c *Component) iGETFor(arg1, arg2 string) error {
	// Setup
	r := mux.NewRouter()
	ctx := context.Background()
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	c.api, err = api.Setup(ctx, cfg, r, &c.RDSInjector.Client)
	if err != nil {
		return err
	}

	req := httptest.NewRequest(http.MethodGet, arg1, nil)
	req.Header.Set(models.AcceptLanguageHeaderName, "en")
	c.response = httptest.NewRecorder()

	c.api.Router.ServeHTTP(c.response, req)
	c.payload, err = io.ReadAll(c.response.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Component) iShouldReceiveTheFollowingJSONResponse(arg1 *godog.DocString) error {
	//returnedBoundary := models.BoundaryDataResults{}
	//err := json.Unmarshal(c.payload, &returnedBoundary)
	//if err != nil {
	//	return err
	//}

	if string(c.payload) != arg1.Content {
		return fmt.Errorf("payload does not match expected response: payload[%s] | expected[%s]",
			string(c.payload), arg1.Content)
	}

	return nil
}

func (c *Component) theHTTPStatusCodeShouldBe(arg1 string) error {
	expectedStatusCode, err := strconv.Atoi(arg1)
	if err != nil {
		return err
	}

	if c.response.Code != expectedStatusCode {
		return fmt.Errorf("expected status code [%v] but received [%v]",
			expectedStatusCode, c.response.Code)
	}

	return nil
}

func (c *Component) theResponseHeaderShouldBe(arg1, arg2 string) error {
	h := c.response.Header().Get(arg1)
	if h != arg2 {
		return fmt.Errorf("expected response header key [%s] value to be [%s] but got [%s]",
			arg1, arg2, h)
	}

	return nil
}
