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
	backingHttp "github.com/rendyfutsuy/base-go/modules/backing/delivery/http"
	"github.com/rendyfutsuy/base-go/modules/backing/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockBackingUsecase struct {
	mock.Mock
}

func (m *mockBackingUsecase) Create(ctx context.Context, req *dto.ReqCreateBacking, userID string) (*models.Backing, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Backing), args.Error(1)
}

func (m *mockBackingUsecase) Update(ctx context.Context, id string, req *dto.ReqUpdateBacking, userID string) (*models.Backing, error) {
	args := m.Called(ctx, id, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Backing), args.Error(1)
}

func (m *mockBackingUsecase) Delete(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *mockBackingUsecase) GetByID(ctx context.Context, id string) (*models.Backing, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Backing), args.Error(1)
}

func (m *mockBackingUsecase) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqBackingIndexFilter) ([]models.Backing, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Backing), args.Int(1), args.Error(2)
}

func (m *mockBackingUsecase) GetAll(ctx context.Context, filter dto.ReqBackingIndexFilter) ([]models.Backing, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Backing), args.Error(1)
}

func (m *mockBackingUsecase) Export(ctx context.Context, filter dto.ReqBackingIndexFilter) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
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

func TestNewBackingHandler_RegisterRoutes(t *testing.T) {
	e := echo.New()
	mockUC := new(mockBackingUsecase)

	backingHttp.NewBackingHandler(e, mockUC, &noopPageRequestMiddleware{}, &noopMiddlewareAuth{}, &noopMiddlewarePermission{})

	require.NotNil(t, findRoute(e.Routes(), http.MethodPost, "/v1/backing"))
	require.NotNil(t, findRoute(e.Routes(), http.MethodGet, "/v1/backing/:id"))
	require.NotNil(t, findRoute(e.Routes(), http.MethodGet, "/v1/backing/export"))
}

func findRoute(routes []*echo.Route, method, path string) *echo.Route {
	for _, r := range routes {
		if r.Method == method && r.Path == path {
			return r
		}
	}
	return nil
}

func TestBackingHandler_CreateSuccess(t *testing.T) {
	e := newEcho()
	payload := fmt.Sprintf(`{"type_id":"%s","name":"BACKING-A"}`, uuid.New().String())
	req := httptest.NewRequest(http.MethodPost, "/v1/backing", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	userID := uuid.New()
	c.Set("user", models.User{ID: userID})

	mockUC := new(mockBackingUsecase)
	handler := &backingHttp.BackingHandler{Usecase: mockUC}

	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*dto.ReqCreateBacking"), userID.String()).
		Return(&models.Backing{ID: uuid.New(), Name: "BACKING-A"}, nil).Once()

	err := handler.Create(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestBackingHandler_CreateValidationError(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodPost, "/v1/backing", strings.NewReader(`{"name":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", models.User{ID: uuid.New()})

	handler := &backingHttp.BackingHandler{Usecase: new(mockBackingUsecase)}

	err := handler.Create(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestBackingHandler_UpdateSuccess(t *testing.T) {
	e := newEcho()
	backingID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPut, "/v1/backing/"+backingID, strings.NewReader(fmt.Sprintf(`{"type_id":"%s","name":"UPDATED"}`, uuid.New().String())))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(backingID)

	mockUC := new(mockBackingUsecase)
	handler := &backingHttp.BackingHandler{Usecase: mockUC}

	mockUC.On("Update", mock.Anything, backingID, mock.AnythingOfType("*dto.ReqUpdateBacking"), "").
		Return(&models.Backing{ID: uuid.MustParse(backingID), Name: "UPDATED"}, nil).Once()

	err := handler.Update(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestBackingHandler_DeleteSuccess(t *testing.T) {
	e := newEcho()
	backingID := uuid.New().String()
	req := httptest.NewRequest(http.MethodDelete, "/v1/backing/"+backingID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(backingID)

	mockUC := new(mockBackingUsecase)
	handler := &backingHttp.BackingHandler{Usecase: mockUC}

	mockUC.On("Delete", mock.Anything, backingID, "").Return(nil).Once()

	err := handler.Delete(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestBackingHandler_GetIndexSuccess(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/backing?page=1&per_page=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	pageReq := request.NewPageRequest(1, 10, "", "", "")
	c.Set("page_request", pageReq)

	mockUC := new(mockBackingUsecase)
	handler := &backingHttp.BackingHandler{Usecase: mockUC}

	mockUC.On("GetIndex", mock.Anything, *pageReq, dto.ReqBackingIndexFilter{}).
		Return([]models.Backing{{ID: uuid.New(), Name: "BACKING"}}, 1, nil).Once()

	err := handler.GetIndex(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestBackingHandler_GetByIDSuccess(t *testing.T) {
	e := newEcho()
	backingID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/v1/backing/"+backingID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(backingID)

	mockUC := new(mockBackingUsecase)
	handler := &backingHttp.BackingHandler{Usecase: mockUC}

	mockUC.On("GetByID", mock.Anything, backingID).
		Return(&models.Backing{ID: uuid.MustParse(backingID), Name: "BACKING"}, nil).Once()

	err := handler.GetByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestBackingHandler_ExportSuccess(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/backing/export", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockBackingUsecase)
	handler := &backingHttp.BackingHandler{Usecase: mockUC}

	mockUC.On("Export", mock.Anything, dto.ReqBackingIndexFilter{}).
		Return([]byte("excel"), nil).Once()

	err := handler.Export(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, constants.ExcelContent, rec.Header().Get(echo.HeaderContentType))
	mockUC.AssertExpectations(t)
}

func TestBackingHandler_GetIndexBindError(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/backing", strings.NewReader("invalid"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("page_request", request.NewPageRequest(1, 10, "", "", ""))

	handler := &backingHttp.BackingHandler{Usecase: new(mockBackingUsecase)}

	err := handler.GetIndex(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestBackingHandler_ExportValidateError(t *testing.T) {
	e := newEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/backing/export", strings.NewReader(`{"type_ids":["not-uuid"]}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &backingHttp.BackingHandler{Usecase: new(mockBackingUsecase)}

	err := handler.Export(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.NotEmpty(t, resp.Message)
}
