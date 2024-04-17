package heartbeats

import (
	"context"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type Storage struct {
	rdb *redis.Client
}

func NewStorage(connectionString string) (*Storage, error) {
	opts, err := redis.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)

	return &Storage{
		rdb: rdb,
	}, nil
}

func (s *Storage) SaveHeartbeat(ctx context.Context, workerID int) error {
	return s.rdb.Set(ctx, strconv.Itoa(workerID), time.Now().Format(time.DateTime), time.Minute).Err()
}
func (s *Storage) GetAllHeartbeats(ctx context.Context) ([]model.Heartbeat, error) {

	keys, err := s.rdb.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	var heartbeats []model.Heartbeat
	for _, key := range keys {
		value, err := s.rdb.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		workerID, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}

		sentAt, err := time.Parse(time.DateTime, value)
		if err != nil {
			return nil, err
		}

		heartbeats = append(heartbeats, model.Heartbeat{
			WorkerID: workerID,
			SentAt:   sentAt,
		})
	}
	return heartbeats, nil
}

func (s *Storage) Close() error {
	return s.rdb.Close()
}
