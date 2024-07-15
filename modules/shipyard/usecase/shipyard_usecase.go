package usecase

import (
	"time"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/shipyard"
)

// shipyardUsecase is a struct that implements the Usecase interface.
type shipyardUsecase struct {
	shipyardRepo   shipyard.Repository // shipyardRepo is the repository for shipyard data.
	contextTimeout time.Duration       // contextTimeout is the timeout duration for the context.
}

// NewShipyardUsecase creates a new shipyard use case with the provided repository and timeout.
func NewShipyardUsecase(r shipyard.Repository, timeout time.Duration) (result shipyard.Usecase) {
	result = &shipyardUsecase{
		shipyardRepo:   r,       // Set the repository for shipyard data.
		contextTimeout: timeout, // Set the timeout duration for the context.
	}
	return
}
