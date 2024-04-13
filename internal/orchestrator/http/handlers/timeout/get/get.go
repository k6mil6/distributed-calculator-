package get

import (
	"context"
	"github.com/go-chi/render"
	"github.com/k6mil6/distributed-calculator/internal/model"
	orchestratorhttp "github.com/k6mil6/distributed-calculator/internal/orchestrator/http"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	"log/slog"
	"net/http"
)

type Response struct {
	Timeouts model.Timeouts `json:"timeouts"`
}

func New(ctx context.Context, log *slog.Logger, timeout orchestratorhttp.Timeout) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.timeouts.actual_timeouts.New"

		log = log.With(
			slog.String("op", op),
		)

		timeouts, err := timeout.GetActualTimeouts(ctx)
		if err != nil {
			log.Error("error getting actual timeouts:", err)

			render.JSON(w, r, resp.Error("error getting actual timeouts"))

			return
		}

		render.JSON(w, r, Response{
			Timeouts: timeouts,
		})
	}
}
