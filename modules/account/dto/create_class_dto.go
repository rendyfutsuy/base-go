package dto

import "github.com/google/uuid"

type ReqCreateAccount struct {
	Name string `json:"name" validate:"required"`
}

type ToDBCreateAccount struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	CreatedByID uuid.UUID
}
