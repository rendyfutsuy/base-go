package dto

type ReqAuthUser struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserProfile struct {
	UserId           string   `json:"id"`
	Name             string   `json:"name"`
	Username         string   `json:"username"`
	Role             string   `json:"role"`
	Email            string   `json:"email"`
	IsFirstTimeLogin bool     `json:"is_first_time_login"`
	Permissions      []string `json:"permissions"`
}

type ReqUpdateProfile struct {
	Name string `form:"name" json:"name" validate:"required"`
}

type ReqUpdatePassword struct {
	OldPassword          string `form:"old_password" json:"old_password" validate:"required"`
	NewPassword          string `form:"new_password" json:"new_password" validate:"required"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" validate:"required,eqfield=NewPassword"`
}

type ReqResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ReqResetPassword struct {
	Password             string `json:"password" validate:"required,min=8,max=25"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}
