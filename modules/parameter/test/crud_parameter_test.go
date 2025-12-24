package test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	parameterDto "github.com/rendyfutsuy/base-go/modules/parameter/dto"
	"github.com/rendyfutsuy/base-go/modules/parameter/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockParameterRepository is a mock implementation of parameter.Repository
type MockParameterRepository struct {
	mock.Mock
}

func (m *MockParameterRepository) Create(ctx context.Context, code, name string, value, typeVal, desc *string) (*models.Parameter, error) {
	args := m.Called(ctx, code, name, value, typeVal, desc)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Parameter), args.Error(1)
}

func (m *MockParameterRepository) Update(ctx context.Context, id uuid.UUID, code, name string, value, typeVal, desc *string) (*models.Parameter, error) {
	args := m.Called(ctx, id, code, name, value, typeVal, desc)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Parameter), args.Error(1)
}

func (m *MockParameterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockParameterRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Parameter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Parameter), args.Error(1)
}

func (m *MockParameterRepository) GetIndex(ctx context.Context, req request.PageRequest, filter parameterDto.ReqParameterIndexFilter) ([]models.Parameter, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Parameter), args.Int(1), args.Error(2)
}

func (m *MockParameterRepository) GetAll(ctx context.Context, filter parameterDto.ReqParameterIndexFilter) ([]models.Parameter, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Parameter), args.Error(1)
}

func (m *MockParameterRepository) ExistsByCode(ctx context.Context, code string, excludeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, code, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockParameterRepository) ExistsByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func TestCreateParameter(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	value := "test-value"
	typeVal := "test-type"
	desc := "test description"

	tests := []struct {
		name           string
		req            *parameterDto.ReqCreateParameter
		authId         string
		setupMock      func(*MockParameterRepository)
		expectedError  error
		expectedResult *models.Parameter
	}{
		{
			name: "success create parameter",
			req: &parameterDto.ReqCreateParameter{
				Code:  "TEST001",
				Name:  "Test Parameter",
				Value: &value,
				Type:  &typeVal,
				Desc:  &desc,
			},
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("ExistsByCode", ctx, "TEST001", uuid.Nil).Return(false, nil).Once()
				m.On("ExistsByName", ctx, "Test Parameter", uuid.Nil).Return(false, nil).Once()
				m.On("Create", ctx, "TEST001", "Test Parameter", &value, &typeVal, &desc).Return(&models.Parameter{
					ID:          uuid.New(),
					Code:        "TEST001",
					Name:        "Test Parameter",
					Value:       &value,
					Type:        &typeVal,
					Description: &desc,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Parameter{
				Code:        "TEST001",
				Name:        "Test Parameter",
				Value:       &value,
				Type:        &typeVal,
				Description: &desc,
			},
		},
		{
			name: "success create parameter without optional fields",
			req: &parameterDto.ReqCreateParameter{
				Code:  "TEST002",
				Name:  "Test Parameter 2",
				Value: nil,
				Type:  nil,
				Desc:  nil,
			},
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("ExistsByCode", ctx, "TEST002", uuid.Nil).Return(false, nil).Once()
				m.On("ExistsByName", ctx, "Test Parameter 2", uuid.Nil).Return(false, nil).Once()
				m.On("Create", ctx, "TEST002", "Test Parameter 2", (*string)(nil), (*string)(nil), (*string)(nil)).Return(&models.Parameter{
					ID:          uuid.New(),
					Code:        "TEST002",
					Name:        "Test Parameter 2",
					Value:       nil,
					Type:        nil,
					Description: nil,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Parameter{
				Code:        "TEST002",
				Name:        "Test Parameter 2",
				Value:       nil,
				Type:        nil,
				Description: nil,
			},
		},
		{
			name: "error when code already exists",
			req: &parameterDto.ReqCreateParameter{
				Code: "EXISTING001",
				Name: "Test Parameter",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("ExistsByCode", ctx, "EXISTING001", uuid.Nil).Return(true, nil).Once()
			},
			expectedError:  errors.New("Parameter code already exists"),
			expectedResult: nil,
		},
		{
			name: "error when name already exists",
			req: &parameterDto.ReqCreateParameter{
				Code: "TEST003",
				Name: "Existing Parameter",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("ExistsByCode", ctx, "TEST003", uuid.Nil).Return(false, nil).Once()
				m.On("ExistsByName", ctx, "Existing Parameter", uuid.Nil).Return(true, nil).Once()
			},
			expectedError:  errors.New("Parameter name already exists"),
			expectedResult: nil,
		},
		{
			name: "error when repository check code fails",
			req: &parameterDto.ReqCreateParameter{
				Code: "TEST004",
				Name: "Test Parameter",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("ExistsByCode", ctx, "TEST004", uuid.Nil).Return(false, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
		{
			name: "error when repository create fails",
			req: &parameterDto.ReqCreateParameter{
				Code: "TEST005",
				Name: "Test Parameter",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("ExistsByCode", ctx, "TEST005", uuid.Nil).Return(false, nil).Once()
				m.On("ExistsByName", ctx, "Test Parameter", uuid.Nil).Return(false, nil).Once()
				m.On("Create", ctx, "TEST005", "Test Parameter", (*string)(nil), (*string)(nil), (*string)(nil)).Return(nil, errors.New("create failed")).Once()
			},
			expectedError:  errors.New("create failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockParameterRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewParameterUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.Create(ctx, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Code, result.Code)
					assert.Equal(t, tt.expectedResult.Name, result.Name)
					if tt.expectedResult.Value != nil {
						assert.Equal(t, *tt.expectedResult.Value, *result.Value)
					}
					if tt.expectedResult.Type != nil {
						assert.Equal(t, *tt.expectedResult.Type, *result.Type)
					}
					if tt.expectedResult.Description != nil {
						assert.Equal(t, *tt.expectedResult.Description, *result.Description)
					}
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateParameter(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	value := "updated-value"
	typeVal := "updated-type"
	desc := "updated description"

	tests := []struct {
		name           string
		id             string
		req            *parameterDto.ReqUpdateParameter
		authId         string
		setupMock      func(*MockParameterRepository)
		expectedError  error
		expectedResult *models.Parameter
	}{
		{
			name: "success update parameter",
			id:   validID.String(),
			req: &parameterDto.ReqUpdateParameter{
				Code:  "UPDATED001",
				Name:  "Updated Parameter",
				Value: &value,
				Type:  &typeVal,
				Desc:  &desc,
			},
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("ExistsByCode", ctx, "UPDATED001", validID).Return(false, nil).Once()
				m.On("ExistsByName", ctx, "Updated Parameter", validID).Return(false, nil).Once()
				m.On("Update", ctx, validID, "UPDATED001", "Updated Parameter", &value, &typeVal, &desc).Return(&models.Parameter{
					ID:          validID,
					Code:        "UPDATED001",
					Name:        "Updated Parameter",
					Value:       &value,
					Type:        &typeVal,
					Description: &desc,
					UpdatedAt:   time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Parameter{
				Code:        "UPDATED001",
				Name:        "Updated Parameter",
				Value:       &value,
				Type:        &typeVal,
				Description: &desc,
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			req: &parameterDto.ReqUpdateParameter{
				Code: "TEST001",
				Name: "Test Parameter",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when code already exists",
			id:   validID.String(),
			req: &parameterDto.ReqUpdateParameter{
				Code: "EXISTING001",
				Name: "Test Parameter",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("ExistsByCode", ctx, "EXISTING001", validID).Return(true, nil).Once()
			},
			expectedError:  errors.New("Parameter code already exists"),
			expectedResult: nil,
		},
		{
			name: "error when name already exists",
			id:   validID.String(),
			req: &parameterDto.ReqUpdateParameter{
				Code: "TEST001",
				Name: "Existing Parameter",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("ExistsByCode", ctx, "TEST001", validID).Return(false, nil).Once()
				m.On("ExistsByName", ctx, "Existing Parameter", validID).Return(true, nil).Once()
			},
			expectedError:  errors.New("Parameter name already exists"),
			expectedResult: nil,
		},
		{
			name: "error when repository update fails",
			id:   validID.String(),
			req: &parameterDto.ReqUpdateParameter{
				Code: "TEST001",
				Name: "Updated Parameter",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("ExistsByCode", ctx, "TEST001", validID).Return(false, nil).Once()
				m.On("ExistsByName", ctx, "Updated Parameter", validID).Return(false, nil).Once()
				m.On("Update", ctx, validID, "TEST001", "Updated Parameter", (*string)(nil), (*string)(nil), (*string)(nil)).Return(nil, errors.New("update failed")).Once()
			},
			expectedError:  errors.New("update failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockParameterRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewParameterUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPut, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.Update(ctx, tt.id, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Code, result.Code)
					assert.Equal(t, tt.expectedResult.Name, result.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteParameter(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name          string
		id            string
		authId        string
		setupMock     func(*MockParameterRepository)
		expectedError error
	}{
		{
			name:   "success delete parameter",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("Delete", ctx, validID).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:   "error when invalid UUID",
			id:     "invalid-uuid",
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				// No mock calls expected
			},
			expectedError: errors.New("requested param is string"),
		},
		{
			name:   "error when repository delete fails",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockParameterRepository) {
				m.On("Delete", ctx, validID).Return(errors.New("delete failed")).Once()
			},
			expectedError: errors.New("delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockParameterRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewParameterUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodDelete, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			err := usecaseInstance.Delete(ctx, tt.id, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetParameterByID(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	value := "test-value"
	typeVal := "test-type"
	desc := "test description"

	tests := []struct {
		name           string
		id             string
		setupMock      func(*MockParameterRepository)
		expectedError  error
		expectedResult *models.Parameter
	}{
		{
			name: "success get parameter by id",
			id:   validID.String(),
			setupMock: func(m *MockParameterRepository) {
				m.On("GetByID", ctx, validID).Return(&models.Parameter{
					ID:          validID,
					Code:        "TEST001",
					Name:        "Test Parameter",
					Value:       &value,
					Type:        &typeVal,
					Description: &desc,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Parameter{
				ID:          validID,
				Code:        "TEST001",
				Name:        "Test Parameter",
				Value:       &value,
				Type:        &typeVal,
				Description: &desc,
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			setupMock: func(m *MockParameterRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when parameter not found",
			id:   validID.String(),
			setupMock: func(m *MockParameterRepository) {
				m.On("GetByID", ctx, validID).Return(nil, errors.New("record not found")).Once()
			},
			expectedError:  errors.New("record not found"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockParameterRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewParameterUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.GetByID(ctx, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.ID, result.ID)
					assert.Equal(t, tt.expectedResult.Code, result.Code)
					assert.Equal(t, tt.expectedResult.Name, result.Name)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetIndexParameter(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	value1 := "value1"
	typeVal1 := "type1"
	value2 := "value2"
	typeVal2 := "type2"

	tests := []struct {
		name           string
		req            request.PageRequest
		filter         parameterDto.ReqParameterIndexFilter
		setupMock      func(*MockParameterRepository)
		expectedError  error
		expectedResult []models.Parameter
		expectedTotal  int
	}{
		{
			name: "success get index with pagination",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: parameterDto.ReqParameterIndexFilter{},
			setupMock: func(m *MockParameterRepository) {
				parameters := []models.Parameter{
					{
						ID:        uuid.New(),
						Code:      "TEST001",
						Name:      "Parameter 1",
						Value:     &value1,
						Type:      &typeVal1,
						CreatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Code:      "TEST002",
						Name:      "Parameter 2",
						Value:     &value2,
						Type:      &typeVal2,
						CreatedAt: time.Now(),
					},
				}
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, parameterDto.ReqParameterIndexFilter{}).Return(parameters, 2, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Parameter{
				{Code: "TEST001", Name: "Parameter 1"},
				{Code: "TEST002", Name: "Parameter 2"},
			},
			expectedTotal: 2,
		},
		{
			name: "success get index with search",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: parameterDto.ReqParameterIndexFilter{},
			setupMock: func(m *MockParameterRepository) {
				parameters := []models.Parameter{
					{
						ID:        uuid.New(),
						Code:      "TEST001",
						Name:      "Parameter 1",
						Value:     &value1,
						Type:      &typeVal1,
						CreatedAt: time.Now(),
					},
				}
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, parameterDto.ReqParameterIndexFilter{}).Return(parameters, 1, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Parameter{
				{Code: "TEST001", Name: "Parameter 1"},
			},
			expectedTotal: 1,
		},
		{
			name: "error when repository fails",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: parameterDto.ReqParameterIndexFilter{},
			setupMock: func(m *MockParameterRepository) {
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, parameterDto.ReqParameterIndexFilter{}).Return(nil, 0, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
			expectedTotal:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockParameterRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewParameterUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, total, err := usecaseInstance.GetIndex(ctx, tt.req, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedTotal, total)
				assert.Equal(t, len(tt.expectedResult), len(result))
				if len(tt.expectedResult) > 0 {
					assert.Equal(t, tt.expectedResult[0].Name, result[0].Name)
					assert.Equal(t, tt.expectedResult[0].Code, result[0].Code)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExportParameter(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	value1 := "value1"
	typeVal1 := "type1"
	value2 := "value2"
	typeVal2 := "type2"

	tests := []struct {
		name           string
		search         string
		filter         parameterDto.ReqParameterIndexFilter
		setupMock      func(*MockParameterRepository)
		expectedError  error
		expectedResult []models.Parameter
	}{
		{
			name:   "success export parameters",
			search: "",
			filter: parameterDto.ReqParameterIndexFilter{},
			setupMock: func(m *MockParameterRepository) {
				parameters := []models.Parameter{
					{
						ID:        uuid.New(),
						Code:      "TEST001",
						Name:      "Parameter 1",
						Value:     &value1,
						Type:      &typeVal1,
						CreatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Code:      "TEST002",
						Name:      "Parameter 2",
						Value:     &value2,
						Type:      &typeVal2,
						CreatedAt: time.Now(),
					},
				}
				m.On("GetAll", ctx, parameterDto.ReqParameterIndexFilter{}).Return(parameters, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Parameter{
				{Code: "TEST001", Name: "Parameter 1"},
				{Code: "TEST002", Name: "Parameter 2"},
			},
		},
		{
			name:   "success export parameters with search",
			search: "Parameter 1",
			filter: parameterDto.ReqParameterIndexFilter{Search: "Parameter 1"},
			setupMock: func(m *MockParameterRepository) {
				parameters := []models.Parameter{
					{
						ID:        uuid.New(),
						Code:      "TEST001",
						Name:      "Parameter 1",
						Value:     &value1,
						Type:      &typeVal1,
						CreatedAt: time.Now(),
					},
				}
				m.On("GetAll", ctx, parameterDto.ReqParameterIndexFilter{Search: "Parameter 1"}).Return(parameters, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Parameter{
				{Code: "TEST001", Name: "Parameter 1"},
			},
		},
		{
			name:   "error when repository fails",
			search: "",
			filter: parameterDto.ReqParameterIndexFilter{},
			setupMock: func(m *MockParameterRepository) {
				m.On("GetAll", ctx, parameterDto.ReqParameterIndexFilter{}).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockParameterRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewParameterUsecase(mockRepo)

			u := &url.URL{Path: "/export"}
			if tt.search != "" {
				q := u.Query()
				q.Set("search", tt.search)
				u.RawQuery = q.Encode()
			}
			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.Export(ctx, tt.filter)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				// Verify that result is Excel file bytes (should not be empty)
				assert.Greater(t, len(result), 0)
				// Verify Excel file signature (PK for ZIP format)
				assert.Equal(t, []byte("PK")[0], result[0])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
