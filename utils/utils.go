package utils

import (
	"context"
	"fmt"
	"strconv"

	errs "github.com/ONSdigital/dp-areas-api/apierrors"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/jackc/pgx/v4/pgxpool"
)

// ValidatePositiveInt obtains the positive int value of query
func ValidatePositiveInt(parameter string) (val int, err error) {
	val, err = strconv.Atoi(parameter)
	if err != nil {
		return -1, errs.ErrInvalidQueryParameter
	}
	if val < 0 {
		return -1, errs.ErrInvalidQueryParameter
	}
	return val, nil
}

// gGenerateRDSHandle returns an open rds handle
func GenerateTestRDSHandle(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
    dbName := cfg.RDSDBName
    dbUser := cfg.RDSDBUser
    dbHost := cfg.RDSDBHost
    dbPort := cfg.RDSDBPort
    dbEndpoint := fmt.Sprintf("%s:%s", dbHost, dbPort)
    region := cfg.AWSRegion

    creds := credentials.NewEnvCredentials()
    authToken, err := rdsutils.BuildAuthToken(dbEndpoint, region, dbUser, creds)
    if err != nil {
		log.Error(ctx, "error building auth token for rds connection", err)
        return nil, err
    }

   connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
        dbHost, dbPort, dbUser, authToken, dbName,
    )

	// generate the rds connection
	rdsConn, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		log.Error(ctx, "error connecting to rds instance", err)
		return nil, err
	}

	return rdsConn, nil
}
