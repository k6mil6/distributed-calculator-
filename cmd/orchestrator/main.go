package main

import (
	"context"
	"github.com/k6mil6/distributed-calculator/internal/config"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/app"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/checker"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/fetcher"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/finaliser"
	"github.com/k6mil6/distributed-calculator/internal/storage"
	"github.com/k6mil6/distributed-calculator/lib/logger"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Get()
	log := logger.SetupLogger(cfg.Env).With(slog.String("env", cfg.Env))
	log.Debug("logger debug mode enabled")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	storages, err := storage.New(cfg.PostgresDatabaseDSN, cfg.RedisDatabaseDSN, cfg.DBRetriesNumber, cfg.DBRetryCooldown)
	if err != nil {
		log.Error("failed to connect to database", err)

		return
	}

	defer func() {
		if err := storages.CloseAll(); err != nil {
			log.Error("failed to close storages", err)
		}
	}()

	f := fetcher.New(storages.ExpressionsStorage, storages.SubexpressionsStorage, cfg.FetcherInterval, log)
	c := checker.New(storages.SubexpressionsStorage, cfg.CheckerInterval, log)
	fin := finaliser.New(log, storages.SubexpressionsStorage, storages.ExpressionsStorage)

	ch := make(chan bool, 1)

	go f.Start(ctx)
	go c.Start(ctx)
	go fin.Start(ctx, ch)

	application := app.New(
		ctx,
		log,
		cfg.GrpcPort,
		storages,
		cfg.TokenTTL,
		cfg.Secret,
		cfg.HttpPort,
		ch,
	)

	go func() {
		application.GRPCServer.MustRun()
	}()

	<-ctx.Done()

	application.GRPCServer.Stop()
	log.Info("gRPC server stopped")
	log.Info("server stopped")
}
