package get

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	orchestratorhttp "github.com/k6mil6/distributed-calculator/internal/orchestrator/http"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http/middleware/user/identity"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	expressionService "github.com/k6mil6/distributed-calculator/internal/orchestrator/service/expression"
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

		userID := identity.GetUserID(r.Context())

		expression, err := expression.Get(ctx, id, userID)
		if err != nil {
			if errors.Is(err, expressionService.ErrExpressionNotBelongsToUser) {
				log.Info("expression not belongs to user", slog.Any("id", id))

				render.JSON(w, r, resp.Error("expression does not belongs to user"))

				return
			}

			log.Error("error getting expression:", err)

			render.JSON(w, r, resp.Error("error getting expression"))

			return
		}

		render.JSON(w, r, expression)
	}
}
