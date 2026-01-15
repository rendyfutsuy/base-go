package storage

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rendyfutsuy/base-go/utils"
)

// MinIOStorage implements the Storage interface for MinIO.
type MinIOStorage struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	ctx        context.Context
}

func (m *MinIOStorage) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exists, err := m.client.BucketExists(ctx, m.bucketName)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("bucket %s does not exist", m.bucketName)
	}

	return nil
}

// NewMinIOStorage initializes a new MinIOStorage instance.
func NewMinIOStorage() (*MinIOStorage, error) {
	endpoint := utils.ConfigVars.String("minio.endpoint")
	accessKeyID := utils.ConfigVars.String("minio.access_key_id")
	secretAccessKey := utils.ConfigVars.String("minio.secret_access_key")
	useSSL := utils.ConfigVars.Bool("minio.use_ssl")

	// Validate the endpoint and credentials
	if endpoint == "" || accessKeyID == "" || secretAccessKey == "" {
		return nil, fmt.Errorf("MinIO configuration is missing required fields")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	bucketName := utils.ConfigVars.String("minio.bucket_name")
	if bucketName == "" {
		return nil, fmt.Errorf("MinIO bucket name is missing")
	}

	return &MinIOStorage{
		client:     client,
		bucketName: bucketName,
		endpoint:   utils.ConfigVars.String("minio.origin_endpoint"),
		ctx:        context.Background(),
	}, nil
}

func (s *MinIOStorage) GetFullURL(path string) string {
	endpoint := utils.ConfigVars.String("minio.origin_endpoint")
	bucket := utils.ConfigVars.String("minio.bucket_name")

	return fmt.Sprintf("%s/%s/%s", endpoint, bucket, path)
}

// UploadFile uploads a file to MinIO and returns the URL.
func (m *MinIOStorage) UploadFile(buf bytes.Buffer, fileName string, destinatedPath string) (string, error) {
	if m.client == nil || m.bucketName == "" || m.endpoint == "" {
		return "", fmt.Errorf("MinIO storage is not properly initialized")
	}

	objectName := fmt.Sprintf("%s/%s", destinatedPath, fileName)

	// Upload the file
	_, err := m.client.PutObject(m.ctx, m.bucketName, objectName, bytes.NewReader(buf.Bytes()), int64(buf.Len()), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		utils.Logger.Error("Uploading file to MinIO - Error:" + err.Error())
		return "", err
	}

	// Construct the full URL
	fileURL := fmt.Sprintf("%s/%s/%s", m.endpoint, m.bucketName, objectName)

	return fileURL, nil
}

// DeleteFile deletes a file from MinIO given the URL of the file.
func (m *MinIOStorage) DeleteFile(fileURL string) error {
	u, err := url.Parse(fileURL)
	if err != nil {
		utils.Logger.Error("Parsing URL - Error:" + err.Error())
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	// Ensure the bucket name and object key are correctly parsed
	if !strings.HasPrefix(u.Path, "/"+m.bucketName+"/") {
		return fmt.Errorf("URL does not match the expected format with bucket name")
	}

	// Extract the object key from the URL path
	objectKey := strings.TrimPrefix(u.Path, "/"+m.bucketName+"/")

	// Proceed to delete the object
	err = m.client.RemoveObject(m.ctx, m.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		utils.Logger.Error("Deleting file from MinIO - Error:" + err.Error())
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GeneratePresignedURL generates a presigned URL for accessing the uploaded file.
func (m *MinIOStorage) GeneratePresignedURL(fullURL string) (string, error) {
	// Parse the full URL
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return "", fmt.Errorf("parsing URL - %v", err)
	}

	// Extract object key from the URL
	objectName := strings.TrimPrefix(parsedURL.Path, fmt.Sprintf("/%s/", m.bucketName))

	// Generate a presigned URL with a validity of 1 hour
	presignedURL, err := m.client.PresignedGetObject(m.ctx, m.bucketName, objectName, time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("generating presigned URL - %v", err)
	}

	return presignedURL.String(), nil
}

// GeneratePresignedURL generates a presigned URL for accessing the uploaded file, with Preview.
func (m *MinIOStorage) GeneratePresignedURLWithPreview(fullURL string) (string, error) {
	// Parse the full URL
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return "", fmt.Errorf("parsing URL - %v", err)
	}

	// Extract object key from the URL
	objectName := strings.TrimPrefix(parsedURL.Path, fmt.Sprintf("/%s/", m.bucketName))

	params := url.Values{}
	params.Set("response-content-type", m.GetContentType(objectName))

	// Generate a presigned URL with a validity of 1 hour
	presignedURL, err := m.client.PresignedGetObject(m.ctx, m.bucketName, objectName, time.Hour, params)
	if err != nil {
		return "", fmt.Errorf("generating presigned URL - %v", err)
	}

	return presignedURL.String(), nil
}

func (m *MinIOStorage) GetContentType(filename string) string {
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream" // Default jika tidak terdeteksi
	}
	return mimeType
}

func (m *MinIOStorage) DownloadFile(fileURL string) (*bytes.Buffer, error) {
	u, err := url.Parse(fileURL)
	if err != nil {
		utils.Logger.Error("Parsing URL - Error:" + err.Error())
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// Ensure the bucket name and object key are correctly parsed
	if !strings.HasPrefix(u.Path, "/"+m.bucketName+"/") {
		u.Path = fmt.Sprint("/"+m.bucketName+"/", u.Path)
	}

	// Extract the object key from the URL path
	objectKey := strings.TrimPrefix(u.Path, "/"+m.bucketName+"/")

	// Proceed to download the object
	object, err := m.client.GetObject(m.ctx, m.bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		utils.Logger.Error("Downloading file from MinIO - Error:" + err.Error())
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(object)

	return buf, nil
}

func (m *MinIOStorage) CopyFile(originalFileURL string, overideName *string) (string, error) {
	// Parse original URL
	u, err := url.Parse(originalFileURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse original file URL: %w", err)
	}

	// Pastikan path sesuai dengan bucket
	if !strings.HasPrefix(u.Path, "/"+m.bucketName+"/") {
		return "", fmt.Errorf("URL does not match expected format with bucket name")
	}

	// Extract object key
	srcObjectName := strings.TrimPrefix(u.Path, "/"+m.bucketName+"/")

	// Tentukan nama baru untuk file hasil copy

	newFileName := ""
	if overideName != nil {
		newFileName = *overideName
	} else {
		timeFormat := "20060201150405"
		createdAtString := time.Now().UTC().Format(timeFormat)
		ext := filepath.Ext(srcObjectName)
		nameWithoutExt := strings.TrimSuffix(filepath.Base(srcObjectName), ext)
		newFileName = nameWithoutExt + "_copied_at_" + createdAtString + ext
	}

	newObjectName := filepath.Join(filepath.Dir(srcObjectName), newFileName)

	// Setup sumber
	src := minio.CopySrcOptions{
		Bucket: m.bucketName,
		Object: srcObjectName,
	}

	// Setup tujuan
	dst := minio.CopyDestOptions{
		Bucket: m.bucketName,
		Object: newObjectName,
	}

	// Lakukan copy
	_, err = m.client.CopyObject(m.ctx, dst, src)
	if err != nil {
		return "", fmt.Errorf("failed to copy object: %w", err)
	}

	// Bangun URL file hasil copy
	copiedFileURL := fmt.Sprintf("%s/%s/%s", m.endpoint, m.bucketName, newObjectName)
	return copiedFileURL, nil
}
