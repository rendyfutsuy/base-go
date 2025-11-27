package repository

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/modules/regency/dto"
	"gorm.io/gorm"
)

// applyCityFilters applies all filters from ReqCityIndexFilter to the query
func applyCityFilters(query *gorm.DB, filter dto.ReqCityIndexFilter) *gorm.DB {
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
func applyDistrictFilters(query *gorm.DB, filter dto.ReqDistrictIndexFilter) *gorm.DB {
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
func applySubdistrictFilters(query *gorm.DB, filter dto.ReqSubdistrictIndexFilter) *gorm.DB {
	// Apply filters with multiple values support
	if filter.DistrictID != uuid.Nil {
		query = query.Where("s.district_id = ?", filter.DistrictID)
	}
	if len(filter.Names) > 0 {
		query = query.Where("s.name IN (?)", filter.Names)
	}
	return query
}

// ApplyFilters applies filters to the query
// Implements NeedFilterPredefine interface
// This method routes to the appropriate filter function based on the filter type
func (r *regencyRepository) ApplyFilters(query *gorm.DB, filter interface{}) *gorm.DB {
	// Try to cast to each filter type and apply appropriate filters
	if cityFilter, ok := filter.(dto.ReqCityIndexFilter); ok {
		return applyCityFilters(query, cityFilter)
	}
	if districtFilter, ok := filter.(dto.ReqDistrictIndexFilter); ok {
		return applyDistrictFilters(query, districtFilter)
	}
	if subdistrictFilter, ok := filter.(dto.ReqSubdistrictIndexFilter); ok {
		return applySubdistrictFilters(query, subdistrictFilter)
	}
	// If filter type doesn't match any known type, return query unchanged
	return query
}

// Compile-time check to ensure regencyRepository implements NeedFilterPredefine interface
var _ request.NeedFilterPredefine = (*regencyRepository)(nil)
