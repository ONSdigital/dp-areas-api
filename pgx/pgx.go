package pgx

import (
	"context"

	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ONSdigital/log.go/v2/log"
)

//go:generate moq -out mock/pgx.go -pkg mock . PGXPool

// PGXPool interface is abstraction of pgxpool package
type PGXPool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Close()
}

type PGX struct {
	Pool PGXPool
}

func NewPGXHandler(ctx context.Context, cfg *config.Config) (*PGX, error) {
	var connectionString string
	if cfg.DPPostgresLocal {
		connectionString = cfg.GetLocalDBConnectionString()
	} else {
		authToken, err := rdsutils.BuildAuthToken(cfg.GetDBEndpoint(), cfg.AWSRegion, cfg.RDSDBUser, credentials.NewEnvCredentials())
		if err != nil {
			log.Error(ctx, "error building auth token for rds connection", err)
			return nil, err
		}
		connectionString = cfg.GetRemoteDBConnectionString(authToken)
	}
	// generate the rds connection
	rdsConn, err := pgxpool.Connect(ctx, connectionString)
	if err != nil {
		log.Error(ctx, "error connecting to rds instance", err)
		return nil, err
	}
	return &PGX{
		Pool: rdsConn,
	}, nil
}
