package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"gorm.io/gorm"
)

type Post struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id"`
	CreatedBy        uuid.UUID `gorm:"column:created_by;type:uuid;not null" json:"created_by"`
	Title            string    `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Description      string    `gorm:"column:description;type:text;not null" json:"description"`
	ShortDescription string    `gorm:"column:short_description;type:varchar(255);not null" json:"short_description"`
	Price            float64   `gorm:"column:price;type:numeric(18,2);not null" json:"price"`
	DiscountRate     float64   `gorm:"column:discount_rate;type:numeric(5,2);not null" json:"discount_rate"`
	ThumbnailURL     *string   `gorm:"column:deletable;<-:false" json:"thumbnail_url"`
	Files            []File    `gorm:"many2many:files_to_module;joinForeignKey:ID;joinReferences:FileID" json:"files"`
	CreatedAt        time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (Post) TableName() string {
	return "posts"
}

// AfterFind computes ThumbnailURL from files_to_module pivot with type = "thumbnail"
func (p *Post) AfterFind(tx *gorm.DB) (err error) {
	if p == nil || p.ID == uuid.Nil {
		return nil
	}
	var filePath *string
	err = tx.Table("files_to_module ftm").
		Select("f.file_path").
		Joins("JOIN files f ON f.id = ftm.file_id AND f.deleted_at IS NULL").
		Where("ftm.module_type = ? AND ftm.module_id = ? AND ftm.type = ?", constants.ModuleTypePost, p.ID, constants.FileTypeThumbnail).
		Order("ftm.created_at DESC").
		Limit(1).
		Scan(&filePath).Error
	if err != nil {
		return nil // do not block read on error
	}
	if filePath != nil {
		p.ThumbnailURL = filePath
	}
	return nil
}
