package test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
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
	httpHandler "github.com/rendyfutsuy/base-go/modules/user_management/delivery/http"
	"github.com/rendyfutsuy/base-go/modules/user_management/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockUserManagementUsecase struct {
	mock.Mock
}

func (m *mockUserManagementUsecase) CreateUser(ctx context.Context, req *dto.ReqCreateUser, userID string) (*models.User, error) {
	args := m.Called(ctx, req, userID)
	if user := args.Get(0); user != nil {
		return user.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserManagementUsecase) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if user := args.Get(0); user != nil {
		return user.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserManagementUsecase) GetAllUser(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	if users := args.Get(0); users != nil {
		return users.([]models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserManagementUsecase) GetIndexUser(ctx context.Context, req request.PageRequest, filter dto.ReqUserIndexFilter) ([]models.User, int, error) {
	args := m.Called(ctx, req, filter)
	if users := args.Get(0); users != nil {
		return users.([]models.User), args.Int(1), args.Error(2)
	}
	return nil, 0, args.Error(2)
}

func (m *mockUserManagementUsecase) UpdateUser(ctx context.Context, id string, req *dto.ReqUpdateUser, userID string) (*models.User, error) {
	args := m.Called(ctx, id, req, userID)
	if user := args.Get(0); user != nil {
		return user.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserManagementUsecase) SoftDeleteUser(ctx context.Context, id string, userID string) (*models.User, error) {
	args := m.Called(ctx, id, userID)
	if user := args.Get(0); user != nil {
		return user.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserManagementUsecase) UserNameIsNotDuplicated(ctx context.Context, name string, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, name, id)
	if user := args.Get(0); user != nil {
		return user.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserManagementUsecase) BlockUser(ctx context.Context, id string, req *dto.ReqBlockUser) (*models.User, error) {
	args := m.Called(ctx, id, req)
	if user := args.Get(0); user != nil {
		return user.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserManagementUsecase) ActivateUser(ctx context.Context, id string, req *dto.ReqActivateUser) (*models.User, error) {
	args := m.Called(ctx, id, req)
	if user := args.Get(0); user != nil {
		return user.(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserManagementUsecase) UpdateUserPassword(ctx context.Context, userId string, passwordChunks *dto.ReqUpdateUserPassword) error {
	args := m.Called(ctx, userId, passwordChunks)
	return args.Error(0)
}

func (m *mockUserManagementUsecase) UpdateUserPasswordNoCheckRequired(ctx context.Context, userId string, passwordChunks *dto.ReqUpdateUserPassword) error {
	args := m.Called(ctx, userId, passwordChunks)
	return args.Error(0)
}

func (m *mockUserManagementUsecase) AssertCurrentUserPassword(ctx context.Context, id string, inputtedPassword string) error {
	args := m.Called(ctx, id, inputtedPassword)
	return args.Error(0)
}

func (m *mockUserManagementUsecase) ImportUsersFromExcel(ctx context.Context, filePath string) (*dto.ResImportUsers, error) {
	args := m.Called(ctx, filePath)
	if res := args.Get(0); res != nil {
		return res.(*dto.ResImportUsers), args.Error(1)
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
		if r.Path == path && r.Method == method {
			return true
		}
	}
	return false
}

func createMultipartRequest(t *testing.T, fieldName, filename string, content []byte) (*http.Request, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filename)
	require.NoError(t, err)
	_, err = part.Write(content)
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/v1/user-management/user/import", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	return req, writer.FormDataContentType()
}

func TestNewUserManagementHandler_RegisterRoutes(t *testing.T) {
	e := echo.New()
	mockUC := new(mockUserManagementUsecase)

	authMw := &noopMiddlewareAuth{}
	permMw := &noopMiddlewarePermission{}
	pageMw := &noopPageRequestMiddleware{}

	httpHandler.NewUserManagementHandler(e, mockUC, pageMw, authMw, permMw)

	require.True(t, routeExists(e.Routes(), http.MethodPost, "/v1/user-management/user"))
	require.True(t, routeExists(e.Routes(), http.MethodGet, "/v1/user-management/user"))
	require.True(t, routeExists(e.Routes(), http.MethodPatch, "/v1/user-management/user/:id/password"))
	require.True(t, routeExists(e.Routes(), http.MethodPost, "/v1/user-management/user/import"))
}

func TestUserHandler_CreateUserSuccess(t *testing.T) {
	e := newEcho()
	roleID := uuid.New()
	body := `{"name":"JOHN DOE","username":"JDOE","role_id":"` + roleID.String() + `","email":"test@example.com","nik":"1234567890","password":"PASSWORD!1","password_confirmation":"PASSWORD!1"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/user-management/user", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	authUser := models.User{ID: uuid.New()}
	c.Set("user", authUser)

	mockUC := new(mockUserManagementUsecase)
	handler := &httpHandler.UserManagementHandler{UserUseCase: mockUC}

	mockUC.On("CreateUser", mock.Anything, mock.AnythingOfType("*dto.ReqCreateUser"), authUser.ID.String()).
		Return(&models.User{ID: uuid.New(), FullName: "JOHN DOE"}, nil).Once()

	err := handler.CreateUser(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestUserHandler_CreateUserValidationError(t *testing.T) {
	e := newEcho()
	body := `{"name":"","username":"","role_id":"","email":"","nik":"","password":"","password_confirmation":""}`
	req := httptest.NewRequest(http.MethodPost, "/v1/user-management/user", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", models.User{ID: uuid.New()})

	handler := &httpHandler.UserManagementHandler{UserUseCase: new(mockUserManagementUsecase)}

	err := handler.CreateUser(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_GetIndexUserSuccess(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/user-management/user?page=1&per_page=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	pageReq := &request.PageRequest{Page: 1, PerPage: 10}
	c.Set("page_request", pageReq)

	mockUC := new(mockUserManagementUsecase)
	handler := &httpHandler.UserManagementHandler{UserUseCase: mockUC}

	users := []models.User{{ID: uuid.New(), FullName: "TEST USER"}}
	mockUC.On("GetIndexUser", mock.Anything, *pageReq, mock.AnythingOfType("dto.ReqUserIndexFilter")).
		Return(users, len(users), nil).Once()

	err := handler.GetIndexUser(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestUserHandler_GetUserByIDInvalidUUID(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/user-management/user/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	handler := &httpHandler.UserManagementHandler{UserUseCase: new(mockUserManagementUsecase)}

	err := handler.GetUserByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserHandler_UpdateUserSuperAdminPassword(t *testing.T) {
	e := newEcho()
	userID := uuid.New()
	body := `{"name":"JOHN DOE","username":"JDOE","role_id":"` + uuid.New().String() + `","email":"test@example.com","is_active":true,"gender":"M","password":"PASSWORD!1","password_confirmation":"PASSWORD!1"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/user-management/user/"+userID.String(), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	authUser := models.User{ID: uuid.New(), RoleName: constants.AuthRoleSuperAdmin}
	c.Set("user", authUser)

	mockUC := new(mockUserManagementUsecase)
	handler := &httpHandler.UserManagementHandler{UserUseCase: mockUC}

	mockUC.On("UpdateUser", mock.Anything, userID.String(), mock.AnythingOfType("*dto.ReqUpdateUser"), authUser.ID.String()).
		Return(&models.User{ID: userID, FullName: "JOHN DOE"}, nil).Once()
	mockUC.On("UpdateUserPasswordNoCheckRequired", mock.Anything, userID.String(), mock.AnythingOfType("*dto.ReqUpdateUserPassword")).
		Return(nil).Once()

	err := handler.UpdateUser(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestUserHandler_GetDuplicatedUserNotFound(t *testing.T) {
	e := newEcho()
	body := `{"username":"JDOE"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/user-management/user/check-name", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockUserManagementUsecase)
	handler := &httpHandler.UserManagementHandler{UserUseCase: mockUC}

	mockUC.On("UserNameIsNotDuplicated", mock.Anything, "JDOE", uuid.Nil).Return(nil, nil).Once()

	err := handler.GetDuplicatedUser(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "User Info with such name is not found", resp.Message)
	mockUC.AssertExpectations(t)
}

func TestUserHandler_GetDuplicatedUserDuplicate(t *testing.T) {
	e := newEcho()
	body := `{"username":"JDOE"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/user-management/user/check-name", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockUserManagementUsecase)
	handler := &httpHandler.UserManagementHandler{UserUseCase: mockUC}

	mockUC.On("UserNameIsNotDuplicated", mock.Anything, "JDOE", uuid.Nil).
		Return(&models.User{ID: uuid.New(), FullName: "JOHN DOE"}, nil).Once()

	err := handler.GetDuplicatedUser(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestUserImportHandler_MissingFile(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodPost, "/v1/user-management/user/import", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &httpHandler.UserManagementHandler{UserUseCase: new(mockUserManagementUsecase)}

	err := handler.ImportUsersFromExcel(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserImportHandler_InvalidExtension(t *testing.T) {
	e := newEcho()
	req, contentType := createMultipartRequest(t, "file", "users.txt", []byte("dummy"))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, contentType)
	c := e.NewContext(req, rec)

	handler := &httpHandler.UserManagementHandler{UserUseCase: new(mockUserManagementUsecase)}

	err := handler.ImportUsersFromExcel(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserImportHandler_Success(t *testing.T) {
	e := newEcho()
	req, contentType := createMultipartRequest(t, "file", "users.xlsx", []byte("dummy"))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, contentType)
	c := e.NewContext(req, rec)

	mockUC := new(mockUserManagementUsecase)
	handler := &httpHandler.UserManagementHandler{UserUseCase: mockUC}

	res := &dto.ResImportUsers{FailedCount: 0}
	mockUC.On("ImportUsersFromExcel", mock.Anything, mock.AnythingOfType("string")).
		Return(res, nil).Once()

	err := handler.ImportUsersFromExcel(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestUserImportHandler_PartialFailure(t *testing.T) {
	e := newEcho()
	req, contentType := createMultipartRequest(t, "file", "users.xls", []byte("dummy"))
	rec := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, contentType)
	c := e.NewContext(req, rec)

	mockUC := new(mockUserManagementUsecase)
	handler := &httpHandler.UserManagementHandler{UserUseCase: mockUC}

	res := &dto.ResImportUsers{FailedCount: 2}
	mockUC.On("ImportUsersFromExcel", mock.Anything, mock.AnythingOfType("string")).
		Return(res, nil).Once()

	err := handler.ImportUsersFromExcel(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestDownloadUserImportTemplate(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/user-management/user/import/template", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &httpHandler.UserManagementHandler{UserUseCase: new(mockUserManagementUsecase)}

	err := handler.DownloadUserImportTemplate(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, constants.ExcelContent, rec.Header().Get(constants.FieldContentType))
}

func TestUpdateUserPasswordSuccess(t *testing.T) {
	e := newEcho()
	userID := uuid.New()
	body := `{"old_password":"SECRET","new_password":"PASSWORD!1","password_confirmation":"PASSWORD!1"}`
	req := httptest.NewRequest(http.MethodPatch, "/v1/user-management/user/"+userID.String()+"/password", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(userID.String())

	mockUC := new(mockUserManagementUsecase)
	handler := &httpHandler.UserManagementHandler{UserUseCase: mockUC}

	mockUC.On("UpdateUserPassword", mock.Anything, userID.String(), mock.AnythingOfType("*dto.ReqUpdateUserPassword")).Return(nil).Once()
	mockUC.On("GetUserByID", mock.Anything, userID.String()).Return(&models.User{ID: userID, FullName: "JOHN"}, nil).Once()

	err := handler.UpdateUserPassword(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestUpdateUserPasswordInvalidUUID(t *testing.T) {
	e := newEcho()
	body := `{"old_password":"SECRET","new_password":"PASSWORD!1","password_confirmation":"PASSWORD!1"}`
	req := httptest.NewRequest(http.MethodPatch, "/v1/user-management/user/invalid/password", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	handler := &httpHandler.UserManagementHandler{UserUseCase: new(mockUserManagementUsecase)}

	err := handler.UpdateUserPassword(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestConfirmCurrentUserPasswordSuccess(t *testing.T) {
	e := newEcho()
	body := `{"password":"PASSWORD!1"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/user-management/user/password-confirmation", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userID := uuid.New()
	c.Set("user", models.User{ID: userID})

	mockUC := new(mockUserManagementUsecase)
	handler := &httpHandler.UserManagementHandler{UserUseCase: mockUC}

	mockUC.On("AssertCurrentUserPassword", mock.Anything, userID.String(), "PASSWORD!1").Return(nil).Once()
	mockUC.On("GetUserByID", mock.Anything, userID.String()).Return(&models.User{ID: userID, FullName: "JOHN"}, nil).Once()

	err := handler.ConfirmCurrentUserPassword(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}
