package config

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewMinio(config *viper.Viper, log *logrus.Logger) *minio.Client {
	minioClient, err := minio.New(config.GetString("minio.servers"), &minio.Options{
		Creds:  credentials.NewStaticV4(config.GetString("minio.accessKey"), config.GetString("minio.secret"), ""),
		Secure: config.GetBool("minio.useSSL"),
	})

	if err != nil {
		log.Fatalf("failed to connect minio: %v", err)
	}

	_, err = minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Fatalf("failed to connect minio: %v", err)
	}

	return minioClient
}
