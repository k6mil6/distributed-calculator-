package orchestrator

import (
	"context"
	"fmt"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"github.com/k6mil6/distributed-calculator/pkg/subexpression_remaker"
	"log/slog"
)

type Orchestrator struct {
	log *slog.Logger

	subexpressionProvider SubexpressionProvider
	subexpressionSaver    SubexpressionSaver
}

type SubexpressionSaver interface {
	TakeSubexpression(context context.Context, id int) (int, error)
	SubexpressionIsDone(context context.Context, id int, result float64) error
	Save(context context.Context, subExpression model.Subexpression) error
	Delete(context context.Context, id int) error
}

type SubexpressionProvider interface {
	NonTakenSubexpressions(context context.Context) ([]model.Subexpression, error)
	DoneSubexpressions(context context.Context) ([]model.Subexpression, error)
	SubexpressionByDependableId(context context.Context, id int) ([]model.Subexpression, error)
}

func New(log *slog.Logger, subexpressionProvider SubexpressionProvider, subexpressionSaver SubexpressionSaver) *Orchestrator {
	return &Orchestrator{
		log:                   log,
		subexpressionProvider: subexpressionProvider,
		subexpressionSaver:    subexpressionSaver,
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
		return model.Subexpression{}, nil
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

	log.Info("expression result saved", slog.Int("id", subexpressionID))

	doneSubexpressions, err := o.subexpressionProvider.DoneSubexpressions(ctx)
	if err != nil {
		log.Error("error getting done subexpressions", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	for _, subexpression := range doneSubexpressions {
		dependableSubexpressions, err := o.subexpressionProvider.SubexpressionByDependableId(ctx, subexpression.ID)
		if err != nil {
			log.Error("error getting subexpression", err)

			return 0, fmt.Errorf("%s: %w", op, err)

		}

		for _, dependableSubexpression := range dependableSubexpressions {
			remadeSubexpression := subexpression_remaker.Remake(dependableSubexpression.Subexpression, subexpression.ID, subexpression.Result)
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

			if err := o.subexpressionSaver.Delete(ctx, dependableSubexpression.ID); err != nil {
				log.Error("error deleting subexpression:", err)

				return 0, fmt.Errorf("%s: %w", op, err)

			}

			if err := o.subexpressionSaver.Save(ctx, dependableSubexpression); err != nil {
				log.Error("error saving subexpression:", err)

				return 0, fmt.Errorf("%s: %w", op, err)
			}
		}

	}

	return subexpressionID, nil
}
