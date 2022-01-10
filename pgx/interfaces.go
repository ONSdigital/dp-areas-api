package pgx

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate moq -out mock/pgx.go -pkg mock . PGXPool

type PGXPool interface {
	Connect(ctx context.Context, connString string) (*pgxpool.Pool, error)
	Close()
}
