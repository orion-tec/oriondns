package db

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/orion-tec/oriondns/config"
)

type DB struct {
	*pgxpool.Pool
}

func New(cfg *config.Config) *DB {

	// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10&pool_max_conn_lifetime=1h30m
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?",
		cfg.DB.User,
		url.QueryEscape(os.Getenv("DB_PASSWORD")),
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	return &DB{
		conn,
	}
}

func NewWithPool(pool *pgxpool.Pool) *DB {
	return &DB{
		pool,
	}
}
