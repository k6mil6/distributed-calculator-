package http

import (
	"context"
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/internal/model"
)

type Auth interface {
	Login(ctx context.Context, login, password string) (string, error)
	Register(ctx context.Context, login, password string) (int, error)
}

type Expression interface {
	Save(ctx context.Context, expression model.Expression) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID, userID int64) (model.Expression, error)
	AllExpressions(ctx context.Context, userID int64) ([]model.Expression, error)
}

type Timeout interface {
	Save(ctx context.Context, timeouts model.Timeouts) (int, error)
	GetActualTimeouts(ctx context.Context, userID int64) (model.Timeouts, error)
}
