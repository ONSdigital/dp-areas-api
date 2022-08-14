package healthcheck

import (
	"context"
	"net/http"
	"regexp"

	"github.com/ONSdigital/dp-areas-api/api"
	"github.com/ONSdigital/dp-areas-api/config"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	RDSHealthy = "RDS Healthy"
	authError  = "SQLSTATE 28000"
)

var (
	authErrorRegex = regexp.MustCompile(authError)
)

func RDSHealthCheck(ctx context.Context, cfg *config.Config, rds api.RDSAreaStore) health.Checker {
	return func(ctx context.Context, state *health.CheckState) error {
		err := rds.Ping(ctx)
		log.Info(context.Background(), "Checking rds connection status...")

		// special case for auth errors
		if err != nil && authErrorRegex.MatchString(err.Error()) {
			log.Error(context.Background(), "Attempting to re-connect to rds instance", err)
			rds.Init(ctx, cfg)
			pingErr := rds.Ping(ctx)
			if pingErr != nil {
				log.Error(context.Background(), "Attempting to re-connect to rds instance on PAM fail", pingErr)
			}
		}

		if err != nil {
			if stateErr := state.Update(health.StatusCritical, err.Error(), http.StatusBadGateway); stateErr != nil {
				log.Error(context.Background(), "Error updating state during area service healthcheck", stateErr)
			}
			// log the error
			log.Error(context.Background(), "Error running area service healthcheck", err)
			return err
		}

		if stateErr := state.Update(health.StatusOK, RDSHealthy, http.StatusOK); stateErr != nil {
			log.Error(context.Background(), "Error updating state during area service healthcheck", stateErr)
		}

		return nil
	}
}
