package all

import (
	"context"
	"github.com/go-chi/render"
	orchestratorhttp "github.com/k6mil6/distributed-calculator/internal/orchestrator/http"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	"log/slog"
	"net/http"
)

func New(ctx context.Context, logger *slog.Logger, expression orchestratorhttp.Expression) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.expression.all_expressions.New"

		logger = logger.With(
			slog.String("op", op),
		)

		expressions, err := expression.AllExpressions(ctx)
		if err != nil {
			logger.Error("error getting all expressions:", err)

			render.JSON(w, r, resp.Error("error getting all expressions"))

			return
		}

		render.JSON(w, r, expressions)
	}
}
