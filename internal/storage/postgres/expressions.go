package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/internal/model"
	errs "github.com/k6mil6/distributed-calculator/internal/storage/errors"
	"github.com/k6mil6/distributed-calculator/internal/timeout"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"time"
)

type ExpressionsStorage struct {
	db *sqlx.DB
}

func NewExpressionStorage(db *sqlx.DB) *ExpressionsStorage {
	return &ExpressionsStorage{
		db: db,
	}
}

func (s *ExpressionsStorage) Save(context context.Context, expression model.Expression) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	var timeoutsID int

	tx, err := conn.BeginTx(context, nil)
	if err != nil {
		return err
	}

	if expression.Timeouts == nil {
		err = tx.QueryRowContext(
			context,
			`SELECT id FROM timeouts WHERE user_id = $1 ORDER BY id DESC LIMIT 1`,
			expression.UserID,
		).Scan(&timeoutsID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				tx.Rollback()
				return errs.ErrTimeoutNotFound
			}
			tx.Rollback()
			return err
		}
	} else {
		timeouts, err := json.Marshal(expression.Timeouts)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.QueryRowContext(
			context,
			`INSERT INTO timeouts (timeouts_values, user_id) VALUES ($1, $2) RETURNING id`,
			timeouts,
			expression.UserID,
		).Scan(&timeoutsID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if _, err := tx.ExecContext(
		context,
		`INSERT INTO expressions (id, user_id, expression, timeouts_id) VALUES ($1, $2, $3, $4)`,
		expression.ID,
		expression.UserID,
		expression.Expression,
		timeoutsID,
	); err != nil {
		tx.Rollback()
		return handlePQError(err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func handlePQError(err error) error {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return errs.ErrExpressionInProgress
		}
	}
	return err
}

func (s *ExpressionsStorage) Get(context context.Context, id uuid.UUID) (model.Expression, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return model.Expression{}, err
	}
	defer conn.Close()

	var expression dbExpression

	query := `SELECT e.id, e.user_id, e.expression, e.created_at, e.is_taken, t.timeouts_values, e.is_done, e.result
              FROM expressions AS e
              LEFT JOIN timeouts AS t ON e.timeouts_id = t.id
              WHERE e.id = $1
              ORDER BY e.created_at`

	if err := conn.GetContext(context, &expression, query, id); err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			return model.Expression{}, errs.ErrExpressionNotFound
		}
		return model.Expression{}, err
	}

	return model.Expression(expression), nil
}

func (s *ExpressionsStorage) NonTakenExpressions(context context.Context) ([]model.Expression, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var expressions []dbExpression

	query := `SELECT e.id, e.expression, e.created_at, e.is_taken, t.timeouts_values
              FROM expressions AS e
              LEFT JOIN timeouts AS t ON e.timeouts_id = t.id
              WHERE e.is_taken = false
              ORDER BY e.created_at`

	if err := conn.SelectContext(context, &expressions, query); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return lo.Map(expressions, func(expression dbExpression, _ int) model.Expression {
		return model.Expression(expression)
	}), nil
}

func (s *ExpressionsStorage) TakeExpression(context context.Context, id uuid.UUID) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(
		context,
		`UPDATE expressions SET is_taken = true WHERE id = $1`,
		id,
	)

	return err
}

func (s *ExpressionsStorage) AllExpressionsByUser(context context.Context, userID int64) ([]model.Expression, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var expressions []dbExpression

	if err := conn.SelectContext(context, &expressions, `SELECT id, expression, created_at, is_taken, result, is_done FROM expressions WHERE user_id = $1 ORDER BY created_at`, userID); err != nil {
		return nil, err
	}

	return lo.Map(expressions, func(expression dbExpression, _ int) model.Expression {
		return model.Expression(expression)
	}), nil
}

func (s *ExpressionsStorage) AllExpressions(context context.Context) ([]model.Expression, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var expressions []dbExpression

	if err := conn.SelectContext(context, &expressions, `SELECT id, expression, created_at, is_taken, result, is_done FROM expressions ORDER BY created_at`); err != nil {
		return nil, err
	}

	return lo.Map(expressions, func(expression dbExpression, _ int) model.Expression {
		return model.Expression(expression)
	}), nil
}

func (s *ExpressionsStorage) UpdateResult(context context.Context, id uuid.UUID, result float64) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(
		context,
		`UPDATE expressions SET result = $1, is_done = true WHERE id = $2`,
		result,
		id,
	)

	return err
}

func (s *ExpressionsStorage) Close() error {
	return s.db.Close()
}

type dbExpression struct {
	ID         uuid.UUID       `db:"id"`
	UserID     int64           `db:"user_id"`
	Expression string          `db:"expression"`
	CreatedAt  time.Time       `db:"created_at"`
	Timeouts   timeout.Timeout `db:"timeouts_values"`
	IsTaken    bool            `db:"is_taken"`
	Result     float64         `db:"result"`
	IsDone     bool            `db:"is_done"`
}
