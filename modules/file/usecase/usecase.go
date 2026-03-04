package usecase

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
	filemod "github.com/rendyfutsuy/base-go/modules/file"
	"github.com/rendyfutsuy/base-go/modules/file/dto"
	"github.com/rendyfutsuy/base-go/utils"
	utilsServices "github.com/rendyfutsuy/base-go/utils/services"
)

type Usecase interface {
	Upload(ctx context.Context, input dto.UploadInput) (*models.File, error)
	AssignFiles(ctx context.Context, input dto.AssignFilesToModule) error
	UnassignFiles(ctx context.Context, input dto.UnassignFilesFromModule) error
	Delete(ctx context.Context, id string) error
}

type fileUsecase struct {
	repo filemod.Repository
}

func NewFileUsecase(repo filemod.Repository) Usecase {
	return &fileUsecase{repo: repo}
}

func (u *fileUsecase) Upload(ctx context.Context, input dto.UploadInput) (*models.File, error) {
	dest := strings.Trim(input.DestRoot, "/")
	if input.ExtraPath != nil && *input.ExtraPath != "" {
		dest = dest + "/" + strings.Trim(*input.ExtraPath, "/")
	}

	ext := filepath.Ext(input.OriginalFileName)
	if ext == "" {
		ext = ""
	}
	newName := uuid.NewString() + ext

	var buf bytes.Buffer
	buf.Write(input.Data)
	fileURL, err := utilsServices.UploadFile(buf, newName, dest)
	if err != nil {
		return nil, err
	}

	desc := "-"
	if input.Description != nil && *input.Description != "" {
		desc = *input.Description
	}

	return u.repo.Create(ctx, dto.ToDBFile{
		Name:        input.OriginalFileName,
		FilePath:    &fileURL,
		Description: &desc,
	})
}

func (u *fileUsecase) AssignFiles(ctx context.Context, input dto.AssignFilesToModule) error {
	return u.repo.AssignFilesToModule(ctx, input)
}

func (u *fileUsecase) UnassignFiles(ctx context.Context, input dto.UnassignFilesFromModule) error {
	return u.repo.RemoveFilesFromModule(ctx, input.ModuleType, input.ModuleID)
}

func (u *fileUsecase) Delete(ctx context.Context, id string) error {
	uid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, uid)
}
