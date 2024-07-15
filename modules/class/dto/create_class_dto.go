package dto

type ReqCreateClass struct {
	Name string `json:"name" validate:"required"`
}

func (r *ReqCreateClass) ToDBCreateClass(code, authId string) ToDBCreateClass {
	return ToDBCreateClass{
		Name:        r.Name,
		Code:        code,
		CreatedByID: authId,
	}
}

type ToDBCreateClass struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	CreatedByID string
}
