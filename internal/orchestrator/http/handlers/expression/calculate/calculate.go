package calculate

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/internal/model"
	orchestratorhttp "github.com/k6mil6/distributed-calculator/internal/orchestrator/http"
	"github.com/k6mil6/distributed-calculator/internal/orchestrator/http/middleware/user/identity"
	resp "github.com/k6mil6/distributed-calculator/internal/orchestrator/response"
	expressionService "github.com/k6mil6/distributed-calculator/internal/orchestrator/service/expression"
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

		err := handleEmptyParameters(req)
		if err != nil {
			log.Error("error validating parameters:", err)

			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		userID := identity.GetUserID(r.Context())

		id, err := expression.Save(ctx, model.Expression{
			ID:         req.Id,
			UserID:     userID,
			Expression: req.Expression,
			Timeouts:   req.Timeouts,
		})

		if err != nil {
			handleErrors(w, r, err, log)

			fmt.Println(err)

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

func handleEmptyParameters(req Request) error {
	if req.Id == uuid.Nil {
		return errors.New("id is empty")
	}

	if req.Expression == "" {
		return errors.New("expression is empty")
	}

	return nil
}

func handleErrors(w http.ResponseWriter, r *http.Request, err error, log *slog.Logger) {
	if errors.Is(err, expressionService.ErrExpressionInProgress) {
		log.Error("expression is in progress")

		render.JSON(w, r, resp.Error("expression is in progress"))

		return
	}

	if errors.Is(err, expressionService.ErrTimeoutNotFound) {
		log.Error("timeout not found")

		render.JSON(w, r, resp.Error("timeout not found, add it in the request"))

		return
	}
}
