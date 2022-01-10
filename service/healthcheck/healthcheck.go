package healthcheck

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-areas-api/rds"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
)

var (
	dbInstanceName = "develop-area-profiles-postgres-one"
)

const RDSHealthy = "RDS Healthy"

func RDSHealthCheck(client rds.Client) health.Checker {
	return func(ctx context.Context, state *health.CheckState) error {
		dbInstanceRequest := models.BuildDescibeDBInstancesRequest(&dbInstanceName)
		_, err := client.DescribeDBInstances(dbInstanceRequest)

		if err != nil {
			if stateErr := state.Update(health.StatusCritical, err.Error(), http.StatusTooManyRequests); stateErr != nil {
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
