package dto

import (
	"github.com/rendyfutsuy/base-go/utils"
)

type ReqAuthUser struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserProfile struct {
	UserId string           `json:"user_id"`
	Name   string           `json:"name"`
	Role   utils.NullString `json:"role"`
	Email  string           `json:"email"`
	Status string           `json:"status"`
	Gender string           `json:"gender"`
}

type ReqUpdateProfile struct {
	Name string `json:"name" validate:"required"`
}

type ReqUpdatePassword struct {
	OldPassword          string `json:"old_password" validate:"required"`
	NewPassword          string `json:"new_password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=NewPassword"`
}

type ReqResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ReqResetPassword struct {
	Password             string `json:"password" validate:"required,min=8,max=25"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}
