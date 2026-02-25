package usecase

import (
	"bytes"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	paramMod "github.com/rendyfutsuy/base-go/modules/parameter"
	"github.com/rendyfutsuy/base-go/modules/post"
	"github.com/rendyfutsuy/base-go/modules/post/dto"
	"github.com/rendyfutsuy/base-go/utils"
	utilsServices "github.com/rendyfutsuy/base-go/utils/services"
)

type postUsecase struct {
	repo      post.Repository
	paramRepo paramMod.Repository
}

func NewPostUsecase(repo post.Repository, paramRepo paramMod.Repository) post.Usecase {
	return &postUsecase{repo: repo, paramRepo: paramRepo}
}

func (u *postUsecase) Create(ctx context.Context, req *dto.ReqCreatePost, authId string, thumbnailData []byte, thumbnailName string) (*models.Post, error) {
	// Validate parameter types
	if err := u.validateParameterType(ctx, req.LangID, "lang"); err != nil {
		return nil, err
	}
	for _, tid := range req.TopicIDs {
		if err := u.validateParameterType(ctx, tid, "topic"); err != nil {
			return nil, err
		}
	}

	createdBy := uuid.Nil
	if authId != "" {
		if uid, err := utils.StringToUUID(authId); err == nil {
			createdBy = uid
		}
	}

	// Upload thumbnail first if provided, so URL can be saved in a single create
	var uploadedURL *string
	if len(thumbnailData) > 0 && thumbnailName != "" {
		var buf bytes.Buffer
		buf.Write(thumbnailData)
		key := uuid.NewString()
		destinatedPath := "posts/thumbnails/" + key
		url, err := utilsServices.UploadFile(buf, thumbnailName, destinatedPath)
		if err != nil {
			return nil, errors.New("Failed to upload thumbnail file")
		}
		uploadedURL = &url
	}

	c, err := u.repo.Create(ctx, createdBy, dto.ToDBPost{
		Title:            req.Title,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Price:            req.Price,
		DiscountRate:     req.DiscountRate,
		LangID:           req.LangID,
		TopicIDs:         req.TopicIDs,
		ThumbnailURL: func() *string {
			if uploadedURL != nil {
				return uploadedURL
			}
			return req.ThumbnailURL
		}(),
	})
	if err != nil {
		return nil, err
	}

	// Assign relations via parameters_to_module
	if err := u.paramRepo.AssignParametersToModule(ctx, "post", c.ID, []uuid.UUID{req.LangID}); err != nil {
		return nil, err
	}
	if len(req.TopicIDs) > 0 {
		if err := u.paramRepo.AssignParametersToModule(ctx, "post", c.ID, req.TopicIDs); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (u *postUsecase) Update(ctx context.Context, id string, req *dto.ReqUpdatePost, authId string, thumbnailData []byte, thumbnailName string) (*models.Post, error) {
	// Validate parameter types
	if err := u.validateParameterType(ctx, req.LangID, "lang"); err != nil {
		return nil, err
	}
	for _, tid := range req.TopicIDs {
		if err := u.validateParameterType(ctx, tid, "topic"); err != nil {
			return nil, err
		}
	}

	cid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	// Upload thumbnail first if provided
	var uploadedURL *string
	if len(thumbnailData) > 0 && thumbnailName != "" {
		var buf bytes.Buffer
		buf.Write(thumbnailData)
		destinatedPath := "posts/thumbnails/" + id
		url, err := utilsServices.UploadFile(buf, thumbnailName, destinatedPath)
		if err != nil {
			return nil, errors.New("Failed to upload thumbnail file")
		}
		uploadedURL = &url
	}

	c, err := u.repo.Update(ctx, cid, dto.ToDBPost{
		Title:            req.Title,
		Description:      req.Description,
		ShortDescription: req.ShortDescription,
		Price:            req.Price,
		DiscountRate:     req.DiscountRate,
		RemoveThumbnail:  req.RemoveThumbnail,
		LangID:           req.LangID,
		TopicIDs:         req.TopicIDs,
		ThumbnailURL: func() *string {
			if uploadedURL != nil {
				return uploadedURL
			}
			return req.ThumbnailURL
		}(),
	})
	if err != nil {
		return nil, err
	}

	// Clear existing relations
	if err := u.paramRepo.ClearParametersForModule(ctx, "post", c.ID); err != nil {
		return nil, err
	}

	// Re-assign relations: for simplicity, append new assignments (idempotency relies on unique checks if needed)
	if err := u.paramRepo.AssignParametersToModule(ctx, "post", c.ID, []uuid.UUID{req.LangID}); err != nil {
		return nil, err
	}
	if len(req.TopicIDs) > 0 {
		if err := u.paramRepo.AssignParametersToModule(ctx, "post", c.ID, req.TopicIDs); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (u *postUsecase) Delete(ctx context.Context, id string, authId string) error {
	cid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.Delete(ctx, cid)
}

func (u *postUsecase) GetByID(ctx context.Context, id string) (*models.Post, error) {
	cid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, cid)
}

func (u *postUsecase) GetParameterReferences(ctx context.Context, id string) (*dto.ReferenceObject, *dto.ReferenceObject, []dto.ReferenceObject, error) {
	cid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, nil, nil, err
	}
	params, err := u.paramRepo.GetByModule(ctx, "post", cid)
	if err != nil {
		return nil, nil, nil, err
	}
	var level *dto.ReferenceObject
	var lang *dto.ReferenceObject
	topics := make([]dto.ReferenceObject, 0)
	for _, p := range params {
		if p.Type == nil {
			continue
		}
		switch *p.Type {
		case "lang":
			lang = &dto.ReferenceObject{ID: p.ID, Name: p.Name}
		case "topic":
			topics = append(topics, dto.ReferenceObject{ID: p.ID, Name: p.Name})
		}
	}
	return level, lang, topics, nil
}
func (u *postUsecase) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqPostIndexFilter) ([]models.Post, int, error) {
	return u.repo.GetIndex(ctx, req, filter)
}

func (u *postUsecase) GetAll(ctx context.Context, filter dto.ReqPostIndexFilter) ([]models.Post, error) {
	return u.repo.GetAll(ctx, filter)
}

func (u *postUsecase) validateParameterType(ctx context.Context, id uuid.UUID, expectedType string) error {
	p, err := u.paramRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil || p.Type == nil || *p.Type != expectedType {
		return errors.New("invalid parameter type for " + expectedType)
	}
	return nil
}
