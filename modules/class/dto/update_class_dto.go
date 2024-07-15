package dto

type ReqUpdateClass struct {
	Name string `json:"name" validate:"required"`
}

func (r *ReqUpdateClass) ToDBUpdateClass(authId string) ToDBUpdateClass {
	return ToDBUpdateClass{
		Name:        r.Name,
		UpdatedByID: authId,
	}
}

type ToDBUpdateClass struct {
	Name        string `json:"name"`
	UpdatedByID string
}
