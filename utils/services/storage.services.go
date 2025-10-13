package services

import (
	"bytes"
	"fmt"

	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/services/storage"
	"go.uber.org/zap"
)

const (
	AWS   = "s3"
	MINIO = "minio"
)

func GetStorage(driver string) (storage.Storage, error) {
	switch driver {
	case AWS:
		return storage.NewS3Storage()
	case MINIO:
		return storage.NewMinIOStorage()
	default:
		zap.S().Errorf("unsupported file driver: %s", driver)
		return nil, fmt.Errorf("unsupported file driver: %s", driver)
	}
}

// GetFullURL return full URL of a relative path based on storage driver.
//
// It takes the document relative path as parameter.
//
// It returns a string representing the URL of the uploaded file, and an error.
func GetFullURL(path string) (string, error) {
	driver := utils.ConfigVars.String("file.driver")
	storage, err := GetStorage(driver)
	if err != nil {
		utils.Logger.Error(err.Error())
		return "", err
	}

	return storage.GetFullURL(path), nil
}

// UploadFile uploads a file to the configured storage provider.
//
// It takes a bytes.Buffer containing the file content, a fileName, and a
// destinatedPath as input parameters.
//
// It returns a string representing the URL of the uploaded file, and an error.
func UploadFile(buf bytes.Buffer, fileName string, destinatedPath string) (string, error) {
	driver := utils.ConfigVars.String("file.driver")
	storage, err := GetStorage(driver)
	if err != nil {
		utils.Logger.Error(err.Error())
		return "", err
	}
	return storage.UploadFile(buf, fileName, destinatedPath)
}

// DeleteFile deletes a file from the configured storage provider.
//
// It takes a string representing the URL of the file as input parameter.
//
// It returns an error.
func DeleteFile(fileURL string) error {
	driver := utils.ConfigVars.String("file.driver")
	storage, err := GetStorage(driver)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	return storage.DeleteFile(fileURL)
}

// GeneratePresignedURL generates a presigned URL for accessing the uploaded file.
//
// It takes a string representing the URL of the uploaded file as input parameter.
//
// It returns a string representing the presigned URL for accessing the uploaded file, and an error.
func GeneratePresignedURL(fullURL string) (string, error) {
	// Generate presigned URL for accessing the uploaded file

	// only generate presigned URL if fullURL is not empty, dont return error if empty
	if fullURL == "" {
		return "", nil
	}

	driver := utils.ConfigVars.String("file.driver")
	storage, err := GetStorage(driver)
	if err != nil {
		utils.Logger.Error(err.Error())
		return "", err
	}
	return storage.GeneratePresignedURL(fullURL)
}

// GeneratePresignedURL generates a presigned URL for accessing the uploaded file. With Preview
//
// It takes a string representing the URL of the uploaded file as input parameter.
//
// It returns a string representing the presigned URL for accessing the uploaded file, and an error.
func GeneratePresignedURLWithPreview(fullURL string) (string, error) {
	// Generate presigned URL for accessing the uploaded file

	// only generate presigned URL if fullURL is not empty, dont return error if empty
	if fullURL == "" {
		return "", nil
	}

	driver := utils.ConfigVars.String("file.driver")
	storage, err := GetStorage(driver)
	if err != nil {
		utils.Logger.Error(err.Error())
		return "", err
	}
	return storage.GeneratePresignedURLWithPreview(fullURL)
}

func DownloadFile(fileURL string) (*bytes.Buffer, error) {
	driver := utils.ConfigVars.String("file.driver")
	storage, err := GetStorage(driver)
	if err != nil {
		utils.Logger.Error(err.Error())
		return nil, err
	}
	return storage.DownloadFile(fileURL)
}

func CopyFile(path string, overrideName *string) (string, error) {
	driver := utils.ConfigVars.String("file.driver")
	storage, err := GetStorage(driver)
	if err != nil {
		utils.Logger.Error(err.Error())
		return "", err
	}

	return storage.CopyFile(path, overrideName)
}
