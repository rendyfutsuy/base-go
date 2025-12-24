package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	mod "github.com/rendyfutsuy/base-go/modules/regency"
	"github.com/rendyfutsuy/base-go/modules/regency/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type regencyUsecase struct {
	repo mod.Repository
}

func NewRegencyUsecase(repo mod.Repository) mod.Usecase {
	return &regencyUsecase{repo: repo}
}

// Province Usecase
func (u *regencyUsecase) CreateProvince(ctx context.Context, reqBody *dto.ReqCreateProvince, userID string) (*models.Province, error) {
	exists, err := u.repo.ExistsProvinceByName(ctx, reqBody.Name, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.ProvinceNameAlreadyExists)
	}

	return u.repo.CreateProvince(ctx, reqBody.Name)
}

func (u *regencyUsecase) UpdateProvince(ctx context.Context, id string, reqBody *dto.ReqUpdateProvince, userID string) (*models.Province, error) {
	pid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	exists, err := u.repo.ExistsProvinceByName(ctx, reqBody.Name, pid)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.ProvinceNameAlreadyExists)
	}

	res, err := u.repo.UpdateProvince(ctx, pid, reqBody.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.ProvinceNotFound, id)
		}
		return nil, err
	}
	return res, nil
}

func (u *regencyUsecase) DeleteProvince(ctx context.Context, id string, userID string) error {
	pid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteProvince(ctx, pid)
}

func (u *regencyUsecase) GetProvinceByID(ctx context.Context, id string) (*models.Province, error) {
	pid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetProvinceByID(ctx, pid)
}

func (u *regencyUsecase) GetProvinceIndex(ctx context.Context, req request.PageRequest, filter dto.ReqProvinceIndexFilter) ([]models.Province, int, error) {
	return u.repo.GetProvinceIndex(ctx, req, filter)
}

func (u *regencyUsecase) GetAllProvince(ctx context.Context, filter dto.ReqProvinceIndexFilter) ([]models.Province, error) {
	return u.repo.GetAllProvince(ctx, filter)
}

func (u *regencyUsecase) ExportProvince(ctx context.Context, filter dto.ReqProvinceIndexFilter) ([]byte, error) {
	list, err := u.repo.GetAllProvince(ctx, filter)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Provinces"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "Name")
	f.SetCellValue(sheet, "B1", "Update Date")

	for i, p := range list {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), p.Name)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), p.UpdatedAt.Local().Format("2006/01/02"))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// City Usecase
func (u *regencyUsecase) CreateCity(ctx context.Context, reqBody *dto.ReqCreateCity, userID string) (*models.City, error) {
	// Check if province_id exists
	if reqBody.ProvinceID != uuid.Nil {
		provinceObject, err := u.repo.GetProvinceByID(ctx, reqBody.ProvinceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(constants.CityProvinceNotFound)
			}
			return nil, err
		}
		// Additional check: ensure provinceObject is valid
		if provinceObject == nil || provinceObject.ID == uuid.Nil {
			return nil, errors.New(constants.CityProvinceNotFound)
		}
	}

	exists, err := u.repo.ExistsCityByName(ctx, reqBody.ProvinceID, reqBody.Name, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.CityNameAlreadyExists)
	}

	return u.repo.CreateCity(ctx, reqBody.ProvinceID, reqBody.Name, reqBody.AreaCode)
}

func (u *regencyUsecase) UpdateCity(ctx context.Context, id string, reqBody *dto.ReqUpdateCity, userID string) (*models.City, error) {
	cid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	// Check if province_id exists
	if reqBody.ProvinceID != uuid.Nil {
		provinceObject, err := u.repo.GetProvinceByID(ctx, reqBody.ProvinceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(constants.CityProvinceNotFound)
			}
			return nil, err
		}
		// Additional check: ensure provinceObject is valid
		if provinceObject == nil || provinceObject.ID == uuid.Nil {
			return nil, errors.New(constants.CityProvinceNotFound)
		}
	}

	exists, err := u.repo.ExistsCityByName(ctx, reqBody.ProvinceID, reqBody.Name, cid)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.CityNameAlreadyExists)
	}

	res, err := u.repo.UpdateCity(ctx, cid, reqBody.ProvinceID, reqBody.Name, reqBody.AreaCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.CityNotFound, id)
		}
		return nil, err
	}
	return res, nil
}

func (u *regencyUsecase) DeleteCity(ctx context.Context, id string, userID string) error {
	cid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteCity(ctx, cid)
}

func (u *regencyUsecase) GetCityByID(ctx context.Context, id string) (*models.City, error) {
	cid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetCityByID(ctx, cid)
}

func (u *regencyUsecase) GetCityIndex(ctx context.Context, req request.PageRequest, filter dto.ReqCityIndexFilter) ([]models.City, int, error) {
	return u.repo.GetCityIndex(ctx, req, filter)
}

func (u *regencyUsecase) GetAllCity(ctx context.Context, filter dto.ReqCityIndexFilter) ([]models.City, error) {
	return u.repo.GetAllCity(ctx, filter)
}

func (u *regencyUsecase) ExportCity(ctx context.Context, filter dto.ReqCityIndexFilter) ([]byte, error) {
	list, err := u.repo.GetAllCity(ctx, filter)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Cities"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "Province ID")
	f.SetCellValue(sheet, "B1", "Area Code")
	f.SetCellValue(sheet, "C1", "Name")
	f.SetCellValue(sheet, "D1", "Update Date")

	for i, c := range list {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), c.ProvinceID.String())
		if c.AreaCode != nil {
			f.SetCellValue(sheet, "B"+strconv.Itoa(row), *c.AreaCode)
		} else {
			f.SetCellValue(sheet, "B"+strconv.Itoa(row), "")
		}
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), c.Name)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), c.UpdatedAt.Local().Format("2006/01/02"))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (u *regencyUsecase) GetCityAreaCodes(ctx context.Context, search string) ([]string, error) {
	return u.repo.GetCityAreaCodes(ctx, search)
}

// District Usecase
func (u *regencyUsecase) CreateDistrict(ctx context.Context, reqBody *dto.ReqCreateDistrict, userID string) (*models.District, error) {
	// Check if city_id exists
	if reqBody.CityID != uuid.Nil {
		cityObject, err := u.repo.GetCityByID(ctx, reqBody.CityID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(constants.DistrictCityNotFound)
			}
			return nil, err
		}
		// Additional check: ensure cityObject is valid
		if cityObject == nil || cityObject.ID == uuid.Nil {
			return nil, errors.New(constants.DistrictCityNotFound)
		}
	}

	exists, err := u.repo.ExistsDistrictByName(ctx, reqBody.CityID, reqBody.Name, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.DistrictNameAlreadyExists)
	}

	return u.repo.CreateDistrict(ctx, reqBody.CityID, reqBody.Name)
}

func (u *regencyUsecase) UpdateDistrict(ctx context.Context, id string, reqBody *dto.ReqUpdateDistrict, userID string) (*models.District, error) {
	did, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	// Check if city_id exists
	if reqBody.CityID != uuid.Nil {
		cityObject, err := u.repo.GetCityByID(ctx, reqBody.CityID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(constants.DistrictCityNotFound)
			}
			return nil, err
		}
		// Additional check: ensure cityObject is valid
		if cityObject == nil || cityObject.ID == uuid.Nil {
			return nil, errors.New(constants.DistrictCityNotFound)
		}
	}

	exists, err := u.repo.ExistsDistrictByName(ctx, reqBody.CityID, reqBody.Name, did)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.DistrictNameAlreadyExists)
	}

	res, err := u.repo.UpdateDistrict(ctx, did, reqBody.CityID, reqBody.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.DistrictNotFound, id)
		}
		return nil, err
	}
	return res, nil
}

func (u *regencyUsecase) DeleteDistrict(ctx context.Context, id string, userID string) error {
	did, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteDistrict(ctx, did)
}

func (u *regencyUsecase) GetDistrictByID(ctx context.Context, id string) (*models.District, error) {
	did, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetDistrictByID(ctx, did)
}

func (u *regencyUsecase) GetDistrictIndex(ctx context.Context, req request.PageRequest, filter dto.ReqDistrictIndexFilter) ([]models.District, int, error) {
	return u.repo.GetDistrictIndex(ctx, req, filter)
}

func (u *regencyUsecase) GetAllDistrict(ctx context.Context, filter dto.ReqDistrictIndexFilter) ([]models.District, error) {
	return u.repo.GetAllDistrict(ctx, filter)
}

func (u *regencyUsecase) ExportDistrict(ctx context.Context, filter dto.ReqDistrictIndexFilter) ([]byte, error) {
	list, err := u.repo.GetAllDistrict(ctx, filter)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Districts"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "City ID")
	f.SetCellValue(sheet, "B1", "Name")
	f.SetCellValue(sheet, "C1", "Update Date")

	for i, d := range list {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), d.CityID.String())
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), d.Name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), d.UpdatedAt.Local().Format("2006/01/02"))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Subdistrict Usecase
func (u *regencyUsecase) CreateSubdistrict(ctx context.Context, reqBody *dto.ReqCreateSubdistrict, userID string) (*models.Subdistrict, error) {
	// Check if district_id exists
	if reqBody.DistrictID != uuid.Nil {
		districtObject, err := u.repo.GetDistrictByID(ctx, reqBody.DistrictID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(constants.SubdistrictDistrictNotFound)
			}
			return nil, err
		}
		// Additional check: ensure districtObject is valid
		if districtObject == nil || districtObject.ID == uuid.Nil {
			return nil, errors.New(constants.SubdistrictDistrictNotFound)
		}
	}

	exists, err := u.repo.ExistsSubdistrictByName(ctx, reqBody.DistrictID, reqBody.Name, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.SubdistrictNameAlreadyExists)
	}

	return u.repo.CreateSubdistrict(ctx, reqBody.DistrictID, reqBody.Name)
}

func (u *regencyUsecase) UpdateSubdistrict(ctx context.Context, id string, reqBody *dto.ReqUpdateSubdistrict, userID string) (*models.Subdistrict, error) {
	sid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}

	// Check if district_id exists
	if reqBody.DistrictID != uuid.Nil {
		districtObject, err := u.repo.GetDistrictByID(ctx, reqBody.DistrictID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New(constants.SubdistrictDistrictNotFound)
			}
			return nil, err
		}
		// Additional check: ensure districtObject is valid
		if districtObject == nil || districtObject.ID == uuid.Nil {
			return nil, errors.New(constants.SubdistrictDistrictNotFound)
		}
	}

	exists, err := u.repo.ExistsSubdistrictByName(ctx, reqBody.DistrictID, reqBody.Name, sid)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.SubdistrictNameAlreadyExists)
	}

	res, err := u.repo.UpdateSubdistrict(ctx, sid, reqBody.DistrictID, reqBody.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.SubdistrictNotFound, id)
		}
		return nil, err
	}
	return res, nil
}

func (u *regencyUsecase) DeleteSubdistrict(ctx context.Context, id string, userID string) error {
	sid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteSubdistrict(ctx, sid)
}

func (u *regencyUsecase) GetSubdistrictByID(ctx context.Context, id string) (*models.Subdistrict, error) {
	sid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetSubdistrictByID(ctx, sid)
}

func (u *regencyUsecase) GetSubdistrictIndex(ctx context.Context, req request.PageRequest, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, int, error) {
	return u.repo.GetSubdistrictIndex(ctx, req, filter)
}

func (u *regencyUsecase) GetAllSubdistrict(ctx context.Context, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, error) {
	return u.repo.GetAllSubdistrict(ctx, filter)
}

func (u *regencyUsecase) ExportSubdistrict(ctx context.Context, filter dto.ReqSubdistrictIndexFilter) ([]byte, error) {
	list, err := u.repo.GetAllSubdistrict(ctx, filter)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Subdistricts"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "District ID")
	f.SetCellValue(sheet, "B1", "Name")
	f.SetCellValue(sheet, "C1", "Update Date")

	for i, s := range list {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), s.DistrictID.String())
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), s.Name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), s.UpdatedAt.Local().Format("2006/01/02"))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
