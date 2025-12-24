package filters

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/modules/regency/dto"
	"gorm.io/gorm"
)

// applyCityFilters applies all filters from ReqCityIndexFilter to the query
func ApplyCityFilters(query *gorm.DB, filter dto.ReqCityIndexFilter) *gorm.DB {
	// Apply filters with multiple values support
	if filter.ProvinceID != uuid.Nil {
		query = query.Where("c.province_id = ?", filter.ProvinceID)
	}
	if len(filter.Names) > 0 {
		query = query.Where("c.name IN (?)", filter.Names)
	}
	return query
}

// applyDistrictFilters applies all filters from ReqDistrictIndexFilter to the query
func ApplyDistrictFilters(query *gorm.DB, filter dto.ReqDistrictIndexFilter) *gorm.DB {
	// Apply filters with multiple values support
	if filter.CityID != uuid.Nil {
		query = query.Where("d.city_id = ?", filter.CityID)
	}
	if len(filter.Names) > 0 {
		query = query.Where("d.name IN (?)", filter.Names)
	}
	return query
}

// applySubdistrictFilters applies all filters from ReqSubdistrictIndexFilter to the query
func ApplySubdistrictFilters(query *gorm.DB, filter dto.ReqSubdistrictIndexFilter) *gorm.DB {
	// Apply filters with multiple values support
	if filter.DistrictID != uuid.Nil {
		query = query.Where("s.district_id = ?", filter.DistrictID)
	}
	if len(filter.Names) > 0 {
		query = query.Where("s.name IN (?)", filter.Names)
	}
	return query
}
