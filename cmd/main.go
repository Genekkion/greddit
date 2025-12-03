package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	authdb "greddit/internal/infra/db/postgres/auth"

	servicesauth "greddit/internal/services/auth"

	"greddit/internal/infra/auth/local/hs256"

	httpserver "greddit/internal/infra/http/server"

	"greddit/internal/infra/http/routing"

	"greddit/internal/env"
	"greddit/internal/infra/db/postgres"
	"greddit/internal/infra/log"
)

var (
	isDev          = env.GetBoolEnvDef("IS_DEV", false)
	httpAddr       = env.GetStringEnvDef("HTTP_ADDR", "127.0.0.1:3000")
	allowedOrigins = strings.TrimSpace(env.GetStringEnvDef("ALLOWED_ORIGINS", "*"))

	pgConnStr = env.GetStringEnvOrFatal("PGSQL_CONN_STR")
)

func main() {
	defer log.CleanupDefaultLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	routingParam := routing.RouterParams{
		IsDev: isDev,
	}

	logger := log.NewLogger(log.NewHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	routingParam.Logger = logger

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

	{
		secret, err := hs256.NewSecret()
		if err != nil {
			logger.Error("Error generating secret",
				"error",
				err,
			)
			os.Exit(exitKeyFailure)
		}
		jwkSource, err := hs256.NewSource(secret)
		if err != nil {
			logger.Error("Error creating jwk source",
				"error",
				err,
			)
			os.Exit(exitKeyFailure)
		}
		users := authdb.NewUsersRepo(pool)
		ser := servicesauth.NewService(logger, jwkSource, users)
		routingParam.AuthSer = &ser
	}

	err = routingParam.Validate()
	if err != nil {
		logger.Error("Error validating router params",
			"error",
			err,
		)
		os.Exit(exitRoutingParamValidationFailure)
	}

	wg := sync.WaitGroup{}
	errCh := make(chan error, 10)

	wg.Go(func() {
		opts := []httpserver.Option{
			httpserver.WithAddress(httpAddr),
			httpserver.WithAllowedOrigins(allowedOrigins),
		}
		svr, err := httpserver.New(routingParam, opts...)
		if err != nil {
			errCh <- err
			return
		}
		errCh <- svr.Start(ctx)
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.Info("Received os signal to terminate, shutting down services",
			"signal", sig,
		)
	case err := <-errCh:
		logger.Info("Received error from a service, shutting down all services",
			"error", err,
		)
	}

	cancel()
	wg.Wait()
}

const (
	exitUnknownError = iota
	exitDbFailure
	exitRoutingParamValidationFailure
	exitKeyFailure
)
