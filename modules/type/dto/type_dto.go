package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
)

type ReqCreateType struct {
	SubgroupID uuid.UUID `form:"subgroup_id" json:"subgroup_id" validate:"required"`
	Name       string    `form:"name" json:"name" validate:"required,max=255"`
}

type ReqUpdateType struct {
	SubgroupID uuid.UUID `form:"subgroup_id" json:"subgroup_id" validate:"required"`
	Name       string    `form:"name" json:"name" validate:"required,max=255"`
}

type RespType struct {
	ID           uuid.UUID `json:"id"`
	SubgroupID   uuid.UUID `json:"subgroup_id"`
	SubgroupName string    `json:"subgroup_name"`
	GoodsGroupID uuid.UUID `json:"goods_group_id"`
	TypeCode     string    `json:"type_code"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToRespType(m models.Type) RespType {
	return RespType{
		ID:           m.ID,
		SubgroupID:   m.SubgroupID,
		SubgroupName: m.SubgroupName,
		GoodsGroupID: m.GoodsGroupID,
		TypeCode:     m.TypeCode,
		Name:         m.Name,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

type RespTypeIndex struct {
	ID             uuid.UUID `json:"id"`
	SubgroupID     uuid.UUID `json:"subgroup_id"`
	SubgroupName   string    `json:"subgroup_name"`
	GoodsGroupName string    `json:"goods_group_name"`
	TypeCode       string    `json:"type_code"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ToRespTypeIndex(m models.Type) RespTypeIndex {
	return RespTypeIndex{
		ID:             m.ID,
		SubgroupID:     m.SubgroupID,
		SubgroupName:   m.SubgroupName,
		GoodsGroupName: m.GoodsGroupName,
		TypeCode:       m.TypeCode,
		Name:           m.Name,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

// ReqTypeIndexFilter for filtering type index
type ReqTypeIndexFilter struct {
	Search      string   `query:"search"`       // Search keyword for filtering by type_code, name
	TypeCodes   []string `query:"type_codes"`   // Filter by type codes (multiple values)
	SubgroupIDs []string `query:"subgroup_ids"` // Filter by subgroup IDs (multiple values, UUIDs as strings)
	Names       []string `query:"names"`        // Filter by names (multiple values)
}
