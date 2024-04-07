package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/internal/config"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/app"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/checker"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/fetcher"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/finaliser"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http_server/handlers/expression/all_expressions"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http_server/handlers/expression/calculate"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http_server/handlers/expression/expression"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http_server/handlers/timeouts/actual_timeouts"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http_server/handlers/timeouts/set_timeouts"
	mwlogger "github.com/k6mil6/distributed-calculator/internal/orchestrator/http_server/middleware/logger"
	"github.com/k6mil6/distributed-calculator/internal/storage/postgres"
	"github.com/k6mil6/distributed-calculator/pkg/logger"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Get()
	log := logger.SetupLogger(cfg.Env).With(slog.String("env", cfg.Env))
	log.Debug("logger debug mode enabled")

	db, err := sqlx.Connect("postgres", cfg.DatabaseDSN)
	if err != nil {
		log.Error("failed to connect to database", err)
		return
	}
	defer db.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	expressionStorage := postgres.NewExpressionStorage(db)
	subexpressionStorage := postgres.NewSubexpressionStorage(db)
	timeoutsStorage := postgres.NewTimeoutsStorage(db)

	router := chi.NewRouter()

	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/calculate", calculate.New(log, expressionStorage, ctx))
	router.Post("/set_timeouts", set_timeouts.New(log, timeoutsStorage, ctx))

	router.Get("/all_expressions", all_expressions.New(log, expressionStorage, ctx))
	router.Get("/expression/{id}", expression.New(log, expressionStorage, ctx))
	router.Get("/actual_timeouts", actual_timeouts.New(log, timeoutsStorage, ctx))

	f := fetcher.New(expressionStorage, subexpressionStorage, cfg.FetcherInterval, log)
	c := checker.New(subexpressionStorage, cfg.CheckerInterval, log)
	fin := finaliser.New(log, subexpressionStorage, expressionStorage)

	ch := make(chan bool, 1)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go f.Start(ctx)
	go c.Start(ctx)
	go fin.Start(ctx, ch)

	application := app.New(log, cfg.GrpcPort, subexpressionStorage, ch)

	go func() {
		application.GRPCServer.MustRun()
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", err)
		}
	}()

	log.Info("server started", slog.String("address", srv.Addr))

	<-ctx.Done()

	application.GRPCServer.Stop()
	log.Info("gRPC server stopped")
	log.Info("server stopped")
}
