package repository

import (
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/modules/expedition/dto"
	"gorm.io/gorm"
)

// applyExpeditionFilters applies all filters from ReqExpeditionIndexFilter to the query
func applyExpeditionFilters(query *gorm.DB, filter dto.ReqExpeditionIndexFilter) *gorm.DB {
	// Apply filters with multiple values support
	// Priority: ExpeditionCodesOrNames > ExpeditionCodes/ExpeditionNames
	if len(filter.ExpeditionCodesOrNames) > 0 {
		// Use OR condition to match either code or name (for bulk import scenarios)
		query = query.Where("(e.expedition_code IN (?) OR e.expedition_name IN (?))", filter.ExpeditionCodesOrNames, filter.ExpeditionCodesOrNames)
	} else {
		// Use OR condition if both codes and names are provided (since we don't know which one the value is)
		if len(filter.ExpeditionCodes) > 0 && len(filter.ExpeditionNames) > 0 {
			// Check if codes and names are the same (likely from import where we don't know which is which)
			// Use OR condition to match either code or name
			query = query.Where("(e.expedition_code IN (?) OR e.expedition_name IN (?))", filter.ExpeditionCodes, filter.ExpeditionNames)
		} else {
			if len(filter.ExpeditionCodes) > 0 {
				query = query.Where("e.expedition_code IN (?)", filter.ExpeditionCodes)
			}
			if len(filter.ExpeditionNames) > 0 {
				query = query.Where("e.expedition_name IN (?)", filter.ExpeditionNames)
			}
		}
	}
	if len(filter.Addresses) > 0 {
		query = query.Where("e.address IN (?)", filter.Addresses)
	}
	if len(filter.TelpNumbers) > 0 {
		query = query.Where("EXISTS (SELECT 1 FROM expedition_contacts ec WHERE ec.expedition_id = e.id AND ec.deleted_at IS NULL AND ec.phone_type = 'telp' AND ec.phone_number IN (?))", filter.TelpNumbers)
	}
	if len(filter.PhoneNumbers) > 0 {
		query = query.Where("EXISTS (SELECT 1 FROM expedition_contacts ec WHERE ec.expedition_id = e.id AND ec.deleted_at IS NULL AND ec.phone_type = 'hp' AND ec.phone_number IN (?))", filter.PhoneNumbers)
	}
	return query
}

// ApplyFilters applies filters to the query
// Implements NeedFilterPredefine interface
func (r *expeditionRepository) ApplyFilters(query *gorm.DB, filter interface{}) *gorm.DB {
	expeditionFilter, ok := filter.(dto.ReqExpeditionIndexFilter)
	if !ok {
		return query
	}

	return applyExpeditionFilters(query, expeditionFilter)
}

// Compile-time check to ensure expeditionRepository implements NeedFilterPredefine interface
var _ request.NeedFilterPredefine = (*expeditionRepository)(nil)
