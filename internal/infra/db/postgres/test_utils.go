package postgres

import (
	"fmt"
	"greddit/internal/test"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	postgresEnv = map[string]string{
		"POSTGRES_USER":     "greddit",
		"POSTGRES_PASSWORD": "greddit",
		"POSTGRES_DB":       "greddit",
	}
	tables = []string{
		"auth_users",
	}
)

// NewTestPool creates a new pgxpool.Pool for testing.
func NewTestPool(t *testing.T) (pool *pgxpool.Pool, cleanup func()) {
	t.Helper()

	ctx := t.Context()
	req := testcontainers.ContainerRequest{
		Image: "postgres:18.0-trixie",
		Env:   postgresEnv,
		ExposedPorts: []string{
			"5432/tcp",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(
				"5432/tcp",
			).WithStartupTimeout(30*time.Second),
			wait.ForLog(
				"database system is ready to accept connections",
			),
		),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	test.NilErr(t, err)

	cleanup = func() {
		testcontainers.CleanupContainer(t, container)
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
