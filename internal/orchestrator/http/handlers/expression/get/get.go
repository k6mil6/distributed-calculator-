package get

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	orchestratorhttp "github.com/k6mil6/distributed-calculator/internal/orchestrator/http"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	"log/slog"
	"net/http"
)

func New(ctx context.Context, log *slog.Logger, expression orchestratorhttp.Expression) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.expression.all_expressions.New"

		log = log.With(
			slog.String("op", op),
		)

		urlParam := chi.URLParam(r, "id")
		if urlParam == "" {
			log.Error("no id")

			render.JSON(w, r, resp.Error("no id"))

			return
		}
		id, err := uuid.Parse(urlParam)
		if err != nil {
			log.Error("invalid id")

			render.JSON(w, r, resp.Error("invalid id"))

			return
		}

		expression, err := expression.Get(ctx, id)
		if err != nil {
			log.Error("error getting expression:", err)

			render.JSON(w, r, resp.Error("error getting expression"))

			return
		}

		render.JSON(w, r, expression)
	}
}
