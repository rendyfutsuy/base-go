package repository

import (
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/modules/regency/dto"
	"github.com/rendyfutsuy/base-go/modules/regency/repository/filters"
	"gorm.io/gorm"
)

// ApplyFilters applies filters to the query
// Implements NeedFilterPredefine interface
// This method routes to the appropriate filter function based on the filter type
func (r *regencyRepository) ApplyFilters(query *gorm.DB, filter interface{}) *gorm.DB {
	// Try to cast to each filter type and apply appropriate filters
	if cityFilter, ok := filter.(dto.ReqCityIndexFilter); ok {
		return filters.ApplyCityFilters(query, cityFilter)
	}
	if districtFilter, ok := filter.(dto.ReqDistrictIndexFilter); ok {
		return filters.ApplyDistrictFilters(query, districtFilter)
	}
	if subdistrictFilter, ok := filter.(dto.ReqSubdistrictIndexFilter); ok {
		return filters.ApplySubdistrictFilters(query, subdistrictFilter)
	}
	// If filter type doesn't match any known type, return query unchanged
	return query
}

// Compile-time check to ensure regencyRepository implements NeedFilterPredefine interface
var _ request.NeedFilterPredefine = (*regencyRepository)(nil)
