package elasticsearch

import (
	"context"
	"techmind/pkg/config"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/fx"
)

func New(lc fx.Lifecycle, config *config.Config) *elasticsearch.Client {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{config.Elasticsearch.URL},
		Username:  config.Elasticsearch.Username,
		Password:  config.Elasticsearch.Password,
	})
	if err != nil {
		panic(err)
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return client
}
