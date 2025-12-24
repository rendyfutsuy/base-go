package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/regency"
	"github.com/rendyfutsuy/base-go/modules/regency/dto"
	"github.com/rendyfutsuy/base-go/modules/regency/repository/searches"
	"gorm.io/gorm"
)

type regencyRepository struct {
	DB *gorm.DB
}

func NewRegencyRepository(db *gorm.DB) regency.Repository {
	return &regencyRepository{DB: db}
}

// Province Repository Implementation
func (r *regencyRepository) CreateProvince(ctx context.Context, name string) (*models.Province, error) {
	now := time.Now().UTC()
	p := &models.Province{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := r.DB.WithContext(ctx).Create(p).Error; err != nil {
		return nil, err
	}
	if p.ID == uuid.Nil {
		return nil, errors.New(constants.ProvinceCreateFailedIDNotSet)
	}
	return p, nil
}

func (r *regencyRepository) UpdateProvince(ctx context.Context, id uuid.UUID, name string) (*models.Province, error) {
	updates := map[string]interface{}{
		"name":       name,
		"updated_at": time.Now().UTC(),
	}
	p := &models.Province{}
	err := r.DB.WithContext(ctx).Model(&models.Province{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).
		First(p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return p, nil
}

func (r *regencyRepository) DeleteProvince(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Delete(&models.Province{}).Error
}

func (r *regencyRepository) GetProvinceByID(ctx context.Context, id uuid.UUID) (*models.Province, error) {
	p := &models.Province{}
	if err := r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *regencyRepository) ExistsProvinceByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	q := r.DB.WithContext(ctx).Unscoped().Model(&models.Province{}).Where("name = ?", name)
	if excludeID != uuid.Nil {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *regencyRepository) GetProvinceIndex(ctx context.Context, req request.PageRequest, filter dto.ReqProvinceIndexFilter) ([]models.Province, int, error) {
	var provinces []models.Province
	query := r.DB.WithContext(ctx).Table("provinces p").Select("p.id, p.name, p.created_at, p.updated_at").
		Where("p.deleted_at IS NULL")

	searchQuery := req.Search
	query = request.ApplySearchConditionFromInterface(query, searchQuery, searches.NewProvinceSearchHelper())

	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:      "p.created_at",
		DefaultSortOrder:   "DESC",
		AllowedColumns:     []string{"id", "name", "created_at", "updated_at"},
		ColumnPrefix:       "p.",
		MaxPerPage:         100,
		NaturalSortColumns: []string{"p.name"}, // Enable natural sorting for p.name
	}, &provinces)
	if err != nil {
		return nil, 0, err
	}
	return provinces, total, nil
}

func (r *regencyRepository) GetAllProvince(ctx context.Context, filter dto.ReqProvinceIndexFilter) ([]models.Province, error) {
	var provinces []models.Province
	query := r.DB.WithContext(ctx).Table("provinces p").Select("p.id, p.name, p.created_at, p.updated_at").
		Where("p.deleted_at IS NULL")

	query = request.ApplySearchConditionFromInterface(query, filter.Search, searches.NewProvinceSearchHelper())

	if err := query.Order("p.created_at DESC").Find(&provinces).Error; err != nil {
		return nil, err
	}
	return provinces, nil
}

// City Repository Implementation
func (r *regencyRepository) CreateCity(ctx context.Context, provinceID uuid.UUID, name string, areaCode *string) (*models.City, error) {
	now := time.Now().UTC()
	c := &models.City{
		ProvinceID: provinceID,
		Name:       name,
		AreaCode:   areaCode,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := r.DB.WithContext(ctx).Create(c).Error; err != nil {
		return nil, err
	}
	if c.ID == uuid.Nil {
		return nil, errors.New(constants.CityCreateFailedIDNotSet)
	}
	return c, nil
}

func (r *regencyRepository) UpdateCity(ctx context.Context, id uuid.UUID, provinceID uuid.UUID, name string, areaCode *string) (*models.City, error) {
	updates := map[string]interface{}{
		"province_id": provinceID,
		"name":        name,
		"updated_at":  time.Now().UTC(),
	}
	if areaCode != nil {
		updates["area_code"] = areaCode
	}
	c := &models.City{}
	err := r.DB.WithContext(ctx).Model(&models.City{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).
		First(c).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return c, nil
}

func (r *regencyRepository) DeleteCity(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Delete(&models.City{}).Error
}

func (r *regencyRepository) GetCityByID(ctx context.Context, id uuid.UUID) (*models.City, error) {
	c := &models.City{}
	if err := r.DB.WithContext(ctx).Preload("Province").Where("id = ? AND deleted_at IS NULL", id).First(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *regencyRepository) ExistsCityByName(ctx context.Context, provinceID uuid.UUID, name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	q := r.DB.WithContext(ctx).Unscoped().Model(&models.City{}).Where("province_id = ? AND name = ?", provinceID, name)
	if excludeID != uuid.Nil {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *regencyRepository) GetCityIndex(ctx context.Context, req request.PageRequest, filter dto.ReqCityIndexFilter) ([]models.City, int, error) {
	var cities []models.City
	query := r.DB.WithContext(ctx).Table("cities c").Select("c.id, c.province_id, c.name, c.area_code, c.created_at, c.updated_at").
		Where("c.deleted_at IS NULL")

	searchQuery := req.Search
	query = request.ApplySearchConditionFromInterface(query, searchQuery, searches.NewCitySearchHelper())

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:      "c.created_at",
		DefaultSortOrder:   "DESC",
		AllowedColumns:     []string{"id", "province_id", "name", "area_code", "created_at", "updated_at"},
		ColumnPrefix:       "c.",
		MaxPerPage:         100,
		NaturalSortColumns: []string{"c.name"}, // Enable natural sorting for c.name
	}, &cities)
	if err != nil {
		return nil, 0, err
	}
	return cities, total, nil
}

func (r *regencyRepository) GetAllCity(ctx context.Context, filter dto.ReqCityIndexFilter) ([]models.City, error) {
	var cities []models.City
	query := r.DB.WithContext(ctx).Table("cities c").Select("c.id, c.province_id, c.name, c.area_code, c.created_at, c.updated_at").
		Where("c.deleted_at IS NULL")

	query = request.ApplySearchConditionFromInterface(query, filter.Search, searches.NewCitySearchHelper())

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	if err := query.Order("c.created_at DESC").Find(&cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}

func (r *regencyRepository) GetCityAreaCodes(ctx context.Context, search string) ([]string, error) {
	var areaCodes []string

	query := r.DB.WithContext(ctx).Table("cities c").
		Where("c.deleted_at IS NULL").
		Where("c.area_code IS NOT NULL AND c.area_code <> ''")

	if search != "" {
		query = query.Where("LOWER(c.area_code) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

	if err := query.Distinct("c.area_code").
		Order("c.area_code ASC").
		Pluck("c.area_code", &areaCodes).Error; err != nil {
		return nil, err
	}

	return areaCodes, nil
}

// District Repository Implementation
func (r *regencyRepository) CreateDistrict(ctx context.Context, cityID uuid.UUID, name string) (*models.District, error) {
	now := time.Now().UTC()
	d := &models.District{
		CityID:    cityID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := r.DB.WithContext(ctx).Create(d).Error; err != nil {
		return nil, err
	}
	if d.ID == uuid.Nil {
		return nil, errors.New(constants.DistrictCreateFailedIDNotSet)
	}
	return d, nil
}

func (r *regencyRepository) UpdateDistrict(ctx context.Context, id uuid.UUID, cityID uuid.UUID, name string) (*models.District, error) {
	updates := map[string]interface{}{
		"city_id":    cityID,
		"name":       name,
		"updated_at": time.Now().UTC(),
	}
	d := &models.District{}
	err := r.DB.WithContext(ctx).Model(&models.District{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).
		First(d).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return d, nil
}

func (r *regencyRepository) DeleteDistrict(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Delete(&models.District{}).Error
}

func (r *regencyRepository) GetDistrictByID(ctx context.Context, id uuid.UUID) (*models.District, error) {
	d := &models.District{}
	if err := r.DB.WithContext(ctx).Preload("City").Preload("City.Province").Where("id = ? AND deleted_at IS NULL", id).First(d).Error; err != nil {
		return nil, err
	}
	return d, nil
}

func (r *regencyRepository) ExistsDistrictByName(ctx context.Context, cityID uuid.UUID, name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	q := r.DB.WithContext(ctx).Unscoped().Model(&models.District{}).Where("city_id = ? AND name = ?", cityID, name)
	if excludeID != uuid.Nil {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *regencyRepository) GetDistrictIndex(ctx context.Context, req request.PageRequest, filter dto.ReqDistrictIndexFilter) ([]models.District, int, error) {
	var districts []models.District
	query := r.DB.WithContext(ctx).Table("districts d").Select("d.id, d.city_id, d.name, d.created_at, d.updated_at").
		Where("d.deleted_at IS NULL")

	searchQuery := req.Search
	query = request.ApplySearchConditionFromInterface(query, searchQuery, searches.NewDistrictSearchHelper())

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:      "d.created_at",
		DefaultSortOrder:   "DESC",
		AllowedColumns:     []string{"id", "city_id", "name", "created_at", "updated_at"},
		ColumnPrefix:       "d.",
		MaxPerPage:         100,
		NaturalSortColumns: []string{"d.name"}, // Enable natural sorting for d.name
	}, &districts)
	if err != nil {
		return nil, 0, err
	}
	return districts, total, nil
}

func (r *regencyRepository) GetAllDistrict(ctx context.Context, filter dto.ReqDistrictIndexFilter) ([]models.District, error) {
	var districts []models.District
	query := r.DB.WithContext(ctx).Table("districts d").Select("d.id, d.city_id, d.name, d.created_at, d.updated_at").
		Where("d.deleted_at IS NULL")

	query = request.ApplySearchConditionFromInterface(query, filter.Search, searches.NewDistrictSearchHelper())

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	if err := query.Order("d.created_at DESC").Find(&districts).Error; err != nil {
		return nil, err
	}
	return districts, nil
}

// Subdistrict Repository Implementation
func (r *regencyRepository) CreateSubdistrict(ctx context.Context, districtID uuid.UUID, name string) (*models.Subdistrict, error) {
	now := time.Now().UTC()
	s := &models.Subdistrict{
		DistrictID: districtID,
		Name:       name,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := r.DB.WithContext(ctx).Create(s).Error; err != nil {
		return nil, err
	}
	if s.ID == uuid.Nil {
		return nil, errors.New(constants.SubdistrictCreateFailedIDNotSet)
	}
	return s, nil
}

func (r *regencyRepository) UpdateSubdistrict(ctx context.Context, id uuid.UUID, districtID uuid.UUID, name string) (*models.Subdistrict, error) {
	updates := map[string]interface{}{
		"district_id": districtID,
		"name":        name,
		"updated_at":  time.Now().UTC(),
	}
	s := &models.Subdistrict{}
	err := r.DB.WithContext(ctx).Model(&models.Subdistrict{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).
		First(s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return s, nil
}

func (r *regencyRepository) DeleteSubdistrict(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).Delete(&models.Subdistrict{}).Error
}

func (r *regencyRepository) GetSubdistrictByID(ctx context.Context, id uuid.UUID) (*models.Subdistrict, error) {
	s := &models.Subdistrict{}
	if err := r.DB.WithContext(ctx).Preload("District").Preload("District.City").Preload("District.City.Province").Where("id = ? AND deleted_at IS NULL", id).First(s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (r *regencyRepository) ExistsSubdistrictByName(ctx context.Context, districtID uuid.UUID, name string, excludeID uuid.UUID) (bool, error) {
	var count int64
	q := r.DB.WithContext(ctx).Unscoped().Model(&models.Subdistrict{}).Where("district_id = ? AND name = ?", districtID, name)
	if excludeID != uuid.Nil {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *regencyRepository) GetSubdistrictIndex(ctx context.Context, req request.PageRequest, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, int, error) {
	var subdistricts []models.Subdistrict
	query := r.DB.WithContext(ctx).Table("subdistricts s").Select("s.id, s.district_id, s.name, s.created_at, s.updated_at").
		Where("s.deleted_at IS NULL")

	searchQuery := req.Search
	query = request.ApplySearchConditionFromInterface(query, searchQuery, searches.NewSubdistrictSearchHelper())

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:      "s.created_at",
		DefaultSortOrder:   "DESC",
		AllowedColumns:     []string{"id", "district_id", "name", "created_at", "updated_at"},
		ColumnPrefix:       "s.",
		MaxPerPage:         100,
		NaturalSortColumns: []string{"s.name"}, // Enable natural sorting for s.name
	}, &subdistricts)
	if err != nil {
		return nil, 0, err
	}
	return subdistricts, total, nil
}

func (r *regencyRepository) GetAllSubdistrict(ctx context.Context, filter dto.ReqSubdistrictIndexFilter) ([]models.Subdistrict, error) {
	var subdistricts []models.Subdistrict
	query := r.DB.WithContext(ctx).Table("subdistricts s").Select("s.id, s.district_id, s.name, s.created_at, s.updated_at").
		Where("s.deleted_at IS NULL")

	query = request.ApplySearchConditionFromInterface(query, filter.Search, searches.NewSubdistrictSearchHelper())

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	if err := query.Order("s.created_at DESC").Find(&subdistricts).Error; err != nil {
		return nil, err
	}
	return subdistricts, nil
}
