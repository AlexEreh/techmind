package app

import (
	"techmind/internal/connetions/elasticsearch"
	"techmind/internal/connetions/gotenberg"
	"techmind/internal/connetions/minio"
	"techmind/internal/connetions/postgres"
	"techmind/internal/di"

	"go.uber.org/fx"
)

var App = fx.Options(
	fx.Provide(
		postgres.New,
		minio.New,
		elasticsearch.New,
		gotenberg.New,
	),
	di.Repository,
	di.Service,
	di.Transport,
)
