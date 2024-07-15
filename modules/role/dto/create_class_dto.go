package dto

import "github.com/google/uuid"

type ReqCreateRole struct {
	Name string `json:"name" validate:"required"`
}

type ToDBCreateRole struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	CreatedByID uuid.UUID
}
