package dto

type ReqUpdateConveyance struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required"`
}

func (r *ReqUpdateConveyance) ToDBUpdateConveyance(authId string) ToDBUpdateConveyance {
	return ToDBUpdateConveyance{
		Name:        r.Name,
		Type:        r.Type,
		UpdatedByID: authId,
	}
}

type ToDBUpdateConveyance struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	UpdatedByID string
}
