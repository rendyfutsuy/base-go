package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rendyfutsuy/base-go/utils"
)

type LocalStorage struct {
	baseDir        string
	originEndpoint string
}

func NewLocalStorage() (*LocalStorage, error) {
	base := utils.ConfigVars.String("local.base_dir")
	if base == "" {
		base = "public/storage"
	}
	origin := utils.ConfigVars.String("local.origin_endpoint")
	// ensure base directory exists
	if err := os.MkdirAll(base, 0755); err != nil {
		return nil, fmt.Errorf("failed to initialize local storage dir: %w", err)
	}
	return &LocalStorage{baseDir: base, originEndpoint: origin}, nil
}

func (l *LocalStorage) HealthCheck(ctx context.Context) error {
	// Check if baseDir exists and is writable
	testFile := filepath.Join(l.baseDir, fmt.Sprintf(".health_%d", time.Now().UnixNano()))
	if err := os.WriteFile(testFile, []byte("ok"), 0644); err != nil {
		return err
	}
	_ = os.Remove(testFile)
	return nil
}

func (l *LocalStorage) GetFullURL(path string) string {
	// Build full URL from originEndpoint if available, otherwise return relative path under /storage
	if l.originEndpoint != "" {
		return fmt.Sprintf("%s/%s", strings.TrimRight(l.originEndpoint, "/"), strings.TrimLeft(path, "/"))
	}
	return fmt.Sprintf("/storage/%s", strings.TrimLeft(path, "/"))
}

func (l *LocalStorage) UploadFile(buf bytes.Buffer, fileName string, destinatedPath string) (string, error) {
	// Write to local filesystem under baseDir/destinatedPath/fileName
	dir := filepath.Join(l.baseDir, destinatedPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create destination dir: %w", err)
	}
	target := filepath.Join(dir, fileName)
	if err := os.WriteFile(target, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}
	// Return full URL-like path consistent with GetFullURL scheme
	relPath := fmt.Sprintf("%s/%s", destinatedPath, fileName)
	return l.GetFullURL(relPath), nil
}

func (l *LocalStorage) DeleteFile(fileURL string) error {
	// Map URL back to local file path
	path := l.resolveLocalPath(fileURL)
	if path == "" {
		return fmt.Errorf("invalid file URL")
	}
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}

func (l *LocalStorage) GeneratePresignedURL(fullURL string) (string, error) {
	// For local storage, just return the same URL
	return "http://localhost:" + utils.ConfigVars.String("app_port") + fullURL, nil
}

func (l *LocalStorage) GeneratePresignedURLWithPreview(fullURL string) (string, error) {
	// For local storage, just return the same URL
	return "http://localhost:" + utils.ConfigVars.String("app_port") + fullURL, nil
}

func (l *LocalStorage) DownloadFile(fileURL string) (*bytes.Buffer, error) {
	path := l.resolveLocalPath(fileURL)
	if path == "" {
		return nil, fmt.Errorf("invalid file URL")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var b bytes.Buffer
	if _, err := io.Copy(&b, f); err != nil {
		return nil, err
	}
	return &b, nil
}

func (l *LocalStorage) CopyFile(originalFileURL string, overrideName *string) (string, error) {
	srcPath := l.resolveLocalPath(originalFileURL)
	if srcPath == "" {
		return "", fmt.Errorf("invalid original file URL")
	}
	// Determine destination file name
	srcDir := filepath.Dir(srcPath)
	srcBase := filepath.Base(srcPath)
	dstName := srcBase
	if overrideName != nil && *overrideName != "" {
		dstName = *overrideName
	}
	dstPath := filepath.Join(srcDir, dstName)
	src, err := os.Open(srcPath)
	if err != nil {
		return "", err
	}
	defer src.Close()
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}
	// Build URL
	rel := strings.TrimPrefix(dstPath, l.baseDir+string(os.PathSeparator))
	rel = filepath.ToSlash(rel)
	return l.GetFullURL(rel), nil
}

func (l *LocalStorage) resolveLocalPath(fileURL string) string {
	// Handle full URLs and relative URLs
	if strings.HasPrefix(fileURL, "http://") || strings.HasPrefix(fileURL, "https://") {
		u, err := url.Parse(fileURL)
		if err != nil {
			return ""
		}
		// Assume path after origin maps to /storage/<...> or configured origin base
		p := strings.TrimPrefix(u.Path, "/")
		// If originEndpoint is set, it may include no /storage prefix; we accept any path
		return filepath.Join(l.baseDir, filepath.FromSlash(p))
	}
	// Relative path starting with /storage or direct relative
	p := strings.TrimPrefix(fileURL, "/")
	if strings.HasPrefix(p, "storage/") {
		p = strings.TrimPrefix(p, "storage/")
	}
	return filepath.Join(l.baseDir, filepath.FromSlash(p))
}
