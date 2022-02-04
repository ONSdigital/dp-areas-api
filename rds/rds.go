package rds

import (
	"context"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-areas-api/pgx"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RDS struct {
	conn pgx.PGXPool
}

func (r *RDS) Init(ctx context.Context, cfg *config.Config) error {
	var connectionString string
	if cfg.DPPostgresLocal {
		connectionString = cfg.GetLocalDBConnectionString()
	} else {
		authToken, err := rdsutils.BuildAuthToken(cfg.GetDBEndpoint(), cfg.AWSRegion, cfg.RDSDBUser, credentials.NewEnvCredentials())
		if err != nil {
			log.Error(ctx, "error building auth token for rds connection", err)
			return err
		}
		connectionString = cfg.GetRemoteDBConnectionString(authToken)
	}

	rdsConn, err := pgxpool.Connect(ctx, connectionString)
	if err != nil {
		log.Error(ctx, "error connecting to rds instance", err)
		return err
	}

	r.conn = rdsConn
	return nil
}

func (r *RDS) Close() {
	r.conn.Close()
}

func (r *RDS) ValidateArea(areaCode string) error {
	var code string
	err := r.conn.QueryRow(context.Background(), getAreaCode, areaCode).Scan(&code)
	return err
}

func (r *RDS) GetArea(areaId string) (*models.AreaDataRDS, error) {
	var (
		code   string
		id     int64
		active bool
	)
	err := r.conn.QueryRow(context.Background(), getBasicArea, areaId).Scan(&id, &code, &active)
	if err != nil {
		return nil, err
	}
	return &models.AreaDataRDS{Id: id, Code: code, Active: active}, nil
}

func (r *RDS) GetRelationships(areaCode string) ([]*models.AreaBasicData, error) {
	var relationships []*models.AreaBasicData
	rows, err := r.conn.Query(context.Background(), getRelationShipAreas, areaCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rs models.AreaBasicData
		rows.Scan(&rs.Code, &rs.Name)
		relationships = append(relationships, &rs)
	}
	return relationships, nil
}

func (r *RDS) BuildTables(ctx context.Context, executionList []string) error {
	for index := range executionList {
		logData := log.Data{"Exceuting Create Table Query": executionList[index]}
		_, err := r.conn.Exec(ctx, executionList[index])
		if err != nil {
			return err
		}
		log.Info(ctx, "table created successfully:", logData)
	}
	return nil
}
