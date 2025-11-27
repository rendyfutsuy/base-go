package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
)

type ReqCreateSubGroup struct {
	GoodsGroupID uuid.UUID `form:"goods_group_id" json:"goods_group_id" validate:"required"`
	Name         string    `form:"name" json:"name" validate:"required,max=255,uppercase_letters"`
}

type ReqUpdateSubGroup struct {
	GoodsGroupID uuid.UUID `form:"goods_group_id" json:"goods_group_id" validate:"required"`
	Name         string    `form:"name" json:"name" validate:"required,max=255,uppercase_letters"`
}

type RespSubGroup struct {
	ID             uuid.UUID `json:"id"`
	GoodsGroupID   uuid.UUID `json:"goods_group_id"`
	GoodsGroupName string    `json:"goods_group_name"`
	SubgroupCode   string    `json:"subgroup_code"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      string    `json:"updated_by"`
	Deletable      bool      `json:"deletable"`
}

func ToRespSubGroup(m models.SubGroup) RespSubGroup {
	return RespSubGroup{
		ID:             m.ID,
		GoodsGroupID:   m.GoodsGroupID,
		GoodsGroupName: m.GoodsGroupName,
		SubgroupCode:   m.SubgroupCode,
		Name:           m.Name,
		CreatedAt:      m.CreatedAt,
		CreatedBy:      m.CreatedBy,
		UpdatedAt:      m.UpdatedAt,
		UpdatedBy:      m.UpdatedBy,
		Deletable:      m.Deletable,
	}
}

type RespSubGroupIndex struct {
	ID             uuid.UUID `json:"id"`
	GoodsGroupID   uuid.UUID `json:"goods_group_id"`
	GoodsGroupName string    `json:"goods_group_name"`
	SubgroupCode   string    `json:"subgroup_code"`
	Name           string    `json:"name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Deletable      bool      `json:"deletable"`
}

func ToRespSubGroupIndex(m models.SubGroup) RespSubGroupIndex {
	return RespSubGroupIndex{
		ID:             m.ID,
		GoodsGroupID:   m.GoodsGroupID,
		GoodsGroupName: m.GoodsGroupName,
		SubgroupCode:   m.SubgroupCode,
		Name:           m.Name,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		Deletable:      m.Deletable,
	}
}

// ReqSubGroupIndexFilter for filtering sub-group index with multiple values support
type ReqSubGroupIndexFilter struct {
	Search        string      `query:"search" json:"search"`                   // Search keyword for filtering by subgroup_code, name
	SubgroupCodes []string    `query:"subgroup_codes" json:"subgroup_codes"`   // Multiple values
	Names         []string    `query:"names" json:"names"`                     // Multiple values
	GoodsGroupIDs []uuid.UUID `query:"goods_group_ids" json:"goods_group_ids"` // Multiple values
	SortBy        string      `query:"sort_by" json:"sort_by"`
	SortOrder     string      `query:"sort_order" json:"sort_order"`
}
