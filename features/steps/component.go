package steps

import (
	"strconv"
	"testing"
	"time"

	"github.com/ONSdigital/dp-areas-api/service"
	componenttest "github.com/ONSdigital/dp-component-test"
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
