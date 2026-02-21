package dto

type ReqVerifyOTP struct {
	Token string `json:"token" validate:"required"`
}
