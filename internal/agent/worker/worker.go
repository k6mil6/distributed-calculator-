package worker

import (
	"context"
	"github.com/k6mil6/distributed-calculator/internal/agent/evaluator"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"log/slog"
	"time"
)

type Communicator interface {
	GetFreeExpressions(ctx context.Context) (model.Subexpression, error)
	SaveResult(ctx context.Context, subexpressionID int, result float64) (int, error)
	SendHeartbeat(ctx context.Context, workerID int) error
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

		ch := make(chan int)
		go func() {
			id := <-ch
			err = communicator.SendHeartbeat(ctx, id)
			if err != nil {
				w.log.Error("error sending heartbeat", err)
			}
		}()

		result, err := evaluator.Evaluate(resp, w.timeout, ch, w.id, w.log)
		if err != nil {
			w.log.Error("error evaluating expression", err)
			time.Sleep(w.timeout)

			continue
		}

		w.log.Info("evaluated expression", slog.Int("id", resp.ID), slog.Float64("result", result.Result))

		_, err = communicator.SaveResult(ctx, resp.ID, result.Result)

		if err != nil {
			w.log.Error("error evaluating expression", err)
			time.Sleep(w.timeout)

			continue
		}
	}
}
