package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/internal/storage/postgres"
	"github.com/k6mil6/distributed-calculator/internal/storage/redis/heartbeats"
	"io"
	"reflect"
)

type Storages struct {
	ExpressionsStorage    *postgres.ExpressionsStorage
	SubexpressionsStorage *postgres.SubexpressionsStorage
	TimeoutsStorage       *postgres.TimeoutsStorage
	UsersStorage          *postgres.UsersStorage
	HeartbeatsStorage     *heartbeats.Storage
}

func New(postgresConnectionString, redisConnectionString string) (Storages, error) {
	db, err := sqlx.Connect("postgres", postgresConnectionString)
	if err != nil {
		return Storages{}, err
	}

	heartbeatsRdb, err := heartbeats.NewStorage(redisConnectionString)
	if err != nil {
		return Storages{}, err
	}

	return Storages{
		ExpressionsStorage:    postgres.NewExpressionStorage(db),
		SubexpressionsStorage: postgres.NewSubexpressionStorage(db),
		TimeoutsStorage:       postgres.NewTimeoutsStorage(db),
		UsersStorage:          postgres.NewUsersStorage(db),
		HeartbeatsStorage:     heartbeatsRdb,
	}, nil
}

func (s *Storages) CloseAll() error {
	val := reflect.ValueOf(s).Elem()
	var errList []error

	for i := 0; i < val.NumField(); i++ {
		storage := val.Field(i).Interface()
		if closer, ok := storage.(io.Closer); ok {
			err := closer.Close()
			if err != nil {
				errList = append(errList, err)
			}
		}
	}

	if len(errList) > 0 {
		return fmt.Errorf("failed to close all storages: %v", errList)
	}

	return nil
}
