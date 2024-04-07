package main

import (
	"context"
	"github.com/k6mil6/distributed-calculator/internal/agent/clients/orchestrator/grpc"
	"github.com/k6mil6/distributed-calculator/internal/agent/worker"
	"github.com/k6mil6/distributed-calculator/internal/config"
	"github.com/k6mil6/distributed-calculator/pkg/logger"
	"log/slog"
	"sync"
	"time"
)

func main() {
	cfg := config.Get()
	log := logger.SetupLogger(cfg.Env)

	log = log.With(slog.String("env", cfg.Env))
	log.Debug("logger debug mode enabled")

	ctx := context.Background()

	client, err := grpc.New(
		ctx,
		log,
		cfg.GRPCServerAddress,
		cfg.GRPCReconnectTimeout,
		cfg.GRPCReconnectRetries,
	)

	if err != nil {
		log.Error("failed to create grpc client", err)
		return
	}

	var wg sync.WaitGroup

	for i := 0; i < cfg.GoroutineNumber; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			w := worker.New(log, cfg.WorkerTimeout)
			w.Start(client, ctx)
		}(i)
		time.Sleep(cfg.WorkerTimeout)
	}

	wg.Wait()
}
