package dto

import (
	"time"

	"github.com/google/uuid"
)

type ReqCreateCob struct {
	CategoryID              string  `json:"category_id" validate:"required"`
	Name                    string  `json:"name" validate:"required"`
	Code                    string  `json:"code" validate:"required"`
	Forms                   *string `json:"forms"`
	ActiveDate              string  `json:"active_date"`
	IsHiddenFromFacultative *bool    `json:"is_hidden_from_facultative"`
	IsInactive              *bool    `json:"is_inactive"`
	IsFromWebCrediit        *bool    `json:"is_from_web_credit"`
}

type ToDBCreateCob struct {
	CategoryID              uuid.UUID  `json:"category_id"`
	Name                    string     `json:"name"`
	Code                    string     `json:"code"`
	Forms                   *string    `json:"forms"`
	ActiveDate              *time.Time `json:"active_date"`
	IsHiddenFromFacultative *bool       `json:"is_hidden_from_facultative"`
	IsInactive              *bool       `json:"is_inactive"`
	IsFromWebCrediit        *bool       `json:"is_from_web_credit"`
	CreatedByID             string
}

func (r *ReqCreateCob) ToDBCreateCob(createdByID string) (*ToDBCreateCob, error) {
	var (
		activeDateNil *time.Time
		err           error
	)
	activeDate, err := time.Parse("2006-01-02", r.ActiveDate)
	if err != nil {
		activeDateNil = nil
	} else {
		activeDateNil = &activeDate
	}

	categoryID, err := uuid.Parse(r.CategoryID)
	if err != nil {
		return nil, err
	}

	return &ToDBCreateCob{
		CategoryID:              categoryID,
		Name:                    r.Name,
		Code:                    r.Code,
		Forms:                   r.Forms,
		ActiveDate:              activeDateNil,
		IsHiddenFromFacultative: r.IsHiddenFromFacultative,
		IsInactive:              r.IsInactive,
		IsFromWebCrediit:        r.IsFromWebCrediit,
		CreatedByID:             createdByID,
	}, err
}

type CobJson struct {
	ID                      MongoID       `json:"_id"`
	Category                CategoryField `json:"category"`
	Parent                  *ParentField  `json:"parent,omitempty"`
	ActiveDate              *time.Time    `json:"activedate,omitempty"`
	Code                    string        `json:"code"`
	Name                    string        `json:"name"`
	Forms                   *string       `json:"forms,omitempty"`
	IsHiddenFromFacultative *bool         `json:"ishiddenfromfacultative"`
	IsInactive              *bool         `json:"isinactive"`
	IsFromWebCrediit        *bool         `json:"is_from_web_credit"`
}

type CategoryField struct {
	ID MongoID `json:"_id"`
}

type ParentField struct {
	ID   *MongoID `json:"_id,omitempty"`
	Name string   `json:"name,omitempty"`
}

func (c *CobJson) JSONToDBCreateCob(createdByID, categoryID string) (*ToDBCreateCob, error) {
	var (
		activeDateNil *time.Time
		err           error
	)
	if c.ActiveDate != nil {
		activeDateNil = c.ActiveDate
	}

	catId, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}

	return &ToDBCreateCob{
		CategoryID:              catId,
		Name:                    c.Name,
		Code:                    c.Code,
		Forms:                   c.Forms,
		ActiveDate:              activeDateNil,
		IsHiddenFromFacultative: c.IsHiddenFromFacultative,
		IsInactive:              c.IsInactive,
		IsFromWebCrediit:        c.IsFromWebCrediit,
		CreatedByID:             createdByID,
	}, err
}

func (c *CobJson) JSONToDBCreateSubcob(createdByID, categoryID, parentID string) (*ToDBCreateSubcob, error) {
	var (
		activeDateNil *time.Time
		err           error
	)
	if c.ActiveDate != nil {
		activeDateNil = c.ActiveDate
	}

	catId, err := uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}

	parentId, err := uuid.Parse(parentID)
	if err != nil {
		return nil, err
	}

	return &ToDBCreateSubcob{
		CategoryID:              catId,
		CobID:                   parentId,
		Name:                    c.Name,
		Code:                    c.Code,
		Forms:                   c.Forms,
		ActiveDate:              activeDateNil,
		IsHiddenFromFacultative: c.IsHiddenFromFacultative,
		IsInactive:              c.IsInactive,
		IsFromWebCrediit:        c.IsFromWebCrediit,
		CreatedByID:             createdByID,
	}, err
}
