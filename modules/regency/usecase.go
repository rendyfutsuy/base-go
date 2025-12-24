package regency

import (
	"context"

	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/regency/dto"
)

type Usecase interface {
	// Province
	CreateProvince(ctx context.Context, req *dto.ReqCreateProvince, authId string) (*models.Province, error)
	UpdateProvince(ctx context.Context, id string, req *dto.ReqUpdateProvince, authId string) (*models.Province, error)
	DeleteProvince(ctx context.Context, id string, authId string) error
	GetProvinceByID(ctx context.Context, id string) (*models.Province, error)
	GetProvinceIndex(ctx context.Context, req request.PageRequest, filter dto.ReqProvinceIndexFilter) ([]models.Province, int, error)
	GetAllProvince(ctx context.Context, filter dto.ReqProvinceIndexFilter) ([]models.Province, error)
	ExportProvince(ctx context.Context, filter dto.ReqProvinceIndexFilter) ([]byte, error)

	// City
	CreateCity(ctx context.Context, req *dto.ReqCreateCity, authId string) (*models.City, error)
	UpdateCity(ctx context.Context, id string, req *dto.ReqUpdateCity, authId string) (*models.City, error)
	DeleteCity(ctx context.Context, id string, authId string) error
	GetCityByID(ctx context.Context, id string) (*models.City, error)
	GetCityIndex(ctx context.Context, req request.PageRequest, filter dto.ReqCityIndexFilter) ([]models.City, int, error)
	GetAllCity(ctx context.Context, filter dto.ReqCityIndexFilter) ([]models.City, error)
	ExportCity(ctx context.Context, filter dto.ReqCityIndexFilter) ([]byte, error)
	GetCityAreaCodes(ctx context.Context, search string) ([]string, error)

	// District
	CreateDistrict(ctx context.Context, req *dto.ReqCreateDistrict, authId string) (*models.District, error)
	UpdateDistrict(ctx context.Context, id string, req *dto.ReqUpdateDistrict, authId string) (*models.District, error)
	DeleteDistrict(ctx context.Context, id string, authId string) error
	GetDistrictByID(ctx context.Context, id string) (*models.District, error)
	GetDistrictIndex(ctx context.Context, req request.PageRequest, filter dto.ReqDistrictIndexFilter) ([]models.District, int, error)
	GetAllDistrict(ctx context.Context, filter dto.ReqDistrictIndexFilter) ([]models.District, error)
	ExportDistrict(ctx context.Context, filter dto.ReqDistrictIndexFilter) ([]byte, error)

	// Subdistrict
	CreateSubdistrict(ctx context.Context, req *dto.ReqCreateSubdistrict, authId string) (*models.Subdistrict, error)
	UpdateSubdistrict(ctx context.Context, id string, req *dto.ReqUpdateSubdistrict, authId string) (*models.Subdistrict, error)
	DeleteSubdistrict(ctx context.Context, id string, authId string) error
	GetSubdistrictByID(ctx context.Context, id string) (*models.Subdistrict, error)
	GetSubdistrictIndex(ctx context.Context, req request.PageRequest, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, int, error)
	GetAllSubdistrict(ctx context.Context, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, error)
	ExportSubdistrict(ctx context.Context, filter dto.ReqSubdistrictIndexFilter) ([]byte, error)
}
