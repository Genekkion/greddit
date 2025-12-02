package main

import (
	"context"
	"greddit/internal/env"
	"greddit/internal/infra/db/postgres"
	"greddit/internal/infra/log"
	"log/slog"
	"os"
)

var (
	pgConnStr = env.GetStringEnvOrFatal("PGSQL_CONN_STR")
)

func main() {
	defer log.CleanupDefaultLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.NewLogger(log.NewHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	pool, err := postgres.New(ctx, pgConnStr)
	if err != nil {
		logger.Error("Error creating postgres pool",
			"error",
			err,
		)
		os.Exit(exitDbFailure)
	}
	defer pool.Close()

	err = postgres.Init(pool, ctx)
	if err != nil {
		logger.Error("Error initializing postgres pool",
			"error",
			err,
		)
		os.Exit(exitDbFailure)
	}

	select {}
}
