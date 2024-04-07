package expression

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/internal/model"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	"log/slog"
	"net/http"
)

type ExpressionsSelector interface {
	Get(context context.Context, id uuid.UUID) (model.Expression, error)
}

func New(logger *slog.Logger, expressionsSelector ExpressionsSelector, context context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.expression.all_expressions.New"

		logger = logger.With(
			slog.String("op", op),
		)

		urlParam := chi.URLParam(r, "id")
		if urlParam == "" {
			logger.Error("no id")

			render.JSON(w, r, resp.Error("no id"))

			return
		}
		id, err := uuid.Parse(urlParam)
		if err != nil {
			logger.Error("invalid id")

			render.JSON(w, r, resp.Error("invalid id"))

			return
		}

		expression, err := expressionsSelector.Get(context, id)
		if err != nil {
			logger.Error("error getting expression:", err)

			render.JSON(w, r, resp.Error("error getting expression"))

			return
		}

		render.JSON(w, r, expression)
	}
}
