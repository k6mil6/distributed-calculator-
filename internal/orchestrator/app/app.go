package app

import (
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/app/grpc"
	orchestratorService "github.com/k6mil6/distributed-calculator/internal/orchestrator/service/orchestrator"
	"github.com/k6mil6/distributed-calculator/internal/storage/postgres"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	subExpressionStorage *postgres.SubexpressionStorage,
	ch chan bool,
) *App {
	orchestrator := orchestratorService.New(log, subExpressionStorage, subExpressionStorage, ch)

	grpcApp := grpcapp.New(log, grpcPort, orchestrator)

	return &App{
		GRPCServer: grpcApp,
	}
}
