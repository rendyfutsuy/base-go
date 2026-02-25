package test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/knadh/koanf/v2"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/file/dto"
	"github.com/rendyfutsuy/base-go/modules/file/usecase"
	"github.com/rendyfutsuy/base-go/utils"
	utilsServices "github.com/rendyfutsuy/base-go/utils/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Create(ctx context.Context, data dto.ToDBFile) (*models.File, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.File), args.Error(1)
}
func (m *MockFileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockFileRepository) AssignFilesToModule(ctx context.Context, input dto.AssignFilesToModule) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}
func (m *MockFileRepository) RemoveFilesFromModule(ctx context.Context, moduleType string, moduleID uuid.UUID) error {
	args := m.Called(ctx, moduleType, moduleID)
	return args.Error(0)
}

func initLocalStorageBaseDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "file-usecase-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	// Ensure koanf exists
	if utils.ConfigVars == nil {
		utils.ConfigVars = koanf.New(".")
	}
	utils.ConfigVars.Set("local.base_dir", dir)
	utils.ConfigVars.Set("local.origin_endpoint", "")
	if err := utilsServices.InitStorage(utilsServices.LOCAL); err != nil {
		t.Fatalf("failed to init local storage: %v", err)
	}
	return dir
}

func TestUpload_DescriptionDefaultAndPathAppend(t *testing.T) {
	base := initLocalStorageBaseDir(t)
	defer os.RemoveAll(base)

	ctx := context.Background()
	mockRepo := new(MockFileRepository)
	uc := usecase.NewFileUsecase(mockRepo)

	data := []byte("hello world")
	orig := "picture.jpg"
	destRoot := "uploads"
	extra := "images"
	var desc *string = nil // trigger default "-"

	mockRepo.On("Create", ctx, mock.MatchedBy(func(d dto.ToDBFile) bool {
		if d.Name != orig {
			return false
		}
		if d.FilePath == nil || !strings.Contains(*d.FilePath, "/storage/uploads/images/") || !strings.HasSuffix(*d.FilePath, ".jpg") {
			return false
		}
		if d.Description == nil || *d.Description != "-" {
			return false
		}
		return true
	})).Return(&models.File{
		ID:        uuid.New(),
		Name:      orig,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil).Once()

	res, err := uc.Upload(ctx, dto.UploadInput{
		Data:             data,
		OriginalFileName: orig,
		DestRoot:         destRoot,
		ExtraPath:        &extra,
		Description:      desc,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestUpload_WithDescriptionAndNoExtraPath(t *testing.T) {
	base := initLocalStorageBaseDir(t)
	defer os.RemoveAll(base)

	ctx := context.Background()
	mockRepo := new(MockFileRepository)
	uc := usecase.NewFileUsecase(mockRepo)

	data := []byte("content")
	orig := "doc.pdf"
	destRoot := "docs"
	desc := "My PDF"

	mockRepo.On("Create", ctx, mock.MatchedBy(func(d dto.ToDBFile) bool {
		return d.Name == orig &&
			d.FilePath != nil &&
			strings.Contains(*d.FilePath, "/storage/docs/") &&
			strings.HasSuffix(*d.FilePath, ".pdf") &&
			d.Description != nil && *d.Description == desc
	})).Return(&models.File{
		ID:        uuid.New(),
		Name:      orig,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil).Once()

	res, err := uc.Upload(ctx, dto.UploadInput{
		Data:             data,
		OriginalFileName: orig,
		DestRoot:         destRoot,
		Description:      &desc,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)

	// additionally check file actually exists on disk
	// compute expected directory under base
	// we can't know the exact filename, but we can glob for *.pdf under base/docs
	glob := filepath.Join(base, "docs", "*.pdf")
	matches, _ := filepath.Glob(glob)
	assert.True(t, len(matches) >= 1)
}
