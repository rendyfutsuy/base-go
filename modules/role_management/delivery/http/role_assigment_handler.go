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
		return c.JSON(http.StatusBadRequest, ResponseError{Message: constants.ErrorUUIDNotRecognized})
	}

	req := new(dto.ReqUpdatePermissionGroupAssignmentToRole)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// Re-assign Permission groups to role
	res, err := handler.RoleUseCase.ReAssignPermissionByGroup(c, id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespRoleDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) AssignUsersToRole(c echo.Context) error {
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: constants.ErrorUUIDNotRecognized})
	}

	req := new(dto.ReqUpdateAssignUsersToRole)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// Re-assign Permission groups to role
	res, err := handler.RoleUseCase.AssignUsersToRole(c, id, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespRoleDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
