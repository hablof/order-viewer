package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connectTimeout time.Duration = 5 * time.Second
)

func NewPostgres(ctx context.Context, url string) (*pgxpool.Pool, error) {

	ctx, cf := context.WithTimeout(ctx, connectTimeout)
	defer cf()

	conn, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
