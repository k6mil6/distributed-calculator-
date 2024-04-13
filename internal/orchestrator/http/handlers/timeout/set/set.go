package set

import (
	"context"
	"github.com/go-chi/render"
	"github.com/k6mil6/distributed-calculator/internal/model"
	orchestratorhttp "github.com/k6mil6/distributed-calculator/internal/orchestrator/http"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	"github.com/k6mil6/distributed-calculator/internal/timeout"
	"log/slog"
	"net/http"
)

type Request struct {
	Timeouts timeout.Timeout `json:"timeouts"`
}

func New(ctx context.Context, log *slog.Logger, timeout orchestratorhttp.Timeout) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.timeouts.set_timeouts.New"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("error decoding JSON request:", err)

			render.JSON(w, r, resp.Error("error decoding JSON request"))

			return
		}

		if _, err := timeout.Save(ctx, model.Timeouts{
			Timeouts: req.Timeouts,
		}); err != nil {
			log.Error("error setting timeouts:", err)

			render.JSON(w, r, resp.Error("error setting timeouts"))

			return

		}

		render.JSON(w, r, resp.OK())
	}
}
