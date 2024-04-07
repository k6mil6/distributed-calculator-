package orchestrator

import (
	"context"
	"fmt"
	"github.com/k6mil6/distributed-calculator/internal/errors"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"log/slog"
)

type Orchestrator struct {
	ch  chan bool
	log *slog.Logger

	subexpressionProvider SubexpressionProvider
	subexpressionSaver    SubexpressionSaver
}

type SubexpressionSaver interface {
	TakeSubexpression(context context.Context, id int) (int, error)
	SubexpressionIsDone(context context.Context, id int, result float64) error
}

type SubexpressionProvider interface {
	NonTakenSubexpressions(context context.Context) ([]model.Subexpression, error)
}

func New(log *slog.Logger, subexpressionProvider SubexpressionProvider, subexpressionSaver SubexpressionSaver, ch chan bool) *Orchestrator {
	return &Orchestrator{
		log:                   log,
		subexpressionProvider: subexpressionProvider,
		subexpressionSaver:    subexpressionSaver,
		ch:                    ch,
	}
}

func (o *Orchestrator) GetFreeExpressions(ctx context.Context) (model.Subexpression, error) {
	op := "Orchestrator.GetFreeExpressions"

	log := o.log.With(
		slog.String("op", op),
	)

	latestSubexpressions, err := o.subexpressionProvider.NonTakenSubexpressions(ctx)
	if err != nil {
		log.Error("error getting non taken subexpressions", err)

		return model.Subexpression{}, fmt.Errorf("%s: %w", op, err)
	}

	if len(latestSubexpressions) == 0 {
		return model.Subexpression{}, errors.ErrSubexpressionNotFound
	}

	subexpression := latestSubexpressions[0]

	log.Info("subexpression found", slog.Int("id", subexpression.ID))

	workerId, err := o.subexpressionSaver.TakeSubexpression(ctx, subexpression.ID)
	if err != nil {
		log.Error("error taking subexpression", err)

		return model.Subexpression{}, fmt.Errorf("%s: %w", op, err)
	}

	subexpression.WorkerId = workerId

	log.Info("subexpression taken", slog.Int("id", subexpression.ID))

	return subexpression, nil
}

func (o *Orchestrator) SaveResult(ctx context.Context, subexpressionID int, result float64) (int, error) {
	op := "Orchestrator.SaveResult"

	log := o.log.With(
		slog.String("op", op),
	)

	err := o.subexpressionSaver.SubexpressionIsDone(ctx, subexpressionID, result)
	if err != nil {
		log.Error("error saving expression result", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	o.ch <- true

	log.Info("expression result saved", slog.Int("id", subexpressionID))

	return subexpressionID, nil
}
