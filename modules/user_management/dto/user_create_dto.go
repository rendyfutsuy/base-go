package dto

import (
	"github.com/google/uuid"
)

type ReqCheckDuplicatedUser struct {
	UserName       string    `json:"username" validate:"required"`
	ExcludedUserId uuid.UUID `json:"excluded_user_id"`
}

type ReqCheckDuplicatedEmail struct {
	Email          string    `json:"email" validate:"required,email"`
	ExcludedUserId uuid.UUID `json:"excluded_user_id"`
}

type ReqCreateUser struct {
	FullName             string    `form:"name" json:"name" validate:"required,max=80"`
	Username             string    `form:"username" json:"username" validate:"required"`
	RoleId               uuid.UUID `form:"role_id" json:"role_id" validate:"required"`
	Email                string    `form:"email" json:"email" validate:"required,email"`
	NIK                  string    `form:"nik" json:"nik" validate:"required"`
	Password             string    `form:"password" json:"password" validate:"required,min=8"`
	PasswordConfirmation string    `form:"password_confirmation" json:"password_confirmation" validate:"required,eqfield=Password"`
}

type ReqRegisterUser struct {
	FullName             string `form:"name" json:"name" validate:"required,max=80"`
	Username             string `form:"username" json:"username" validate:"required"`
	Email                string `form:"email" json:"email" validate:"required,email"`
	NIK                  string `form:"nik" json:"nik" validate:"required"`
	Password             string `form:"password" json:"password" validate:"required,min=8"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" validate:"required,eqfield=Password"`
}

func (r *ReqCreateUser) ToDBCreateUser(code, authId string) ToDBCreateUser {
	return ToDBCreateUser{
		FullName:         r.FullName,
		Username:         r.Username,
		RoleId:           r.RoleId,
		Email:            r.Email,
		Password:         r.Password,
		Nik:              r.NIK,
		IsVerifiedNow:    true,
		IsFirstTimeLogin: true, // Explicitly set to true for new users
	}
}

func (r *ReqRegisterUser) ToDBRegisterUser(code, authId string, roleId uuid.UUID) ToDBCreateUser {
	return ToDBCreateUser{
		FullName:         r.FullName,
		Username:         r.Username,
		RoleId:           roleId, // by default its should be filled with role "USER"
		Email:            r.Email,
		Password:         r.Password,
		Nik:              r.NIK,
		IsVerifiedNow:    false,
		IsFirstTimeLogin: false, // Explicitly set to false for registered users
	}
}

type ToDBCreateUser struct {
	FullName         string    `json:"name"`
	Username         string    `json:"username"`
	RoleId           uuid.UUID `json:"role_id"`
	Email            string    `json:"email"`
	Nik              string    `json:"nik"`
	IsActive         bool      `json:"is_active"`
	Gender           string    `json:"gender"`
	Password         string    `json:"password"`
	IsVerifiedNow    bool      `json:"is_verified_now"`
	IsFirstTimeLogin bool      `json:"is_first_time_login"`
}
