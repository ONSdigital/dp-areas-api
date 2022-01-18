package pgx

import (
	"context"

	"fmt"

	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ONSdigital/log.go/v2/log"
)

//go:generate moq -out mock/pgx.go -pkg mock . PGXPool

type PGXPool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Close()
}

type PGX struct {
	Pool PGXPool
}

func NewPGXHandler(ctx context.Context, cfg *config.Config) (*PGX, error) {
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
	return &PGX{
		Pool: rdsConn,
	}, nil
}
