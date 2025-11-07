package dto

import "github.com/google/uuid"

// Permission Assignment
type ReqUpdatePermissionGroupAssignmentToRole struct {
	PermissionGroupIds []uuid.UUID `form:"permission_groups" json:"permission_groups" validate:"required,min=1"`
}

func (r *ReqUpdatePermissionGroupAssignmentToRole) ToDBUpdatePermissionGroupAssignmentToRole() ToDBUpdatePermissionGroupAssignmentToRole {
	return ToDBUpdatePermissionGroupAssignmentToRole{
		PermissionGroupIds: r.PermissionGroupIds,
	}
}

type ToDBUpdatePermissionGroupAssignmentToRole struct {
	PermissionGroupIds []uuid.UUID `json:"permission_groups"`
}

// User Assignment
type ReqUpdateAssignUsersToRole struct {
	UserIds []uuid.UUID `form:"users" json:"users" validate:"required,min=1"`
}

func (r *ReqUpdateAssignUsersToRole) ToDBUpdateAssignUsersToRole() ToDBUpdateAssignUsersToRole {
	return ToDBUpdateAssignUsersToRole{
		UserIds: r.UserIds,
	}
}

type ToDBUpdateAssignUsersToRole struct {
	UserIds []uuid.UUID `json:"users"`
}
