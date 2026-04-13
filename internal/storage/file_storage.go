package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/falaqmsi/go-example/internal/config"
)

// FileStorage defines the behavior for any file upload storage backend.
type FileStorage interface {
	// SaveFile uploads a file and returns its accessible URL or object key.
	SaveFile(ctx context.Context, file *multipart.FileHeader) (string, error)
}

// NewFileStorage creates the appropriate FileStorage implementation based on config.
func NewFileStorage(cfg config.StorageConfig) (FileStorage, error) {
	switch cfg.Type {
	case "local":
		return NewLocalStorage(cfg.LocalUploadDir)
	case "minio", "s3":
		return NewMinioStorage(cfg)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}
