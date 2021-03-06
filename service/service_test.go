package service_test

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-areas-api/api"
	"net/http"
	"sync"
	"testing"
	"time"

	apiMock "github.com/ONSdigital/dp-areas-api/api/mock"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/service"
	"github.com/ONSdigital/dp-areas-api/service/mock"
	serviceMock "github.com/ONSdigital/dp-areas-api/service/mock"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx           = context.Background()
	testBuildTime = "BuildTime"
	testGitCommit = "GitCommit"
	testVersion   = "Version"
	errServer     = errors.New("HTTP Server error")
)

var (
	errHealthcheck = errors.New("healthCheck error")
	errMongo       = errors.New("mongoDB error")
)

var funcDoGetHealthcheckErr = func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return nil, errHealthcheck
}

var funcDoGetHTTPServerNil = func(bindAddr string, router http.Handler) service.HTTPServer {
	return nil
}

var funcDoGetMongoDBErr = func(ctx context.Context, cfg *config.Config) (api.AreaStore, error) {
	return nil, errMongo
}

func TestRun(t *testing.T) {

	Convey("Having a set of mocked dependencies", t, func() {

		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
		}

		serverWg := &sync.WaitGroup{}
		serverMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return nil
			},
		}

		mongoMock := &apiMock.AreaStoreMock{
			CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
			CloseFunc: func(ctx context.Context) error {
				return nil
			},
		}

		failingServerMock := &serviceMock.HTTPServerMock{
			ListenAndServeFunc: func() error {
				serverWg.Done()
				return errServer
			},
		}

		funcDoGetHealthcheckOk := func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
			return hcMock, nil
		}

		funcDoGetHTTPServer := func(bindAddr string, router http.Handler) service.HTTPServer {
			return serverMock
		}

		funcDoGetFailingHTTPServer := func(bindAddr string, router http.Handler) service.HTTPServer {
			return failingServerMock
		}

		funcDoGetMongoDBOk := func(ctx context.Context, cfg *config.Config) (api.AreaStore, error) {
			return mongoMock, nil
		}

		Convey("Given that initialising mongoDB returns an error", func() {
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:  funcDoGetHTTPServer,
				DoGetMongoDBFunc:     funcDoGetMongoDBErr,
				DoGetHealthCheckFunc: funcDoGetHealthcheckOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the same error and the flag is not set", func() {
				So(err, ShouldResemble, errMongo)
				So(svcList.MongoDB, ShouldBeFalse)
				So(svcList.HealthCheck, ShouldBeFalse)
			})
		})

		Convey("Given that initialising healthcheck returns an error", func() {

			// setup (run before each `Convey` at this scope / indentation):
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:  funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc: funcDoGetHealthcheckErr,
				DoGetMongoDBFunc:     funcDoGetMongoDBOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails with the same error and the flag is not set", func() {
				So(err, ShouldResemble, errHealthcheck)
				So(svcList.HealthCheck, ShouldBeFalse)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("Given that all dependencies are successfully initialised", func() {

			// setup (run before each `Convey` at this scope / indentation):
			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc:  funcDoGetHTTPServer,
				DoGetHealthCheckFunc: funcDoGetHealthcheckOk,
				DoGetMongoDBFunc:     funcDoGetMongoDBOk,
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("When service Run succeeds and all the flags are set", func() {
				So(err, ShouldBeNil)
				So(svcList.HealthCheck, ShouldBeTrue)
				So(svcList.MongoDB, ShouldBeTrue)
			})

			Convey("Then checkers are registered and the healthcheck and http server started", func() {
				So(hcMock.AddCheckCalls(), ShouldHaveLength, 1)
				So(hcMock.AddCheckCalls()[0].Name, ShouldResemble, "Mongo DB")
				So(initMock.DoGetHTTPServerCalls(), ShouldHaveLength, 1)
				So(initMock.DoGetHTTPServerCalls()[0].BindAddr, ShouldEqual, "localhost:25500")
				So(initMock.DoGetMongoDBCalls()[0].Cfg.MongoConfig.BindAddr, ShouldEqual, "localhost:27017")
				So(hcMock.StartCalls(), ShouldHaveLength, 1)
				//!!! a call needed to stop the server, maybe ?
				serverWg.Wait() // Wait for HTTP server go-routine to finish
				So(serverMock.ListenAndServeCalls(), ShouldHaveLength, 1)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		// ADD CODE: put this code in, if you have Checkers to register
		Convey("Given that Checkers cannot be registered", func() {

			// setup (run before each `Convey` at this scope / indentation):
			errAddheckFail := errors.New("Error(s) registering checkers for healthcheck")
			hcMockAddFail := &serviceMock.HealthCheckerMock{
				AddCheckFunc: func(name string, checker healthcheck.Checker) error { return errAddheckFail },
				StartFunc:    func(ctx context.Context) {},
			}

			initMock := &serviceMock.InitialiserMock{
				DoGetHTTPServerFunc: funcDoGetHTTPServerNil,
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMockAddFail, nil
				},
				DoGetMongoDBFunc: funcDoGetMongoDBOk,
				// ADD CODE: add the checkers that you want to register here
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			Convey("Then service Run fails, but all checks try to register", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldResemble, fmt.Sprintf("unable to register checkers: %s", errAddheckFail.Error()))
				So(svcList.HealthCheck, ShouldBeTrue)
				// ADD CODE: add code to confirm checkers exist
				So(hcMockAddFail.AddCheckCalls(), ShouldHaveLength, 1)
				So(hcMockAddFail.AddCheckCalls()[0].Name, ShouldResemble, "Mongo DB") // ADD CODE: change the '0' to the number of checkers you have registered
			})
			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})

		Convey("Given that all dependencies are successfully initialised but the http server fails", func() {

			// setup (run before each `Convey` at this scope / indentation):
			initMock := &serviceMock.InitialiserMock{
				DoGetHealthCheckFunc: funcDoGetHealthcheckOk,
				DoGetHTTPServerFunc:  funcDoGetFailingHTTPServer,
				DoGetMongoDBFunc:     funcDoGetMongoDBOk,
			}
			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			serverWg.Add(1)
			_, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			Convey("Then the error is returned in the error channel", func() {
				sErr := <-svcErrors
				So(sErr.Error(), ShouldResemble, fmt.Sprintf("failure in http listen and serve: %s", errServer.Error()))
				So(failingServerMock.ListenAndServeCalls(), ShouldHaveLength, 1)
			})

			Reset(func() {
				// This reset is run after each `Convey` at the same scope (indentation)
			})
		})
	})
}

func TestClose(t *testing.T) {

	Convey("Having a correctly initialised service", t, func() {

		cfg, err := config.Get()
		So(err, ShouldBeNil)

		hcStopped := false

		// healthcheck Stop does not depend on any other service being closed/stopped
		hcMock := &serviceMock.HealthCheckerMock{
			AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
			StartFunc:    func(ctx context.Context) {},
			StopFunc:     func() { hcStopped = true },
		}

		// server Shutdown will fail if healthcheck is not stopped
		serverMock := &mock.HTTPServerMock{
			ListenAndServeFunc: func() error { return nil },
			ShutdownFunc: func(ctx context.Context) error {
				if !hcStopped {
					return errors.New("Server stopped before healthcheck")
				}
				return nil
			},
		}

		Convey("Closing the service results in all the dependencies being closed in the expected order", func() {

			mongoMock := &apiMock.AreaStoreMock{
				CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
				CloseFunc: func(ctx context.Context) error {
					return nil
				},
			}

			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return serverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetMongoDBFunc: func(ctx context.Context, cfg *config.Config) (api.AreaStore, error) {
					return mongoMock, nil
				},
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)
			So(err, ShouldBeNil)

			err = svc.Close(context.Background())
			So(err, ShouldBeNil)
			So(hcMock.StopCalls(), ShouldHaveLength, 1)
			So(mongoMock.CloseCalls(), ShouldHaveLength, 1)
			So(serverMock.ShutdownCalls(), ShouldHaveLength, 1)
		})

		Convey("If Mongo fails to Close and returns an error", func() {

			mongoMockCloseErr := &apiMock.AreaStoreMock{
				CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
				CloseFunc: func(ctx context.Context) error {
					return errors.New("Closing mongo timed out")
				},
			}

			initMock := &mock.InitialiserMock{
				DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer { return serverMock },
				DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
					return hcMock, nil
				},
				DoGetMongoDBFunc: func(ctx context.Context, cfg *config.Config) (api.AreaStore, error) {
					return mongoMockCloseErr, nil
				},
			}

			svcErrors := make(chan error, 1)
			svcList := service.NewServiceList(initMock)
			svc, err := service.Run(ctx, cfg, svcList, testBuildTime, testGitCommit, testVersion, svcErrors)

			err = svc.Close(context.Background())
			So(err, ShouldBeError, "failed to shutdown gracefully")
			So(svc.ServiceList.MongoDB, ShouldBeTrue)
		})

		Convey("If service times out while shutting down, the Close operation fails with the expected error", func() {
			mongoMock := &apiMock.AreaStoreMock{
				CheckerFunc: func(ctx context.Context, state *healthcheck.CheckState) error { return nil },
				CloseFunc:   func(ctx context.Context) error { return nil },
			}

			cfg.GracefulShutdownTimeout = 100 * time.Millisecond
			timeoutServerMock := &mock.HTTPServerMock{
				ListenAndServeFunc: func() error { return nil },
				ShutdownFunc: func(ctx context.Context) error {
					time.Sleep(200 * time.Millisecond)
					return nil
				},
			}

			svcList := service.NewServiceList(nil)
			svcList.HealthCheck = true
			svc := service.Service{
				Config:      cfg,
				ServiceList: svcList,
				Server:      timeoutServerMock,
				HealthCheck: hcMock,
				MongoDB:     mongoMock,
			}

			err = svc.Close(context.Background())
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, "context deadline exceeded")
			So(hcMock.StopCalls(), ShouldHaveLength, 1)
			So(timeoutServerMock.ShutdownCalls(), ShouldHaveLength, 1)
		})
	})
}
