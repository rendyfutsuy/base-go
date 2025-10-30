package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
)

func (repo *roleRepository) ReAssignPermissionsToPermissionGroup(ctx context.Context, id uuid.UUID, permissions []uuid.UUID) error {
	// reset role permissions group assignment
	query := `
		DELETE FROM permissions_modules
		WHERE permission_group_id = $1
	`
	// Execute the query and delete requested row.
	_, err := repo.Conn.ExecContext(ctx, query, id)

	// Handle the error.
	if err != nil {
		// Print an error message if delete row fails.
		fmt.Println(constants.SQLErrorScanRow, err)
		return err
	}

	// assign permission group to role, by create new pivot entry that on modules_roles
	for _, permissionId := range permissions {
		// assign permission group to role
		_, err := repo.Conn.ExecContext(ctx, `INSERT INTO permissions_modules
				(permission_group_id, permission_id)
			VALUES
				($1, $2)`,
			id,
			permissionId,
		)

		if err != nil {
			fmt.Println("Error scanning row Permission Group:", err)
			return err
		}
	}

	return nil
}
