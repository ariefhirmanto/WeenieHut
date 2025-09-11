package storage

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	S3AccessKeyID     = os.Getenv("S3_ACCESS_KEY_ID")
	S3SecretAccessKey = os.Getenv("S3_SECRET_ACCESS_KEY")
	S3Endpoint        = os.Getenv("S3_ENDPOINT")
	S3Bucket          = os.Getenv("S3_BUCKET")
)

type MinioStorage struct {
	client *minio.Client
}

func New(
	s3Endpoint, s3AccessKeyID, s3SecretAccessKey string,
) *MinioStorage {

	mc, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKeyID, s3SecretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatal(err)
	}

	return &MinioStorage{
		client: mc,
	}
}

func (s *MinioStorage) UploadFile(ctx context.Context, bucket, localPath, remotePath string) (string, error) {
	_, err := s.client.FPutObject(ctx, bucket, remotePath, localPath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	presignedURL, err := s.client.PresignedGetObject(ctx, bucket, remotePath, 10*time.Minute, s.client.EndpointURL().Query())
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
