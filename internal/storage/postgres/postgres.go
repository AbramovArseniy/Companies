package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConn struct {
	pgxpool.Conn
}

func New(dbAddress string) (*pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), dbAddress)
	if err != nil {
		return nil, fmt.Errorf("error while connecting to database: %w", err)
	}
	return dbPool, err
}
