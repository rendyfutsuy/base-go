package utils

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const TEMP_PATH = "./public/temp"

func GetFileSizeFromURL(url string) (int64, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	contentLength := resp.Header.Get("Content-Length")
	var size int64
	_, err = fmt.Sscanf(contentLength, "%d", &size)
	if err != nil {
		return 0, err
	}

	return size, nil
}

// formatFileSize converts bytes into a human-readable string (e.g., "17.4 MB")
func GetFormattedFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func GetFileNameAndExt(fileName string) (string, string) {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))], filepath.Ext(fileName)
}

func SaveFileToLocalStorage(data []byte, filename, path string) (string, error) {
	publicPath := TEMP_PATH

	folder := filepath.Join(publicPath, path)
	newpath := filepath.Join(path, filename)
	localpath := filepath.Join(folder, filename)

	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(localpath, data, 0644); err != nil {
		return "", err
	}

	return newpath, nil
}
