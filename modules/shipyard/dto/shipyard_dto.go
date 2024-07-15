package dto

// ReqShipyard is a struct that represents a request to create or update a shipyard.
type ReqShipyard struct {
	Name string `json:"name" validate:"required"` // Name is the name of the shipyard. It is required.
	Yard string `json:"yard" validate:"required"` // Yard is the yard of the shipyard. It is required.
}
