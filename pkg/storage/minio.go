package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

type MinIO struct {
	client *minio.Client
	bucket string
}

func NewMinIO(ctx context.Context, cfg Config) (*MinIO, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("check bucket: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("create bucket: %w", err)
		}
	}

	return &MinIO{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

func (m *MinIO) Upload(
	ctx context.Context,
	objectName string,
	reader io.Reader,
	size int64,
	contentType string,
) error {
	_, err := m.client.PutObject(
		ctx,
		m.bucket,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return fmt.Errorf("upload object: %w", err)
	}

	return nil
}

func (m *MinIO) Download(
	ctx context.Context,
	objectName string,
) (*minio.Object, error) {
	obj, err := m.client.GetObject(
		ctx,
		m.bucket,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("download object: %w", err)
	}

	return obj, nil
}

func (m *MinIO) Delete(
	ctx context.Context,
	objectName string,
) error {
	err := m.client.RemoveObject(
		ctx,
		m.bucket,
		objectName,
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	return nil
}

func (m *MinIO) Exists(
	ctx context.Context,
	objectName string,
) (bool, error) {
	_, err := m.client.StatObject(
		ctx,
		m.bucket,
		objectName,
		minio.StatObjectOptions{},
	)

	if err == nil {
		return true, nil
	}

	resp := minio.ToErrorResponse(err)

	if resp.Code == "NoSuchKey" {
		return false, nil
	}

	return false, err
}

func (m *MinIO) PresignedURL(
	ctx context.Context,
	objectName string,
	expiry time.Duration,
) (*url.URL, error) {
	u, err := m.client.PresignedGetObject(
		ctx,
		m.bucket,
		objectName,
		expiry,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("generate url: %w", err)
	}

	return u, nil
}
