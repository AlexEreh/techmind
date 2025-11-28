package main

import (
	"fmt"
	"techmind/app"
	"techmind/internal/connetions/postgres"
	"techmind/pkg/config"

	"go.uber.org/fx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}
	postgres.RunMigrations(cfg)

	fx.New(
		fx.Supply(cfg),
		app.App,
	).Run()
}
