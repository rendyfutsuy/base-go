package test

import (
	"context"
	"encoding/json"
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
	groupHttp "github.com/rendyfutsuy/base-go/modules/group/delivery/http"
	"github.com/rendyfutsuy/base-go/modules/group/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockGroupUsecase struct {
	mock.Mock
}

func (m *mockGroupUsecase) Create(ctx context.Context, req *dto.ReqCreateGroup, userID string) (*models.Group, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Group), args.Error(1)
}

func (m *mockGroupUsecase) Update(ctx context.Context, id string, req *dto.ReqUpdateGroup, userID string) (*models.Group, error) {
	args := m.Called(ctx, id, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Group), args.Error(1)
}

func (m *mockGroupUsecase) Delete(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *mockGroupUsecase) GetByID(ctx context.Context, id string) (*models.Group, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Group), args.Error(1)
}

func (m *mockGroupUsecase) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqGroupIndexFilter) ([]models.Group, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Group), args.Int(1), args.Error(2)
}

func (m *mockGroupUsecase) GetAll(ctx context.Context, filter dto.ReqGroupIndexFilter) ([]models.Group, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Group), args.Error(1)
}

func (m *mockGroupUsecase) Export(ctx context.Context, filter dto.ReqGroupIndexFilter) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockGroupUsecase) ExistsInSubGroups(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func newGroupEcho() *echo.Echo {
	e := echo.New()
	v := validator.New()
	utils.RegisterCustomValidator(v)
	e.Validator = &utils.CustomValidator{Validator: v}
	return e
}

func TestGroupHandler_CreateSuccess(t *testing.T) {
	e := newGroupEcho()
	req := httptest.NewRequest(http.MethodPost, "/v1/group", strings.NewReader(`{"name":"GROUP-A"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	userID := uuid.New()
	c.Set("user", models.User{ID: userID})

	mockUC := new(mockGroupUsecase)
	handler := &groupHttp.GroupHandler{Usecase: mockUC}

	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*dto.ReqCreateGroup"), userID.String()).
		Return(&models.Group{ID: uuid.New(), Name: "GROUP-A"}, nil).Once()

	err := handler.Create(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestGroupHandler_CreateValidationError(t *testing.T) {
	e := newGroupEcho()
	req := httptest.NewRequest(http.MethodPost, "/v1/group", strings.NewReader(`{"name":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", models.User{ID: uuid.New()})

	handler := &groupHttp.GroupHandler{Usecase: new(mockGroupUsecase)}

	err := handler.Create(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGroupHandler_UpdateSuccess(t *testing.T) {
	e := newGroupEcho()
	groupID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPut, "/v1/group/"+groupID, strings.NewReader(`{"name":"GROUP-B"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(groupID)

	mockUC := new(mockGroupUsecase)
	handler := &groupHttp.GroupHandler{Usecase: mockUC}

	mockUC.On("Update", mock.Anything, groupID, mock.AnythingOfType("*dto.ReqUpdateGroup"), "").
		Return(&models.Group{ID: uuid.MustParse(groupID), Name: "GROUP-B", Deletable: true}, nil).Once()

	err := handler.Update(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestGroupHandler_DeleteSuccess(t *testing.T) {
	e := newGroupEcho()
	groupID := uuid.New().String()
	req := httptest.NewRequest(http.MethodDelete, "/v1/group/"+groupID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(groupID)

	mockUC := new(mockGroupUsecase)
	handler := &groupHttp.GroupHandler{Usecase: mockUC}

	mockUC.On("Delete", mock.Anything, groupID, "").Return(nil).Once()

	err := handler.Delete(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestGroupHandler_GetIndexSuccess(t *testing.T) {
	e := newGroupEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/group?page=1&per_page=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	pageReq := request.NewPageRequest(1, 5, "", "", "")
	c.Set("page_request", pageReq)

	mockUC := new(mockGroupUsecase)
	handler := &groupHttp.GroupHandler{Usecase: mockUC}

	mockUC.On("GetIndex", mock.Anything, *pageReq, dto.ReqGroupIndexFilter{}).
		Return([]models.Group{{ID: uuid.New(), Name: "GROUP-A"}}, 1, nil).Once()

	err := handler.GetIndex(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestGroupHandler_GetByIDSuccess(t *testing.T) {
	e := newGroupEcho()
	groupID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/v1/group/"+groupID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(groupID)

	mockUC := new(mockGroupUsecase)
	handler := &groupHttp.GroupHandler{Usecase: mockUC}

	mockUC.On("GetByID", mock.Anything, groupID).
		Return(&models.Group{ID: uuid.MustParse(groupID), Name: "GROUP-A"}, nil).Once()

	err := handler.GetByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestGroupHandler_ExportSuccess(t *testing.T) {
	e := newGroupEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/group/export", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockGroupUsecase)
	handler := &groupHttp.GroupHandler{Usecase: mockUC}

	mockUC.On("Export", mock.Anything, dto.ReqGroupIndexFilter{}).
		Return([]byte("excel"), nil).Once()

	err := handler.Export(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, constants.ExcelContent, rec.Header().Get(echo.HeaderContentType))
	mockUC.AssertExpectations(t)
}

func TestGroupHandler_GetIndexBindError(t *testing.T) {
	e := newGroupEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/group", strings.NewReader("invalid"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("page_request", request.NewPageRequest(1, 5, "", "", ""))

	handler := &groupHttp.GroupHandler{Usecase: new(mockGroupUsecase)}

	err := handler.GetIndex(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGroupHandler_ExportValidateError(t *testing.T) {
	e := newGroupEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/group/export", strings.NewReader(`invalid`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &groupHttp.GroupHandler{Usecase: new(mockGroupUsecase)}

	err := handler.Export(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.NotEmpty(t, resp.Message)
}
