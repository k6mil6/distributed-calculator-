package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/backend/internal/config"
	mwlogger "github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/http-server/middleware/logger"
	"github.com/k6mil6/distributed-calculator/backend/internal/storage/migrations"
	"github.com/k6mil6/distributed-calculator/backend/pkg/logger"
	_ "github.com/lib/pq"
	"log/slog"
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

	if err := migrations.Start(cfg.DatabaseDSN); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Error("failed to start migrations", err)
			return
		}
	}

	router := chi.NewRouter()

	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

}
