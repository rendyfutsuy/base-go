package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"gorm.io/gorm"
)

// ReAssignPermissionGroup resets role permissions group assignment and assigns new permissions
// based on the permission group input.
func (repo *roleRepository) ReAssignPermissionGroup(ctx context.Context, id uuid.UUID, permissionGroupReq dto.ToDBUpdatePermissionGroupAssignmentToRole) error {
	// Delete existing assignments using parameter binding
	err := repo.DB.WithContext(ctx).
		Exec("DELETE FROM modules_roles WHERE role_id = ?", id).Error

	if err != nil {
		fmt.Println(constants.SQLErrorScanRow, err)
		return err
	}

	// Insert new assignments
	for _, permissionGroupId := range permissionGroupReq.PermissionGroupIds {
		err := repo.DB.WithContext(ctx).
			Exec("INSERT INTO modules_roles (permission_group_id, role_id) VALUES (?, ?)",
				permissionGroupId, id).Error

		if err != nil {
			fmt.Println("Error inserting Permission Group:", err)
			return err
		}
	}

	return nil
}

// GetTotalUser retrieves the total number of users associated with a given role ID.
func (repo *roleRepository) GetTotalUser(ctx context.Context, id uuid.UUID) (total int, err error) {
	var count int64
	err = repo.DB.WithContext(ctx).
		Table("users usr").
		Joins("JOIN roles role ON usr.role_id = role.id").
		Where("role.id = ? AND role.deleted_at IS NULL AND usr.deleted_at IS NULL", id).
		Count(&count).Error

	if err != nil {
		// If no rows found, return 0
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}

	total = int(count)
	return total, nil
}

// GetPermissionFromRoleId retrieves the permissions associated with a given role ID.
func (repo *roleRepository) GetPermissionFromRoleId(ctx context.Context, id uuid.UUID) (permissions []models.Permission, err error) {
	err = repo.DB.WithContext(ctx).
		Table("permissions ps").
		Select("DISTINCT ps.id", "ps.name", "pg.module AS module").
		Joins("JOIN permissions_modules ppg ON ps.id = ppg.permission_id").
		Joins("JOIN permission_groups pg ON ppg.permission_group_id = pg.id").
		Joins("JOIN modules_roles pgr ON pg.id = pgr.permission_group_id").
		Where("ps.deleted_at IS NULL AND pgr.role_id = ?", id).
		Find(&permissions).Error

	if err != nil {
		return nil, fmt.Errorf(constants.PermissionGroupFetchError)
	}

	return permissions, nil
}

// GetPermissionGroupFromRoleId retrieves the permission groups associated with a given role ID.
func (repo *roleRepository) GetPermissionGroupFromRoleId(ctx context.Context, id uuid.UUID) (permissionGroups []models.PermissionGroup, err error) {
	err = repo.DB.WithContext(ctx).
		Table("permission_groups pg").
		Select("pg.id", "pg.name", "pg.module").
		Joins("JOIN modules_roles pgr ON pg.id = pgr.permission_group_id").
		Where("pgr.role_id = ?", id).
		Find(&permissionGroups).Error

	if err != nil {
		return nil, fmt.Errorf(constants.PermissionGroupFetchError)
	}

	return permissionGroups, nil
}

// AssignUsers updates the role_id of users in the users table based on the provided roleId and userReq.
func (repo *roleRepository) AssignUsers(ctx context.Context, roleId uuid.UUID, userReq []uuid.UUID) error {
	// Update users in batch using parameter binding
	for _, userId := range userReq {
		err := repo.DB.WithContext(ctx).
			Model(&models.User{}).
			Where("id = ?", userId).
			Update("role_id", roleId).Error

		if err != nil {
			fmt.Printf("Error updating user role (role_id: %s, user_id: %s): %v\n", roleId, userId, err)
			return err
		}
	}

	return nil
}

// GetUserByID retrieves a user by ID
// TOD: This Method only temporarily, after user management module is created. use that instead.
func (repo *roleRepository) GetUserByID(ctx context.Context, id uuid.UUID) (user *models.User, err error) {
	user = &models.User{}

	err = repo.DB.WithContext(ctx).
		Select("id", "full_name").
		Where("id = ?", id).
		First(user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}
