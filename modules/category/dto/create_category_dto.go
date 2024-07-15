package dto

type ReqCreateCategory struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (r *ReqCreateCategory) ToDBCreateCategory(code, createdByID string) ToDBCreateCategory {
	return ToDBCreateCategory{
		Name:        r.Name,
		Code:        code,
		Description: r.Description,
		CreatedByID: createdByID,
	}
}

type ToDBCreateCategory struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	CreatedByID string
}
