// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-areas-api/api"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/service"
	"net/http"
	"sync"
)

// Ensure, that InitialiserMock does implement service.Initialiser.
// If this is not the case, regenerate this file with moq.
var _ service.Initialiser = &InitialiserMock{}

// InitialiserMock is a mock implementation of service.Initialiser.
//
// 	func TestSomethingThatUsesInitialiser(t *testing.T) {
//
// 		// make and configure a mocked service.Initialiser
// 		mockedInitialiser := &InitialiserMock{
// 			DoGetHTTPServerFunc: func(bindAddr string, router http.Handler) service.HTTPServer {
// 				panic("mock out the DoGetHTTPServer method")
// 			},
// 			DoGetHealthCheckFunc: func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
// 				panic("mock out the DoGetHealthCheck method")
// 			},
// 			DoGetMongoDBFunc: func(ctx context.Context, cfg config.MongoConfig) (api.AreaStore, error) {
// 				panic("mock out the DoGetMongoDB method")
// 			},
// 		}
//
// 		// use mockedInitialiser in code that requires service.Initialiser
// 		// and then make assertions.
//
// 	}
type InitialiserMock struct {
	// DoGetHTTPServerFunc mocks the DoGetHTTPServer method.
	DoGetHTTPServerFunc func(bindAddr string, router http.Handler) service.HTTPServer

	// DoGetHealthCheckFunc mocks the DoGetHealthCheck method.
	DoGetHealthCheckFunc func(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error)

	// DoGetMongoDBFunc mocks the DoGetMongoDB method.
	DoGetMongoDBFunc func(ctx context.Context, cfg config.MongoConfig) (api.AreaStore, error)

	// calls tracks calls to the methods.
	calls struct {
		// DoGetHTTPServer holds details about calls to the DoGetHTTPServer method.
		DoGetHTTPServer []struct {
			// BindAddr is the bindAddr argument value.
			BindAddr string
			// Router is the router argument value.
			Router http.Handler
		}
		// DoGetHealthCheck holds details about calls to the DoGetHealthCheck method.
		DoGetHealthCheck []struct {
			// Cfg is the cfg argument value.
			Cfg *config.Config
			// BuildTime is the buildTime argument value.
			BuildTime string
			// GitCommit is the gitCommit argument value.
			GitCommit string
			// Version is the version argument value.
			Version string
		}
		// DoGetMongoDB holds details about calls to the DoGetMongoDB method.
		DoGetMongoDB []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Cfg is the cfg argument value.
			Cfg config.MongoConfig
		}
	}
	lockDoGetHTTPServer  sync.RWMutex
	lockDoGetHealthCheck sync.RWMutex
	lockDoGetMongoDB     sync.RWMutex
}

// DoGetHTTPServer calls DoGetHTTPServerFunc.
func (mock *InitialiserMock) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	if mock.DoGetHTTPServerFunc == nil {
		panic("InitialiserMock.DoGetHTTPServerFunc: method is nil but Initialiser.DoGetHTTPServer was just called")
	}
	callInfo := struct {
		BindAddr string
		Router   http.Handler
	}{
		BindAddr: bindAddr,
		Router:   router,
	}
	mock.lockDoGetHTTPServer.Lock()
	mock.calls.DoGetHTTPServer = append(mock.calls.DoGetHTTPServer, callInfo)
	mock.lockDoGetHTTPServer.Unlock()
	return mock.DoGetHTTPServerFunc(bindAddr, router)
}

// DoGetHTTPServerCalls gets all the calls that were made to DoGetHTTPServer.
// Check the length with:
//     len(mockedInitialiser.DoGetHTTPServerCalls())
func (mock *InitialiserMock) DoGetHTTPServerCalls() []struct {
	BindAddr string
	Router   http.Handler
} {
	var calls []struct {
		BindAddr string
		Router   http.Handler
	}
	mock.lockDoGetHTTPServer.RLock()
	calls = mock.calls.DoGetHTTPServer
	mock.lockDoGetHTTPServer.RUnlock()
	return calls
}

// DoGetHealthCheck calls DoGetHealthCheckFunc.
func (mock *InitialiserMock) DoGetHealthCheck(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	if mock.DoGetHealthCheckFunc == nil {
		panic("InitialiserMock.DoGetHealthCheckFunc: method is nil but Initialiser.DoGetHealthCheck was just called")
	}
	callInfo := struct {
		Cfg       *config.Config
		BuildTime string
		GitCommit string
		Version   string
	}{
		Cfg:       cfg,
		BuildTime: buildTime,
		GitCommit: gitCommit,
		Version:   version,
	}
	mock.lockDoGetHealthCheck.Lock()
	mock.calls.DoGetHealthCheck = append(mock.calls.DoGetHealthCheck, callInfo)
	mock.lockDoGetHealthCheck.Unlock()
	return mock.DoGetHealthCheckFunc(cfg, buildTime, gitCommit, version)
}

// DoGetHealthCheckCalls gets all the calls that were made to DoGetHealthCheck.
// Check the length with:
//     len(mockedInitialiser.DoGetHealthCheckCalls())
func (mock *InitialiserMock) DoGetHealthCheckCalls() []struct {
	Cfg       *config.Config
	BuildTime string
	GitCommit string
	Version   string
} {
	var calls []struct {
		Cfg       *config.Config
		BuildTime string
		GitCommit string
		Version   string
	}
	mock.lockDoGetHealthCheck.RLock()
	calls = mock.calls.DoGetHealthCheck
	mock.lockDoGetHealthCheck.RUnlock()
	return calls
}

// DoGetMongoDB calls DoGetMongoDBFunc.
func (mock *InitialiserMock) DoGetMongoDB(ctx context.Context, cfg config.MongoConfig) (api.AreaStore, error) {
	if mock.DoGetMongoDBFunc == nil {
		panic("InitialiserMock.DoGetMongoDBFunc: method is nil but Initialiser.DoGetMongoDB was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Cfg config.MongoConfig
	}{
		Ctx: ctx,
		Cfg: cfg,
	}
	mock.lockDoGetMongoDB.Lock()
	mock.calls.DoGetMongoDB = append(mock.calls.DoGetMongoDB, callInfo)
	mock.lockDoGetMongoDB.Unlock()
	return mock.DoGetMongoDBFunc(ctx, cfg)
}

// DoGetMongoDBCalls gets all the calls that were made to DoGetMongoDB.
// Check the length with:
//     len(mockedInitialiser.DoGetMongoDBCalls())
func (mock *InitialiserMock) DoGetMongoDBCalls() []struct {
	Ctx context.Context
	Cfg config.MongoConfig
} {
	var calls []struct {
		Ctx context.Context
		Cfg config.MongoConfig
	}
	mock.lockDoGetMongoDB.RLock()
	calls = mock.calls.DoGetMongoDB
	mock.lockDoGetMongoDB.RUnlock()
	return calls
}
