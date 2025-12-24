package test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/middleware"
	reqMw "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	expeditionMod "github.com/rendyfutsuy/base-go/modules/expedition"
	expeditionHttp "github.com/rendyfutsuy/base-go/modules/expedition/delivery/http"
	"github.com/rendyfutsuy/base-go/modules/expedition/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockExpeditionUsecase struct {
	mock.Mock
}

func (m *mockExpeditionUsecase) Create(ctx context.Context, req *dto.ReqCreateExpedition, authId string) (*models.Expedition, error) {
	args := m.Called(ctx, req, authId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expedition), args.Error(1)
}

func (m *mockExpeditionUsecase) Update(ctx context.Context, id string, req *dto.ReqUpdateExpedition, authId string) (*models.Expedition, error) {
	args := m.Called(ctx, id, req, authId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expedition), args.Error(1)
}

func (m *mockExpeditionUsecase) Delete(ctx context.Context, id string, authId string) error {
	args := m.Called(ctx, id, authId)
	return args.Error(0)
}

func (m *mockExpeditionUsecase) GetByID(ctx context.Context, id string) (*models.Expedition, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expedition), args.Error(1)
}

func (m *mockExpeditionUsecase) GetContactsByExpeditionID(ctx context.Context, id string) ([]models.ExpeditionContact, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ExpeditionContact), args.Error(1)
}

func (m *mockExpeditionUsecase) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqExpeditionIndexFilter) ([]models.Expedition, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Expedition), args.Int(1), args.Error(2)
}

func (m *mockExpeditionUsecase) GetAll(ctx context.Context, filter dto.ReqExpeditionIndexFilter) ([]models.Expedition, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Expedition), args.Error(1)
}

func (m *mockExpeditionUsecase) Export(ctx context.Context, filter dto.ReqExpeditionIndexFilter) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

type mockMiddlewareAuth struct {
	mock.Mock
}

func (m *mockMiddlewareAuth) AuthorizationCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

type mockMiddlewarePermission struct {
	mock.Mock
}

func (m *mockMiddlewarePermission) PermissionValidation(args []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

type mockMiddlewarePageRequest struct {
	mock.Mock
}

func (m *mockMiddlewarePageRequest) PageRequestCtx(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		pageReq := &request.PageRequest{
			Page:      1,
			PerPage:   10,
			SortBy:    "id",
			SortOrder: "desc",
		}
		c.Set("page_request", pageReq)
		return next(c)
	}
}

func (m *mockMiddlewarePageRequest) PageRequestCtxWithoutLimitation(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return utils.ValidateRequest(i, cv.validator)
}

func newEcho() *echo.Echo {
	e := echo.New()
	v := validator.New()
	utils.RegisterCustomValidator(v)
	e.Validator = &customValidator{validator: v}
	return e
}

func setHandlerValidator(handler *expeditionHttp.ExpeditionHandler, v *validator.Validate) {
	val := reflect.ValueOf(handler).Elem()
	field := val.FieldByName("validator")
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(v))
	}
}

func newExpeditionHandler(mockUC expeditionMod.Usecase, mockAuthMw middleware.IMiddlewareAuth, mockPermMw middleware.IMiddlewarePermission, mockPageReqMw reqMw.IMiddlewarePageRequest) *expeditionHttp.ExpeditionHandler {
	handler := &expeditionHttp.ExpeditionHandler{
		Usecase: mockUC,
	}
	val := reflect.ValueOf(handler).Elem()

	authField := val.FieldByName("middlewareAuth")
	if authField.IsValid() && authField.CanSet() {
		authField.Set(reflect.ValueOf(mockAuthMw))
	}

	permField := val.FieldByName("middlewarePermission")
	if permField.IsValid() && permField.CanSet() {
		permField.Set(reflect.ValueOf(mockPermMw))
	}

	pageReqField := val.FieldByName("mwPageRequest")
	if pageReqField.IsValid() && pageReqField.CanSet() {
		pageReqField.Set(reflect.ValueOf(mockPageReqMw))
	}

	setHandlerValidator(handler, validator.New())
	return handler
}

func TestExpeditionHandler_CreateSuccess(t *testing.T) {
	e := newEcho()
	reqBody := `{"expedition_name":"TEST EXPEDITION","address":"TEST ADDRESS"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/expedition", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", models.User{ID: uuid.New()})

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	expeditionID := uuid.New()
	createdExp := &models.Expedition{
		ID:             expeditionID,
		ExpeditionCode: "EXP001",
		ExpeditionName: "TEST EXPEDITION",
		Address:        "TEST ADDRESS",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*dto.ReqCreateExpedition"), mock.AnythingOfType("string")).
		Return(createdExp, nil).Once()
	mockUC.On("GetContactsByExpeditionID", mock.Anything, expeditionID.String()).
		Return([]models.ExpeditionContact{}, nil).Once()

	err := handler.Create(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Status  int `json:"status"`
		Message string
		Data    dto.RespExpedition `json:"data"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, "TEST EXPEDITION", resp.Data.ExpeditionName)
	mockUC.AssertExpectations(t)
}

func TestExpeditionHandler_CreateValidationError(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodPost, "/v1/expedition", strings.NewReader(`{"expedition_name":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", models.User{ID: uuid.New()})

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	err := handler.Create(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestExpeditionHandler_CreateUsecaseError(t *testing.T) {
	e := newEcho()
	reqBody := `{"expedition_name":"TEST EXPEDITION","address":"TEST ADDRESS"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/expedition", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", models.User{ID: uuid.New()})

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*dto.ReqCreateExpedition"), mock.AnythingOfType("string")).
		Return(nil, errors.New("usecase error")).Once()

	err := handler.Create(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestExpeditionHandler_UpdateSuccess(t *testing.T) {
	e := newEcho()
	expeditionID := uuid.New().String()
	reqBody := `{"expedition_name":"UPDATED EXPEDITION","address":"UPDATED ADDRESS"}`
	req := httptest.NewRequest(http.MethodPut, "/v1/expedition/"+expeditionID, strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(expeditionID)

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	updatedExp := &models.Expedition{
		ID:             uuid.MustParse(expeditionID),
		ExpeditionCode: "EXP001",
		ExpeditionName: "UPDATED EXPEDITION",
		Address:        "UPDATED ADDRESS",
		UpdatedAt:      time.Now(),
	}

	mockUC.On("Update", mock.Anything, expeditionID, mock.AnythingOfType("*dto.ReqUpdateExpedition"), mock.AnythingOfType("string")).
		Return(updatedExp, nil).Once()
	mockUC.On("GetByID", mock.Anything, expeditionID).
		Return(updatedExp, nil).Once()
	mockUC.On("GetContactsByExpeditionID", mock.Anything, expeditionID).
		Return([]models.ExpeditionContact{}, nil).Once()

	err := handler.Update(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestExpeditionHandler_UpdateValidationError(t *testing.T) {
	e := newEcho()
	expeditionID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPut, "/v1/expedition/"+expeditionID, strings.NewReader(`{"expedition_name":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(expeditionID)

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	err := handler.Update(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestExpeditionHandler_DeleteSuccess(t *testing.T) {
	e := newEcho()
	expeditionID := uuid.New().String()
	req := httptest.NewRequest(http.MethodDelete, "/v1/expedition/"+expeditionID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(expeditionID)

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	mockUC.On("Delete", mock.Anything, expeditionID, mock.AnythingOfType("string")).
		Return(nil).Once()

	err := handler.Delete(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestExpeditionHandler_DeleteUsecaseError(t *testing.T) {
	e := newEcho()
	expeditionID := uuid.New().String()
	req := httptest.NewRequest(http.MethodDelete, "/v1/expedition/"+expeditionID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(expeditionID)

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	mockUC.On("Delete", mock.Anything, expeditionID, mock.AnythingOfType("string")).
		Return(errors.New("not found")).Once()

	err := handler.Delete(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestExpeditionHandler_GetIndexSuccess(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/expedition?page=1&per_page=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	mockPageReqMw.PageRequestCtx(func(c echo.Context) error {
		c.Set("page_request", &request.PageRequest{Page: 1, PerPage: 10})
		return nil
	})(c)

	expeditions := []models.Expedition{
		{
			ID:             uuid.New(),
			ExpeditionCode: "EXP001",
			ExpeditionName: "TEST EXPEDITION",
			Address:        "TEST ADDRESS",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	mockUC.On("GetIndex", mock.Anything, mock.AnythingOfType("request.PageRequest"), mock.AnythingOfType("dto.ReqExpeditionIndexFilter")).
		Return(expeditions, 1, nil).Once()

	err := handler.GetIndex(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestExpeditionHandler_GetByIDSuccess(t *testing.T) {
	e := newEcho()
	expeditionID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/v1/expedition/"+expeditionID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(expeditionID)

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	expedition := &models.Expedition{
		ID:             uuid.MustParse(expeditionID),
		ExpeditionCode: "EXP001",
		ExpeditionName: "TEST EXPEDITION",
		Address:        "TEST ADDRESS",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mockUC.On("GetByID", mock.Anything, expeditionID).
		Return(expedition, nil).Once()
	mockUC.On("GetContactsByExpeditionID", mock.Anything, expeditionID).
		Return([]models.ExpeditionContact{}, nil).Once()

	err := handler.GetByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestExpeditionHandler_GetByIDUsecaseError(t *testing.T) {
	e := newEcho()
	expeditionID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/v1/expedition/"+expeditionID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(expeditionID)

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	mockUC.On("GetByID", mock.Anything, expeditionID).
		Return(nil, errors.New("not found")).Once()

	err := handler.GetByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestExpeditionHandler_ExportSuccess(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/expedition/export", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	excelBytes := []byte("PK\x03\x04") // Excel file signature

	mockUC.On("Export", mock.Anything, mock.AnythingOfType("dto.ReqExpeditionIndexFilter")).
		Return(excelBytes, nil).Once()

	err := handler.Export(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, constants.ExcelContent, rec.Header().Get(echo.HeaderContentType))
	mockUC.AssertExpectations(t)
}

func TestExpeditionHandler_ExportUsecaseError(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/expedition/export", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockExpeditionUsecase)
	mockAuthMw := new(mockMiddlewareAuth)
	mockPermMw := new(mockMiddlewarePermission)
	mockPageReqMw := new(mockMiddlewarePageRequest)
	handler := newExpeditionHandler(mockUC, mockAuthMw, mockPermMw, mockPageReqMw)

	mockUC.On("Export", mock.Anything, mock.AnythingOfType("dto.ReqExpeditionIndexFilter")).
		Return(nil, errors.New("export error")).Once()

	err := handler.Export(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockUC.AssertExpectations(t)
}

// TestExpeditionHandler_ExportValidateError is removed because
// expedition_codes filter doesn't have validation that would fail
// The filter accepts any string values, so validation won't fail
