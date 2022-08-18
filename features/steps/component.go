package steps

import (
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

	svc *service.Service
}

func NewComponent(t *testing.T) *Component {
	component := &Component{
		ErrorFeature:        componenttest.ErrorFeature{TB: t},
		AuthServiceInjector: componenttest.NewAuthorizationFeature(),
		RDSInjector:         &RDSFeature{},
	}

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
	ctx.Step(`^I GET "([^"]*)"$`, c.iGET)
	ctx.Step(`^I should receive the following JSON response:$`, c.iShouldReceiveTheFollowingJSONResponse)
	ctx.Step(`^the HTTP status code should be "([^"]*)"$`, c.theHTTPStatusCodeShouldBe)
	ctx.Step(`^the response header "([^"]*)" should be "([^"]*)"$`, c.theResponseHeaderShouldBe)
}

func (c *Component) iGET(arg1 string) error {
	return godog.ErrPending
}

func (c *Component) iShouldReceiveTheFollowingJSONResponse(arg1 *godog.DocString) error {
	return godog.ErrPending
}

func (c *Component) theHTTPStatusCodeShouldBe(arg1 string) error {
	return godog.ErrPending
}

func (c *Component) theResponseHeaderShouldBe(arg1, arg2 string) error {
	return godog.ErrPending
}
