package rds

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-areas-api/pgx"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RDS struct {
	conn pgx.PGXPool
}

func (r *RDS) Init(ctx context.Context, cfg *config.Config) error {
	authToken, err := rdsutils.BuildAuthToken(
		cfg.GetRDSEndpoint(),
		cfg.AWSRegion, cfg.RDSDBUser,
		credentials.NewEnvCredentials(),
	)
	if err != nil {
		return err
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.RDSDBHost, cfg.RDSDBPort, cfg.RDSDBUser, authToken, cfg.RDSDBName,
	)

	rdsConn, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
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
