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
		"POSTGRES_USER":     "bastion",
		"POSTGRES_PASSWORD": "bastion",
		"POSTGRES_DB":       "bastion",
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

	connStr := fmt.Sprintf("postgresql://bastion:bastion@%s:%s/bastion?sslmode=disable", host, port.Port())

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
