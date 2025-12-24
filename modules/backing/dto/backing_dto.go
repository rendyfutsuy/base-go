package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
)

type ReqCreateBacking struct {
	TypeID uuid.UUID `form:"type_id" json:"type_id" validate:"required"`
	Name   string    `form:"name" json:"name" validate:"required,max=255"`
}

type ReqUpdateBacking struct {
	TypeID uuid.UUID `form:"type_id" json:"type_id" validate:"required"`
	Name   string    `form:"name" json:"name" validate:"required,max=255"`
}

type RespBacking struct {
	ID           uuid.UUID `json:"id"`
	TypeID       uuid.UUID `json:"type_id"`
	BackingCode  string    `json:"backing_code,omitempty"` // omit when creating new backing
	Name         string    `json:"name"`
	TypeName     string    `json:"type_name"`
	SubgroupID   uuid.UUID `json:"subgroup_id"`
	SubgroupName string    `json:"subgroup_name"`
	GroupID      uuid.UUID `json:"groups_id"`
	GroupName    string    `json:"groups_name"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    string    `json:"updated_by"`
	Deletable    bool      `json:"deletable"`
}

func ToRespBacking(m models.Backing) RespBacking {
	return RespBacking{
		ID:           m.ID,
		TypeID:       m.TypeID,
		BackingCode:  m.BackingCode,
		Name:         m.Name,
		TypeName:     m.TypeName,
		SubgroupID:   m.SubgroupID,
		SubgroupName: m.SubgroupName,
		GroupID:      m.GroupID,
		GroupName:    m.GroupName,
		CreatedAt:    m.CreatedAt,
		CreatedBy:    m.CreatedBy,
		UpdatedAt:    m.UpdatedAt,
		UpdatedBy:    m.UpdatedBy,
		Deletable:    true, // placeholder
	}
}

type RespBackingIndex struct {
	ID           uuid.UUID `json:"id"`
	TypeID       uuid.UUID `json:"type_id"`
	BackingCode  string    `json:"backing_code"`
	Name         string    `json:"name"`
	TypeName     string    `json:"type_name"`
	SubgroupName string    `json:"subgroup_name"`
	GroupName    string    `json:"group_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Deletable    bool      `json:"deletable"`
}

func ToRespBackingIndex(m models.Backing) RespBackingIndex {
	return RespBackingIndex{
		ID:           m.ID,
		TypeID:       m.TypeID,
		BackingCode:  m.BackingCode,
		Name:         m.Name,
		TypeName:     m.TypeName,
		SubgroupName: m.SubgroupName,
		GroupName:    m.GroupName,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		Deletable:    true, // placeholder
	}
}

// ReqBackingIndexFilter for filtering backing index with multiple values support
type ReqBackingIndexFilter struct {
	Search       string      `query:"search" json:"search"`               // Search keyword for filtering by backing_code, name
	BackingCodes []string    `query:"backing_codes" json:"backing_codes"` // Multiple values
	Names        []string    `query:"names" json:"names"`                 // Multiple values
	TypeIDs      []uuid.UUID `query:"type_ids" json:"type_ids"`           // Multiple values
	SubgroupIDs  []string    `query:"subgroup_ids" json:"subgroup_ids"`   // Filter by subgroup IDs (multiple values, UUIDs as strings)
	GoodGroupIDs []string    `query:"group_ids" json:"group_ids"`         // Filter by goodGroup IDs (multiple values, UUIDs as strings)
	SortBy       string      `query:"sort_by" json:"sort_by"`
	SortOrder    string      `query:"sort_order" json:"sort_order"`
}
