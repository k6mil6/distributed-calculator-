package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"github.com/k6mil6/distributed-calculator/internal/timeout"
)

type TimeoutsStorage struct {
	db *sqlx.DB
}

func NewTimeoutsStorage(db *sqlx.DB) *TimeoutsStorage {
	return &TimeoutsStorage{
		db: db,
	}
}

func (s *TimeoutsStorage) Save(context context.Context, timeouts model.Timeouts) (int, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	query := `INSERT INTO timeouts (user_id, timeouts_values) VALUES ($1, $2) RETURNING id`

	var id int

	if err := conn.QueryRowContext(
		context,
		query,
		timeouts.UserID,
		timeouts.Timeouts,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *TimeoutsStorage) GetActualTimeouts(context context.Context, userID int64) (model.Timeouts, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return model.Timeouts{}, err
	}
	defer conn.Close()

	var timeouts dbTimeouts

	if err := conn.GetContext(
		context,
		&timeouts,
		`SELECT timeouts_values FROM timeouts WHERE user_id = $1 ORDER BY id DESC LIMIT 1`,
		userID,
	); err != nil {
		return model.Timeouts{}, err
	}
	return model.Timeouts(timeouts), nil
}

type dbTimeouts struct {
	ID       int             `db:"id"`
	UserID   int64           `db:"user_id"`
	Timeouts timeout.Timeout `db:"timeouts_values"`
}
