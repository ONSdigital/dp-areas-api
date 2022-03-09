package service

import (
	"context"
	s3 "github.com/ONSdigital/dp-s3/v2"

	"github.com/ONSdigital/dp-areas-api/api"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-areas-api/models/DBRelationalSchema"

	"github.com/ONSdigital/dp-areas-api/pgx"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	health "github.com/ONSdigital/dp-areas-api/service/healthcheck"
)

// Service contains all the configs, server and clients to run the dp-areas-api API
const (
	databaseName     = "dp-areas-api"
	dpPublishingUser = "dp-areas-api-publishing"
)

// Service contains all the configs, server and clients to run the dp-areas-api API
type Service struct {
	Config      *config.Config
	Server      HTTPServer
	Router      *mux.Router
	API         *api.API
	ServiceList *ExternalServiceList
	HealthCheck HealthChecker
	RDS         api.RDSAreaStore
	PGX         *pgx.PGX
}

// Run the service
func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {

	log.Info(ctx, "running service")

	log.Info(ctx, "using service configuration", log.Data{"config": cfg})

	// Get HTTP Server and ... // ADD CODE: Add any middleware that your service requires
	r := mux.NewRouter()

	s := serviceList.GetHTTPServer(cfg.BindAddr, r)

	// Get RDS client
	rds, err := serviceList.getRDSDB(ctx, cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise rds db client", err)
		return nil, err
	}

	// Get S3 client
	s3Client, err := serviceList.getS3Client(cfg)
	if err != nil {
		log.Fatal(ctx, "failed to initialise s3 client", err)
		return nil, err
	}

	// only run publishing user or if pointing at local postgres instance
	if cfg.RDSDBUser == dpPublishingUser || cfg.DPPostgresLocal {
		// rds table schema builder
		rdsSchema := &models.DatabaseSchema{
			DBName:       databaseName,
			SchemaString: DBRelationalSchema.DBSchema,
		}
		err = rdsSchema.BuildDatabaseSchemaModel()
		if err != nil {
			log.Fatal(ctx, "error building database schema model", err)
			return nil, err
		}
		rdsSchema.TableSchemaBuilder()
		if err != nil {
			log.Fatal(ctx, "error building database table schema", err)
			return nil, err
		}
		err = rds.BuildTables(ctx, rdsSchema.ExecutionList)
		if err != nil {
			log.Fatal(ctx, "error building database schema target", err)
			return nil, err
		}
	}

	// Setup the API
	a, _ := api.Setup(ctx, cfg, r, rds)

	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		log.Fatal(ctx, "could not instantiate healthcheck", err)
		return nil, err
	}

	if err := registerCheckers(ctx, cfg, hc, rds, s3Client); err != nil {
		return nil, errors.Wrap(err, "unable to register checkers")
	}

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	hc.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if err := s.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return &Service{
		Config:      cfg,
		Router:      r,
		API:         a,
		HealthCheck: hc,
		ServiceList: serviceList,
		Server:      s,
		RDS:         rds,
	}, nil
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.ServiceList.HealthCheck {
			svc.HealthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.Server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

		// close RDS DB connection
		if svc.RDS != nil {
			svc.RDS.Close()
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Error(ctx, "failed to shutdown gracefully", err)
		return err
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func registerCheckers(ctx context.Context, cfg *config.Config, hc HealthChecker, rds api.RDSAreaStore, s3Client *s3.Client) (err error) {
	hasErrors := false

	if err := hc.AddCheck("S3 healthchecker", s3Client.Checker); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for s3 client", err)
	}
	if err := hc.AddCheck("RDS healthchecker", health.RDSHealthCheck(ctx, cfg, rds)); err != nil {
		hasErrors = true
		log.Error(ctx, "error adding check for rds client", err)
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}
	return nil
}
