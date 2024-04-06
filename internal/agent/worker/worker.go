package worker

import (
	"context"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"log/slog"
	"time"
)

type Communicator interface {
	GetFreeExpressions(ctx context.Context) (model.Subexpression, error)
	SaveResult(ctx context.Context, subexpressionID int, result float64) (int, error)
	SendHeartbeat(ctx context.Context) error
}

type Worker struct {
	id      int
	timeout time.Duration
	log     *slog.Logger
}

func New(log *slog.Logger, timeout time.Duration) *Worker {
	return &Worker{
		log:     log,
		timeout: timeout,
	}
}

func (w *Worker) Start(communicator Communicator, ctx context.Context) {
	for {
		resp, err := communicator.GetFreeExpressions(ctx)
		if err != nil {
			w.log.Error("error getting free expressions", err)
			time.Sleep(w.timeout)
			continue
		}

		w.log.Info("got free expression", slog.Int("id", resp.ID))
	}
}
