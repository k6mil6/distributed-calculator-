package heartbeat

import (
	"context"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"log/slog"
)

type Heartbeat struct {
	log   *slog.Logger
	saver Saver
}

type Saver interface {
	SaveHeartbeat(ctx context.Context, workerID int) error
	GetAllHeartbeats(ctx context.Context) ([]model.Heartbeat, error)
}

func New(log *slog.Logger, saver Saver) *Heartbeat {
	return &Heartbeat{
		log:   log,
		saver: saver,
	}
}

func (h *Heartbeat) SaveHeartbeat(ctx context.Context, workerID int) error {
	return h.saver.SaveHeartbeat(ctx, workerID)
}

func (h *Heartbeat) GetAllHeartbeats(ctx context.Context) ([]model.Heartbeat, error) {
	return h.saver.GetAllHeartbeats(ctx)
}
