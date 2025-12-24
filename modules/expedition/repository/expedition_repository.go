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
	"github.com/rendyfutsuy/base-go/modules/expedition"
	"github.com/rendyfutsuy/base-go/modules/expedition/dto"
	rsearchexpedition "github.com/rendyfutsuy/base-go/modules/expedition/repository/searches"
	"gorm.io/gorm"
)

type expeditionRepository struct {
	DB *gorm.DB
}

func NewExpeditionRepository(db *gorm.DB) *expeditionRepository {
	return &expeditionRepository{
		DB: db,
	}
}

func (r *expeditionRepository) Create(ctx context.Context, params expedition.CreateExpeditionParams) (*models.Expedition, error) {
	now := time.Now().UTC()
	exp := &models.Expedition{
		ExpeditionName: params.ExpeditionName,
		Address:        params.Address,
		Notes:          params.Notes,
		CreatedAt:      now,
		CreatedBy:      params.CreatedBy,
		UpdatedAt:      now,
		UpdatedBy:      params.CreatedBy,
	}

	// Start transaction
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Omit expedition_code to let database generate it using DEFAULT generate_expedition_code()
	if err := tx.Omit("expedition_code").Create(exp).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if exp.ID == uuid.Nil {
		tx.Rollback()
		return nil, errors.New(constants.ExpeditionCreateFailedIDNotSet)
	}

	// Create contacts
	contacts := make([]models.ExpeditionContact, 0)

	// Process TelpNumbers: first telp (index 0) becomes primary
	for i, telp := range params.TelpNumbers {
		if telp.PhoneNumber != "" {
			contacts = append(contacts, models.ExpeditionContact{
				ExpeditionID: exp.ID,
				PhoneType:    constants.ExpeditionContactTypeTelp,
				PhoneNumber:  telp.PhoneNumber,
				AreaCode:     telp.AreaCode,
				IsPrimary:    i == 0, // First telp is always primary
				CreatedAt:    now,
				CreatedBy:    params.CreatedBy,
				UpdatedAt:    now,
				UpdatedBy:    params.CreatedBy,
			})
		}
	}

	// Process PhoneNumbers: first hp (index 0) becomes primary
	for i, phoneNumber := range params.PhoneNumbers {
		if phoneNumber != "" {
			contacts = append(contacts, models.ExpeditionContact{
				ExpeditionID: exp.ID,
				PhoneType:    constants.ExpeditionContactTypePhone,
				PhoneNumber:  phoneNumber,
				IsPrimary:    i == 0, // First hp is always primary
				CreatedAt:    now,
				CreatedBy:    params.CreatedBy,
				UpdatedAt:    now,
				UpdatedBy:    params.CreatedBy,
			})
		}
	}

	if len(contacts) > 0 {
		if err := tx.Create(&contacts).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return exp, nil
}

func (r *expeditionRepository) Update(ctx context.Context, id uuid.UUID, params expedition.UpdateExpeditionParams) (*models.Expedition, error) {
	updates := map[string]interface{}{
		"expedition_name": params.ExpeditionName,
		"address":         params.Address,
		"updated_at":      time.Now().UTC(),
		"updated_by":      params.UpdatedBy,
	}
	if params.Notes != nil {
		updates["notes"] = *params.Notes
	} else {
		updates["notes"] = nil
	}

	// Start transaction
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	exp := &models.Expedition{}
	err := tx.Model(&models.Expedition{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).
		Take(exp).Error
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	// Update contacts: Always hard delete existing contacts before creating new ones
	now := time.Now().UTC()
	if err := tx.Unscoped().Where("expedition_id = ?", id).
		Delete(&models.ExpeditionContact{}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	contacts := make([]models.ExpeditionContact, 0)

	// Process TelpNumbers: first telp (index 0) becomes primary
	for i, telp := range params.TelpNumbers {
		if telp.PhoneNumber != "" {
			contacts = append(contacts, models.ExpeditionContact{
				ExpeditionID: exp.ID,
				PhoneType:    constants.ExpeditionContactTypeTelp,
				PhoneNumber:  telp.PhoneNumber,
				AreaCode:     telp.AreaCode,
				IsPrimary:    i == 0, // First telp is always primary
				CreatedAt:    now,
				CreatedBy:    params.UpdatedBy,
				UpdatedAt:    now,
				UpdatedBy:    params.UpdatedBy,
			})
		}
	}

	// Process PhoneNumbers: first hp (index 0) becomes primary
	for i, phoneNumber := range params.PhoneNumbers {
		if phoneNumber != "" {
			contacts = append(contacts, models.ExpeditionContact{
				ExpeditionID: exp.ID,
				PhoneType:    constants.ExpeditionContactTypePhone,
				PhoneNumber:  phoneNumber,
				IsPrimary:    i == 0, // First hp is always primary
				CreatedAt:    now,
				CreatedBy:    params.UpdatedBy,
				UpdatedAt:    now,
				UpdatedBy:    params.UpdatedBy,
			})
		}
	}

	if len(contacts) > 0 {
		if err := tx.Create(&contacts).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return exp, nil
}

func (r *expeditionRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	updates := map[string]interface{}{
		"deleted_at": time.Now().UTC(),
		"deleted_by": deletedBy,
	}
	return r.DB.WithContext(ctx).Model(&models.Expedition{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

func (r *expeditionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Expedition, error) {
	exp := &models.Expedition{}
	query := r.DB.WithContext(ctx).Table("expeditions e").
		Select(`
			e.*,
			NOT EXISTS (
				SELECT 1 FROM suppliers s
				WHERE s.expedition_arrives_id = e.id AND s.deleted_at IS NULL
			) AND NOT EXISTS (
				SELECT 1 FROM customers c
				WHERE c.expedition_send_id = e.id AND c.deleted_at IS NULL
			) as deletable
		`).
		Where("e.id = ? AND e.deleted_at IS NULL", id)

	err := query.Scan(exp).Error
	if err != nil {
		return nil, err
	}
	// Scan() doesn't return error for record not found, so check if ID is nil
	if exp.ID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}
	return exp, nil
}

func (r *expeditionRepository) ExistsByExpeditionName(ctx context.Context, expeditionName string, excludeID uuid.UUID) (bool, error) {
	var count int64
	q := r.DB.WithContext(ctx).Unscoped().Model(&models.Expedition{}).Where("expedition_name = ?", expeditionName)
	if excludeID != uuid.Nil {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *expeditionRepository) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqExpeditionIndexFilter) ([]models.Expedition, int, error) {
	var expeditions []models.Expedition
	query := r.DB.WithContext(ctx).Table("expeditions e").
		Select(`
			e.id, 
			e.expedition_code, 
			e.expedition_name, 
			e.address, 
			e.created_at, 
			e.updated_at,
			(SELECT CASE 
				WHEN ec_telp.area_code IS NULL OR ec_telp.area_code = '' THEN ec_telp.phone_number 
				ELSE ec_telp.area_code || '-' || ec_telp.phone_number 
			END
			 FROM expedition_contacts ec_telp 
			 WHERE ec_telp.expedition_id = e.id AND ec_telp.phone_type = 'telp' AND ec_telp.is_primary = true AND ec_telp.deleted_at IS NULL 
			 LIMIT 1) as primary_telp_number,
			(SELECT ec_hp.phone_number FROM expedition_contacts ec_hp WHERE ec_hp.expedition_id = e.id AND ec_hp.phone_type = 'hp' AND ec_hp.is_primary = true AND ec_hp.deleted_at IS NULL LIMIT 1) as primary_phone_number,
			NOT EXISTS (
				SELECT 1 FROM suppliers s
				WHERE s.expedition_arrives_id = e.id AND s.deleted_at IS NULL
			) AND NOT EXISTS (
				SELECT 1 FROM customers c
				WHERE c.expedition_send_id = e.id AND c.deleted_at IS NULL
			) as deletable
		`).
		Where("e.deleted_at IS NULL")

	// Apply search from PageRequest
	query = request.ApplySearchConditionFromInterface(query, req.Search, rsearchexpedition.NewExpeditionSearchHelper())

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	// Pagination
	total, err := request.ApplyPagination(query, req, request.PaginationConfig{
		DefaultSortBy:      "e.created_at",
		DefaultSortOrder:   "DESC",
		MaxPerPage:         100,
		SortMapping:        mapExpeditionIndexSortColumn,
		NaturalSortColumns: []string{"e.expedition_name", "address"}, // Enable natural sorting for e.expedition_name
	}, &expeditions)
	if err != nil {
		return nil, 0, err
	}
	return expeditions, total, nil
}

func (r *expeditionRepository) GetAll(ctx context.Context, filter dto.ReqExpeditionIndexFilter) ([]models.Expedition, error) {
	var expeditions []models.Expedition
	query := r.DB.WithContext(ctx).Table("expeditions e").
		Select(`
			e.id, 
			e.expedition_code, 
			e.expedition_name, 
			e.address, 
			e.notes, 
			e.created_at, 
			e.updated_at,
			(SELECT CASE 
				WHEN ec_telp.area_code IS NULL OR ec_telp.area_code = '' THEN ec_telp.phone_number 
				ELSE ec_telp.area_code || '-' || ec_telp.phone_number 
			END
			 FROM expedition_contacts ec_telp 
			 WHERE ec_telp.expedition_id = e.id AND ec_telp.phone_type = 'telp' AND ec_telp.is_primary = true AND ec_telp.deleted_at IS NULL 
			 LIMIT 1) as primary_telp_number,
			(SELECT ec_hp.phone_number FROM expedition_contacts ec_hp WHERE ec_hp.expedition_id = e.id AND ec_hp.phone_type = 'hp' AND ec_hp.is_primary = true AND ec_hp.deleted_at IS NULL LIMIT 1) as primary_phone_number
		`).
		Where("e.deleted_at IS NULL")

	// Apply search from filter
	query = request.ApplySearchConditionFromInterface(query, filter.Search, rsearchexpedition.NewExpeditionSearchHelper())

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	// Determine sorting with natural sorting support
	sortExpression := request.BuildSortExpressionForExport(
		filter.SortBy,
		filter.SortOrder,
		"e.created_at",
		"DESC",
		mapExpeditionIndexSortColumn,
		[]string{"e.expedition_name", "e.address"}, // Enable natural sorting for expedition_name and address
	)

	// Order results
	if err := query.Order(sortExpression).Find(&expeditions).Error; err != nil {
		return nil, err
	}
	return expeditions, nil
}

func (r *expeditionRepository) GetAllForExport(ctx context.Context, filter dto.ReqExpeditionIndexFilter) ([]dto.ExpeditionExport, error) {
	// First, get all expeditions
	type ExpeditionBase struct {
		ID             uuid.UUID
		ExpeditionCode string
		ExpeditionName string
		Address        string
		UpdatedAt      time.Time
	}

	var expeditionsBase []ExpeditionBase
	query := r.DB.WithContext(ctx).Table("expeditions e").
		Select(`
			e.id,
			e.expedition_code,
			e.expedition_name,
			e.address,
			e.updated_at
		`).
		Where("e.deleted_at IS NULL")

	// Apply search from filter
	query = request.ApplySearchConditionFromInterface(query, filter.Search, rsearchexpedition.NewExpeditionSearchHelper())

	// Apply filters with multiple values support
	query = r.ApplyFilters(query, filter)

	// Determine sorting
	sortBy := "e.created_at"
	if mapped := mapExpeditionIndexSortColumn(filter.SortBy); mapped != "" {
		sortBy = mapped
	}

	sortOrder := request.ValidateAndSanitizeSortOrder(filter.SortOrder)
	if sortOrder == "" {
		sortOrder = "DESC"
	}

	// Order results
	if err := query.Order(sortBy + " " + sortOrder).Find(&expeditionsBase).Error; err != nil {
		return nil, err
	}

	// Get expedition IDs
	expeditionIDs := make([]uuid.UUID, len(expeditionsBase))
	for i, exp := range expeditionsBase {
		expeditionIDs[i] = exp.ID
	}

	// Fetch all HP and Telp phone numbers for all expeditions
	var contacts []struct {
		ExpeditionID uuid.UUID
		PhoneType    string
		PhoneNumber  string
		AreaCode     *string
	}
	if len(expeditionIDs) > 0 {
		if err := r.DB.WithContext(ctx).Table("expedition_contacts").
			Select("expedition_id, phone_type, phone_number, area_code").
			Where("expedition_id IN (?) AND deleted_at IS NULL", expeditionIDs).
			Order("is_primary DESC, created_at ASC").
			Find(&contacts).Error; err != nil {
			return nil, err
		}
	}

	// Group phone numbers by expedition_id and phone_type
	phoneNumbersMap := make(map[uuid.UUID][]string)
	telpNumbersMap := make(map[uuid.UUID][]string)
	for _, contact := range contacts {
		formatted := formatContactNumber(contact.AreaCode, contact.PhoneNumber)
		if contact.PhoneType == "hp" {
			phoneNumbersMap[contact.ExpeditionID] = append(phoneNumbersMap[contact.ExpeditionID], formatted)
		} else if contact.PhoneType == "telp" {
			telpNumbersMap[contact.ExpeditionID] = append(telpNumbersMap[contact.ExpeditionID], formatted)
		}
	}

	// Map to ExpeditionExport
	expeditions := make([]dto.ExpeditionExport, len(expeditionsBase))
	for i, exp := range expeditionsBase {
		expeditions[i] = dto.ExpeditionExport{
			ExpeditionCode: exp.ExpeditionCode,
			ExpeditionName: exp.ExpeditionName,
			Address:        exp.Address,
			PhoneNumbers:   phoneNumbersMap[exp.ID],
			TelpNumbers:    telpNumbersMap[exp.ID],
			UpdatedAt:      exp.UpdatedAt,
		}
	}

	return expeditions, nil
}

func (r *expeditionRepository) GetContactsByExpeditionID(ctx context.Context, expeditionID uuid.UUID) ([]models.ExpeditionContact, error) {
	var contacts []models.ExpeditionContact
	err := r.DB.WithContext(ctx).
		Where("expedition_id = ? AND deleted_at IS NULL", expeditionID).
		Order("is_primary DESC, created_at ASC").
		Find(&contacts).Error
	return contacts, err
}

// Implement expedition.Repository interface
var _ expedition.Repository = (*expeditionRepository)(nil)

func formatContactNumber(areaCode *string, phoneNumber string) string {
	if phoneNumber == "" {
		return ""
	}
	if areaCode == nil {
		return phoneNumber
	}
	ac := strings.TrimSpace(*areaCode)
	if ac == "" {
		return phoneNumber
	}
	return ac + "-" + strings.TrimSpace(phoneNumber)
}
