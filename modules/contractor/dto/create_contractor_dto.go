package dto

type ReqCreateContractor struct {
	Name string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
}

func (r *ReqCreateContractor) ToDBCreateContractor(code, authId string) ToDBCreateContractor {
	return ToDBCreateContractor{
		Name:        r.Name,
		Code:        code,
		Address:     r.Address,
		CreatedByID: authId,
	}
}

type ToDBCreateContractor struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Address     string `json:"address"`
	CreatedByID string
}
