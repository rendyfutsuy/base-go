package dto

import (
	"github.com/google/uuid"
)

type ReqUpdateUserPassword struct {
	NewPassword          string `form:"new_password" json:"new_password" validate:"required"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" validate:"required,eqfield=NewPassword"`
}

type ReqBlockUser struct {
	IsBlock bool `form:"is_block" json:"is_block"`
}

type ReqActivateUser struct {
	IsActive bool `form:"is_active" json:"is_active"`
}

type ReqUpdateUser struct {
	FullName             string    `form:"name" json:"name" validate:"required,max=80"`
	Username             string    `form:"username" json:"username" validate:"required"`
	RoleId               uuid.UUID `form:"role_id" json:"role_id" validate:"required"`
	Email                string    `form:"email" json:"email"`
	IsActive             bool      `form:"is_active" json:"is_active"`
	Gender               string    `form:"gender" json:"gender"`
	Password             string    `form:"password" json:"password"`
	PasswordConfirmation string    `form:"password_confirmation" json:"password_confirmation"`
}

func (r *ReqUpdateUser) ToDBUpdateUser(authId string) ToDBUpdateUser {
	return ToDBUpdateUser{
		FullName: r.FullName,
		Username: r.Username,
		RoleId:   r.RoleId,
		Email:    r.Email,
		IsActive: r.IsActive,
		Gender:   r.Gender,
	}
}

type ToDBUpdateUser struct {
	FullName string    `json:"name"`
	Username string    `json:"username"`
	RoleId   uuid.UUID `json:"role_id"`
	Email    string    `json:"email"`
	IsActive bool      `json:"is_active"`
	Gender   string    `json:"gender"`
}
