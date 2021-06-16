package steps

import (
	"context"
	"fmt"
	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/mongo"
	"github.com/ONSdigital/dp-topic-api/service"
	serviceMock "github.com/ONSdigital/dp-topic-api/service/mock"
	"github.com/ONSdigital/dp-topic-api/store"
	"github.com/benweissmann/memongo"
	"github.com/cucumber/godog"
	"net/http"
)

type TopicComponent struct {
	ErrorFeature   componenttest.ErrorFeature
	svc            *service.Service
	errorChan      chan error
	MongoClient    *mongo.Mongo
	Config         *config.Config
	HTTPServer     *http.Server
	ServiceRunning bool
}

func NewTopicComponent(mongoFeature *componenttest.MongoFeature, zebedeeURL string) (*TopicComponent, error) {

	f := &TopicComponent{
		HTTPServer:     &http.Server{},
		errorChan:      make(chan error),
		ServiceRunning: false,
	}

	var err error

	f.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	f.Config.ZebedeeURL = zebedeeURL

	f.Config.EnablePrivateEndpoints = false // for component tests, ensure 'false' to start

	f.Config.EnablePermissionsAuth = false

	getMongoURI := fmt.Sprintf("localhost:%d", mongoFeature.Server.Port())
	mongodb := &mongo.Mongo{
		Database:          memongo.RandomDatabase(),
		URI:               getMongoURI,
		Username:          "",
		Password:          "",
		TopicsCollection:  f.Config.MongoConfig.TopicsCollection,
		ContentCollection: f.Config.MongoConfig.ContentCollection,
		IsSSL:             false,
	}

	if err := mongodb.Init(context.TODO(), false, true); err != nil {
		return nil, err
	}

	f.MongoClient = mongodb

	initMock := &serviceMock.InitialiserMock{
		DoGetMongoDBFunc:     f.DoGetMongoDB,
		DoGetHealthCheckFunc: f.DoGetHealthcheckOk,
		DoGetHTTPServerFunc:  f.DoGetHTTPServer,
	}

	f.svc = service.New(f.Config, service.NewServiceList(initMock))

	return f, nil
}

func (f *TopicComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^private endpoints are enabled$`, f.privateEndpointsAreEnabled)
	ctx.Step(`^I have these topics:$`, f.iHaveTheseTopics)
	ctx.Step(`^I have these contents:$`, f.iHaveTheseContents)
}

func (f *TopicComponent) Reset() *TopicComponent {
	f.MongoClient.Database = memongo.RandomDatabase()
	f.MongoClient.Init(context.TODO(), false, true)
	f.Config.EnablePrivateEndpoints = false
	return f
}

func (f *TopicComponent) Close() error {
	if f.svc != nil && f.ServiceRunning {
		f.svc.Close(context.Background())
		f.ServiceRunning = false
	}
	return nil
}

func (f *TopicComponent) InitialiseService() (http.Handler, error) {
	if err := f.svc.Run(context.Background(), "1", "", "", f.errorChan); err != nil {
		return nil, err
	}
	f.ServiceRunning = true
	return f.HTTPServer.Handler, nil
}

func (f *TopicComponent) DoGetHealthcheckOk(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return &serviceMock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (f *TopicComponent) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
	f.HTTPServer.Addr = bindAddr
	f.HTTPServer.Handler = router
	return f.HTTPServer
}

// DoGetMongoDB returns a MongoDB
func (f *TopicComponent) DoGetMongoDB(ctx context.Context, cfg *config.Config) (store.MongoDB, error) {
	return f.MongoClient, nil
}
