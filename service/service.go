package service

import (
	"context"
	"net/http"

	clientsidentity "github.com/ONSdigital/dp-api-clients-go/identity"
	"github.com/ONSdigital/dp-authorisation/auth"
	"github.com/justinas/alice"

	dphandlers "github.com/ONSdigital/dp-net/handlers"
	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/dp-topic-api/api"
	"github.com/ONSdigital/dp-topic-api/config"
	"github.com/ONSdigital/dp-topic-api/store"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// check that DatsetAPIStore satifies the the store.Storer interface
var _ store.Storer = (*DatsetAPIStore)(nil)

//DatsetAPIStore is a wrapper which embeds (Neo4j) Mongo structs which between them satisfy the store.Storer interface.
type DatsetAPIStore struct {
	store.MongoDB
	//	store.GraphDB
}

// Service contains all the configs, server and clients to run the dp-topic-api API
type Service struct {
	Config         *config.Config
	ServiceList    *ExternalServiceList
	Server         HTTPServer
	Router         *mux.Router
	API            *api.API
	HealthCheck    HealthChecker
	mongoDB        store.MongoDB
	IdentityClient *clientsidentity.Client
}

// New creates a new service
func New(cfg *config.Config, serviceList *ExternalServiceList) *Service {
	svc := &Service{
		Config:      cfg,
		ServiceList: serviceList,
	}
	return svc
}

// Run the service
func (svc *Service) Run(ctx context.Context, buildTime, gitCommit, version string, svcErrors chan error) (err error) {

	// Get MongoDB client
	svc.mongoDB, err = svc.ServiceList.GetMongoDB(ctx, svc.Config)
	if err != nil {
		log.Event(ctx, "failed to initialise mongo DB", log.FATAL, log.Error(err))
		return err
	}
	store := store.DataStore{Backend: DatsetAPIStore{svc.mongoDB}}

	// Get Identity Client (only if private endpoints are enabled)
	if svc.Config.EnablePrivateEndpoints {
		// Only in Publishing ... create client(s):
		svc.IdentityClient = clientsidentity.New(svc.Config.ZebedeeURL)
	}

	// Get HealthCheck
	svc.HealthCheck, err = svc.ServiceList.GetHealthCheck(svc.Config, buildTime, gitCommit, version)
	if err != nil {
		log.Event(ctx, "could not instantiate healthcheck", log.FATAL, log.Error(err))
		return err
	}

	if err := svc.registerCheckers(ctx); err != nil {
		return errors.Wrap(err, "unable to register checkers")
	}

	// Get HTTP router and server with middleware
	router := mux.NewRouter()
	middle := svc.createMiddleware(svc.Config)
	svc.Server = svc.ServiceList.GetHTTPServer(svc.Config.BindAddr, middle.Then(router))

	// Setup the API
	permissions := getAuthorisationHandlers(ctx, svc.Config)
	svc.API = api.Setup(ctx, svc.Config, router, store, permissions)

	svc.HealthCheck.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()

	return nil
}

func getAuthorisationHandlers(ctx context.Context, cfg *config.Config) api.AuthHandler {
	if cfg.EnablePermissionsAuth == false {
		log.Event(ctx, "feature flag not enabled defaulting to nop authZ impl", log.INFO, log.Data{"feature": "ENABLE_PERMISSIONS_AUTHZ"})
		return &auth.NopHandler{}
	}

	log.Event(ctx, "feature flag enabled", log.INFO, log.Data{"feature": "ENABLE_PERMISSIONS_AUTHZ"})

	authClient := auth.NewPermissionsClient(dphttp.NewClient())
	authVerifier := auth.DefaultPermissionsVerifier()

	// for checking caller permissions when we only have a user/service token
	permissions := auth.NewHandler(
		auth.NewPermissionsRequestBuilder(cfg.ZebedeeURL),
		authClient,
		authVerifier,
	)

	return permissions
}

// CreateMiddleware creates an Alice middleware chain of handlers
// to forward collectionID from cookie from header
func (svc *Service) createMiddleware(cfg *config.Config) alice.Chain {

	// healthcheck
	healthcheckHandler := healthcheckMiddleware(svc.HealthCheck.Handler, "/health")
	middleware := alice.New(healthcheckHandler)

	// Only add the identity middleware when running in publishing.
	if cfg.EnablePrivateEndpoints {
		middleware = middleware.Append(dphandlers.IdentityWithHTTPClient(svc.IdentityClient))
	}

	return middleware
}

// healthcheckMiddleware creates a new http.Handler to intercept /health requests.
func healthcheckMiddleware(healthcheckHandler func(http.ResponseWriter, *http.Request), path string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			if req.Method == "GET" && req.URL.Path == path {
				healthcheckHandler(w, req)
				return
			}

			h.ServeHTTP(w, req)
		})
	}
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Event(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout}, log.INFO)
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
			log.Event(ctx, "failed to shutdown http server", log.Error(err), log.ERROR)
			hasShutdownError = true
		}

		// ADD CODE HERE: Close other dependencies, in the expected order

		// close mongoDB
		if svc.ServiceList.MongoDB {
			if err := svc.mongoDB.Close(ctx); err != nil {
				log.Event(ctx, "error closing mongoDB", log.Error(err), log.ERROR)
				hasShutdownError = true
			}
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Event(ctx, "shutdown timed out", log.ERROR, log.Error(ctx.Err()))
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Event(ctx, "failed to shutdown gracefully ", log.ERROR, log.Error(err))
		return err
	}

	log.Event(ctx, "graceful shutdown was successful", log.INFO)
	return nil
}

// registerCheckers - registers functions which are periodically called to validate
//      the health state of external services that this application depends upon.
func (svc *Service) registerCheckers(ctx context.Context) (err error) {

	// ADD CODE: add other health checks here, as per dp-upload-service

	hasErrors := false

	if svc.Config.EnablePrivateEndpoints {
		// Only in Publishing ... add check(s):

		if err = svc.HealthCheck.AddCheck("Zebedee", svc.IdentityClient.Checker); err != nil {
			hasErrors = true
			log.Event(ctx, "error adding check for zebedee", log.ERROR, log.Error(err))
		}
	}

	if err = svc.HealthCheck.AddCheck("Mongo DB", svc.mongoDB.Checker); err != nil {
		hasErrors = true
		log.Event(ctx, "error adding check for mongo db", log.ERROR, log.Error(err))
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}
	return nil
}
