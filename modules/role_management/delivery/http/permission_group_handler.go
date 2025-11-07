package http

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
)

type ModuleFunction struct {
	Name string
	ID   uuid.UUID
}

// permission group scope
// get permission group
// get index permission group
// get all permission group
// export all account based on type

func (handler *RoleManagementHandler) GetIndexPermissionGroup(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.RoleUseCase.GetIndexPermissionGroup(c, *pageRequest)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	respPermissionGroup := []dto.RespPermissionGroupIndex{}

	for _, v := range res {
		respPermissionGroup = append(respPermissionGroup, dto.ToRespPermissionGroupIndex(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respPermissionGroup, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *RoleManagementHandler) GetAllPermissionGroup(c echo.Context) error {

	res, err := handler.RoleUseCase.GetAllPermissionGroup(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	respPermissionGroup := []dto.RespPermissionGroup{}

	for _, v := range res {
		respPermissionGroup = append(respPermissionGroup, dto.ToRespPermissionGroup(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respPermissionGroup)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) GetPermissionGroupByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	res, err := handler.RoleUseCase.GetPermissionGroupByID(c, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resResp := dto.ToRespPermissionGroupDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) GetDuplicatedPermissionGroup(c echo.Context) error {
	req := new(dto.ReqCheckDuplicatedPermissionGroup)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate input
	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// initialize uid
	uid := uuid.Nil

	// PermissionGroupId can be null
	if req.ExcludedPermissionGroupId != uuid.Nil {
		uid = req.ExcludedPermissionGroupId
	}

	res, err := handler.RoleUseCase.PermissionGroupNameIsNotDuplicated(c, req.Name, uid)

	// if name havent been uses by existing account info, return not found error
	if res == nil {
		return c.JSON(http.StatusNotFound, response.SetErrorResponse(http.StatusNotFound, "PermissionGroup Info with such name is not found"))
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// if name already uses by existing account info, return PermissionGroup object
	resResp := dto.ToRespPermissionGroup(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

// buildPermissionGroupsByModule builds permission groups grouped by module with assignment status
// It takes an echo context, optional roleID string, and returns permission groups organized by module
// If roleID is provided, it will mark which permission groups are assigned to that role
func (handler *RoleManagementHandler) buildPermissionGroupsByModule(c echo.Context, roleIDStr string) ([]dto.RespPermissionGroupByModule, error) {
	var role *models.Role
	var rolePermissionGroups map[uuid.UUID]bool

	// If role_id is provided, validate and get role's permission groups
	if roleIDStr != "" {
		roleID, err := uuid.Parse(roleIDStr)
		if err != nil {
			return nil, err
		}

		// Validate role_id exists in database
		role, err = handler.RoleUseCase.GetRoleByID(c, roleID.String())
		if err != nil {
			return nil, err
		}

		// Get permission groups assigned to the role
		rolePermissionGroups = make(map[uuid.UUID]bool)
		if role != nil && role.PermissionGroups != nil {
			for _, pg := range role.PermissionGroups {
				rolePermissionGroups[pg.ID] = true
			}
		}
	} else {
		// If role_id is not provided, create empty map (all values will be false)
		rolePermissionGroups = make(map[uuid.UUID]bool)
	}

	// Get all permission groups
	res, err := handler.RoleUseCase.GetAllPermissionGroup(c)
	if err != nil {
		return nil, err
	}

	modules := make(map[utils.NullString][]dto.RespPermissionGroup)
	moduleNames := []utils.NullString{}
	moduleRead := make(map[utils.NullString]bool)

	for _, group := range res {
		// Check if the module is already in the map, if not initialize it
		if _, ok := modules[group.Module]; !ok {
			modules[group.Module] = []dto.RespPermissionGroup{}
		}

		// Check if this permission group is assigned to the role
		isAssigned := rolePermissionGroups[group.ID]

		// Append to the module slice
		modules[group.Module] = append(modules[group.Module], dto.RespPermissionGroup{
			Name:  group.Name,
			ID:    group.ID,
			Value: isAssigned,
		})

		if moduleRead[group.Module] {
			continue
		}

		moduleNames = append(moduleNames, group.Module)
		moduleRead[group.Module] = true
	}

	respPermissionGroup := []dto.RespPermissionGroupByModule{}

	for _, module := range moduleNames {
		// sort permission name
		sort.Slice(modules[module], func(i, j int) bool {
			return modules[module][i].Name < modules[module][j].Name
		})

		respPermissionGroup = append(respPermissionGroup, dto.RespPermissionGroupByModule{
			Name:             module,
			PermissionGroups: modules[module], // Assign the slice of functions
		})
	}

	return respPermissionGroup, nil
}

// GetAllPermissionGroupByModule godoc
// @Summary		Get all permission groups grouped by module
// @Description	Retrieve all permission groups organized by module. Returns permission groups with a value field indicating if the permission group is assigned to the specified role. The role_id query parameter is optional.
// @Tags			Role Management
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			role_id	query		string	false	"Optional role ID to check permission group assignment"
// @Success		200		{object}	response.NonPaginationResponse{data=[]dto.RespPermissionGroupByModule}	"Successfully retrieved permission groups by module"
// @Failure		400		{object}	response.NonPaginationResponse	"Bad request - validation error or role not found"
// @Failure		401		{object}	response.NonPaginationResponse	"Unauthorized"
// @Router			/v1/role-management/role/module-access [get]
func (handler *RoleManagementHandler) GetAllPermissionGroupByModule(c echo.Context) error {
	// Get role_id from query parameter (optional)
	roleIDStr := c.QueryParam("role_id")

	respPermissionGroup, err := handler.buildPermissionGroupsByModule(c, roleIDStr)
	if err != nil {
		// Handle specific error for UUID parsing
		if roleIDStr != "" {
			if _, parseErr := uuid.Parse(roleIDStr); parseErr != nil {
				return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
			}
		}
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respPermissionGroup)

	return c.JSON(http.StatusOK, resp)
}
