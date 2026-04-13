package service

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/falaqmsi/go-example/internal/storage"
)

// UploadService contains business logic for processing file uploads.
type UploadService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error)
}

type uploadService struct {
	storage storage.FileStorage
}

// NewUploadService constructs a new UploadService using the provided FileStorage.
func NewUploadService(s storage.FileStorage) UploadService {
	return &uploadService{storage: s}
}

// UploadFile delegates the actual persistence to the configured storage backend
// and potentially could encapsulate business rules like file typing, size constraints, etc.
func (s *uploadService) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	// For example, file validation could occur here.
	// max size, mime-type whitelists, etc.

	url, err := s.storage.SaveFile(ctx, file)
	if err != nil {
		return "", fmt.Errorf("uploadService failed to save file: %w", err)
	}

	return url, nil
}
