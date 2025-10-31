package dto

import (
	"github.com/google/uuid"
)

type ReqCheckDuplicatedUser struct {
	FullName       string    `json:"name" validate:"required"`
	ExcludedUserId uuid.UUID `json:"excluded_user_id"`
}

type ReqCreateUser struct {
	FullName             string    `json:"name" validate:"required,max=80"`
	RoleId               uuid.UUID `json:"role_id" validate:"required"`
	Email                string    `json:"email" validate:"required,email,emaildomain"`
	IsActive             bool      `json:"is_active"`
	Gender               string    `json:"gender" validate:"required,oneof='male' 'female'"`
	Password             string    `json:"password" validate:"required"`
	PasswordConfirmation string    `json:"password_confirmation" validate:"required,eqfield=Password"`
}

func (r *ReqCreateUser) ToDBCreateUser(code, authId string) ToDBCreateUser {
	return ToDBCreateUser{
		FullName: r.FullName,
		RoleId:   r.RoleId,
		Email:    r.Email,
		IsActive: r.IsActive,
		Gender:   r.Gender,
		Password: r.Password,
	}
}

type ToDBCreateUser struct {
	FullName string    `json:"name"`
	RoleId   uuid.UUID `json:"role_id"`
	Email    string    `json:"email"`
	IsActive bool      `json:"is_active"`
	Gender   string    `json:"gender"`
	Password string    `json:"password"`
}
