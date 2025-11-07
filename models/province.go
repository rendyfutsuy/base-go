package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Province represents province table
type Province struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id" validate:"required"`
	Name      string         `gorm:"column:name;type:varchar(100);not null;uniqueIndex" json:"name" validate:"required"`
	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Province) TableName() string {
	return "province"
}

// City represents city table
type City struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id" validate:"required"`
	ProvinceID uuid.UUID      `gorm:"column:province_id;type:uuid;not null" json:"province_id" validate:"required"`
	Province   Province       `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
	Name       string         `gorm:"column:name;type:varchar(255);not null" json:"name" validate:"required"`
	CreatedAt  time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (City) TableName() string {
	return "city"
}

// District represents district table
type District struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id" validate:"required"`
	CityID    uuid.UUID      `gorm:"column:city_id;type:uuid;not null" json:"city_id" validate:"required"`
	City      City           `gorm:"foreignKey:CityID" json:"city,omitempty"`
	Name      string         `gorm:"column:name;type:varchar(255);not null" json:"name" validate:"required"`
	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (District) TableName() string {
	return "district"
}

// Subdistrict represents subdistrict table
type Subdistrict struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id" validate:"required"`
	DistrictID uuid.UUID      `gorm:"column:district_id;type:uuid;not null" json:"district_id" validate:"required"`
	District   District       `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Name       string         `gorm:"column:name;type:varchar(255);not null" json:"name" validate:"required"`
	CreatedAt  time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Subdistrict) TableName() string {
	return "subdistrict"
}

