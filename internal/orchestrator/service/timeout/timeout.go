package timeout

import (
	"context"
	"fmt"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"log/slog"
)

type Timeout struct {
	log      *slog.Logger
	saver    Saver
	provider Provider
	secret   string
}

type Saver interface {
	Save(context context.Context, timeouts model.Timeouts) (int, error)
}

type Provider interface {
	GetActualTimeouts(context context.Context, userID int64) (model.Timeouts, error)
}

func New(
	log *slog.Logger,
	saver Saver,
	provider Provider,
) *Timeout {
	return &Timeout{
		log:      log,
		saver:    saver,
		provider: provider,
	}
}

func (t *Timeout) Save(ctx context.Context, timeouts model.Timeouts) (int, error) {
	const op = "Timeout.Save"

	log := t.log.With(
		slog.String("op", op),
	)

	log.Info("saving timeouts")

	id, err := t.saver.Save(ctx, timeouts)
	if err != nil {
		log.Error("failed to save timeouts", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (t *Timeout) GetActualTimeouts(ctx context.Context, userID int64) (model.Timeouts, error) {
	const op = "Timeout.GetActualTimeouts"

	log := t.log.With(
		slog.String("op", op),
	)

	log.Info("getting actual timeouts")

	timeouts, err := t.provider.GetActualTimeouts(ctx, userID)
	if err != nil {
		log.Error("failed to get actual timeouts", err)
		return model.Timeouts{}, fmt.Errorf("%s: %w", op, err)
	}

	return timeouts, nil
}
