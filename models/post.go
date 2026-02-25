package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id"`
	CreatedBy        uuid.UUID `gorm:"column:created_by;type:uuid;not null" json:"created_by"`
	Title            string    `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Description      string    `gorm:"column:description;type:text;not null" json:"description"`
	ShortDescription string    `gorm:"column:short_description;type:varchar(255);not null" json:"short_description"`
	Price            float64   `gorm:"column:price;type:numeric(18,2);not null" json:"price"`
	DiscountRate     float64   `gorm:"column:discount_rate;type:numeric(5,2);not null" json:"discount_rate"`
	ThumbnailURL     *string   `gorm:"column:thumbnail_url;type:text" json:"thumbnail_url"`
	CreatedAt        time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (Post) TableName() string {
	return "posts"
}
