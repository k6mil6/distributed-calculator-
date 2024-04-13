package calculate

import (
	"context"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/internal/model"
	orchestratorhttp "github.com/k6mil6/distributed-calculator/internal/orchestrator/http"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	"github.com/k6mil6/distributed-calculator/internal/timeout"
	"log/slog"
	"net/http"
)

type Request struct {
	Id         uuid.UUID       `json:"id"`
	Expression string          `json:"expression"`
	Timeouts   timeout.Timeout `json:"timeouts,omitempty"`
}

type Response struct {
	resp.Response
	Id uuid.UUID `json:"id"`
}

func New(ctx context.Context, log *slog.Logger, expression orchestratorhttp.Expression) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.expression.calculate.New"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("error decoding JSON request:", err)

			render.JSON(w, r, resp.Error("error decoding JSON request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		id, err := expression.Save(ctx, model.Expression{
			ID:         req.Id,
			Expression: req.Expression,
			Timeouts:   req.Timeouts,
		})
		if err != nil {
			log.Error("error saving expression:", err)

			render.JSON(w, r, resp.Error("error saving expression"))

			return
		}

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Id:       id,
	})
}
