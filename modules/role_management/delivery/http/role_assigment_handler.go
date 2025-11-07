package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
)

// role scope
// re-assign role management
// assign users to role

func (handler *RoleManagementHandler) ReAssignPermissionByGroup(c echo.Context) error {
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	req := new(dto.ReqUpdatePermissionGroupAssignmentToRole)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Re-assign Permission groups to role
	res, err := handler.RoleUseCase.ReAssignPermissionByGroup(c, id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// if name already uses by existing account info, return Role object
	modules := []dto.RespPermissionGroupByModule{}

	modules, err = handler.buildPermissionGroupsByModule(c, id)
	if err != nil {
		// Handle specific error for UUID parsing
		if id != "" {
			if _, parseErr := uuid.Parse(id); parseErr != nil {
				return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
			}
		}
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resResp := dto.ToRespRoleDetail(*res, modules)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) AssignUsersToRole(c echo.Context) error {
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, constants.ErrorUUIDNotRecognized))
	}

	req := new(dto.ReqUpdateAssignUsersToRole)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// Re-assign Permission groups to role
	res, err := handler.RoleUseCase.AssignUsersToRole(c, id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	// if name already uses by existing account info, return Role object
	modules := []dto.RespPermissionGroupByModule{}

	resResp := dto.ToRespRoleDetail(*res, modules)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
