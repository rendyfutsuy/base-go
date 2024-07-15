package shipyard

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/shipyard/dto"
)

// Usecase is an interface that defines the methods a shipyard use case must implement.
type Usecase interface {
	FindShipyardByUUIDOrCode(requestID interface{}) (result models.Shipyard, err error)       // FindShipyardByUUIDOrCode finds a shipyard by its UUID or code.
	FetchAllActiveShipyards() (result []*models.Shipyard, total int, lastPage int, err error) // FetchAllActiveShipyards fetches all active shipyards.

	StoreShipyard(data dto.ReqShipyard) (result models.Shipyard, err error)                  // StoreShipyard stores a new shipyard.
	UpdateShipyard(id interface{}, data dto.ReqShipyard) (result models.Shipyard, err error) // UpdateShipyard updates an existing shipyard.

	DeleteShipyard(data models.Shipyard) (err error) // DeleteShipyard deletes a shipyard.

	FetchShipyards(req request.PageRequest) (result []*models.Shipyard, total int, lastPage int, err error) // FetchShipyards fetches shipyards based on a PageRequest. It returns a slice of shipyards, the total number of shipyards, the last page number, and any error encountered.
}
