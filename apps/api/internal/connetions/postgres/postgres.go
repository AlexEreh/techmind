package postgres

import (
	"context"
	"database/sql"
	"techmind/pkg/config"
	"techmind/schema/ent"

	"github.com/pressly/goose/v3"
	"go.uber.org/fx"

	_ "github.com/lib/pq"
)

func RunMigrations(config *config.Config) {
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	db, err := sql.Open("postgres", config.Postgres.Conn)
	if err != nil {
		panic(err)
	}
	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}

func New(lc fx.Lifecycle, config *config.Config) *ent.Client {
	client, err := ent.Open("postgres", config.Postgres.Conn)
	if err != nil {
		panic(err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})

	return client
}
