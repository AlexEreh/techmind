package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
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

	// Initialize index with Russian analyzer
	initializeElasticsearchIndex(client)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return client
}

func initializeElasticsearchIndex(client *elasticsearch.Client) {
	indexName := "documents"

	// Check if index exists
	existsRes, err := client.Indices.Exists([]string{indexName})
	if err != nil {
		panic(err)
	}
	defer existsRes.Body.Close()

	// If index exists, skip initialization
	if existsRes.StatusCode == 200 {
		return
	}

	// Create index with Russian analyzer settings
	settings := map[string]interface{}{
		"settings": map[string]interface{}{
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"rebuilt_russian": map[string]interface{}{
						"tokenizer": "standard",
						"filter": []string{
							"lowercase",
							"russian_stop",
							"russian_stemmer",
						},
					},
				},
				"filter": map[string]interface{}{
					"russian_stop": map[string]interface{}{
						"type":      "stop",
						"stopwords": "_russian_",
					},
					"russian_stemmer": map[string]interface{}{
						"type":     "stemmer",
						"language": "russian",
					},
				},
			},
		},
	}

	body, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}

	res, err := client.Indices.Create(indexName, client.Indices.Create.WithBody(bytes.NewReader(body)))
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		panic("Failed to create Elasticsearch index")
	}
}
