package usecase

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/cob-subcob/dto"
	"github.com/google/uuid"
)

func (uc *cobsubcobUsecase) InsertCategoryCobSubcob(cat []dto.CategoryJson, cob []dto.CobJson, createdByID string) (err error) {
	catCobMap := make(map[string]map[string][]dto.CobJson)

	// fill category
	for _, c := range cat {
		if _, ok := catCobMap[c.ID.OID]; !ok {
			catCobMap[c.ID.OID] = make(map[string][]dto.CobJson)
		}
	}

	// fill cob
	for _, c := range cob {
		if c.Parent != nil {
			continue
		}

		if _, ok := catCobMap[c.Category.ID.OID]; !ok {
			continue
		}

		catCobMap[c.Category.ID.OID][c.ID.OID] = []dto.CobJson{}
	}

	// fill subcob
	for _, c := range cob {
		if c.Parent == nil {
			continue
		}

		if _, ok := catCobMap[c.Category.ID.OID]; !ok {
			continue
		}

		if _, ok := catCobMap[c.Category.ID.OID][c.Parent.ID.OID]; !ok {
			continue
		}

		catCobMap[c.Category.ID.OID][c.Parent.ID.OID] = append(catCobMap[c.Category.ID.OID][c.Parent.ID.OID], c)
	}

	// query
	trx, err := uc.cobsubcobRepo.StartTransaction()
	if err != nil {
		return err
	}

	// insert category
	catMap := make(map[string]dto.CategoryJson)
	for _, c := range cat {
		catMap[c.ID.OID] = c
	}

	for _, c := range catMap {
		cToDB := c.ToDBCategory(createdByID)

		cRes, errC := uc.categoryRepo.CreateCategory(trx, cToDB)
		if errC != nil {
			err = errC
			break
		}

		// update key
		for k, v := range catCobMap {
			if k == c.ID.OID {
				delete(catCobMap, k)
				catCobMap[cRes.ID.String()] = v
				break
			}
		}
	}

	if err != nil {
		if err := trx.Rollback(); err != nil {
			return err
		}
		return err
	}

	// insert cob
	cobMap := make(map[string]dto.CobJson)
	for _, c := range cob {
		cobMap[c.ID.OID] = c
	}

	for catID, cMap := range catCobMap {

		oldIdArr := []string{}
		newIdArr := []string{}

		for cobId := range cMap {
			oldIdArr = append(oldIdArr, cobId)
			cob := cobMap[cobId]

			cobToDB, errC := cob.JSONToDBCreateCob(createdByID, catID)
			if errC != nil {
				err = errC
				break
			}

			cobRes, errC := uc.cobsubcobRepo.CreateCob(trx, *cobToDB)
			if errC != nil {
				err = errC
				break
			}
			newIdArr = append(newIdArr, cobRes.ID.String())

		}

		// update key
		for i := range oldIdArr {
			catCobMap[catID][newIdArr[i]] = catCobMap[catID][oldIdArr[i]]
			delete(catCobMap[catID], oldIdArr[i])
		}
	}

	if err != nil {
		if err := trx.Rollback(); err != nil {
			return err
		}
		return err
	}

	// insert subcob
	for catID, cMap := range catCobMap {
		for cobID, subcobs := range cMap {
			for _, subcob := range subcobs {
				subcobToDB, errC := subcob.JSONToDBCreateSubcob(createdByID, catID, cobID)
				if errC != nil {
					err = errC
					break
				}

				_, errC = uc.cobsubcobRepo.CreateSubcob(trx, *subcobToDB)
				if errC != nil {
					err = errC
					break
				}
			}
		}
	}

	if err != nil {
		if err := trx.Rollback(); err != nil {
			return err
		}
		return err
	}

	err = trx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (uc *cobsubcobUsecase) GetCobByID(id string) (cobRes *models.Cob, err error) {

	uId := uuid.MustParse(id)

	return uc.cobsubcobRepo.GetCobByID(uId)
}

func (uc *cobsubcobUsecase) GetSubcobByID(id string) (subcobRes *models.Subcob, err error) {

	uId := uuid.MustParse(id)

	return uc.cobsubcobRepo.GetSubcobByID(uId)
}

func (uc *cobsubcobUsecase) GetIndexCob(req *request.PageRequest) (cobs []models.Cob, total int, err error) {
	return uc.cobsubcobRepo.GetIndexCob(*req)
}

func (uc *cobsubcobUsecase) GetIndexSubcob(req *request.PageRequest) (subcobs []models.Subcob, total int, err error) {
	return uc.cobsubcobRepo.GetIndexSubcob(*req)
}

func (uc *cobsubcobUsecase) GetAllCob() (cobs []models.Cob, err error) {
	return uc.cobsubcobRepo.GetAllCob()
}

func (uc *cobsubcobUsecase) GetAllSubcob() (subcobs []models.Subcob, err error) {
	return uc.cobsubcobRepo.GetAllSubcob()
}
