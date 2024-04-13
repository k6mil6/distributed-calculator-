package finaliser

import (
	"context"
	"github.com/google/uuid"
	"github.com/k6mil6/distributed-calculator/internal/model"
	sb "github.com/k6mil6/distributed-calculator/lib/subexpressions/remaker"
	"log/slog"
)

type SubexpressionResultSaver interface {
	DoneSubexpressions(context context.Context) ([]model.Subexpression, error)
	SubexpressionByDependableId(context context.Context, id int) ([]model.Subexpression, error)
	Delete(context context.Context, id int) error
	Save(context context.Context, subExpression model.Subexpression) error
	CompleteSubexpression(context context.Context, id uuid.UUID) (model.Subexpression, error)
}

type ExpressionProvider interface {
	AllExpressions(context context.Context) ([]model.Expression, error)
	UpdateResult(context context.Context, id uuid.UUID, result float64) error
}

type Finaliser struct {
	logger                   *slog.Logger
	subexpressionResultSaver SubexpressionResultSaver
	expressionProvider       ExpressionProvider
}

func New(logger *slog.Logger, subexpressionResultSaver SubexpressionResultSaver, expressionProvider ExpressionProvider) *Finaliser {
	return &Finaliser{
		logger:                   logger,
		subexpressionResultSaver: subexpressionResultSaver,
		expressionProvider:       expressionProvider,
	}
}

func (f *Finaliser) Start(ctx context.Context, ch chan bool) {
	f.logger.Info("finaliser started")

	for {
		select {
		case <-ch:
			f.finalise(ctx)
		case <-ctx.Done():
			f.logger.Info("finaliser stopped")
			return
		}
	}
}

func (f *Finaliser) finalise(ctx context.Context) {
	op := "finaliser.finalise"
	log := f.logger.With(
		slog.String("op", op),
	)

	log.Info("finalising")

	doneSubexpressions, err := f.subexpressionResultSaver.DoneSubexpressions(ctx)
	if err != nil {
		log.Error("error getting done subexpressions:", err)

		return
	}

	log.Info("done subexpressions found", doneSubexpressions)

	for _, subexpression := range doneSubexpressions {
		dependableSubexpressions, err := f.subexpressionResultSaver.SubexpressionByDependableId(ctx, subexpression.ID)
		if err != nil {
			log.Error("error getting subexpression:", err)

			return
		}

		log.Info("dependable subexpressions found", dependableSubexpressions)

		for _, dependableSubexpression := range dependableSubexpressions {
			remadeSubexpression := sb.Remake(dependableSubexpression.Subexpression, subexpression.ID, subexpression.Result)
			dependableSubexpression.Subexpression = remadeSubexpression

			indexToDelete := -1
			for i, value := range dependableSubexpression.DependsOn {
				if value == subexpression.ID {
					indexToDelete = i
					break
				}
			}

			if indexToDelete != -1 {
				dependableSubexpression.DependsOn = append(dependableSubexpression.DependsOn[:indexToDelete], dependableSubexpression.DependsOn[indexToDelete+1:]...)
			}

			if len(dependableSubexpression.DependsOn) == 0 {
				dependableSubexpression.DependsOn = nil
			}

			log.Info("subexpression remade", slog.Any("subexpression", dependableSubexpression.Subexpression), slog.Any("depends on", dependableSubexpression.DependsOn))

			if err := f.subexpressionResultSaver.Delete(ctx, dependableSubexpression.ID); err != nil {
				log.Error("error deleting subexpression:", err)

				return
			}

			if err := f.subexpressionResultSaver.Save(ctx, dependableSubexpression); err != nil {
				log.Error("error saving subexpression:", err)

				return
			}
		}

		f.convertSubexpressionToExpression(ctx)
	}
}

func (f *Finaliser) convertSubexpressionToExpression(ctx context.Context) {
	op := "finaliser.convertSubexpressionToExpression"

	log := f.logger.With(
		slog.String("op", op),
	)

	expressions, err := f.expressionProvider.AllExpressions(ctx)
	if err != nil {
		log.Error("error getting all expressions:", err)

		return
	}

	for _, expression := range expressions {
		if expression.IsDone {
			continue
		}

		subexp, err := f.subexpressionResultSaver.CompleteSubexpression(ctx, expression.ID)
		if err != nil {
			log.Error("error completing subexpression", err)

			continue
		}

		if !subexp.IsDone {
			log.Info("expression is not done", slog.Any("expression", expression.ID), slog.Any("is_done", subexp.IsDone))

			continue
		}

		expression.Result = subexp.Result
		expression.IsDone = true

		if err := f.expressionProvider.UpdateResult(ctx, expression.ID, expression.Result); err != nil {
			log.Error("error updating expression result:", err)

			return
		}

		log.Info("expression result updated", slog.Any("expression", expression))
	}
}
