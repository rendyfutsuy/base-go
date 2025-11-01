package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"gorm.io/gorm"
)

// CreateUser creates a new user information entry in the database.
//
// It takes a ToDBCreateUser parameter and returns an User pointer and an error.
func (repo *userRepository) CreateUser(ctx context.Context, userReq dto.ToDBCreateUser) (userRes *models.User, err error) {
	now := time.Now().UTC()
	expiredAt := now.AddDate(0, 3, 0)

	// Get password template from config (default to "temp" if not configured)
	passwordTemplate := "temp"
	if utils.ConfigVars.Exists("user.default_password_template") {
		passwordTemplate = utils.ConfigVars.String("user.default_password_template")
	}

	userRes = &models.User{
		FullName:          userReq.FullName,
		Username:          userReq.Username,
		Email:             userReq.Email,
		RoleId:            userReq.RoleId,
		Nik:               userReq.Nik,
		IsActive:          userReq.IsActive,
		Gender:            userReq.Gender,
		Password:          passwordTemplate,
		CreatedAt:         now,
		UpdatedAt:         now,
		PasswordExpiredAt: expiredAt,
	}

	// Create user - GORM will insert all fields from struct
	err = repo.DB.WithContext(ctx).Create(userRes).Error

	if err != nil {
		return nil, err
	}

	// Reload only the fields we need to return
	err = repo.DB.WithContext(ctx).
		Select("id", "full_name", "created_at", "updated_at", "deleted_at").
		Where("id = ?", userRes.ID).
		First(userRes).Error

	if err != nil {
		return nil, err
	}

	return userRes, nil
}

// GetUserByID retrieves an user information entry by ID from the database.
//
// It takes a uuid.UUID parameter representing the ID and returns an User pointer and an error.
func (repo *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (user *models.User, err error) {
	user = &models.User{}

	err = repo.DB.WithContext(ctx).
		Table("users usr").
		Select(`
			usr.id,
			usr.full_name,
			usr.email,
			usr.created_at,
			usr.updated_at,
			usr.deleted_at,
			usr.role_id,
			usr.is_active,
			rl.name AS role_name,
			usr.gender,
			CASE 
				WHEN usr.is_active THEN 'active'
				ELSE 'inactive'
			END AS active_status,
			CASE 
				WHEN usr.counter >= 3 THEN true
				ELSE false
			END AS is_blocked
		`).
		Joins("JOIN roles rl ON rl.id = usr.role_id").
		Where("usr.id = ? AND usr.deleted_at IS NULL", id).
		Scan(user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetIndexUser retrieves a paginated list of user information from the database.
//
// It takes a PageRequest parameter and returns a slice of User, the total number of
// user information entries, and an error.
// its can search by user name, user code, user alias_1, user alias_2, user alias_3, user alias_4, user address, user email, user phone_number, type name
func (repo *userRepository) GetIndexUser(ctx context.Context, req request.PageRequest, filter dto.ReqUserIndexFilter) (users []models.User, total int, err error) {
	searchQuery := req.Search

	// Build base query with joins
	query := repo.DB.WithContext(ctx).
		Table("users usr").
		Select(`
			usr.id,
			usr.full_name,
			usr.email,
			usr.gender,
			usr.is_active,
			usr.counter,
			usr.created_at,
			usr.updated_at,
			usr.deleted_at,
			CASE 
				WHEN usr.is_active THEN 'active'
				ELSE 'inactive'
			END AS active_status,
			CASE 
				WHEN usr.counter >= 3 THEN true
				ELSE false
			END AS is_blocked,
			rl.name AS role_name
		`).
		Joins("JOIN roles rl ON rl.id = usr.role_id").
		Where("usr.deleted_at IS NULL")

	// Apply search query with parameter binding
	query = request.ApplySearchCondition(query, searchQuery, []string{
		"usr.full_name",
		"usr.gender",
		"usr.email",
		"rl.name",
	})

	// Apply role IDs filter
	if len(filter.RoleIds) > 0 {
		query = query.Where("rl.id = ANY(?)", pq.Array(filter.RoleIds))
	}

	// Apply role name filter
	if filter.RoleName != "" {
		query = query.Where("rl.name = ?", filter.RoleName)
	}

	// Apply pagination using generic function
	config := request.PaginationConfig{
		DefaultSortBy:    "usr.created_at",
		DefaultSortOrder: "DESC",
		MaxPerPage:       100,
		SortMapping:      repo.SortColumnMapping,
	}

	total, err = request.ApplyPagination(query, req, config, &users)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetAllUser retrieves all user information entries from the database.
//
// Returns a slice of models.User and an error.
func (repo *userRepository) GetAllUser(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := repo.DB.WithContext(ctx).
		Select("id", "full_name", "created_at").
		Where("deleted_at IS NULL").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates an existing user information entry in the database.
//
// It takes an ID of the user information and a ToDBUpdateUser parameter.
// It returns an User pointer and an error.
func (repo *userRepository) UpdateUser(ctx context.Context, id uuid.UUID, userReq dto.ToDBUpdateUser) (userRes *models.User, err error) {
	updates := map[string]interface{}{
		"full_name":  userReq.FullName,
		"email":      userReq.Email,
		"gender":     userReq.Gender,
		"is_active":  userReq.IsActive,
		"role_id":    userReq.RoleId,
		"updated_at": time.Now().UTC(),
	}

	userRes = &models.User{}
	err = repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).
		Select("id", "full_name", "created_at", "updated_at", "deleted_at").
		First(userRes).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}
		return nil, err
	}

	return userRes, nil
}

// SoftDeleteUser soft deletes an user user entry in the database.
//
// It takes an id of type uuid.UUID and an userReq of type dto.ToDBDeleteUser as parameters.
// It returns the soft deleted user user entry of type models.User and an error.
func (repo *userRepository) SoftDeleteUser(ctx context.Context, id uuid.UUID, userReq dto.ToDBDeleteUser) (userRes *models.User, err error) {
	userRes = &models.User{}

	// GORM soft delete automatically sets deleted_at
	err = repo.DB.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Delete(&models.User{}).Error

	if err != nil {
		return nil, err
	}

	// Get the deleted user (with Unscoped to include soft deleted)
	err = repo.DB.WithContext(ctx).
		Unscoped().
		Select("id", "full_name", "created_at", "updated_at", "deleted_at").
		Where("id = ?", id).
		First(userRes).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}
		return nil, err
	}

	return userRes, nil
}

func (repo *userRepository) BlockUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error) {
	userRes = &models.User{}

	err = repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("counter", 4).
		Select("id", "full_name", "counter", "created_at", "updated_at", "deleted_at").
		First(userRes).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}
		return nil, err
	}

	return userRes, nil
}

func (repo *userRepository) UnBlockUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error) {
	userRes = &models.User{}

	err = repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("counter", 0).
		Select("id", "full_name", "counter", "created_at", "updated_at", "deleted_at").
		First(userRes).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}
		return nil, err
	}

	return userRes, nil
}

func (repo *userRepository) ActivateUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error) {
	userRes = &models.User{}

	err = repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", true).
		Select("id", "full_name", "is_active", "created_at", "updated_at", "deleted_at").
		First(userRes).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}
		return nil, err
	}

	return userRes, nil
}

func (repo *userRepository) DisActivateUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error) {
	userRes = &models.User{}

	err = repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", false).
		Select("id", "full_name", "is_active", "created_at", "updated_at", "deleted_at").
		First(userRes).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.UserIDNotFound, id)
		}
		return nil, err
	}

	return userRes, nil
}

// CountUser retrieves the count of user information entries from the database.
//
// Returns a pointer to an integer and an error.
func (repo *userRepository) CountUser(ctx context.Context) (count *int, err error) {
	var result int64
	err = repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Count(&result).Error

	if err != nil {
		return nil, err
	}

	resultInt := int(result)
	count = &resultInt
	return count, nil
}

// EmailIsNotDuplicated checks if an email is not duplicated in the users table, excluding a specific ID if provided.
//
// Parameters:
// - email: the email to check for duplication.
// - excludedId: the ID to exclude from the check. If set to uuid.Nil, no exclusion is applied.
//
// Returns:
// - bool: true if the email is not duplicated, false otherwise.
// - error: an error if the check fails.
func (repo *userRepository) EmailIsNotDuplicated(ctx context.Context, email string, excludedId uuid.UUID) (bool, error) {
	var count int64
	query := repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ? AND deleted_at IS NULL", email)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// UserNameIsNotDuplicated checks if the provided user name is not duplicated in the database.
//
// It takes a name string and an excludedId UUID as parameters.
// It returns a boolean indicating whether the name is not duplicated and an error.
func (repo *userRepository) UserNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	var count int64
	query := repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("full_name = ? AND deleted_at IS NULL", name)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// GetDuplicatedUser retrieves the user information with the given name and excluded ID from the database.
//
// Parameters:
// - name: the name of the user information to retrieve.
// - excludedId: the ID of the user information to exclude from the result.
//
// Returns:
// - user: a pointer to the retrieved user information.
// - err: an error if there was a problem retrieving the user information.
func (repo *userRepository) GetDuplicatedUser(ctx context.Context, name string, excludedId uuid.UUID) (user *models.User, err error) {
	user = &models.User{}

	query := repo.DB.WithContext(ctx).
		Select("id", "full_name", "created_at", "updated_at").
		Where("full_name = ? AND deleted_at IS NULL", name)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err = query.First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// UserNameIsNotDuplicatedOnSoftDeleted checks if the provided user name is not duplicated in the database.
//
// It takes a name string and an excludedId UUID as parameters.
// It returns a boolean indicating whether the name is not duplicated and an error.
func (repo *userRepository) UserNameIsNotDuplicatedOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	var count int64
	query := repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Unscoped().
		Where("full_name = ?", name)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// GetDuplicatedUserOnSoftDeleted retrieves the user information with the given name and excluded ID from the database.
//
// Parameters:
// - name: the name of the user information to retrieve.
// - excludedId: the ID of the user information to exclude from the result.
//
// Returns:
// - user: a pointer to the retrieved user information.
// - err: an error if there was a problem retrieving the user information.
func (repo *userRepository) GetDuplicatedUserOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (user *models.User, err error) {
	user = &models.User{}

	query := repo.DB.WithContext(ctx).
		Unscoped().
		Select("id", "full_name", "created_at", "updated_at").
		Where("full_name = ?", name)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err = query.First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// UsernameIsNotDuplicated checks if a username is not duplicated in the users table, excluding a specific ID if provided.
//
// Parameters:
// - username: the username to check for duplication.
// - excludedId: the ID to exclude from the check. If set to uuid.Nil, no exclusion is applied.
//
// Returns:
// - bool: true if the username is not duplicated, false otherwise.
// - error: an error if the check fails.
func (repo *userRepository) UsernameIsNotDuplicated(ctx context.Context, username string, excludedId uuid.UUID) (bool, error) {
	var count int64
	query := repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ? AND deleted_at IS NULL", username)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

// NikIsNotDuplicated checks if a NIK is not duplicated in the users table, excluding a specific ID if provided.
//
// Parameters:
// - nik: the NIK to check for duplication.
// - excludedId: the ID to exclude from the check. If set to uuid.Nil, no exclusion is applied.
//
// Returns:
// - bool: true if the NIK is not duplicated, false otherwise.
// - error: an error if the check fails.
func (repo *userRepository) NikIsNotDuplicated(ctx context.Context, nik string, excludedId uuid.UUID) (bool, error) {
	var count int64
	query := repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("nik = ? AND deleted_at IS NULL", nik)

	if excludedId != uuid.Nil {
		query = query.Where("id <> ?", excludedId)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (repo *userRepository) SortColumnMapping(selectedSortLabel string) string {
	response := ""
	sortLabels := map[string][]string{
		"id": {
			"id",
		},
		"full_name": {
			"full_name",
			"name",
		},
		"email": {
			"email",
		},
		"gender": {
			"gender",
		},
		"is_active": {
			"is_active",
		},
		"counter": {
			"counter",
		},
		"created_at": {
			"created_at",
		},
		"updated_at": {
			"updated_at",
		},
		"deleted_at": {
			"deleted_at",
		},
		"active_status": {
			"active_status",
		},
		"is_blocked": {
			"is_blocked",
		},
		"role_name": {
			"role_name",
		},
	}

	// Loop through the map
	for DBcolumn, labels := range sortLabels {
		for _, sortLabel := range labels {
			// Check if the current sort label matches the selected sort label
			if sortLabel == selectedSortLabel {
				response = DBcolumn
			}
		}
	}

	return response
}
