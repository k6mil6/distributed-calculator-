package all

import (
	"context"
	"github.com/go-chi/render"
	orchestratorhttp "github.com/k6mil6/distributed-calculator/internal/orchestrator/http"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http/middleware/user/identity"
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

		userID := identity.GetUserID(r.Context())

		expressions, err := expression.AllExpressions(ctx, userID)
		if err != nil {
			log.Error("error getting all expressions:", err)

			render.JSON(w, r, resp.Error("error getting all expressions"))

			return
		}

		render.JSON(w, r, expressions)
	}
}
