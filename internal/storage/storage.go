package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/internal/storage/postgres"
	"github.com/k6mil6/distributed-calculator/internal/storage/redis/heartbeats"
	"io"
	"reflect"
	"time"
)

type Storages struct {
	ExpressionsStorage    *postgres.ExpressionsStorage
	SubexpressionsStorage *postgres.SubexpressionsStorage
	TimeoutsStorage       *postgres.TimeoutsStorage
	UsersStorage          *postgres.UsersStorage
	HeartbeatsStorage     *heartbeats.Storage
}

func New(postgresConnectionString, redisConnectionString string, maxRetries int, retryCooldown time.Duration) (Storages, error) {
	var db *sqlx.DB
	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", postgresConnectionString)
		if err == nil {
			break
		}
		time.Sleep(retryCooldown)
	}
	if err != nil {
		return Storages{}, fmt.Errorf("failed to connect to postgres after %d retries: %w", maxRetries, err)
	}

	var heartbeatsRdb *heartbeats.Storage
	for i := 0; i < maxRetries; i++ {
		heartbeatsRdb, err = heartbeats.NewStorage(redisConnectionString)
		if err == nil {
			break
		}
		time.Sleep(retryCooldown)
	}
	if err != nil {
		return Storages{}, fmt.Errorf("failed to connect to redis after %d retries: %w", maxRetries, err)
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
