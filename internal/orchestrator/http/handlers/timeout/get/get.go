package get

import (
	"context"
	"github.com/go-chi/render"
	orchestratorhttp "github.com/k6mil6/distributed-calculator/internal/orchestrator/http"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http/middleware/user/identity"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	"log/slog"
	"net/http"
)

func New(ctx context.Context, log *slog.Logger, timeout orchestratorhttp.Timeout) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.timeouts.actual_timeouts.New"

		log = log.With(
			slog.String("op", op),
		)

		userID := identity.GetUserID(r.Context())

		timeouts, err := timeout.GetActualTimeouts(ctx, userID)
		if err != nil {
			log.Error("error getting actual timeouts:", err)

			render.JSON(w, r, resp.Error("error getting actual timeouts"))

			return
		}

		render.JSON(w, r, timeouts)
	}
}
