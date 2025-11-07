package dto

import (
	"github.com/google/uuid"
)

type ReqCheckDuplicatedUser struct {
	FullName       string    `json:"name" validate:"required"`
	ExcludedUserId uuid.UUID `json:"excluded_user_id"`
}

type ReqCreateUser struct {
	FullName             string    `form:"name" json:"name" validate:"required,max=80,uppercase_letters"`
	Username             string    `form:"username" json:"username" validate:"required,uppercase_letters"`
	RoleId               uuid.UUID `form:"role_id" json:"role_id" validate:"required"`
	Email                string    `form:"email" json:"email"`
	IsActive             bool      `form:"is_active" json:"is_active"`
	Gender               string    `form:"gender" json:"gender"`
	Password             string    `form:"password" json:"password" validate:"required,min=8,password_uppercase"`
	PasswordConfirmation string    `form:"password_confirmation" json:"password_confirmation" validate:"required,eqfield=Password"`
}

func (r *ReqCreateUser) ToDBCreateUser(code, authId string) ToDBCreateUser {
	return ToDBCreateUser{
		FullName: r.FullName,
		Username: r.Username,
		RoleId:   r.RoleId,
		Email:    r.Email,
		IsActive: r.IsActive,
		Gender:   r.Gender,
		Password: r.Password,
	}
}

type ToDBCreateUser struct {
	FullName string    `json:"name"`
	Username string    `json:"username"`
	RoleId   uuid.UUID `json:"role_id"`
	Email    string    `json:"email"`
	Nik      string    `json:"nik"`
	IsActive bool      `json:"is_active"`
	Gender   string    `json:"gender"`
	Password string    `json:"password"`
}
