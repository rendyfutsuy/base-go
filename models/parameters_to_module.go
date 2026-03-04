package models

import (
	"time"

	"github.com/google/uuid"
)

type ParametersToModule struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v7()" json:"id"`
	ParameterID uuid.UUID `gorm:"column:parameter_id;type:uuid;not null" json:"parameter_id"`
	ModuleType  string    `gorm:"column:module_type;type:varchar(255);not null" json:"module_type"`
	ModuleID    uuid.UUID `gorm:"column:module_id;type:uuid;not null" json:"module_id"`
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (ParametersToModule) TableName() string {
	return "parameters_to_module"
}
