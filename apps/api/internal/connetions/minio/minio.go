package minio

import (
	"context"
	"techmind/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
)

func New(lc fx.Lifecycle, config *config.Config) *minio.Client {
	client, err := minio.New(config.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinIO.AccessKey, config.MinIO.SecretKey, ""),
		Secure: config.MinIO.UseSSL,
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
