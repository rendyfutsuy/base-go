package dto

type ReqUpdateCarriage struct {
	Name string `json:"name" validate:"required"`
}

func (r *ReqUpdateCarriage) ToDBUpdateCarriage(authId string) ToDBUpdateCarriage {
	return ToDBUpdateCarriage{
		Name:        r.Name,
		UpdatedByID: authId,
	}
}

type ToDBUpdateCarriage struct {
	Name        string `json:"name"`
	UpdatedByID string
}
