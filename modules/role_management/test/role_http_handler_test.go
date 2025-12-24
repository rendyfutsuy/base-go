package test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	roleHttp "github.com/rendyfutsuy/base-go/modules/role_management/delivery/http"
	"github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRoleManagementUsecase struct {
	mock.Mock
}

func (m *mockRoleManagementUsecase) CreateRole(ctx context.Context, req *dto.ReqCreateRole, userID string) (*models.Role, error) {
	args := m.Called(ctx, req, userID)
	if role := args.Get(0); role != nil {
		return role.(*models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) GetRoleByID(ctx context.Context, id string) (*models.Role, error) {
	args := m.Called(ctx, id)
	if role := args.Get(0); role != nil {
		return role.(*models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) GetAllRole(ctx context.Context) ([]models.Role, error) {
	args := m.Called(ctx)
	if roles := args.Get(0); roles != nil {
		return roles.([]models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) GetIndexRole(ctx context.Context, req request.PageRequest) ([]models.Role, int, error) {
	args := m.Called(ctx, req)
	if roles := args.Get(0); roles != nil {
		return roles.([]models.Role), args.Int(1), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *mockRoleManagementUsecase) UpdateRole(ctx context.Context, id string, req *dto.ReqUpdateRole, userID string) (*models.Role, error) {
	args := m.Called(ctx, id, req, userID)
	if role := args.Get(0); role != nil {
		return role.(*models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) SoftDeleteRole(ctx context.Context, id string, userID string) (*models.Role, error) {
	args := m.Called(ctx, id, userID)
	if role := args.Get(0); role != nil {
		return role.(*models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) RoleNameIsNotDuplicated(ctx context.Context, name string, id uuid.UUID) (*models.Role, error) {
	args := m.Called(ctx, name, id)
	if role := args.Get(0); role != nil {
		return role.(*models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) MyPermissionsByUserToken(ctx context.Context, token string) (*models.Role, error) {
	args := m.Called(ctx, token)
	if role := args.Get(0); role != nil {
		return role.(*models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) ReAssignPermissionByGroup(ctx context.Context, roleId string, req *dto.ReqUpdatePermissionGroupAssignmentToRole) (*models.Role, error) {
	args := m.Called(ctx, roleId, req)
	if role := args.Get(0); role != nil {
		return role.(*models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) AssignUsersToRole(ctx context.Context, roleId string, req *dto.ReqUpdateAssignUsersToRole) (*models.Role, error) {
	args := m.Called(ctx, roleId, req)
	if role := args.Get(0); role != nil {
		return role.(*models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) GetPermissionGroupByID(ctx context.Context, id string) (*models.PermissionGroup, error) {
	args := m.Called(ctx, id)
	if pg := args.Get(0); pg != nil {
		return pg.(*models.PermissionGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) GetAllPermissionGroup(ctx context.Context) ([]models.PermissionGroup, error) {
	args := m.Called(ctx)
	if groups := args.Get(0); groups != nil {
		return groups.([]models.PermissionGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) GetIndexPermissionGroup(ctx context.Context, req request.PageRequest) ([]models.PermissionGroup, int, error) {
	args := m.Called(ctx, req)
	if groups := args.Get(0); groups != nil {
		return groups.([]models.PermissionGroup), args.Int(1), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *mockRoleManagementUsecase) PermissionGroupNameIsNotDuplicated(ctx context.Context, name string, id uuid.UUID) (*models.PermissionGroup, error) {
	args := m.Called(ctx, name, id)
	if pg := args.Get(0); pg != nil {
		return pg.(*models.PermissionGroup), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) GetPermissionByID(ctx context.Context, id string) (*models.Permission, error) {
	args := m.Called(ctx, id)
	if perm := args.Get(0); perm != nil {
		return perm.(*models.Permission), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) GetAllPermission(ctx context.Context) ([]models.Permission, error) {
	args := m.Called(ctx)
	if perms := args.Get(0); perms != nil {
		return perms.([]models.Permission), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleManagementUsecase) GetIndexPermission(ctx context.Context, req request.PageRequest) ([]models.Permission, int, error) {
	args := m.Called(ctx, req)
	if perms := args.Get(0); perms != nil {
		return perms.([]models.Permission), args.Int(1), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *mockRoleManagementUsecase) PermissionNameIsNotDuplicated(ctx context.Context, name string, id uuid.UUID) (*models.Permission, error) {
	args := m.Called(ctx, name, id)
	if perm := args.Get(0); perm != nil {
		return perm.(*models.Permission), args.Error(1)
	}
	return nil, args.Error(1)
}

type noopMiddlewareAuth struct{}

func (n *noopMiddlewareAuth) AuthorizationCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

type noopMiddlewarePermission struct{}

func (n *noopMiddlewarePermission) PermissionValidation(args []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

type noopPageRequestMiddleware struct{}

func (n *noopPageRequestMiddleware) PageRequestCtx(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func (n *noopPageRequestMiddleware) PageRequestCtxWithoutLimitation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

func newEcho() *echo.Echo {
	e := echo.New()
	v := validator.New()
	utils.RegisterCustomValidator(v)
	e.Validator = &utils.CustomValidator{Validator: v}
	return e
}

func routeExists(routes []*echo.Route, method, path string) bool {
	for _, r := range routes {
		if r.Method == method && r.Path == path {
			return true
		}
	}
	return false
}

func TestNewRoleManagementHandler_RegisterRoutes(t *testing.T) {
	e := echo.New()
	mockUC := new(mockRoleManagementUsecase)
	authMw := &noopMiddlewareAuth{}
	permMw := &noopMiddlewarePermission{}
	pageMw := &noopPageRequestMiddleware{}

	roleHttp.NewRoleManagementHandler(e, mockUC, pageMw, authMw, permMw)

	require.True(t, routeExists(e.Routes(), http.MethodPost, "/v1/role-management/role"))
	require.True(t, routeExists(e.Routes(), http.MethodGet, "/v1/role-management/role/:id"))
	require.True(t, routeExists(e.Routes(), http.MethodPatch, "/v1/role-management/role/:id/re-assign-permission-groups"))
	require.True(t, routeExists(e.Routes(), http.MethodPatch, "/v1/role-management/role/:id/assign-users"))
}

func TestRoleHandler_CreateRoleSuccess(t *testing.T) {
	e := newEcho()
	body := fmt.Sprintf(`{"role_name":"ADMIN","description":"desc","accesses":["%s"]}`, uuid.New().String())
	req := httptest.NewRequest(http.MethodPost, "/v1/role-management/role", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	authUser := models.User{ID: uuid.New()}
	c.Set("user", authUser)

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	mockUC.On("CreateRole", mock.Anything, mock.AnythingOfType("*dto.ReqCreateRole"), authUser.ID.String()).
		Return(&models.Role{ID: uuid.New(), Name: "ADMIN"}, nil).Once()

	err := handler.CreateRole(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestRoleHandler_CreateRoleValidationError(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodPost, "/v1/role-management/role", strings.NewReader(`{"role_name":"","accesses":[]}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", models.User{ID: uuid.New()})

	handler := &roleHttp.RoleManagementHandler{RoleUseCase: new(mockRoleManagementUsecase)}

	err := handler.CreateRole(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRoleHandler_GetIndexRoleSuccess(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/role-management/role?page=1&per_page=20", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	pageReq := request.NewPageRequest(1, 20, "", "", "")
	c.Set("page_request", pageReq)

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	mockUC.On("GetIndexRole", mock.Anything, *pageReq).
		Return([]models.Role{{ID: uuid.New(), Name: "ADMIN"}}, 1, nil).Once()

	err := handler.GetIndexRole(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestRoleHandler_GetRoleByIDInvalidUUID(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/role-management/role/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	handler := &roleHttp.RoleManagementHandler{RoleUseCase: new(mockRoleManagementUsecase)}

	err := handler.GetRoleByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRoleHandler_DeleteRolePrevented(t *testing.T) {
	e := newEcho()
	roleID := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/v1/role-management/role/"+roleID.String(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(roleID.String())
	authUser := models.User{ID: uuid.New()}
	c.Set("user", authUser)

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	mockUC.On("SoftDeleteRole", mock.Anything, roleID.String(), authUser.ID.String()).
		Return(&models.Role{ID: roleID, TotalUser: 3}, nil).Once()

	err := handler.DeleteRole(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, constants.ErrorCannotDeleteRole, resp.Message)
	mockUC.AssertExpectations(t)
}

func TestRoleHandler_GetDuplicatedRoleNotFound(t *testing.T) {
	e := newEcho()
	body := `{"role_name":"ADMIN"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/role-management/role/check-name", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	mockUC.On("RoleNameIsNotDuplicated", mock.Anything, "ADMIN", uuid.Nil).Return(nil, nil).Once()

	err := handler.GetDuplicatedRole(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockUC.AssertExpectations(t)
}

func TestRoleHandler_GetDuplicatedRoleConflict(t *testing.T) {
	e := newEcho()
	body := `{"role_name":"ADMIN"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/role-management/role/check-name", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	mockUC.On("RoleNameIsNotDuplicated", mock.Anything, "ADMIN", uuid.Nil).
		Return(&models.Role{ID: uuid.New(), Name: "ADMIN"}, nil).Once()

	err := handler.GetDuplicatedRole(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestRoleHandler_GetMyPermissionsSuccess(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/role-management/role/module-access", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("token", "token-123")

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	roleID := uuid.New()
	assignedPG := uuid.New()
	mockUC.On("MyPermissionsByUserToken", mock.Anything, "token-123").
		Return(&models.Role{ID: roleID, Name: "MANAGER"}, nil).Once()
	mockUC.On("GetRoleByID", mock.Anything, roleID.String()).
		Return(&models.Role{
			ID:   roleID,
			Name: "MANAGER",
			PermissionGroups: []models.PermissionGroup{
				{ID: assignedPG},
			},
		}, nil).Once()
	mockUC.On("GetAllPermissionGroup", mock.Anything).
		Return([]models.PermissionGroup{
			{ID: assignedPG, Name: "GENERAL", Module: utils.NullString{String: "General", Valid: true}},
			{ID: uuid.New(), Name: constants.UserPermissionNameCreate, Module: utils.NullString{String: constants.UserPermissionModuleName, Valid: true}},
		}, nil).Once()

	err := handler.GetMyPermissions(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestRoleAssignment_ReAssignPermissionByGroupSuccess(t *testing.T) {
	e := newEcho()
	roleID := uuid.New()
	body := fmt.Sprintf(`{"permission_groups":["%s"]}`, uuid.New().String())
	req := httptest.NewRequest(http.MethodPatch, "/v1/role-management/role/"+roleID.String()+"/re-assign-permission-groups", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(roleID.String())

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	mockUC.On("ReAssignPermissionByGroup", mock.Anything, roleID.String(), mock.AnythingOfType("*dto.ReqUpdatePermissionGroupAssignmentToRole")).
		Return(&models.Role{ID: roleID}, nil).Once()
	mockUC.On("GetRoleByID", mock.Anything, roleID.String()).
		Return(&models.Role{ID: roleID}, nil).Once()
	mockUC.On("GetAllPermissionGroup", mock.Anything).
		Return([]models.PermissionGroup{}, nil).Once()

	err := handler.ReAssignPermissionByGroup(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestRoleAssignment_AssignUsersToRoleSuccess(t *testing.T) {
	e := newEcho()
	roleID := uuid.New()
	body := fmt.Sprintf(`{"users":["%s"]}`, uuid.New().String())
	req := httptest.NewRequest(http.MethodPatch, "/v1/role-management/role/"+roleID.String()+"/assign-users", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(roleID.String())

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	mockUC.On("AssignUsersToRole", mock.Anything, roleID.String(), mock.AnythingOfType("*dto.ReqUpdateAssignUsersToRole")).
		Return(&models.Role{ID: roleID}, nil).Once()

	err := handler.AssignUsersToRole(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestRoleAssignment_ReAssignPermissionByGroupInvalidUUID(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodPatch, "/v1/role-management/role/invalid/re-assign-permission-groups", strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	handler := &roleHttp.RoleManagementHandler{RoleUseCase: new(mockRoleManagementUsecase)}

	err := handler.ReAssignPermissionByGroup(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPermissionHandler_GetIndexPermissionSuccess(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/role-management/permission?page=1&per_page=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	pageReq := request.NewPageRequest(1, 5, "", "", "")
	c.Set("page_request", pageReq)

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	mockUC.On("GetIndexPermission", mock.Anything, *pageReq).
		Return([]models.Permission{{ID: uuid.New(), Name: "PERMISSION_VIEW"}}, 1, nil).Once()

	err := handler.GetIndexPermission(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestPermissionHandler_GetDuplicatedPermissionNotFound(t *testing.T) {
	e := newEcho()
	body := `{"name":"PERMISSION_CREATE"}`
	req := httptest.NewRequest(http.MethodGet, "/v1/role-management/permission/check-name", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	mockUC.On("PermissionNameIsNotDuplicated", mock.Anything, "PERMISSION_CREATE", uuid.Nil).Return(nil, nil).Once()

	err := handler.GetDuplicatedPermission(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestPermissionGroupHandler_GetIndexSuccess(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/role-management/permission-group?page=1&per_page=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	pageReq := request.NewPageRequest(1, 5, "", "", "")
	c.Set("page_request", pageReq)

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	mockUC.On("GetIndexPermissionGroup", mock.Anything, *pageReq).
		Return([]models.PermissionGroup{{ID: uuid.New(), Name: "General"}}, 1, nil).Once()

	err := handler.GetIndexPermissionGroup(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestPermissionGroupHandler_GetAllPermissionGroupByModuleFiltersUserOps(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/role-management/role/module-access", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockRoleManagementUsecase)
	handler := &roleHttp.RoleManagementHandler{RoleUseCase: mockUC}

	generalModule := utils.NullString{String: "General", Valid: true}
	userModule := utils.NullString{String: constants.UserPermissionModuleName, Valid: true}
	mockUC.On("GetAllPermissionGroup", mock.Anything).
		Return([]models.PermissionGroup{
			{ID: uuid.New(), Name: "General Access", Module: generalModule},
			{ID: uuid.New(), Name: constants.UserPermissionNameCreate, Module: userModule},
			{ID: uuid.New(), Name: constants.UserPermissionNameDelete, Module: userModule},
		}, nil).Once()

	err := handler.GetAllPermissionGroupByModule(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Data []dto.RespPermissionGroupByModule `json:"data"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Len(t, resp.Data, 1)
	require.Equal(t, "General", resp.Data[0].Name)
	assert.Len(t, resp.Data[0].PermissionGroups, 1)
	mockUC.AssertExpectations(t)
}
