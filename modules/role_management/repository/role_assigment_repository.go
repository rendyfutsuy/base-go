package repository

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
)

// ReAssignPermissionGroup resets role permissions group assignment and assigns new permissions
// based on the permission group input.
//
// Parameters:
// - id: The UUID of the role to be reassigned.
// - permissionGroupReq: The request containing the permission group IDs to be assigned to the role.
//
// Returns:
// - error: An error if there was a problem executing the queries.
func (repo *roleRepository) ReAssignPermissionGroup(id uuid.UUID, permissionGroupReq dto.ToDBUpdatePermissionGroupAssignmentToRole) error {
	// reset role permissions group assignment
	query := `
		DELETE FROM modules_roles
		WHERE role_id = $1
	`
	// Execute the query and delete requested row.
	_, err := repo.Conn.Exec(query, id)

	// Handle the error.
	if err != nil {
		// Print an error message if delete row fails.
		fmt.Println(constants.SQLErrorScanRow, err)
		return err
	}

	// get permission based on permission group input, with query permission that permission group have by pivot table permissions_modules
	// construct variable
	// Construct the placeholders for the query
	placeholders := make([]string, len(permissionGroupReq.PermissionGroupIds))
	args := make([]interface{}, len(permissionGroupReq.PermissionGroupIds))

	// set how many arguments for query syntax based on passed permission group ids
	for i, id := range permissionGroupReq.PermissionGroupIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	// setup query syntax to get permission IDs
	query = fmt.Sprintf(`
		SELECT 
			DISTINCT permission_id
		FROM 
			permissions_modules
		WHERE
			permission_group_id IN (%s)`, strings.Join(placeholders, ", "))

	// Execute the query with all the arguments
	rows, err := repo.Conn.Query(query, args...)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	defer rows.Close()

	// assign permission group to role, by create new pivot entry that on modules_roles
	for _, permissionGroupId := range permissionGroupReq.PermissionGroupIds {
		// assign permission group to role
		_, err := repo.Conn.Exec(`INSERT INTO modules_roles
				(permission_group_id, role_id)
			VALUES
				($1, $2)`,
			permissionGroupId,
			id,
		)

		if err != nil {
			fmt.Println("Error scanning row Permission Group:", err)
			return err
		}
	}

	return nil
}

// GetTotalUser retrieves the total number of users associated with a given role ID.
//
// Parameters:
// - id: The UUID of the role.
//
// Returns:
// - total: The total number of users associated with the role.
// - err: An error if there was a problem executing the query.
func (repo *roleRepository) GetTotalUser(id uuid.UUID) (total int, err error) {
	// setup query syntax to get total user
	err = repo.Conn.QueryRow(
		`SELECT 
			COUNT(*) AS total_user
		FROM 
			users usr
		JOIN
			roles role
		ON
			usr.role_id = role.id
		WHERE 
			role.id = $1 AND role.deleted_at IS NULL
		GROUP BY
			role.id`,
		id,
	).Scan(
		&total,
	)

	// if no error and total is zero
	if total == 0 {
		return 0, nil
	}

	// if error occurs, return error
	return total, nil
}

// GetPermissionFromRoleId retrieves the permissions associated with a given role ID.
//
// Parameters:
// - id: The UUID of the role.
//
// Returns:
// - permissions: A slice of models.Permission representing the permissions assigned to the role.
// - err: An error if there was a problem fetching the permissions.
func (repo *roleRepository) GetPermissionFromRoleId(id uuid.UUID) (permissions []models.Permission, err error) {
	// Fetch and assign permissions that role has
	permissionQuery := `SELECT
			DISTINCT ps.id,
			ps.name
		FROM
			permissions ps
		JOIN
			permissions_modules ppg
		ON
			ps.id = ppg.permission_id
		JOIN
			modules_roles pgr
		ON
			ppg.permission_group_id = pgr.permission_group_id
		WHERE
			ps.deleted_at IS NULL
		AND
			pgr.role_id = $1`

	rows, err := repo.Conn.Query(permissionQuery, id)

	if err != nil {
		return nil, fmt.Errorf("Something Wrong when fetching permission group..")
	}

	defer rows.Close()

	// assign Permission to temp variable
	for rows.Next() {
		var permission models.Permission
		err = rows.Scan(
			&permission.ID,
			&permission.Name,
		)

		if err != nil {
			continue
		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// GetPermissionGroupFromRoleId retrieves the permission groups associated with a given role ID.
//
// Parameters:
// - id: The UUID of the role.
//
// Returns:
// - permissionGroups: A slice of models.PermissionGroup representing the permission groups assigned to the role.
// - err: An error if there was a problem fetching the permission groups.
func (repo *roleRepository) GetPermissionGroupFromRoleId(id uuid.UUID) (permissionGroups []models.PermissionGroup, err error) {
	permissionGroupsQuery := `SELECT
			pg.id,
			pg.name,
			pg.module
		FROM
			permission_groups pg
		JOIN
			modules_roles pgr
		ON
			pg.id = pgr.permission_group_id
		WHERE
			pgr.role_id = $1`

	groupRows, err := repo.Conn.Query(permissionGroupsQuery, id)

	if err != nil {
		return nil, fmt.Errorf("Something Wrong when fetching permission group..")
	}
	defer groupRows.Close()

	for groupRows.Next() {
		var permissionGroup models.PermissionGroup
		err = groupRows.Scan(
			&permissionGroup.ID,
			&permissionGroup.Name,
			&permissionGroup.Module,
		)

		if err != nil {
			continue
		}

		permissionGroups = append(permissionGroups, permissionGroup)
	}
	return permissionGroups, nil
}

// AssignUsers updates the role_id of users in the users table based on the provided roleId and userReq.
//
// Parameters:
// - roleId: the UUID of the role to assign to the users.
// - userReq: a slice of UUIDs representing the users to assign the role to.
//
// Returns:
// - error: an error if there was a problem updating the users' role_id in the database.
func (repo *roleRepository) AssignUsers(roleId uuid.UUID, userReq []uuid.UUID) error {
	// Assign the role to users by updating the role_id in the users table
	for _, userId := range userReq {
		// Update the user's role_id in the users table
		_, err := repo.Conn.Exec(`UPDATE users SET role_id = $1 WHERE id = $2`,
			roleId,
			userId,
		)

		// if error occurs, return error
		if err != nil {
			fmt.Printf("Error updating user role (role_id: %s, user_id: %s): %v\n", roleId, userId, err)
			return err
		}
	}

	return nil
}

// TOD: This Method only temporarily, after user management module is created. use that instead.
// this is only to fill the gap of user management module
func (repo *roleRepository) GetUserByID(id uuid.UUID) (user *models.User, err error) {
	// Initialize the user variable
	user = &models.User{}

	err = repo.Conn.QueryRow(
		`SELECT 
			id, full_name
		FROM 
			users
		WHERE 
			id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.FullName,
	)

	// If error occurs, return error
	if err != nil {
		return nil, err
	}

	// If no error, return user
	return user, nil
}
