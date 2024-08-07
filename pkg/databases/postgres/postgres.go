package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

var (
	timeout = 10 * time.Second
)

func NewClient(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	fmt.Println("About to ping: ", dsn, "Conn string: ", conn)

	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	log.Println("Connected to Postgres~!")

	return conn, nil
}
