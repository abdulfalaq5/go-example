package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/falaqmsi/go-example/internal/config"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioStorage struct {
	client *minio.Client
	bucket string
}

// NewMinioStorage initializes a new MinIO/S3 compatible storage backend.
func NewMinioStorage(cfg config.StorageConfig) (FileStorage, error) {
	// Initialize minio client object
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minioStorage: failed to initialize client: %w", err)
	}

	// Ensure bucket exists
	ctx := context.Background()
	exists, err := minioClient.BucketExists(ctx, cfg.MinioBucket)
	if err != nil {
		return nil, fmt.Errorf("minioStorage: check bucket failed: %w", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("minioStorage: create bucket %q failed: %w", cfg.MinioBucket, err)
		}
	}

	return &minioStorage{
		client: minioClient,
		bucket: cfg.MinioBucket,
	}, nil
}

func (s *minioStorage) SaveFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("minioStorage: could not open uploaded file: %w", err)
	}
	defer src.Close()

	ext := filepath.Ext(fileHeader.Filename)
	uniqueName := fmt.Sprintf("%d-%s%s", time.Now().Unix(), uuid.New().String(), ext)

	// Note: We use the header's content type or fallback to octet-stream
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	info, err := s.client.PutObject(ctx, s.bucket, uniqueName, src, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("minioStorage: upload failed: %w", err)
	}

	// For MinIO, we return the object key or construct a presumed public URL
	// If public read access is enabled on the bucket, the URL is accessible directly.
	protocol := "http"
	if s.client.EndpointURL().Scheme == "https" {
		protocol = "https"
	}
	
	fileURL := fmt.Sprintf("%s://%s/%s/%s", protocol, s.client.EndpointURL().Host, s.bucket, info.Key)
	return fileURL, nil
}
