package repository

import (
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/modules/parameter/dto"
	"gorm.io/gorm"
)

// applyParameterFilters applies all filters from ReqParameterIndexFilter to the query
func applyParameterFilters(query *gorm.DB, filter dto.ReqParameterIndexFilter) *gorm.DB {
	// Apply filters with multiple values support
	if len(filter.Types) > 0 {
		query = query.Where("p.type IN (?)", filter.Types)
	}
	if len(filter.Names) > 0 {
		query = query.Where("p.name IN (?)", filter.Names)
	}
	if len(filter.IDs) > 0 {
		query = query.Where("p.id IN (?)", filter.IDs)
	}
	return query
}

// ApplyFilters applies filters to the query
// Implements NeedFilterPredefine interface
func (r *parameterRepository) ApplyFilters(query *gorm.DB, filter interface{}) *gorm.DB {
	parameterFilter, ok := filter.(dto.ReqParameterIndexFilter)
	if !ok {
		return query
	}
	return applyParameterFilters(query, parameterFilter)
}

// Compile-time check to ensure parameterRepository implements NeedFilterPredefine interface
var _ request.NeedFilterPredefine = (*parameterRepository)(nil)
