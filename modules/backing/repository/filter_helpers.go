package repository

import (
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/modules/backing/dto"
	"gorm.io/gorm"
)

// applyBackingFilters applies all filters from ReqBackingIndexFilter to the query
func applyBackingFilters(query *gorm.DB, filter dto.ReqBackingIndexFilter) *gorm.DB {
	// Apply filters with multiple values support
	if len(filter.BackingCodes) > 0 {
		query = query.Where("b.backing_code IN (?)", filter.BackingCodes)
	}
	if len(filter.Names) > 0 {
		query = query.Where("b.name IN (?)", filter.Names)
	}
	if len(filter.TypeIDs) > 0 {
		query = query.Where("b.type_id IN (?)", filter.TypeIDs)
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
	if len(filter.GoodGroupIDs) > 0 {
		// Convert string UUIDs to uuid.UUID for query
		groupUUIDs := make([]uuid.UUID, 0, len(filter.GoodGroupIDs))
		for _, idStr := range filter.GoodGroupIDs {
			if id, err := uuid.Parse(idStr); err == nil {
				groupUUIDs = append(groupUUIDs, id)
			}
		}
		if len(groupUUIDs) > 0 {
			query = query.Where("sg.groups_id IN (?)", groupUUIDs)
		}
	}
	return query
}

// ApplyFilters applies filters to the query
// Implements NeedFilterPredefine interface
func (r *backingRepository) ApplyFilters(query *gorm.DB, filter interface{}) *gorm.DB {
	backingFilter, ok := filter.(dto.ReqBackingIndexFilter)
	if !ok {
		return query
	}
	return applyBackingFilters(query, backingFilter)
}

// Compile-time check to ensure backingRepository implements NeedFilterPredefine interface
var _ request.NeedFilterPredefine = (*backingRepository)(nil)
