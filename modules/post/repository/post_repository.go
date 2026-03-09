package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/post"
	"github.com/rendyfutsuy/base-go/modules/post/dto"
	csearch "github.com/rendyfutsuy/base-go/modules/post/repository/searches"
	"gorm.io/gorm"
)

type postRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) post.Repository {
	return &postRepository{DB: db}
}

func (r *postRepository) Create(ctx context.Context, createdBy uuid.UUID, data dto.ToDBPost) (*models.Post, error) {
	now := time.Now().UTC()
	c := &models.Post{
		CreatedBy:        createdBy,
		Title:            data.Title,
		Description:      data.Description,
		ShortDescription: data.ShortDescription,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := r.DB.WithContext(ctx).Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *postRepository) Update(ctx context.Context, id uuid.UUID, data dto.ToDBPost) (*models.Post, error) {
	updates := map[string]interface{}{
		"title":             data.Title,
		"description":       data.Description,
		"short_description": data.ShortDescription,
		"updated_at":        time.Now().UTC(),
	}

	c := &models.Post{}
	err := r.DB.WithContext(ctx).Model(&models.Post{}).
		Where("id = ?", id).
		Updates(updates).
		First(c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ?", id).Delete(&models.Post{}).Error
}

func (r *postRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	c := &models.Post{}
	if err := r.DB.WithContext(ctx).
		Table("posts c").
		Select(`c.id, c.created_by, c.title, c.description, c.short_description,
			(SELECT f.file_path FROM files_to_module ftm
				JOIN files f ON f.id = ftm.file_id AND f.deleted_at IS NULL
				WHERE ftm.module_type = ? AND ftm.module_id = c.id AND ftm.type = ?
				ORDER BY ftm.created_at DESC
				LIMIT 1
			) AS thumbnail_url,
			c.created_at, c.updated_at`,
			constants.ModuleTypePost, constants.FileTypeThumbnail,
		).
		Where("c.id = ?", id).
		First(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *postRepository) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqPostIndexFilter) ([]models.Post, int, error) {
	var posts []models.Post
	query := r.DB.WithContext(ctx).
		Table("posts c").
		Select(`c.id, c.title, c.short_description,
			(SELECT f.file_path FROM files_to_module ftm
				JOIN files f ON f.id = ftm.file_id AND f.deleted_at IS NULL
				WHERE ftm.module_type = ? AND ftm.module_id = c.id AND ftm.type = ?
				ORDER BY ftm.created_at DESC
				LIMIT 1
			) AS thumbnail_url,
			c.created_at`,
			constants.ModuleTypePost, constants.FileTypeThumbnail,
		).
		Where("1=1")

	// Search support
	query = request.ApplySearchConditionFromInterface(query, req.Search, csearch.NewPostSearchHelper())

	// Filters by parameter relations
	if len(filter.LangIDs) > 0 {
		query = query.Where(`
			EXISTS (
				SELECT 1 FROM parameters_to_module ptm
				JOIN parameters p ON p.id = ptm.parameter_id
				WHERE ptm.module_type = ?
				  AND ptm.module_id = c.id
				  AND p.type = 'lang'
				  AND p.id IN (?)
			)
		`, constants.ModuleTypePost, filter.LangIDs)
	}
	if len(filter.TopicIDs) > 0 {
		query = query.Where(`
			EXISTS (
				SELECT 1 FROM parameters_to_module ptm
				JOIN parameters p ON p.id = ptm.parameter_id
				WHERE ptm.module_type = ?
				  AND ptm.module_id = c.id
				  AND p.type = 'topic'
				  AND p.id IN (?)
			)
		`, constants.ModuleTypePost, filter.TopicIDs)
	}

	// Pagination
	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:    "c.created_at",
		DefaultSortOrder: "DESC",
		MaxPerPage:       100,
		SortMapping: func(s string) string {
			switch s {
			case "id":
				return "c.id"
			case "title":
				return "c.title"
			case "short_description":
				return "c.short_description"
			case "created_at":
				return "c.created_at"
			default:
				return ""
			}
		},
		NaturalSortColumns: []string{"c.title"},
	}, &posts)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (r *postRepository) GetAll(ctx context.Context, filter dto.ReqPostIndexFilter) ([]models.Post, error) {
	var posts []models.Post
	query := r.DB.WithContext(ctx).
		Table("posts c").
		Select(`c.id, c.title, c.short_description,
			(SELECT f.file_path FROM files_to_module ftm
				JOIN files f ON f.id = ftm.file_id AND f.deleted_at IS NULL
				WHERE ftm.module_type = ? AND ftm.module_id = c.id AND ftm.type = ?
				ORDER BY ftm.created_at DESC
				LIMIT 1
			) AS thumbnail_url,
			c.created_at`,
			constants.ModuleTypePost, constants.FileTypeThumbnail,
		).
		Where("1=1")

	// Search support
	query = request.ApplySearchConditionFromInterface(query, filter.Search, csearch.NewPostSearchHelper())

	// Filters (same as index)
	if len(filter.LangIDs) > 0 {
		query = query.Where(`
			EXISTS (
				SELECT 1 FROM parameters_to_module ptm
				JOIN parameters p ON p.id = ptm.parameter_id
				WHERE ptm.module_type = ?
				  AND ptm.module_id = c.id
				  AND p.type = 'lang'
				  AND p.id IN (?)
			)
		`, constants.ModuleTypePost, filter.LangIDs)
	}
	if len(filter.TopicIDs) > 0 {
		query = query.Where(`
			EXISTS (
				SELECT 1 FROM parameters_to_module ptm
				JOIN parameters p ON p.id = ptm.parameter_id
				WHERE ptm.module_type = ?
				  AND ptm.module_id = c.id
				  AND p.type = 'topic'
				  AND p.id IN (?)
			)
		`, constants.ModuleTypePost, filter.TopicIDs)
	}

	if err := query.Order("c.created_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
