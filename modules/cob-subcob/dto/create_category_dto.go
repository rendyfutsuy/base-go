package dto

import "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/dto"

// Define a struct for the MongoDB ObjectId
type MongoID struct {
	OID string `json:"$oid"`
}

// Define a struct for the MongoDB Date
type MongoDate struct {
	Date string `json:"$date"`
}

// Define the main struct for the JSON document
type CategoryJson struct {
	ID          MongoID `json:"_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Code        string  `json:"code"`
}

func (c *CategoryJson) ToDBCategory(createdByID string) dto.ToDBCreateCategory {
	return dto.ToDBCreateCategory{
		Name:        c.Name,
		Description: c.Description,
		Code:        c.Code,
		CreatedByID: createdByID,
	}
}
