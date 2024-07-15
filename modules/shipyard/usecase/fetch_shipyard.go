package usecase

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
)

// FindShipyardByUUIDOrCode is a method of the shipyardUsecase struct.
// It finds a shipyard by its UUID or code.
func (uc *shipyardUsecase) FindShipyardByUUIDOrCode(requestID interface{}) (result models.Shipyard, err error) {
	// Call the FindShipyardByUUIDOrCode method of the shipyard repository.
	return uc.shipyardRepo.FindShipyardByUUIDOrCode(requestID)
}

// FetchAllActiveShipyards is a method of the shipyardUsecase struct.
// It fetches all active shipyards.
func (uc *shipyardUsecase) FetchAllActiveShipyards() (result []*models.Shipyard, total int, lastPage int, err error) {
	// Call the FetchShipyards method of the shipyard repository with a condition to only fetch active shipyards.
	return uc.shipyardRepo.FetchShipyards(&fetchRepoRequest{
		Condition: "deleted_at IS NULL",
	})
}

// FetchShipyards is a function in the shipyardUsecase struct.
// It gets a list of shipyards based on the page request.
func (uc *shipyardUsecase) FetchShipyards(req request.PageRequest) (result []*models.Shipyard, total int, lastPage int, err error) {
	request := fetchRepoRequesFromPageRequest(req)  // Turns the page request into a repository request
	return uc.shipyardRepo.FetchShipyards(&request) // Gets the shipyards from the repository using the request
}
