package usecase

import (
	"fmt"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	_shipyardDTO "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/shipyard/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	// errorUUIDNotRecognized is an error that is returned when a requested parameter is not recognized.
	errorUUIDNotRecognized = fmt.Errorf("Requested param is not recognized")
)

// StoreShipyard is a method of the shipyardUsecase struct.
// It stores a new shipyard.
func (uc *shipyardUsecase) StoreShipyard(data _shipyardDTO.ReqShipyard) (result models.Shipyard, err error) {
	// Create a new Shipyard object from the provided data.
	shipyard := models.Shipyard{
		Name: data.Name,
		Yard: data.Yard,
	}
	// Call the CreateShipyard method of the shipyard repository to store the new shipyard.
	err = uc.shipyardRepo.CreateShipyard(&shipyard)
	if err != nil {
		// If an error occurred, log it and return.
		zap.S().Error(err)
		return
	}
	// Set the result to the new shipyard.
	result = shipyard
	// Return the result and any error that occurred.
	return
}

// UpdateShipyard is a method of the shipyardUsecase struct.
// It updates an existing shipyard.
func (uc *shipyardUsecase) UpdateShipyard(id interface{}, data _shipyardDTO.ReqShipyard) (result models.Shipyard, err error) {
	// Convert the id to a UUID.
	uuid, err := StringToUUID(id)
	if err != nil {
		// If an error occurred, log it and return.
		zap.S().Error(err)
		return
	}
	// Create a new Shipyard object from the provided data and UUID.
	shipyard := models.Shipyard{
		ID:   uuid,
		Name: data.Name,
		Yard: data.Yard,
	}

	// Call the UpdateShipyard method of the shipyard repository to update the shipyard.
	err = uc.shipyardRepo.UpdateShipyard(&shipyard)
	if err != nil {
		// If an error occurred, log it and return.
		zap.S().Error(err)
		return
	}

	// Set the result to the updated shipyard.
	result = shipyard

	// Return the result and any error that occurred.
	return
}

// StringToUUID converts a string to a UUID.
func StringToUUID(request interface{}) (result uuid.UUID, err error) {
	// Try to assert the request as a string.
	stringUUID, ok := request.(string)
	if !ok {
		// If it's not a string, try to assert it as a UUID.
		uuid, ok := request.(uuid.UUID)
		if !ok {
			// If it's neither, return an error.
			err = errorUUIDNotRecognized
			return
		}
		// If it's a UUID, return it.
		result = uuid
	} else {
		// If it's a string, try to parse it as a UUID.
		uuid, parseErr := uuid.Parse(stringUUID)
		if parseErr != nil {
			// If the parsing fails, return an error.
			err = fmt.Errorf("requested param is string")
			return
		}
		// If the parsing succeeds, return the parsed UUID.
		result = uuid
	}
	return
}
