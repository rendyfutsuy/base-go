package repository

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/modules/type/dto"
	"gorm.io/gorm"
)

// applyTypeFilters applies all filters from ReqTypeIndexFilter to the query
func applyTypeFilters(query *gorm.DB, filter dto.ReqTypeIndexFilter) *gorm.DB {
	// Apply filters with multiple values support
	if len(filter.TypeCodes) > 0 {
		query = query.Where("t.type_code IN (?)", filter.TypeCodes)
	}
	if len(filter.Names) > 0 {
		query = query.Where("t.name IN (?)", filter.Names)
	}
	if len(filter.SubgroupIDs) > 0 {
		// Convert string UUIDs to uuid.UUID for query
		subgroupUUIDs := make([]uuid.UUID, 0, len(filter.SubgroupIDs))
		for _, idStr := range filter.SubgroupIDs {
			if id, err := uuid.Parse(idStr); err == nil {
				subgroupUUIDs = append(subgroupUUIDs, id)
			}
		}
		if len(subgroupUUIDs) > 0 {
			query = query.Where("t.subgroup_id IN (?)", subgroupUUIDs)
		}
	}
	return query
}

// ApplyFilters applies filters to the query
// Implements NeedFilterPredefine interface
func (r *typeRepository) ApplyFilters(query *gorm.DB, filter interface{}) *gorm.DB {
	typeFilter, ok := filter.(dto.ReqTypeIndexFilter)
	if !ok {
		return query
	}

	if len(typeFilter.TypeCodes) > 0 {
		query = query.Where("t.type_code IN (?)", typeFilter.TypeCodes)
	}
	if len(typeFilter.Names) > 0 {
		query = query.Where("t.name IN (?)", typeFilter.Names)
	}
	if len(typeFilter.SubgroupIDs) > 0 {
		// Convert string UUIDs to uuid.UUID for query
		subgroupUUIDs := make([]uuid.UUID, 0, len(typeFilter.SubgroupIDs))
		for _, idStr := range typeFilter.SubgroupIDs {
			if id, err := uuid.Parse(idStr); err == nil {
				subgroupUUIDs = append(subgroupUUIDs, id)
			}
		}
		if len(subgroupUUIDs) > 0 {
			query = query.Where("t.subgroup_id IN (?)", subgroupUUIDs)
		}
	}
	if len(typeFilter.GoodGroupIDs) > 0 {
		// Convert string UUIDs to uuid.UUID for query
		groupUUIDs := make([]uuid.UUID, 0, len(typeFilter.GoodGroupIDs))
		for _, idStr := range typeFilter.GoodGroupIDs {
			if id, err := uuid.Parse(idStr); err == nil {
				groupUUIDs = append(groupUUIDs, id)
			}
		}
		if len(groupUUIDs) > 0 {
			query = query.Where("sg.goods_group_id IN (?)", groupUUIDs)
		}
	}
	return applyTypeFilters(query, typeFilter)
}

// Compile-time check to ensure typeRepository implements NeedFilterPredefine interface
var _ request.NeedFilterPredefine = (*typeRepository)(nil)
