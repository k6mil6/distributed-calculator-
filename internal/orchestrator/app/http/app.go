package httpapp

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	authHttp "github.com/k6mil6/distributed-calculator/internal/orchestrator/http"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http/handlers/expression/calculate"
	expressionGet "github.com/k6mil6/distributed-calculator/internal/orchestrator/http/handlers/expression/get"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http/handlers/expression/get/all"
	timeoutGet "github.com/k6mil6/distributed-calculator/internal/orchestrator/http/handlers/timeout/get"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http/handlers/timeout/set"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http/handlers/user/login"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http/handlers/user/register"
	mwlogger "github.com/k6mil6/distributed-calculator/internal/orchestrator/http/middleware/logger"
	"log/slog"
	"net/http"
)

type App struct {
	log    *slog.Logger
	router *chi.Mux
	server *http.Server
}

func New(ctx context.Context, log *slog.Logger, port int, auth authHttp.Auth, expression authHttp.Expression, timeout authHttp.Timeout) *App {
	router := chi.NewRouter()

	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/calculate", calculate.New(ctx, log, expression))
	router.Post("/set_timeouts", set.New(ctx, log, timeout))
	router.Post("/register", register.New(ctx, log, auth))
	router.Post("/login", login.New(ctx, log, auth))

	router.Get("/all_expressions", all.New(ctx, log, expression))
	router.Get("/expression/{id}", expressionGet.New(ctx, log, expression))
	router.Get("/actual_timeouts", timeoutGet.New(ctx, log, timeout))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return &App{
		log:    log,
		router: router,
		server: server,
	}
}

func (a *App) Run() error {
	a.log.Info("starting server", slog.String("address", a.server.Addr))
	return a.server.ListenAndServe()
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
