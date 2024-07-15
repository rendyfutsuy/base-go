package dto

type ReqCreateRole struct {
	Name      string                                `json:"name" validate:"required"`
	Deletable bool                                  `json:"deletable"`
}
