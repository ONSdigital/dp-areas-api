package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-areas-api/api"
	"github.com/ONSdigital/dp-areas-api/rds"
	"github.com/ONSdigital/dp-areas-api/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/mongo"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"

	aurorards "github.com/aws/aws-sdk-go/service/rds"
)

// ExternalServiceList holds the initialiser and initialisation state of external services.
type ExternalServiceList struct {
	HealthCheck bool
	Init        Initialiser
	MongoDB     bool
}

// NewServiceList creates a new service list with the provided initialiser
func NewServiceList(initialiser Initialiser) *ExternalServiceList {
	return &ExternalServiceList{
		HealthCheck: false,
		Init:        initialiser,
		MongoDB:     false,
	}
}

// Init implements the Initialiser interface to initialise dependencies
type Init struct{}

// GetHTTPServer creates an http server
func (e *ExternalServiceList) GetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := e.Init.DoGetHTTPServer(bindAddr, router)
	return s
}

// GetMongoDB creates a mongoDB client and sets the Mongo flag to true
func (e *ExternalServiceList) GetMongoDB(ctx context.Context, cfg config.MongoConfig) (api.AreaStore, error) {
	mongoDB, err := e.Init.DoGetMongoDB(ctx, cfg)
	if err != nil {
		log.Error(ctx, "failed to create mongodb client", err)
		return nil, err
	}
	e.MongoDB = true
	return mongoDB, nil
}

// GetHealthCheck creates a healthcheck with versionInfo and sets teh HealthCheck flag to true
func (e *ExternalServiceList) GetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	hc, err := e.Init.DoGetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	e.HealthCheck = true
	return hc, nil
}

// GetRDSClient creates aurora rds client
func (e *ExternalServiceList) GetRDSClient(region string) rds.Client {
	client := e.Init.DoGetRDSClient(region)
	return client
}

// GetRDSClient creates aurora rds client
func (e *ExternalServiceList) GetPGXPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pgxConn, err := e.Init.DoGetPGXPool(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return pgxConn, nil
}

// DoGetHTTPServer creates an HTTP Server with the provided bind address and router
func (e *Init) DoGetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

// DoGetMongoDB returns a MongoDB
func (e *Init) DoGetMongoDB(ctx context.Context, cfg config.MongoConfig) (api.AreaStore, error) {
	mongoDB, err := mongo.NewMongoStore(ctx, cfg)
	if err != nil {
		log.Error(ctx, "failed to intialise mongo", err)
		return nil, err
	}

	return mongoDB, nil
}

// DoGetHealthCheck creates a healthcheck with versionInfo
func (e *Init) DoGetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, err
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

// DoGetRDSClient creates a cognito client
func (e *Init) DoGetRDSClient(region string) rds.Client {
	client := aurorards.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})), &aws.Config{Region: &region})
	return client
}

// DoGetPGXPool creates a pgx pool connector
func (e *Init) DoGetPGXPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pgxConn, err := utils.GenerateTestRDSHandle(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return pgxConn, nil
}
