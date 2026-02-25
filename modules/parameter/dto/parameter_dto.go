package dto

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
)

type ReqCreateParameter struct {
	Code     string     `form:"code" json:"code" validate:"required,max=255"`
	Name     string     `form:"name" json:"name" validate:"required,max=255"`
	Value    *string    `form:"value" json:"value,omitempty"`
	Type     *string    `form:"type" json:"type,omitempty"`
	Desc     *string    `form:"desc" json:"desc,omitempty"`
	ParentId *uuid.UUID `form:"parent_id" json:"parent_id,omitempty"`
}

type ReqUpdateParameter struct {
	Code     string     `form:"code" json:"code" validate:"required,max=255"`
	Name     string     `form:"name" json:"name" validate:"required,max=255"`
	Value    *string    `form:"value" json:"value,omitempty"`
	Type     *string    `form:"type" json:"type,omitempty"`
	Desc     *string    `form:"desc" json:"desc,omitempty"`
	ParentId *uuid.UUID `form:"parent_id" json:"parent_id,omitempty"`
}

type RespParameter struct {
	ID        uuid.UUID            `json:"id"`
	Code      string               `json:"code"`
	Name      string               `json:"name"`
	Value     *string              `json:"value,omitempty"`
	Type      *string              `json:"type,omitempty"`
	Desc      *string              `json:"desc,omitempty"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
	Parent    *RespParameterParent `json:"parent"`
	Deletable bool                 `json:"deletable"`
}

type RespParameterParent struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func ToRespParameter(m models.Parameter) RespParameter {
	value := m.Value
	typeVal := m.Type
	desc := m.Description
	var parent *RespParameterParent
	if m.ParentID != uuid.Nil && m.ParentName != "" {
		parent = &RespParameterParent{ID: m.ParentID, Name: m.ParentName}
	}
	return RespParameter{
		ID:        m.ID,
		Code:      m.Code,
		Name:      m.Name,
		Value:     value,
		Type:      typeVal,
		Desc:      desc,
		CreatedAt: m.CreatedAt.Format(constants.FormatDateTimeISO8601),
		UpdatedAt: m.UpdatedAt.Format(constants.FormatDateTimeISO8601),
		Parent:    parent,
		Deletable: m.Deletable,
	}
}

type RespParameterIndex struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Value     *string   `json:"value,omitempty"`
	Type      *string   `json:"type,omitempty"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Deletable bool      `json:"deletable"`
}

func ToRespParameterIndex(m models.Parameter) RespParameterIndex {
	value := m.Value
	typeVal := m.Type
	return RespParameterIndex{
		ID:        m.ID,
		Code:      m.Code,
		Name:      m.Name,
		Value:     value,
		Type:      typeVal,
		CreatedAt: m.CreatedAt.Format(constants.FormatDateTimeISO8601),
		UpdatedAt: m.UpdatedAt.Format(constants.FormatDateTimeISO8601),
		Deletable: m.Deletable,
	}
}

// ReqParameterIndexFilter for filtering parameter index (prepared for future use)
type ReqParameterIndexFilter struct {
	Search    string      `query:"search"` // Search keyword for filtering by name and code
	Types     []string    `query:"types"`
	Names     []string    `query:"names"`
	IDs       []uuid.UUID `query:"ids"` // Filter by parameter IDs
	SortBy    string      `query:"sort_by"`
	SortOrder string      `query:"sort_order"`
}
