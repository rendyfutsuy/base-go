package dto

import (
	"time"

	"github.com/google/uuid"
)

type ReqCreateSubcob struct {
	CategoryID              string `json:"category_id"`
	CobID                   string `json:"cob_id"`
	Name                    string `json:"name"`
	Code                    string `json:"code"`
	Forms                   *string `json:"forms"`
	ActiveDate              string `json:"active_date"`
	IsHiddenFromFacultative *bool   `json:"is_hidden_from_facultative"`
	IsInactive              *bool   `json:"is_inactive"`
	IsFromWebCrediit        *bool   `json:"is_from_web_credit"`
}

type ToDBCreateSubcob struct {
	CategoryID              uuid.UUID  `json:"category_id"`
	CobID                   uuid.UUID  `json:"cob_id"`
	Name                    string     `json:"name"`
	Code                    string     `json:"code"`
	Forms                   *string     `json:"forms"`
	ActiveDate              *time.Time `json:"active_date"`
	IsHiddenFromFacultative *bool       `json:"is_hidden_from_facultative"`
	IsInactive              *bool       `json:"is_inactive"`
	IsFromWebCrediit        *bool       `json:"is_from_web_credit"`
	CreatedByID             string
}

func (r *ReqCreateSubcob) ToDBCreateSubcob(createdByID string) (*ToDBCreateSubcob, error) {
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

	cobID, err := uuid.Parse(r.CobID)
	if err != nil {
		return nil, err
	}

	return &ToDBCreateSubcob{
		CategoryID:              categoryID,
		CobID:                   cobID,
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