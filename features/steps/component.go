package steps

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/service"
	mocks "github.com/ONSdigital/dp-areas-api/service/mock"
	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

var (
	BuildTime = strconv.Itoa(time.Now().Nanosecond())
	GitCommit = "component test commit"
	Version   = "component test version"
)

type Component struct {
	componenttest.ErrorFeature
	AuthFeature *componenttest.AuthorizationFeature
	APIFeature  *componenttest.APIFeature
	RDSFeature  *RDSFeature

	svc *service.Service
}

func NewComponent(t *testing.T) *Component {
	cfg, err := config.Get()
	if err != nil {
		return nil
	}

	component := &Component{
		ErrorFeature: componenttest.ErrorFeature{TB: t},
		AuthFeature:  componenttest.NewAuthorizationFeature(),
		RDSFeature:   NewRDSFeature(t, cfg),
	}
	component.APIFeature = componenttest.NewAPIFeature(component.ServiceAPIRouter)

	standardInit := service.Init{}
	initFunctions := &mocks.InitialiserMock{
		DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer {
			return &http.Server{Addr: bindAddr, Handler: router}
		},
		DoGetHealthCheckFunc: func(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
			return &mocks.HealthCheckerMock{
				AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
				StartFunc:    func(ctx context.Context) {},
				StopFunc:     func() {},
			}, nil
		},
		DoGetRDSDBFunc: standardInit.DoGetRDSDB,
	}
	//component.RDSFeature.esServer.NewHandler().Get("/elasticsearch/_cluster/health").Reply(200).Body([]byte(""))

	serviceList := service.NewServiceList(initFunctions)
	component.svc, err = service.Run(context.Background(), cfg, serviceList, BuildTime, GitCommit, Version, make(chan error, 1))
	if err != nil {
		t.Fatalf("service failed to run: %s", err)
	}

	return component
}

// Reset re-initialises the service under test and the api mocks.
func (c *Component) Reset() {
	c.AuthFeature.Reset()
	c.APIFeature.Reset()
	c.RDSFeature.Reset()
}

func (c *Component) Close() {
	c.AuthFeature.Close()
	c.RDSFeature.Close()
	if c.svc != nil {
		_ = c.svc.Close(context.Background())
		c.svc = nil
	}
}

func (c *Component) ServiceAPIRouter() (http.Handler, error) {
	return c.svc.API.Router, nil
}
