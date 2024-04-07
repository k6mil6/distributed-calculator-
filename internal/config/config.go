package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
	"log"
	"sync"
	"time"
)

type Config struct {
	Env             string        `hcl:"env" env:"ENV" default:"local"`
	DatabaseDSN     string        `hcl:"database_dsn" env:"DB_DSN" default:"postgres://postgres:postgres@localhost:5442/calc_db?sslmode=disable"`
	FetcherInterval time.Duration `hcl:"fetcher_interval" env:"FETCHER_INTERVAL" default:"10s"`
	CheckerInterval time.Duration `hcl:"checker_interval" env:"CHECKER_INTERVAL" default:"10s"`
	MigrationPath   string        `hcl:"migration_path" env:"MIGRATION_PATH" default:"./internal/storage/migrations"`
	GrpcPort        int           `hcl:"grpc_port" env:"GRPC_PORT" default:"50051"`

	HeartbeatTimeout     time.Duration `hcl:"heartbeat_timeout" env:"HEARTBEAT_TIMEOUT" default:"30s"`
	GoroutineNumber      int           `hcl:"goroutine_number" env:"GOROUTINE_NUMBER" default:"5"`
	WorkerTimeout        time.Duration `hcl:"worker_timeout" env:"WORKER_TIMEOUT" default:"1m"`
	GRPCServerAddress    string        `hcl:"grpc_server_address" env:"GRPC_SERVER_ADDRESS" default:"localhost:50051"`
	GRPCReconnectTimeout time.Duration `hcl:"grpc_reconnect_timeout" env:"GRPC_RECONNECT_TIMEOUT" default:"5s"`
	GRPCReconnectRetries int           `hcl:"grpc_reconnect_retries" env:"GRPC_RECONNECT_RETRIES" default:"5"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "NFB",
			Files:     []string{"./config.hcl", "./config.local.hcl"},
			FileDecoders: map[string]aconfig.FileDecoder{
				".hcl": aconfighcl.New(),
			},
		})

		if err := loader.Load(); err != nil {
			log.Printf("failed to load config: %v", err)
		}
	})
	return cfg
}
