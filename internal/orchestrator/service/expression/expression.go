package expression

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/internal/model"
	errs "github.com/k6mil6/distributed-calculator/internal/storage/errors"
	"github.com/k6mil6/distributed-calculator/lib/validation"
	"log/slog"
)

var (
	ErrExpressionNotValid   = errors.New("expression not valid")
	ErrExpressionInProgress = errors.New("expression in progress")
)

type Expression struct {
	log      *slog.Logger
	saver    Saver
	provider Provider
}

type Saver interface {
	Save(context context.Context, expression model.Expression) error
}

type Provider interface {
	Get(context context.Context, id uuid.UUID) (model.Expression, error)
	AllExpressions(context context.Context) ([]model.Expression, error)
}

func New(log *slog.Logger, saver Saver, provider Provider) *Expression {
	return &Expression{
		log:      log,
		saver:    saver,
		provider: provider,
	}
}

func (e *Expression) Save(ctx context.Context, expression model.Expression) (uuid.UUID, error) {
	const op = "Expression.Save"

	log := e.log.With(
		slog.String("op", op),
	)

	if !validation.IsMathExpressionValid(expression.Expression) {
		return uuid.UUID{}, ErrExpressionNotValid
	}

	log.Info("math expression is valid")

	err := e.saver.Save(ctx, model.Expression{
		ID:         expression.ID,
		Expression: expression.Expression,
		Timeouts:   expression.Timeouts,
	})

	if err != nil {
		if errors.Is(err, errs.ErrExpressionInProgress) {
			return uuid.UUID{}, ErrExpressionInProgress
		}
		log.Error("error saving expression:", err)

		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("expression saved successfully: ", expression.ID)

	return expression.ID, nil
}

func (e *Expression) Get(ctx context.Context, id uuid.UUID) (model.Expression, error) {
	const op = "Expression.Get"
	log := e.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to get expression")

	expression, err := e.provider.Get(ctx, id)
	if err != nil {
		log.Error("error getting expression:", err)
		return model.Expression{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("expression found")

	return expression, nil
}

func (e *Expression) AllExpressions(ctx context.Context) ([]model.Expression, error) {
	const op = "Expression.AllExpressions"
	log := e.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to get all expressions")

	expressions, err := e.provider.AllExpressions(ctx)
	if err != nil {
		log.Error("error getting all expressions:", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("all expressions found")

	return expressions, nil
}
