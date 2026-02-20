package dto

type ReqSendVerification struct {
	Email string `json:"email" validate:"required,email"`
}

type ReqVerifyOTP struct {
	Email string `json:"email" validate:"required,email"`
	Token string `json:"token" validate:"required"`
}
