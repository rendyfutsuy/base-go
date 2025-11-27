package test

import (
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
	typeHttp "github.com/rendyfutsuy/base-go/modules/type/delivery/http"
	"github.com/rendyfutsuy/base-go/modules/type/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTypeUsecase struct {
	mock.Mock
}

func (m *mockTypeUsecase) Create(c echo.Context, req *dto.ReqCreateType, authId string) (*models.Type, error) {
	args := m.Called(c, req, authId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Type), args.Error(1)
}

func (m *mockTypeUsecase) Update(c echo.Context, id string, req *dto.ReqUpdateType, authId string) (*models.Type, error) {
	args := m.Called(c, id, req, authId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Type), args.Error(1)
}

func (m *mockTypeUsecase) Delete(c echo.Context, id string, authId string) error {
	args := m.Called(c, id, authId)
	return args.Error(0)
}

func (m *mockTypeUsecase) GetByID(c echo.Context, id string) (*models.Type, error) {
	args := m.Called(c, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Type), args.Error(1)
}

func (m *mockTypeUsecase) GetIndex(c echo.Context, req request.PageRequest, filter dto.ReqTypeIndexFilter) ([]models.Type, int, error) {
	args := m.Called(c, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Type), args.Int(1), args.Error(2)
}

func (m *mockTypeUsecase) GetAll(c echo.Context, filter dto.ReqTypeIndexFilter) ([]models.Type, error) {
	args := m.Called(c, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Type), args.Error(1)
}

func (m *mockTypeUsecase) Export(c echo.Context, filter dto.ReqTypeIndexFilter) ([]byte, error) {
	args := m.Called(c, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

func newTypeEcho() *echo.Echo {
	e := echo.New()
	v := validator.New()
	utils.RegisterCustomValidator(v)
	e.Validator = &utils.CustomValidator{Validator: v}
	return e
}

func TestTypeHandler_CreateSuccess(t *testing.T) {
	e := newTypeEcho()
	subGroupID := uuid.New()
	body := fmt.Sprintf(`{"subgroup_id":"%s","name":"TYPE-A"}`, subGroupID.String())
	req := httptest.NewRequest(http.MethodPost, "/v1/type", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", models.User{ID: uuid.New()})

	mockUC := new(mockTypeUsecase)
	handler := &typeHttp.TypeHandler{Usecase: mockUC}

	created := &models.Type{ID: uuid.New(), Name: "TYPE-A"}
	detailed := &models.Type{ID: created.ID, Name: created.Name, SubgroupName: "SUB-GROUP"}
	mockUC.On("Create", mock.Anything, mock.AnythingOfType("*dto.ReqCreateType"), "").
		Return(created, nil).Once()
	mockUC.On("GetByID", mock.Anything, created.ID.String()).
		Return(detailed, nil).Once()

	err := handler.Create(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestTypeHandler_CreateValidationError(t *testing.T) {
	e := newTypeEcho()
	req := httptest.NewRequest(http.MethodPost, "/v1/type", strings.NewReader(`{"name":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", models.User{ID: uuid.New()})

	handler := &typeHttp.TypeHandler{Usecase: new(mockTypeUsecase)}

	err := handler.Create(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTypeHandler_UpdateSuccess(t *testing.T) {
	e := newTypeEcho()
	typeID := uuid.New().String()
	req := httptest.NewRequest(http.MethodPut, "/v1/type/"+typeID, strings.NewReader(fmt.Sprintf(`{"subgroup_id":"%s","name":"UPDATED"}`, uuid.New().String())))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(typeID)

	mockUC := new(mockTypeUsecase)
	handler := &typeHttp.TypeHandler{Usecase: mockUC}

	mockUC.On("Update", mock.Anything, typeID, mock.AnythingOfType("*dto.ReqUpdateType"), "").
		Return(&models.Type{ID: uuid.MustParse(typeID), Name: "UPDATED"}, nil).Once()
	mockUC.On("GetByID", mock.Anything, typeID).
		Return(&models.Type{ID: uuid.MustParse(typeID), Name: "UPDATED", SubgroupName: "SUB"}, nil).Once()

	err := handler.Update(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestTypeHandler_DeleteSuccess(t *testing.T) {
	e := newTypeEcho()
	typeID := uuid.New().String()
	req := httptest.NewRequest(http.MethodDelete, "/v1/type/"+typeID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(typeID)

	mockUC := new(mockTypeUsecase)
	handler := &typeHttp.TypeHandler{Usecase: mockUC}

	mockUC.On("Delete", mock.Anything, typeID, "").Return(nil).Once()

	err := handler.Delete(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestTypeHandler_GetIndexSuccess(t *testing.T) {
	e := newTypeEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/type?page=1&per_page=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	pageReq := request.NewPageRequest(1, 5, "", "", "")
	c.Set("page_request", pageReq)

	mockUC := new(mockTypeUsecase)
	handler := &typeHttp.TypeHandler{Usecase: mockUC}

	mockUC.On("GetIndex", mock.Anything, *pageReq, dto.ReqTypeIndexFilter{}).
		Return([]models.Type{{ID: uuid.New(), Name: "TYPE-A"}}, 1, nil).Once()

	err := handler.GetIndex(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestTypeHandler_GetByIDSuccess(t *testing.T) {
	e := newTypeEcho()
	typeID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/v1/type/"+typeID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(typeID)

	mockUC := new(mockTypeUsecase)
	handler := &typeHttp.TypeHandler{Usecase: mockUC}

	mockUC.On("GetByID", mock.Anything, typeID).
		Return(&models.Type{ID: uuid.MustParse(typeID), Name: "TYPE-A"}, nil).Once()

	err := handler.GetByID(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestTypeHandler_ExportSuccess(t *testing.T) {
	e := newTypeEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/type/export", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUC := new(mockTypeUsecase)
	handler := &typeHttp.TypeHandler{Usecase: mockUC}

	mockUC.On("Export", mock.Anything, dto.ReqTypeIndexFilter{}).
		Return([]byte("excel"), nil).Once()

	err := handler.Export(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, constants.ExcelContent, rec.Header().Get(echo.HeaderContentType))
	mockUC.AssertExpectations(t)
}

func TestTypeHandler_GetIndexBindError(t *testing.T) {
	e := newTypeEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/type", strings.NewReader("invalid"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("page_request", request.NewPageRequest(1, 5, "", "", ""))

	handler := &typeHttp.TypeHandler{Usecase: new(mockTypeUsecase)}

	err := handler.GetIndex(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTypeHandler_ExportValidateError(t *testing.T) {
	e := newTypeEcho()
	req := httptest.NewRequest(http.MethodGet, "/v1/type/export", strings.NewReader(`invalid`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &typeHttp.TypeHandler{Usecase: new(mockTypeUsecase)}

	err := handler.Export(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var resp struct {
		Message string `json:"message"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.NotEmpty(t, resp.Message)
}
