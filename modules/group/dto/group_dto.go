package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
)

type ReqCreateGroup struct {
	Name string `form:"name" json:"name" validate:"required,max=255,uppercase_letters"`
}

type ReqUpdateGroup struct {
	Name string `form:"name" json:"name" validate:"required,max=255,uppercase_letters"`
}

type RespGroup struct {
	ID        uuid.UUID `json:"id"`
	GroupCode string    `json:"group_code,omitempty"` // omit when creating new group
	Name      string    `json:"name"`
	Deletable bool      `json:"deletable"` // true if group is not used in any active sub-group
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToRespGroup(m models.Group) RespGroup {
	return RespGroup{
		ID:        m.ID,
		GroupCode: m.GroupCode,
		Name:      m.Name,
		Deletable: m.Deletable,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type RespGroupIndex struct {
	ID        uuid.UUID `json:"id"`
	GroupCode string    `json:"group_code"`
	Name      string    `json:"name"`
	Deletable bool      `json:"deletable"` // true if group is not used in any active sub-group
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToRespGroupIndex(m models.Group) RespGroupIndex {
	return RespGroupIndex{
		ID:        m.ID,
		GroupCode: m.GroupCode,
		Name:      m.Name,
		Deletable: m.Deletable,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// ReqGroupIndexFilter for filtering group index (prepared for future use)
type ReqGroupIndexFilter struct {
	Search    string `query:"search" json:"search"` // Search keyword for filtering by name and group_code
	SortBy    string `query:"sort_by" json:"sort_by"`
	SortOrder string `query:"sort_order" json:"sort_order"`
}
