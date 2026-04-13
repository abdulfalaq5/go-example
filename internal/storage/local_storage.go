package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type localStorage struct {
	baseDir string
}

// NewLocalStorage creates a file storage implementation that saves files to disk.
func NewLocalStorage(baseDir string) (FileStorage, error) {
	// Ensure the base directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("localStorage: failed to create upload directory %q: %w", baseDir, err)
	}
	return &localStorage{baseDir: baseDir}, nil
}

func (s *localStorage) SaveFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("localStorage: could not open uploaded file: %w", err)
	}
	defer src.Close()

	// Generate a unique filename using UUID to prevent collisions
	ext := filepath.Ext(fileHeader.Filename)
	uniqueName := fmt.Sprintf("%d-%s%s", time.Now().Unix(), uuid.New().String(), ext)
	
	// Create destination path
	destPath := filepath.Join(s.baseDir, uniqueName)
	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("localStorage: could not create destination file: %w", err)
	}
	defer dst.Close()

	// Copy content
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("localStorage: could not save file content: %w", err)
	}

	// For local storage, the exact URL mapping typically depends on how Gin routes static files.
	// We'll return the relative URL path assuming the app serves this directory under /uploads
	return fmt.Sprintf("/uploads/%s", uniqueName), nil
}
