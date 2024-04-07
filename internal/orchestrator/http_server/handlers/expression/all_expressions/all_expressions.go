package all_expressions

import (
	"context"
	"github.com/go-chi/render"
	"github.com/k6mil6/distributed-calculator/internal/model"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	"log/slog"
	"net/http"
)

type ExpressionsSelector interface {
	AllExpressions(context context.Context) ([]model.Expression, error)
}

func New(logger *slog.Logger, expressionsSelector ExpressionsSelector, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.expression.all_expressions.New"

		logger = logger.With(
			slog.String("op", op),
		)

		expressions, err := expressionsSelector.AllExpressions(ctx)
		if err != nil {
			logger.Error("error getting all expressions:", err)

			render.JSON(w, r, resp.Error("error getting all expressions"))

			return
		}

		render.JSON(w, r, expressions)
	}
}
