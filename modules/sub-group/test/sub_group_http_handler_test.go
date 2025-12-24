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
	subgroupHttp "github.com/rendyfutsuy/base-go/modules/sub-group/delivery/http"
	"github.com/rendyfutsuy/base-go/modules/sub-group/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockSubGroupUsecase struct {
	mock.Mock
}

func (m *mockSubGroupUsecase) Create(ctx context.Context, req *dto.ReqCreateSubGroup, userID string) (*models.SubGroup, error) {
	args := m.Called(ctx, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubGroup), args.Error(1)
}

func (m *mockSubGroupUsecase) Update(ctx context.Context, id string, req *dto.ReqUpdateSubGroup, userID string) (*models.SubGroup, error) {
	args := m.Called(ctx, id, req, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubGroup), args.Error(1)
}

func (m *mockSubGroupUsecase) Delete(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *mockSubGroupUsecase) GetByID(ctx context.Context, id string) (*models.SubGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubGroup), args.Error(1)
}

func (m *mockSubGroupUsecase) GetIndex(ctx context.Context, req request.PageRequest, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.SubGroup), args.Int(1), args.Error(2)
}

func (m *mockSubGroupUsecase) GetAll(ctx context.Context, filter dto.ReqSubGroupIndexFilter) ([]models.SubGroup, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SubGroup), args.Error(1)
}

func (m *mockSubGroupUsecase) Export(ctx context.Context, filter dto.ReqSubGroupIndexFilter) ([]byte, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockSubGroupUsecase) ExistsInTypes(ctx context.Context, subGroupID string) (bool, error) {
	args := m.Called(ctx, subGroupID)
	return args.Bool(0), args.Error(1)
}

func newSubGroupEcho() *echo.Echo {
	e := echo.New()
	v := validator.New()
	utils.RegisterCustomValidator(v)
	e.Validator = &utils.CustomValidator{Validator: v}
	return e
}

func TestSubGroupHandler_CreateSuccess(t *testing.T) {
	e := newSubGroupEcho()
	body := fmt.Sprintf(`{"groups_id":"%s","name":"SUBGROUP"}`, uuid.New().String())
	req := httptest.NewRequest(http.MethodPost, "/v1/sub-group", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	userID := uuid.New()
	c.Set("user", models.User{ID: userID})

	mockUC := new(mockSubGroupUsecase)
	handler := &subgroupHttp.SubGroupHandler{Usecase: mockUC}

	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*dto.ReqCreateSubGroup"), userID.String()).
		Return(&models.SubGroup{ID: uuid.New(), Name: "SUBGROUP"}, nil).Once()

	err := handler.Create(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestSubGroupHandler_CreateValidationError(t *testing.T) {
	e := newSubGroupEcho()
	req := httptest.NewRequest(http.MethodPost, "/v1/sub-group", strings.NewReader(`{"name":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", models.User{ID: uuid.New()})

	handler := &subgroupHttp.SubGroupHandler{Usecase: new(mockSubGroupUsecase)}

	err := handler.Create(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSubGroupHandler_UpdateSuccess(t *testing.T) {
	e := newSubGroupEcho()
	subGroupID := uuid.New().String()
	body := fmt.Sprintf(`{"groups_id":"%s","name":"UPDATED"}`, uuid.New().String())
	req := httptest.NewRequest(http.MethodPut, "/v1/sub-group/"+subGroupID, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(subGroupID)

	mockUC := new(mockSubGroupUsecase)
	handler := &subgroupHttp.SubGroupHandler{Usecase: mockUC}

	mockUC.On("Update", mock.Anything, subGroupID, mock.AnythingOfType("*dto.ReqUpdateSubGroup"), "").
		Return(&models.SubGroup{ID: uuid.MustParse(subGroupID), Name: "UPDATED"}, nil).Once()

	err := handler.Update(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestSubGroupHandler_DeleteSuccess(t *testing.T) {
	e := newSubGroupEcho()
	subGroupID := uuid.New().String()
	req := httptest.NewRequest(http.MethodDelete, "/v1/sub-group/"+subGroupID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(subGroupID)

	mockUC := new(mockSubGroupUsecase)
	handler := &subgroupHttp.SubGroupHandler{Usecase: mockUC}

	mockUC.On("Delete", mock.Anything, subGroupID, "").Return(nil).Once()

	err := handler.Delete(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestSubGroupHandler_GetIndexSuccess(t *testing.T) {
	e := newSubGroupEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/sub-group?page=1&per_page=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	pageReq := request.NewPageRequest(1, 5, "", "", "")
	c.Set("page_request", pageReq)

	mockUC := new(mockSubGroupUsecase)
	handler := &subgroupHttp.SubGroupHandler{Usecase: mockUC}

	mockUC.On("GetIndex", mock.Anything, *pageReq, dto.ReqSubGroupIndexFilter{}).
		Return([]models.SubGroup{{ID: uuid.New(), Name: "SUBGROUP"}}, 1, nil).Once()

	err := handler.GetIndex(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestSubGroupHandler_GetByIDSuccess(t *testing.T) {
	e := newSubGroupEcho()
	subGroupID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/v1/sub-group/"+subGroupID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(subGroupID)

	mockUC := new(mockSubGroupUsecase)
	handler := &subgroupHttp.SubGroupHandler{Usecase: mockUC}

	mockUC.On("GetByID", mock.Anything, subGroupID).
		Return(&models.SubGroup{ID: uuid.MustParse(subGroupID), Name: "SUBGROUP"}, nil).Once()

	err := handler.GetByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestSubGroupHandler_ExportSuccess(t *testing.T) {
	e := newSubGroupEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/sub-group/export", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockSubGroupUsecase)
	handler := &subgroupHttp.SubGroupHandler{Usecase: mockUC}

	mockUC.On("Export", mock.Anything, dto.ReqSubGroupIndexFilter{}).
		Return([]byte("excel"), nil).Once()

	err := handler.Export(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, constants.ExcelContent, rec.Header().Get(echo.HeaderContentType))
	mockUC.AssertExpectations(t)
}

func TestSubGroupHandler_GetIndexBindError(t *testing.T) {
	e := newSubGroupEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/sub-group", strings.NewReader("invalid"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("page_request", request.NewPageRequest(1, 5, "", "", ""))

	handler := &subgroupHttp.SubGroupHandler{Usecase: new(mockSubGroupUsecase)}

	err := handler.GetIndex(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSubGroupHandler_ExportValidateError(t *testing.T) {
	e := newSubGroupEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/sub-group/export", strings.NewReader(`{"groups_ids":["invalid"]}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &subgroupHttp.SubGroupHandler{Usecase: new(mockSubGroupUsecase)}

	err := handler.Export(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.NotEmpty(t, resp.Message)
}
