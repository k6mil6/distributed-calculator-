package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/internal/model"
	errs "github.com/k6mil6/distributed-calculator/internal/storage/errors"
	"github.com/lib/pq"
)

type UsersStorage struct {
	db *sqlx.DB
}

func NewUsersStorage(db *sqlx.DB) *UsersStorage {
	return &UsersStorage{
		db: db,
	}
}

func (s *UsersStorage) Save(context context.Context, user model.User) (int, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	query := `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`

	var id int

	if err := conn.QueryRowContext(
		context,
		query,
		user.Login,
		user.PasswordHash,
	).Scan(&id); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, errs.ErrUserExists
		}
		return 0, err
	}

	return id, nil
}

func (s *UsersStorage) GetByLogin(context context.Context, login string) (model.User, error) {
	conn, err := s.db.Connx(context)
	if err != nil {
		return model.User{}, err
	}
	defer conn.Close()

	var user dbUser

	if err := conn.GetContext(
		context,
		&user,
		`SELECT id, login, password_hash FROM users WHERE login = $1`,
		login,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, errs.ErrUserNotFound
		}
		return model.User{}, err
	}

	return model.User(user), nil
}

func (s *UsersStorage) Close() error {
	return s.db.Close()
}

type dbUser struct {
	ID           int64  `db:"id"`
	Login        string `db:"login"`
	PasswordHash []byte `db:"password_hash"`
}
