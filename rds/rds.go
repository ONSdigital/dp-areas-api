package rds

import (
	"context"
	"fmt"

	"github.com/ONSdigital/dp-areas-api/config"
	"github.com/ONSdigital/dp-areas-api/models"
	"github.com/ONSdigital/dp-areas-api/models/DBRelationalData"
	"github.com/ONSdigital/dp-areas-api/pgx"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/rds/rdsutils"
	"github.com/jackc/pgx/v4/pgxpool"
)

type RDS struct {
	conn             pgx.PGXPool
	useLocalPostgres bool
}

func (r *RDS) Init(ctx context.Context, cfg *config.Config) error {
	var connectionString string
	if cfg.DPPostgresLocal {
		connectionString = cfg.GetLocalDBConnectionString()
		r.useLocalPostgres = true
	} else {
		authToken, err := rdsutils.BuildAuthToken(cfg.GetDBEndpoint(), cfg.AWSRegion, cfg.RDSDBUser, credentials.NewEnvCredentials())
		if err != nil {
			log.Error(ctx, "error building auth token for rds connection", err)
			return err
		}
		connectionString = cfg.GetRemoteDBConnectionString(authToken)
		r.useLocalPostgres = false
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

func (r *RDS) GetArea(ctx context.Context, areaId string) (*models.AreasDataResults, error) {
	area := models.AreasDataResults{}
	err := r.conn.QueryRow(ctx, getArea, areaId).Scan(&area.Code, &area.Name, &area.GeometricData, &area.Visible, &area.AreaType)
	if err != nil {
		return nil, err
	}
	return &area, nil
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
	var err error
	for index := range executionList {
		logData := log.Data{"exceuting query": executionList[index]}
		_, err = r.conn.Exec(ctx, executionList[index])
		if err != nil {
			return err
		}
		log.Info(ctx, "query executed successfully:", logData)
	}
	//  seed local instance with test data
	if r.useLocalPostgres {
		err = r.insertAreaTypeTestData(ctx)
		if err != nil {
			return err
		}
		err = r.upsertAreaTestData(ctx)
		if err != nil {
			return err
		}
		err = r.insertRelationshipTypeTestData(ctx)
		if err != nil {
			return err
		}
		err = r.insertAreaNameTestData(ctx)
		if err != nil {
			return err
		}
		err = r.insertAreaRelationshipTestData(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RDS) insertAreaTypeTestData(ctx context.Context) error {
	areaTypeData := DBRelationalData.AreaTypeData
	executionList := make([]string, len(areaTypeData))
	// build queries in order - only insert if doesn't already exist in area_type table
	for key := range areaTypeData {
		executionList[areaTypeData[key]["creation_order"].(int)] = areaTypeData[key]["values"].(string)
	}
	// execute queries
	for index := range executionList {
		logData := log.Data{"exceuting query": executionList[index]}
		rows, err := r.conn.Query(
			ctx,
			areaTypeInsertTransaction,
			executionList[index],
			executionList[index],
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		log.Info(ctx, "area_type table query executed successfully:", logData)
	}
	return nil
}

func (r *RDS) upsertAreaTestData(ctx context.Context) error {
	areaData := DBRelationalData.AreaData
	// build queries
	for code := range areaData {
		queryValues := areaData[code]["values"].(map[string]interface{})
		logData := log.Data{"exceuting query": code}

		// handle scenario where dates not set => pointer to sql null
		var active_from *string
		if queryValues["active_from"].(string) != "" {
			af := queryValues["active_from"].(string)
			active_from = &af
		}
		var active_to *string
		if queryValues["active_to"].(string) != "" {
			at := queryValues["active_to"].(string)
			active_to = &at
		}

		rows, err := r.conn.Query(
			ctx,
			areaInsertTransaction,
			code,
			active_from,
			active_to,
			queryValues["area_type_id"].(int),
			queryValues["geometric_area"].(string),
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		log.Info(ctx, "area table query executed successfully:", logData)
	}
	return nil
}

func (r *RDS) Ping(ctx context.Context) error {
	return r.conn.Ping(ctx)
}

func (r *RDS) UpsertArea(ctx context.Context, area models.AreaParams) (bool, error) {
	var areaTypeId int
	var isInserted bool
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return isInserted, fmt.Errorf("failed to start transaction: %+v", err)
	}

	err = tx.QueryRow(ctx, getAreaType, area.AreaType).Scan(&areaTypeId)
	if err != nil {
		return isInserted, fmt.Errorf("failed to get area type: %+v", err)
	}
	areaDetails := []interface{}{area.Code, area.ActiveFrom, area.ActiveTo, area.GeometricData, areaTypeId, area.Visible}

	err = tx.QueryRow(ctx, upsertArea, areaDetails...).Scan(&isInserted)

	if err != nil {
		tx.Rollback(ctx)
		return isInserted, fmt.Errorf("failed to upsert into area: %+v", err)
	}

	_, err = tx.Exec(ctx, upsertAreaName, area.Code, area.AreaName.Name, area.AreaName.ActiveFrom, area.AreaName.ActiveTo)

	if err != nil {
		tx.Rollback(ctx)
		return isInserted, fmt.Errorf("failed to upsert into area_name: %+v", err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		tx.Rollback(ctx)
		return isInserted, fmt.Errorf("failed to commit: %+v", err)
	}
	return isInserted, nil
}

func (r *RDS) insertRelationshipTypeTestData(ctx context.Context) error {
	relationshipTypeData := DBRelationalData.RelationshipTypeData
	executionList := make([]string, len(relationshipTypeData))
	for key := range relationshipTypeData {
		executionList[relationshipTypeData[key]["creation_order"].(int)] = relationshipTypeData[key]["values"].(string)
	}
	// execute queries
	for index := range executionList {
		logData := log.Data{"exceuting query": executionList[index]}
		rows, err := r.conn.Query(
			ctx,
			relationshipTypeInsertTransaction,
			executionList[index],
			executionList[index],
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		log.Info(ctx, "relationship_type table query executed successfully:", logData)
	}
	return nil
}

func (r *RDS) insertAreaNameTestData(ctx context.Context) error {
	areaNameData := DBRelationalData.AreaNameData
	// build queries
	for name := range areaNameData {
		queryValues := areaNameData[name]["values"].(map[string]interface{})
		logData := log.Data{"exceuting query": name}

		// handle scenario where dates not set => pointer to sql null
		var active_from *string
		if queryValues["active_from"].(string) != "" {
			af := queryValues["active_from"].(string)
			active_from = &af
		}
		var active_to *string
		if queryValues["active_to"].(string) != "" {
			at := queryValues["active_to"].(string)
			active_to = &at
		}

		rows, err := r.conn.Query(
			ctx,
			areaNameInsertTransaction,
			queryValues["area_code"],
			name,
			active_from,
			active_to,
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		log.Info(ctx, "area_name table query executed successfully:", logData)
	}
	return nil
}

func (r *RDS) insertAreaRelationshipTestData(ctx context.Context) error {
	areaRelationshipData := DBRelationalData.AreaRelationshipData
	// build queries
	for code := range areaRelationshipData {
		queryValues := areaRelationshipData[code]["values"].(map[string]interface{})
		logData := log.Data{"exceuting query": code}

		rows, err := r.conn.Query(
			ctx,
			areaRelationshipInsertTransaction,
			code,
			queryValues["rel_area_code"],
			queryValues["rel_type_id"],
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		log.Info(ctx, "area_relationship table query executed successfully:", logData)
	}
	return nil
}
