package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	connectTimeout time.Duration = 5 * time.Second
)

func NewPostgres(ctx context.Context, url string) (*pgx.Conn, error) {

	ctx, cf := context.WithTimeout(ctx, connectTimeout)
	defer cf()

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
