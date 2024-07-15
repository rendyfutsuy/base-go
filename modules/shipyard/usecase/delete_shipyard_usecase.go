package usecase

import "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"

func (uc *shipyardUsecase) DeleteShipyard(data models.Shipyard) (err error) {

	return uc.shipyardRepo.DeleteShipyard(&data)
}
