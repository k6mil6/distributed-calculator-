package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/k6mil6/distributed-calculator/internal/config"
)

func main() {
	cfg := config.Get()

	m, err := migrate.New(
		"file://"+cfg.MigrationPath,
		cfg.DatabaseDSN,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied")
}
