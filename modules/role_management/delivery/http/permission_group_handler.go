package http

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/helpers/response"
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

	res, total, err := handler.RoleUseCase.GetIndexPermissionGroup(*pageRequest)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	respPermissionGroup := []dto.RespPermissionGroupIndex{}

	for _, v := range res {
		respPermissionGroup = append(respPermissionGroup, dto.ToRespPermissionGroupIndex(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respPermissionGroup, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *RoleManagementHandler) GetAllPermissionGroup(c echo.Context) error {

	res, err := handler.RoleUseCase.GetAllPermissionGroup()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
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
		return c.JSON(http.StatusBadRequest, ResponseError{Message: constants.ErrorUUIDNotRecognized})
	}

	res, err := handler.RoleUseCase.GetPermissionGroupByID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespPermissionGroupDetail(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) GetDuplicatedPermissionGroup(c echo.Context) error {
	req := new(dto.ReqCheckDuplicatedPermissionGroup)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// validate input
	// validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// initialize uid
	uid := uuid.Nil

	// PermissionGroupId can be null
	if req.ExcludedPermissionGroupId != uuid.Nil {
		uid = req.ExcludedPermissionGroupId
	}

	res, err := handler.RoleUseCase.PermissionGroupNameIsNotDuplicated(req.Name, uid)

	// if name havent been uses by existing account info, return not found error
	if res == nil {
		return c.JSON(http.StatusNotFound, ResponseError{Message: "PermissionGroup Info with such name is not found"})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	// if name already uses by existing account info, return PermissionGroup object
	resResp := dto.ToRespPermissionGroup(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleManagementHandler) GetAllPermissionGroupByModule(c echo.Context) error {
	res, err := handler.RoleUseCase.GetAllPermissionGroup()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	modules := make(map[utils.NullString][]dto.RespPermissionGroup)
	moduleNames := []utils.NullString{}
	moduleRead := make(map[utils.NullString]bool)

	for _, group := range res {
		// Check if the module is already in the map, if not initialize it
		if _, ok := modules[group.Module]; !ok {
			modules[group.Module] = []dto.RespPermissionGroup{}
		}

		// Append to the module slice
		modules[group.Module] = append(modules[group.Module], dto.RespPermissionGroup{
			Name: group.Name,
			ID:   group.ID,
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

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respPermissionGroup)

	return c.JSON(http.StatusOK, resp)
}
