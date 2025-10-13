package dto

import "github.com/google/uuid"

type ReqCreatePermission struct {
	Name string `json:"name" validate:"required"`
}

type ToDBCreatePermission struct {
	Name        string `json:"name"`
	CreatedByID uuid.UUID
}
