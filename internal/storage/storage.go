package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/internal/storage/postgres"
)

type Storages struct {
	ExpressionsStorage    *postgres.ExpressionsStorage
	SubexpressionsStorage *postgres.SubexpressionsStorage
	TimeoutsStorage       *postgres.TimeoutsStorage
	UsersStorage          *postgres.UsersStorage
}

func New(db *sqlx.DB) Storages {
	return Storages{
		ExpressionsStorage:    postgres.NewExpressionStorage(db),
		SubexpressionsStorage: postgres.NewSubexpressionStorage(db),
		TimeoutsStorage:       postgres.NewTimeoutsStorage(db),
		UsersStorage:          postgres.NewUsersStorage(db),
	}
}
