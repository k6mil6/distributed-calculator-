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
	AllExpressions(ctx context.Context) ([]model.Expression, error)
	Get(ctx context.Context, id uuid.UUID) (model.Expression, error)
}

type Timeout interface {
	GetActualTimeouts(context context.Context) (model.Timeouts, error)
	Save(context context.Context, timeouts model.Timeouts) (int, error)
}
