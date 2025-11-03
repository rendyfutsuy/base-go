package constants

var (
	RoleErrorRoleNotFound  = "Role Name already exists"
	RoleUnderwriterSupport = "Underwriter Support"
	RoleUnderwriter        = "Underwriter"
	RoleGroupHeadTreaty    = "Group Head Treaty"
)

const (
	// Role errors
	RoleNotFound                = "Role Not Found"
	RoleNotFoundWithID          = "Role with ID `%s` is not Found.."
	RoleNotFoundWithName        = "Role with name %s not found"
	RoleNotFoundWithIDRepo      = "role role with id %s not found"
	RoleHasUsersCannotDelete    = "Role has user. Can't be deleted"
	RoleAssignUsersError        = "Something went wrong when assigning users to role, please check if role and users exist"
	RoleNotExist                = "Not Such Role Exist"
	
	// Permission Group errors
	PermissionGroupNotFoundWithID      = "Function with ID `%s` is not Found.."
	PermissionGroupNotFoundWithIDAlt  = "Permission Group with ID `%s` is not Found.."
	PermissionGroupNotFoundRepo        = "permission_group permission_group with id %s not found"
	PermissionGroupAssignError         = "Something Wrong when assigning Permission Group to Role"
	PermissionGroupFetchError          = "Something Wrong when fetching permission group.."
	
	// User errors (in role context)
	UserNotFoundWithID = "User with ID `%s` is not Found.."
	
	// Role assignment errors
	RoleFetchTotalUserError = "Something Wrong when fetching total user"
	
	// Permission errors
	PermissionNotFoundWithID = "permission permission with id %s not found"
)
