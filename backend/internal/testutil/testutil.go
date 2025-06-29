package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

const (
	TestDBName = "oriondns_test"
	TestDBUser = "postgres"
	TestDBHost = "localhost"
	TestDBPort = "5432"
)

func GetTestDBConfig() *pgxpool.Config {
	dbURL := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable",
		TestDBUser, TestDBHost, TestDBPort, TestDBName)

	if envURL := os.Getenv("TEST_DATABASE_URL"); envURL != "" {
		dbURL = envURL
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse test database config: %v", err))
	}

	return config
}

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	config := GetTestDBConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err, "Failed to connect to test database")

	err = pool.Ping(ctx)
	require.NoError(t, err, "Failed to ping test database")

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TruncateAllTables(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()

	tables := []string{
		"stats",
		"blocked_domains",
		"categories",
		"domains",
	}

	for _, table := range tables {
		_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		require.NoError(t, err, "Failed to truncate table %s", table)
	}
}

func CreateTestDB(t *testing.T) {
	adminURL := fmt.Sprintf("postgres://%s@%s:%s/postgres?sslmode=disable",
		TestDBUser, TestDBHost, TestDBPort)

	if envURL := os.Getenv("TEST_DATABASE_URL"); envURL != "" {
		adminURL = envURL
	}

	ctx := context.Background()
	conn, err := pgconn.Connect(ctx, adminURL)
	require.NoError(t, err, "Failed to connect to PostgreSQL")
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", TestDBName)).ReadAll()
	require.NoError(t, err, "Failed to drop test database")

	_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", TestDBName)).ReadAll()
	require.NoError(t, err, "Failed to create test database")
}

func SkipIfNoDatabase(t *testing.T) {
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Skipping database tests (SKIP_DB_TESTS=true)")
	}
}
