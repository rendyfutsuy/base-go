package regency

import (
	"context"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/regency/dto"
)

// Combined Repository Interface
type Repository interface {
	// Province methods
	CreateProvince(ctx context.Context, name string) (*models.Province, error)
	UpdateProvince(ctx context.Context, id uuid.UUID, name string) (*models.Province, error)
	DeleteProvince(ctx context.Context, id uuid.UUID) error
	GetProvinceByID(ctx context.Context, id uuid.UUID) (*models.Province, error)
	GetProvinceIndex(ctx context.Context, req request.PageRequest, filter dto.ReqProvinceIndexFilter) ([]models.Province, int, error)
	GetAllProvince(ctx context.Context, filter dto.ReqProvinceIndexFilter) ([]models.Province, error)
	ExistsProvinceByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error)

	// City methods
	CreateCity(ctx context.Context, provinceID uuid.UUID, name string, areaCode *string) (*models.City, error)
	UpdateCity(ctx context.Context, id uuid.UUID, provinceID uuid.UUID, name string, areaCode *string) (*models.City, error)
	DeleteCity(ctx context.Context, id uuid.UUID) error
	GetCityByID(ctx context.Context, id uuid.UUID) (*models.City, error)
	GetCityIndex(ctx context.Context, req request.PageRequest, filter dto.ReqCityIndexFilter) ([]models.City, int, error)
	GetAllCity(ctx context.Context, filter dto.ReqCityIndexFilter) ([]models.City, error)
	ExistsCityByName(ctx context.Context, provinceID uuid.UUID, name string, excludeID uuid.UUID) (bool, error)
	GetCityAreaCodes(ctx context.Context, search string) ([]string, error)

	// District methods
	CreateDistrict(ctx context.Context, cityID uuid.UUID, name string) (*models.District, error)
	UpdateDistrict(ctx context.Context, id uuid.UUID, cityID uuid.UUID, name string) (*models.District, error)
	DeleteDistrict(ctx context.Context, id uuid.UUID) error
	GetDistrictByID(ctx context.Context, id uuid.UUID) (*models.District, error)
	GetDistrictIndex(ctx context.Context, req request.PageRequest, filter dto.ReqDistrictIndexFilter) ([]models.District, int, error)
	GetAllDistrict(ctx context.Context, filter dto.ReqDistrictIndexFilter) ([]models.District, error)
	ExistsDistrictByName(ctx context.Context, cityID uuid.UUID, name string, excludeID uuid.UUID) (bool, error)

	// Subdistrict methods
	CreateSubdistrict(ctx context.Context, districtID uuid.UUID, name string) (*models.Subdistrict, error)
	UpdateSubdistrict(ctx context.Context, id uuid.UUID, districtID uuid.UUID, name string) (*models.Subdistrict, error)
	DeleteSubdistrict(ctx context.Context, id uuid.UUID) error
	GetSubdistrictByID(ctx context.Context, id uuid.UUID) (*models.Subdistrict, error)
	GetSubdistrictIndex(ctx context.Context, req request.PageRequest, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, int, error)
	GetAllSubdistrict(ctx context.Context, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, error)
	ExistsSubdistrictByName(ctx context.Context, districtID uuid.UUID, name string, excludeID uuid.UUID) (bool, error)
}
