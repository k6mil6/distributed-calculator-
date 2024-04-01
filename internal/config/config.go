package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
	"log"
	"sync"
	"time"
)

type Config struct {
	Env              string        `hcl:"env" env:"ENV" default:"local"`
	DatabaseDSN      string        `hcl:"database_dsn" env:"DB_DSN" default:"postgres://postgres:postgres@localhost:5442/calc_db?sslmode=disable"`
	GoroutineNumber  int           `hcl:"goroutine_number" env:"GOROUTINE_NUMBER" default:"5"`
	OrchestratorURL  string        `hcl:"orchestrator_url" env:"ORCHESTRATOR_URL" default:"http://localhost:5441"`
	WorkerTimeout    time.Duration `hcl:"worker_timeout" env:"WORKER_TIMEOUT" default:"1m"`
	HeartbeatTimeout time.Duration `hcl:"heartbeat_timeout" env:"HEARTBEAT_TIMEOUT" default:"30s"`
	FetcherInterval  time.Duration `hcl:"fetcher_interval" env:"FETCHER_INTERVAL" default:"10s"`
	CheckerInterval  time.Duration `hcl:"checker_interval" env:"CHECKER_INTERVAL" default:"10s"`
	MigrationPath    string        `hcl:"migration_path" env:"MIGRATION_PATH" default:"./internal/storage/migrations"`
	GrpcPort         int           `hcl:"grpc_port" env:"GRPC_PORT" default:"50051"`
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
