package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/internal/model"
)

type UsersStorage struct {
	db *sqlx.DB
}

func NewUsersStorage(db *sqlx.DB) *UsersStorage {
	return &UsersStorage{
		db: db,
	}
}

func (s *UsersStorage) Save(context context.Context, user model.User) error {
	conn, err := s.db.Connx(context)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		context,
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)`,
		user.Login,
		user.PasswordHash,
	); err != nil {
		return err
	}

	return nil
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
		`SELECT id, username, password_hash FROM users WHERE username = $1`,
		login,
	); err != nil {
		return model.User{}, err
	}

	return model.User(user), nil
}

type dbUser struct {
	ID           int64  `db:"id"`
	Login        string `db:"login"`
	PasswordHash string `db:"password_hash"`
}
