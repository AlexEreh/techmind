package minio

import (
	"context"
	"techmind/pkg/config"

	"log"

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

	// Auto-create documents bucket if it doesn't exist
	ctx := context.Background()
	bucketName := "documents"

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Printf("Error checking if bucket exists: %v", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Error creating bucket: %v", err)
		} else {
			log.Printf("Bucket '%s' created successfully", bucketName)
		}
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return client
}
