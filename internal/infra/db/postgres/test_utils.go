package postgres

import (
	"fmt"
	"testing"
	"time"

	"greddit/internal/test"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	postgresEnv = map[string]string{
		"POSTGRES_USER":     "greddit",
		"POSTGRES_PASSWORD": "greddit",
		"POSTGRES_DB":       "greddit",
	}
	tables = []string{
		"auth_users",
		"forum_communities",
	}
)

// NewTestPool creates a new pgxpool.Pool for testing.
func NewTestPool(t *testing.T) (pool *pgxpool.Pool, cleanup func()) {
	t.Helper()

	ctx := t.Context()
	container, err := pg.Run(ctx,
		"postgres:18.1-alpine",
		pg.WithUsername(postgresEnv["POSTGRES_USER"]),
		pg.WithPassword(postgresEnv["POSTGRES_PASSWORD"]),
		pg.WithDatabase(postgresEnv["POSTGRES_DB"]),
		pg.BasicWaitStrategies(),
	)
	test.NilErr(t, err)

	cleanup = func() {
		err := testcontainers.TerminateContainer(container)
		if err != nil {
			fmt.Println("failed to terminate container", err)
		}
	}

	host, err := container.Host(ctx)
	if err != nil {
		cleanup()
		t.Fatalf("failed to get host")
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		cleanup()
		t.Fatalf("failed to get port")
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		postgresEnv["POSTGRES_USER"], postgresEnv["POSTGRES_PASSWORD"],
		host, port.Port(), postgresEnv["POSTGRES_DB"],
	)

	timeout := time.After(30 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

loop:
	for {
		select {
		case <-timeout:
			cleanup()
			t.Fatalf("failed to connect to postgres")
		case <-tick:
			fmt.Println("connecting to postgres")
			pool, err = New(ctx, connStr)
			if err != nil {
				fmt.Println("failed to connect to postgres", err)
				continue
			}

			err = Init(pool, ctx)
			if err == nil {
				fmt.Println("successfully connected to postgres")
				break loop
			} else {
				fmt.Println("failed to init postgres", err)
			}
		}
	}

	cleanup = func() {
		pool.Close()
		testcontainers.CleanupContainer(t, container)
	}
	return pool, cleanup
}

// ClearAllTables clears all tables in the database.
func ClearAllTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	for _, table := range tables {
		stmt := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)

		_, err := pool.Exec(t.Context(), stmt)
		test.NilErr(t, err)
	}
}
