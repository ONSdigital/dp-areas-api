package steps

import (
	"context"
	"net/http"
	"time"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-component-test/utils"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/mongo"
	"github.com/ONSdigital/dp-topic-api/service"
	serviceMock "github.com/ONSdigital/dp-topic-api/service/mock"
	"github.com/ONSdigital/dp-topic-api/store"
	"github.com/cucumber/godog"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
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

func NewTopicComponent(mongoURL, zebedeeURL string) (*TopicComponent, error) {

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
	f.Config.EnablePrivateEndpoints = false

	f.Config.MongoConfig.ClusterEndpoint = mongoURL
	f.Config.MongoConfig.Database = utils.RandomDatabase()
	f.Config.MongoConfig.Username, f.Config.MongoConfig.Password = createCredsInDB(f.Config.MongoConfig.ClusterEndpoint, f.Config.MongoConfig.Database)

	f.MongoClient, err = mongo.NewDBConnection(context.TODO(), f.Config.MongoConfig)
	if err != nil {
		return nil, err
	}

	initMock := &serviceMock.InitialiserMock{
		DoGetMongoDBFunc:     f.DoGetMongoDB,
		DoGetHealthCheckFunc: f.DoGetHealthcheckOk,
		DoGetHTTPServerFunc:  f.DoGetHTTPServer,
	}

	f.svc = service.New(f.Config, service.NewServiceList(initMock))

	return f, nil
}

func createCredsInDB(getMongoURI string, databaseName string) (string, string) {
	username := "admin"
	password, _ := uuid.NewV4()
	mongoConnectionConfig := &dpMongoDriver.MongoDriverConfig{
		TLSConnectionConfig: dpMongoDriver.TLSConnectionConfig{
			IsSSL: false,
		},
		ConnectTimeout: 15 * time.Second,
		QueryTimeout:   15 * time.Second,

		Username:        "",
		Password:        "",
		ClusterEndpoint: getMongoURI,
		Database:        databaseName,
	}
	mongoConnection, err := dpMongoDriver.Open(mongoConnectionConfig)
	if err != nil {
		panic("expected db connection to be opened")
	}

	createCollectionResponse := mongoConnection.RunCommand(context.TODO(), bson.D{{"create", "test"}})

	if createCollectionResponse != nil {
		panic("expected collection creation to go through")
	}
	userCreationResponse := mongoConnection.RunCommand(context.TODO(), bson.D{
		{Key: "createUser", Value: username},
		{Key: "pwd", Value: password.String()},
		{Key: "roles", Value: []bson.M{
			{"role": "root", "db": "admin"},
		}},
	})
	if userCreationResponse != nil {
		panic("expected user creation to go through")
	}
	return username, password.String()
}

func (f *TopicComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^private endpoints are enabled$`, f.privateEndpointsAreEnabled)
	ctx.Step(`^I have these topics:$`, f.iHaveTheseTopics)
	ctx.Step(`^I have these contents:$`, f.iHaveTheseContents)
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
func (f *TopicComponent) DoGetMongoDB(ctx context.Context, cfg config.MongoConfig) (store.MongoDB, error) {
	return f.MongoClient, nil
}
