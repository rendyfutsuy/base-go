package storage

import (
	"bytes"
	"fmt"
	"mime"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	aws_service "github.com/rendyfutsuy/base-go/helpers/aws"
	"github.com/rendyfutsuy/base-go/utils"
)

type S3Storage struct {
	svc *s3.S3
}

func NewS3Storage() (*S3Storage, error) {
	svc := aws_service.StartS3()
	return &S3Storage{svc}, nil
}

func (s *S3Storage) GetFullURL(path string) string {
	url := utils.ConfigVars.String("aws.aws_origin_endpoint")

	return fmt.Sprintf("%s/%s", url, path)
}

func (s *S3Storage) UploadFile(buf bytes.Buffer, fileName string, destinatedPath string) (string, error) {
	svc := s.svc
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(utils.ConfigVars.String("aws.aws_bucket")),
		Key:    aws.String(fmt.Sprintf("%s/%s", destinatedPath, fileName)),
		Body:   bytes.NewReader(buf.Bytes()),
		// ACL:    aws.String("public-read"), // set ACL to public-read, only set it to public read if need it
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s/%s", utils.ConfigVars.String("aws.aws_origin_endpoint"), destinatedPath, fileName)

	return url, nil
}

func (s *S3Storage) DeleteFile(fileURL string) error {
	u, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	objectKey := strings.TrimPrefix(u.Path, "/")
	svc := s.svc
	err = applyLifecycleRule(svc, utils.ConfigVars.String("aws.aws_bucket"))
	if err != nil {
		return fmt.Errorf("failed to apply lifecycle rule: %w", err)
	}

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(utils.ConfigVars.String("aws.aws_bucket")),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(utils.ConfigVars.String("aws.aws_bucket")),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to wait for file deletion: %w", err)
	}

	return nil
}

// GeneratePresignedURL generates a presigned URL for accessing the uploaded file.
func (s *S3Storage) GeneratePresignedURL(fileURL string) (string, error) {
	// Parse the input URL to get the S3 object key (file path)
	u, err := url.Parse(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Extract the object key from the URL path
	objectKey := strings.TrimPrefix(u.Path, "/")

	// Initialize the S3 client
	svc := s.svc
	// Define the request for a presigned URL
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(utils.ConfigVars.String("aws.aws_bucket")),
		Key:    aws.String(objectKey),
	})

	// Generate the presigned URL (valid for 60 minutes)
	presignedURL, err := req.Presign(60 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL, nil
}

func (s *S3Storage) GeneratePresignedURLWithPreview(fullURL string) (string, error) {
	// Parse the input URL to get the S3 object key (file path)
	u, err := url.Parse(fullURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	// Extract the object key from the URL path
	objectKey := strings.TrimPrefix(u.Path, "/")

	// Initialize the S3 client
	svc := s.svc
	// Define the request for a presigned URL
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(utils.ConfigVars.String("aws.aws_bucket")),
		Key:    aws.String(objectKey),
	})

	// Generate the presigned URL (valid for 60 minutes)
	presignedURL, err := req.Presign(60 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL, nil
}

func (s *S3Storage) GetContentType(filename string) string {
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream" // Default jika tidak terdeteksi
	}
	return mimeType
}

// extra function to apply lifecycle rule
// applyLifecycleRule ensures that the bucket has the 30-day deletion rule applied
func applyLifecycleRule(svc *s3.S3, bucket string) error {
	// Get the current lifecycle rules for the bucket
	days, err := strconv.ParseInt(utils.ConfigVars.String("aws.aws_delete_cycle_days"), 10, 64)

	if err != nil {
		return fmt.Errorf("error parsing aws_delete_cycle_days: %w", err)
	}

	// Define the lifecycle rule
	rule := &s3.LifecycleRule{
		ID:     aws.String(fmt.Sprintf("Delete after %d days", days)),
		Filter: &s3.LifecycleRuleFilter{},
		Status: aws.String("Enabled"),
		Expiration: &s3.LifecycleExpiration{
			Days: aws.Int64(days),
		},
	}

	// Apply the lifecycle configuration to the bucket
	input := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
		LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
			Rules: []*s3.LifecycleRule{rule},
		},
	}

	_, err = svc.PutBucketLifecycleConfiguration(input)
	if err != nil {
		return fmt.Errorf("error applying lifecycle configuration: %w", err)
	}

	return nil
}

func (s *S3Storage) DownloadFile(fileURL string) (*bytes.Buffer, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *S3Storage) CopyFile(originalFileURL string, overrideName *string) (string, error) {
	return "", nil
}
