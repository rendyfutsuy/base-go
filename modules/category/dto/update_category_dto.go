package dto

type ReqUpdateCategory struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (r *ReqUpdateCategory) ToDBUpdateCategory(updatedByID string) ToDBUpdateCategory {
	return ToDBUpdateCategory{
		Name:        r.Name,
		Description: r.Description,
		UpdatedByID: updatedByID,
	}
}

type ToDBUpdateCategory struct {
	Name        string `json:"name"`
	Description string `json:"description" validate:"required"`
	UpdatedByID string
}
