package storage

import (
	"bytes"
)

// Storage defines the interface that any storage provider (local, AWS S3, MinIO) must implement.
type Storage interface {
	GetFullURL(path string) string
	UploadFile(buf bytes.Buffer, fileName string, destinatedPath string) (string, error)
	DeleteFile(fileURL string) error
	GeneratePresignedURL(fullURL string) (string, error)
	GeneratePresignedURLWithPreview(fullURL string) (string, error)
	DownloadFile(fileURL string) (*bytes.Buffer, error)
	CopyFile(originalFileURL string, overrideName *string) (string, error)
}
