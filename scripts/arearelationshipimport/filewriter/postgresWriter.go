package filewriter

import (
	"arearelationshipimport/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

type PostgresWriter struct {
	config *config.ImportConfig
}

func NewPostgresWriter(config *config.ImportConfig) FileWriter {
	return &PostgresWriter{config: config}
}

func (p PostgresWriter) Write(ctx context.Context, file *os.File, tbl string) (int64, error) {
	dbconn, err := pgx.Connect(ctx, p.config.GetDatabaseURI())
	if err != nil {
		return 0, err
	}
	defer dbconn.Close(ctx)

	query := fmt.Sprintf("COPY %s FROM STDIN DELIMITER '|' CSV HEADER", tbl)
	res, err := dbconn.PgConn().CopyFrom(ctx, file, query)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}
