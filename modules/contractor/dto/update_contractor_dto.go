package dto

type ReqUpdateContractor struct {
	Name string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
}

func (r *ReqUpdateContractor) ToDBUpdateContractor(authId string) ToDBUpdateContractor {
	return ToDBUpdateContractor{
		Name:        r.Name,
		Address:     r.Address,
		UpdatedByID: authId,
	}
}

type ToDBUpdateContractor struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	UpdatedByID string
}
