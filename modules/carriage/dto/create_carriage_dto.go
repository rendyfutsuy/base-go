package dto

type ReqCreateCarriage struct {
	Name string `json:"name" validate:"required"`
}

func (r *ReqCreateCarriage) ToDBCreateCarriage(code, authId string) ToDBCreateCarriage {
	return ToDBCreateCarriage{
		Name:        r.Name,
		Code:        code,
		CreatedByID: authId,
	}
}

type ToDBCreateCarriage struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	CreatedByID string
}
