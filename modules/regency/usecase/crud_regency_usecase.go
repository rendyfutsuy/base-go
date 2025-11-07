package usecase

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
func (u *regencyUsecase) CreateProvince(c echo.Context, reqBody *dto.ReqCreateProvince, authId string) (*models.Province, error) {
	ctx := c.Request().Context()

	exists, err := u.repo.ExistsProvinceByName(ctx, reqBody.Name, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(constants.ProvinceNameAlreadyExists)
	}

	return u.repo.CreateProvince(ctx, reqBody.Name)
}

func (u *regencyUsecase) UpdateProvince(c echo.Context, id string, reqBody *dto.ReqUpdateProvince, authId string) (*models.Province, error) {
	ctx := c.Request().Context()
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

	return u.repo.UpdateProvince(ctx, pid, reqBody.Name)
}

func (u *regencyUsecase) DeleteProvince(c echo.Context, id string, authId string) error {
	ctx := c.Request().Context()
	pid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteProvince(ctx, pid)
}

func (u *regencyUsecase) GetProvinceByID(c echo.Context, id string) (*models.Province, error) {
	ctx := c.Request().Context()
	pid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetProvinceByID(ctx, pid)
}

func (u *regencyUsecase) GetProvinceIndex(c echo.Context, req request.PageRequest, filter dto.ReqProvinceIndexFilter) ([]models.Province, int, error) {
	ctx := c.Request().Context()
	return u.repo.GetProvinceIndex(ctx, req, filter)
}

func (u *regencyUsecase) GetAllProvince(c echo.Context, filter dto.ReqProvinceIndexFilter) ([]models.Province, error) {
	ctx := c.Request().Context()
	return u.repo.GetAllProvince(ctx, filter)
}

func (u *regencyUsecase) ExportProvince(c echo.Context, filter dto.ReqProvinceIndexFilter) ([]byte, error) {
	ctx := c.Request().Context()
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
func (u *regencyUsecase) CreateCity(c echo.Context, reqBody *dto.ReqCreateCity, authId string) (*models.City, error) {
	ctx := c.Request().Context()

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

	return u.repo.CreateCity(ctx, reqBody.ProvinceID, reqBody.Name)
}

func (u *regencyUsecase) UpdateCity(c echo.Context, id string, reqBody *dto.ReqUpdateCity, authId string) (*models.City, error) {
	ctx := c.Request().Context()
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

	return u.repo.UpdateCity(ctx, cid, reqBody.ProvinceID, reqBody.Name)
}

func (u *regencyUsecase) DeleteCity(c echo.Context, id string, authId string) error {
	ctx := c.Request().Context()
	cid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteCity(ctx, cid)
}

func (u *regencyUsecase) GetCityByID(c echo.Context, id string) (*models.City, error) {
	ctx := c.Request().Context()
	cid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetCityByID(ctx, cid)
}

func (u *regencyUsecase) GetCityIndex(c echo.Context, req request.PageRequest, filter dto.ReqCityIndexFilter) ([]models.City, int, error) {
	ctx := c.Request().Context()
	return u.repo.GetCityIndex(ctx, req, filter)
}

func (u *regencyUsecase) GetAllCity(c echo.Context, filter dto.ReqCityIndexFilter) ([]models.City, error) {
	ctx := c.Request().Context()
	return u.repo.GetAllCity(ctx, filter)
}

func (u *regencyUsecase) ExportCity(c echo.Context, filter dto.ReqCityIndexFilter) ([]byte, error) {
	ctx := c.Request().Context()
	list, err := u.repo.GetAllCity(ctx, filter)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Cities"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "Province ID")
	f.SetCellValue(sheet, "B1", "Name")
	f.SetCellValue(sheet, "C1", "Update Date")

	for i, c := range list {
		row := i + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), c.ProvinceID.String())
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), c.Name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), c.UpdatedAt.Local().Format("2006/01/02"))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// District Usecase
func (u *regencyUsecase) CreateDistrict(c echo.Context, reqBody *dto.ReqCreateDistrict, authId string) (*models.District, error) {
	ctx := c.Request().Context()

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

func (u *regencyUsecase) UpdateDistrict(c echo.Context, id string, reqBody *dto.ReqUpdateDistrict, authId string) (*models.District, error) {
	ctx := c.Request().Context()
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

	return u.repo.UpdateDistrict(ctx, did, reqBody.CityID, reqBody.Name)
}

func (u *regencyUsecase) DeleteDistrict(c echo.Context, id string, authId string) error {
	ctx := c.Request().Context()
	did, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteDistrict(ctx, did)
}

func (u *regencyUsecase) GetDistrictByID(c echo.Context, id string) (*models.District, error) {
	ctx := c.Request().Context()
	did, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetDistrictByID(ctx, did)
}

func (u *regencyUsecase) GetDistrictIndex(c echo.Context, req request.PageRequest, filter dto.ReqDistrictIndexFilter) ([]models.District, int, error) {
	ctx := c.Request().Context()
	return u.repo.GetDistrictIndex(ctx, req, filter)
}

func (u *regencyUsecase) GetAllDistrict(c echo.Context, filter dto.ReqDistrictIndexFilter) ([]models.District, error) {
	ctx := c.Request().Context()
	return u.repo.GetAllDistrict(ctx, filter)
}

func (u *regencyUsecase) ExportDistrict(c echo.Context, filter dto.ReqDistrictIndexFilter) ([]byte, error) {
	ctx := c.Request().Context()
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
func (u *regencyUsecase) CreateSubdistrict(c echo.Context, reqBody *dto.ReqCreateSubdistrict, authId string) (*models.Subdistrict, error) {
	ctx := c.Request().Context()

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

func (u *regencyUsecase) UpdateSubdistrict(c echo.Context, id string, reqBody *dto.ReqUpdateSubdistrict, authId string) (*models.Subdistrict, error) {
	ctx := c.Request().Context()
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

	return u.repo.UpdateSubdistrict(ctx, sid, reqBody.DistrictID, reqBody.Name)
}

func (u *regencyUsecase) DeleteSubdistrict(c echo.Context, id string, authId string) error {
	ctx := c.Request().Context()
	sid, err := utils.StringToUUID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteSubdistrict(ctx, sid)
}

func (u *regencyUsecase) GetSubdistrictByID(c echo.Context, id string) (*models.Subdistrict, error) {
	ctx := c.Request().Context()
	sid, err := utils.StringToUUID(id)
	if err != nil {
		return nil, err
	}
	return u.repo.GetSubdistrictByID(ctx, sid)
}

func (u *regencyUsecase) GetSubdistrictIndex(c echo.Context, req request.PageRequest, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, int, error) {
	ctx := c.Request().Context()
	return u.repo.GetSubdistrictIndex(ctx, req, filter)
}

func (u *regencyUsecase) GetAllSubdistrict(c echo.Context, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, error) {
	ctx := c.Request().Context()
	return u.repo.GetAllSubdistrict(ctx, filter)
}

func (u *regencyUsecase) ExportSubdistrict(c echo.Context, filter dto.ReqSubdistrictIndexFilter) ([]byte, error) {
	ctx := c.Request().Context()
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
