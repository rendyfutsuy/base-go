package dto

type ReqCreateConveyance struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required"`
}

func (r *ReqCreateConveyance) ToDBCreateConveyance(code, authId string) ToDBCreateConveyance {
	return ToDBCreateConveyance{
		Name:        r.Name,
		Code:        code,
		Type:        r.Type,
		CreatedByID: authId,
	}
}

type ToDBCreateConveyance struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Type        string `json:"type"`
	CreatedByID string
}
