package regency

import (
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/regency/dto"
)

type Usecase interface {
	// Province
	CreateProvince(c echo.Context, req *dto.ReqCreateProvince, authId string) (*models.Province, error)
	UpdateProvince(c echo.Context, id string, req *dto.ReqUpdateProvince, authId string) (*models.Province, error)
	DeleteProvince(c echo.Context, id string, authId string) error
	GetProvinceByID(c echo.Context, id string) (*models.Province, error)
	GetProvinceIndex(c echo.Context, req request.PageRequest, filter dto.ReqProvinceIndexFilter) ([]models.Province, int, error)
	GetAllProvince(c echo.Context, filter dto.ReqProvinceIndexFilter) ([]models.Province, error)
	ExportProvince(c echo.Context, filter dto.ReqProvinceIndexFilter) ([]byte, error)

	// City
	CreateCity(c echo.Context, req *dto.ReqCreateCity, authId string) (*models.City, error)
	UpdateCity(c echo.Context, id string, req *dto.ReqUpdateCity, authId string) (*models.City, error)
	DeleteCity(c echo.Context, id string, authId string) error
	GetCityByID(c echo.Context, id string) (*models.City, error)
	GetCityIndex(c echo.Context, req request.PageRequest, filter dto.ReqCityIndexFilter) ([]models.City, int, error)
	GetAllCity(c echo.Context, filter dto.ReqCityIndexFilter) ([]models.City, error)
	ExportCity(c echo.Context, filter dto.ReqCityIndexFilter) ([]byte, error)

	// District
	CreateDistrict(c echo.Context, req *dto.ReqCreateDistrict, authId string) (*models.District, error)
	UpdateDistrict(c echo.Context, id string, req *dto.ReqUpdateDistrict, authId string) (*models.District, error)
	DeleteDistrict(c echo.Context, id string, authId string) error
	GetDistrictByID(c echo.Context, id string) (*models.District, error)
	GetDistrictIndex(c echo.Context, req request.PageRequest, filter dto.ReqDistrictIndexFilter) ([]models.District, int, error)
	GetAllDistrict(c echo.Context, filter dto.ReqDistrictIndexFilter) ([]models.District, error)
	ExportDistrict(c echo.Context, filter dto.ReqDistrictIndexFilter) ([]byte, error)

	// Subdistrict
	CreateSubdistrict(c echo.Context, req *dto.ReqCreateSubdistrict, authId string) (*models.Subdistrict, error)
	UpdateSubdistrict(c echo.Context, id string, req *dto.ReqUpdateSubdistrict, authId string) (*models.Subdistrict, error)
	DeleteSubdistrict(c echo.Context, id string, authId string) error
	GetSubdistrictByID(c echo.Context, id string) (*models.Subdistrict, error)
	GetSubdistrictIndex(c echo.Context, req request.PageRequest, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, int, error)
	GetAllSubdistrict(c echo.Context, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, error)
	ExportSubdistrict(c echo.Context, filter dto.ReqSubdistrictIndexFilter) ([]byte, error)
}
