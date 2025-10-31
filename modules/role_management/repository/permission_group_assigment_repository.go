package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
)

func (repo *roleRepository) ReAssignPermissionsToPermissionGroup(ctx context.Context, id uuid.UUID, permissions []uuid.UUID) error {
	// Delete existing assignments using parameter binding
	err := repo.DB.WithContext(ctx).
		Exec("DELETE FROM permissions_modules WHERE permission_group_id = ?", id).Error

	if err != nil {
		fmt.Println(constants.SQLErrorScanRow, err)
		return err
	}

	// Insert new assignments using parameter binding
	for _, permissionId := range permissions {
		err := repo.DB.WithContext(ctx).
			Exec("INSERT INTO permissions_modules (permission_group_id, permission_id) VALUES (?, ?)",
				id, permissionId).Error

		if err != nil {
			fmt.Println("Error inserting Permission Group:", err)
			return err
		}
	}

	return nil
}
