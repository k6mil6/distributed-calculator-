package app

import (
	"context"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/app/grpc"
	httpapp "github.com/k6mil6/distributed-calculator/internal/orchestrator/app/http"
	authService "github.com/k6mil6/distributed-calculator/internal/orchestrator/service/auth"
	expressionService "github.com/k6mil6/distributed-calculator/internal/orchestrator/service/expression"
	orchestratorService "github.com/k6mil6/distributed-calculator/internal/orchestrator/service/orchestrator"
	timeoutService "github.com/k6mil6/distributed-calculator/internal/orchestrator/service/timeout"
	"github.com/k6mil6/distributed-calculator/internal/storage"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
	HTTPServer *httpapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	grpcPort int,
	storages storage.Storages,
	tokenTTL time.Duration,
	secret string,
	port int,
	ch chan bool,
) *App {
	orchestrator := orchestratorService.New(log, storages.SubexpressionsStorage, storages.SubexpressionsStorage, ch)

	auth := authService.New(log, storages.UsersStorage, storages.UsersStorage, tokenTTL, secret)
	expression := expressionService.New(log, storages.ExpressionsStorage, storages.ExpressionsStorage)
	timeout := timeoutService.New(log, storages.TimeoutsStorage, storages.TimeoutsStorage)

	httpApp := httpapp.New(ctx, log, port, auth, expression, timeout)

	grpcApp := grpcapp.New(log, grpcPort, orchestrator)

	return &App{
		GRPCServer: grpcApp,
		HTTPServer: httpApp,
	}
}
