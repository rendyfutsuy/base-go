package dto

import (
	"github.com/google/uuid"
)

type ReqUpdateUserPassword struct {
	OldPassword          string `json:"old_password" validate:"required"`
	NewPassword          string `json:"new_password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=NewPassword"`
}

type ReqBlockUser struct {
	IsBlock bool `json:"is_block"`
}

type ReqActivateUser struct {
	IsActive bool `json:"is_active"`
}

type ReqUpdateUser struct {
	FullName string    `json:"name" validate:"required,max=80"`
	RoleId   uuid.UUID `json:"role_id" validate:"required"`
	Email    string    `json:"email" validate:"required,email,emaildomain"`
	IsActive bool      `json:"is_active"`
	Gender   string    `json:"gender" validate:"required,oneof='male' 'female'"`
}

func (r *ReqUpdateUser) ToDBUpdateUser(authId string) ToDBUpdateUser {
	return ToDBUpdateUser{
		FullName: r.FullName,
		RoleId:   r.RoleId,
		Email:    r.Email,
		IsActive: r.IsActive,
		Gender:   r.Gender,
	}
}

type ToDBUpdateUser struct {
	FullName string    `json:"name"`
	RoleId   uuid.UUID `json:"role_id"`
	Email    string    `json:"email"`
	IsActive bool      `json:"is_active"`
	Gender   string    `json:"gender"`
}
