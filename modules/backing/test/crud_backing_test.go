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
	backingDto "github.com/rendyfutsuy/base-go/modules/backing/dto"
	"github.com/rendyfutsuy/base-go/modules/backing/usecase"
	typeDto "github.com/rendyfutsuy/base-go/modules/type/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBackingRepository is a mock implementation of backing.Repository
type MockBackingRepository struct {
	mock.Mock
}

func (m *MockBackingRepository) Create(ctx context.Context, typeID uuid.UUID, name string, createdBy string) (*models.Backing, error) {
	args := m.Called(ctx, typeID, name, createdBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Backing), args.Error(1)
}

func (m *MockBackingRepository) Update(ctx context.Context, id uuid.UUID, typeID uuid.UUID, name string, updatedBy string) (*models.Backing, error) {
	args := m.Called(ctx, id, typeID, name, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Backing), args.Error(1)
}

func (m *MockBackingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBackingRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Backing, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Backing), args.Error(1)
}

func (m *MockBackingRepository) GetIndex(ctx context.Context, req request.PageRequest, filter backingDto.ReqBackingIndexFilter) ([]models.Backing, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Backing), args.Int(1), args.Error(2)
}

func (m *MockBackingRepository) GetAll(ctx context.Context, filter backingDto.ReqBackingIndexFilter) ([]models.Backing, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Backing), args.Error(1)
}

func (m *MockBackingRepository) ExistsByNameInType(ctx context.Context, typeID uuid.UUID, name string, excludeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, typeID, name, excludeID)
	return args.Bool(0), args.Error(1)
}

// MockTypeRepository is a mock implementation of type.Repository
type MockTypeRepository struct {
	mock.Mock
}

func (m *MockTypeRepository) Create(ctx context.Context, subgroupID uuid.UUID, name string, createdBy string) (*models.Type, error) {
	args := m.Called(ctx, subgroupID, name, createdBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Type), args.Error(1)
}

func (m *MockTypeRepository) Update(ctx context.Context, id uuid.UUID, subgroupID uuid.UUID, name string, updatedBy string) (*models.Type, error) {
	args := m.Called(ctx, id, subgroupID, name, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Type), args.Error(1)
}

func (m *MockTypeRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *MockTypeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Type, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Type), args.Error(1)
}

func (m *MockTypeRepository) GetIndex(ctx context.Context, req request.PageRequest, filter typeDto.ReqTypeIndexFilter) ([]models.Type, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Type), args.Int(1), args.Error(2)
}

func (m *MockTypeRepository) GetAll(ctx context.Context, filter typeDto.ReqTypeIndexFilter) ([]models.Type, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Type), args.Error(1)
}

func (m *MockTypeRepository) ExistsByNameInSubgroup(ctx context.Context, subgroupID uuid.UUID, name string, excludeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, subgroupID, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTypeRepository) ExistsInBackings(ctx context.Context, typeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, typeID)
	return args.Bool(0), args.Error(1)
}

func TestCreateBacking(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	typeID := uuid.New()
	testUserID := uuid.New()

	tests := []struct {
		name           string
		req            *backingDto.ReqCreateBacking
		authId         string
		setupMock      func(*MockBackingRepository, *MockTypeRepository)
		expectedError  error
		expectedResult *models.Backing
	}{
		{
			name: "success create backing",
			req: &backingDto.ReqCreateBacking{
				TypeID: typeID,
				Name:   "Test Backing",
			},
			authId: testUserID.String(),
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				mt.On("GetByID", ctx, typeID).Return(&models.Type{
					ID:   typeID,
					Name: "Test Type",
				}, nil).Once()
				m.On("ExistsByNameInType", ctx, typeID, "Test Backing", uuid.Nil).Return(false, nil).Once()
				m.On("Create", ctx, typeID, "Test Backing", testUserID.String()).Return(&models.Backing{
					ID:          uuid.New(),
					TypeID:      typeID,
					BackingCode: "01",
					Name:        "Test Backing",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Backing{
				BackingCode: "01",
				Name:        "Test Backing",
			},
		},
		{
			name: "error when name already exists in type",
			req: &backingDto.ReqCreateBacking{
				TypeID: typeID,
				Name:   "Existing Backing",
			},
			authId: testUserID.String(),
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				mt.On("GetByID", ctx, typeID).Return(&models.Type{
					ID:   typeID,
					Name: "Test Type",
				}, nil).Once()
				m.On("ExistsByNameInType", ctx, typeID, "Existing Backing", uuid.Nil).Return(true, nil).Once()
			},
			expectedError:  errors.New("Backing name already exists in this type"),
			expectedResult: nil,
		},
		{
			name: "error when repository check exists fails",
			req: &backingDto.ReqCreateBacking{
				TypeID: typeID,
				Name:   "Test Backing",
			},
			authId: testUserID.String(),
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				mt.On("GetByID", ctx, typeID).Return(&models.Type{
					ID:   typeID,
					Name: "Test Type",
				}, nil).Once()
				m.On("ExistsByNameInType", ctx, typeID, "Test Backing", uuid.Nil).Return(false, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
		{
			name: "error when repository create fails",
			req: &backingDto.ReqCreateBacking{
				TypeID: typeID,
				Name:   "Test Backing",
			},
			authId: testUserID.String(),
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				mt.On("GetByID", ctx, typeID).Return(&models.Type{
					ID:   typeID,
					Name: "Test Type",
				}, nil).Once()
				m.On("ExistsByNameInType", ctx, typeID, "Test Backing", uuid.Nil).Return(false, nil).Once()
				m.On("Create", ctx, typeID, "Test Backing", testUserID.String()).Return(nil, errors.New("create failed")).Once()
			},
			expectedError:  errors.New("create failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBackingRepository)
			mockTypeRepo := new(MockTypeRepository)
			tt.setupMock(mockRepo, mockTypeRepo)

			usecaseInstance := usecase.NewBackingUsecase(mockRepo, mockTypeRepo)

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
					assert.Equal(t, tt.expectedResult.Name, result.Name)
					assert.Equal(t, tt.expectedResult.BackingCode, result.BackingCode)
				}
			}

			mockRepo.AssertExpectations(t)
			mockTypeRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateBacking(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	typeID := uuid.New()
	testUserID := uuid.New()

	tests := []struct {
		name           string
		id             string
		req            *backingDto.ReqUpdateBacking
		authId         string
		setupMock      func(*MockBackingRepository, *MockTypeRepository)
		expectedError  error
		expectedResult *models.Backing
	}{
		{
			name: "success update backing",
			id:   validID.String(),
			req: &backingDto.ReqUpdateBacking{
				TypeID: typeID,
				Name:   "Updated Backing",
			},
			authId: testUserID.String(),
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				mt.On("GetByID", ctx, typeID).Return(&models.Type{
					ID:   typeID,
					Name: "Test Type",
				}, nil).Once()
				m.On("ExistsByNameInType", ctx, typeID, "Updated Backing", validID).Return(false, nil).Once()
				m.On("Update", ctx, validID, typeID, "Updated Backing", testUserID.String()).Return(&models.Backing{
					ID:          validID,
					TypeID:      typeID,
					BackingCode: "02",
					Name:        "Updated Backing",
					UpdatedAt:   time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Backing{
				ID:          validID,
				BackingCode: "02",
				Name:        "Updated Backing",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			req: &backingDto.ReqUpdateBacking{
				TypeID: typeID,
				Name:   "Updated Backing",
			},
			authId: testUserID.String(),
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when name already exists in type",
			id:   validID.String(),
			req: &backingDto.ReqUpdateBacking{
				TypeID: typeID,
				Name:   "Existing Backing",
			},
			authId: testUserID.String(),
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				mt.On("GetByID", ctx, typeID).Return(&models.Type{
					ID:   typeID,
					Name: "Test Type",
				}, nil).Once()
				m.On("ExistsByNameInType", ctx, typeID, "Existing Backing", validID).Return(true, nil).Once()
			},
			expectedError:  errors.New("Backing name already exists in this type"),
			expectedResult: nil,
		},
		{
			name: "error when repository update fails",
			id:   validID.String(),
			req: &backingDto.ReqUpdateBacking{
				TypeID: typeID,
				Name:   "Updated Backing",
			},
			authId: testUserID.String(),
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				mt.On("GetByID", ctx, typeID).Return(&models.Type{
					ID:   typeID,
					Name: "Test Type",
				}, nil).Once()
				m.On("ExistsByNameInType", ctx, typeID, "Updated Backing", validID).Return(false, nil).Once()
				m.On("Update", ctx, validID, typeID, "Updated Backing", testUserID.String()).Return(nil, errors.New("update failed")).Once()
			},
			expectedError:  errors.New("update failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBackingRepository)
			mockTypeRepo := new(MockTypeRepository)
			tt.setupMock(mockRepo, mockTypeRepo)

			usecaseInstance := usecase.NewBackingUsecase(mockRepo, mockTypeRepo)

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
					assert.Equal(t, tt.expectedResult.Name, result.Name)
					assert.Equal(t, tt.expectedResult.BackingCode, result.BackingCode)
				}
			}

			mockRepo.AssertExpectations(t)
			mockTypeRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteBacking(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name          string
		id            string
		authId        string
		setupMock     func(*MockBackingRepository, *MockTypeRepository)
		expectedError error
	}{
		{
			name:   "success delete backing",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				m.On("Delete", ctx, validID).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:   "error when invalid UUID",
			id:     "invalid-uuid",
			authId: "test-auth-id",
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				// No mock calls expected
			},
			expectedError: errors.New("requested param is string"),
		},
		{
			name:   "error when repository delete fails",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				m.On("Delete", ctx, validID).Return(errors.New("delete failed")).Once()
			},
			expectedError: errors.New("delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBackingRepository)
			mockTypeRepo := new(MockTypeRepository)
			tt.setupMock(mockRepo, mockTypeRepo)

			usecaseInstance := usecase.NewBackingUsecase(mockRepo, mockTypeRepo)

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
			mockTypeRepo.AssertExpectations(t)
		})
	}
}

func TestGetBackingByID(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	typeID := uuid.New()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*MockBackingRepository, *MockTypeRepository)
		expectedError  error
		expectedResult *models.Backing
	}{
		{
			name: "success get backing by id",
			id:   validID.String(),
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				m.On("GetByID", ctx, validID).Return(&models.Backing{
					ID:          validID,
					TypeID:      typeID,
					BackingCode: "01",
					Name:        "Test Backing",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Backing{
				ID:          validID,
				BackingCode: "01",
				Name:        "Test Backing",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when backing not found",
			id:   validID.String(),
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				m.On("GetByID", ctx, validID).Return(nil, errors.New("record not found")).Once()
			},
			expectedError:  errors.New("record not found"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBackingRepository)
			mockTypeRepo := new(MockTypeRepository)
			tt.setupMock(mockRepo, mockTypeRepo)

			usecaseInstance := usecase.NewBackingUsecase(mockRepo, mockTypeRepo)

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
					assert.Equal(t, tt.expectedResult.Name, result.Name)
					assert.Equal(t, tt.expectedResult.BackingCode, result.BackingCode)
				}
			}

			mockRepo.AssertExpectations(t)
			mockTypeRepo.AssertExpectations(t)
		})
	}
}

func TestGetIndexBacking(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		req            request.PageRequest
		filter         backingDto.ReqBackingIndexFilter
		setupMock      func(*MockBackingRepository, *MockTypeRepository)
		expectedError  error
		expectedResult []models.Backing
		expectedTotal  int
	}{
		{
			name: "success get index with pagination",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: backingDto.ReqBackingIndexFilter{},
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				backings := []models.Backing{
					{
						ID:          uuid.New(),
						BackingCode: "01",
						Name:        "Backing 1",
						CreatedAt:   time.Now(),
					},
					{
						ID:          uuid.New(),
						BackingCode: "02",
						Name:        "Backing 2",
						CreatedAt:   time.Now(),
					},
				}
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, backingDto.ReqBackingIndexFilter{}).Return(backings, 2, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Backing{
				{BackingCode: "01", Name: "Backing 1"},
				{BackingCode: "02", Name: "Backing 2"},
			},
			expectedTotal: 2,
		},
		{
			name: "success get index with search",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: backingDto.ReqBackingIndexFilter{},
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				backings := []models.Backing{
					{
						ID:          uuid.New(),
						BackingCode: "01",
						Name:        "Backing 1",
						CreatedAt:   time.Now(),
					},
				}
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, backingDto.ReqBackingIndexFilter{}).Return(backings, 1, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Backing{
				{BackingCode: "01", Name: "Backing 1"},
			},
			expectedTotal: 1,
		},
		{
			name: "error when repository fails",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: backingDto.ReqBackingIndexFilter{},
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, backingDto.ReqBackingIndexFilter{}).Return(nil, 0, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
			expectedTotal:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBackingRepository)
			mockTypeRepo := new(MockTypeRepository)
			tt.setupMock(mockRepo, mockTypeRepo)

			usecaseInstance := usecase.NewBackingUsecase(mockRepo, mockTypeRepo)

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
					assert.Equal(t, tt.expectedResult[0].BackingCode, result[0].BackingCode)
				}
			}

			mockRepo.AssertExpectations(t)
			mockTypeRepo.AssertExpectations(t)
		})
	}
}

func TestExportBacking(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		search         string
		filter         backingDto.ReqBackingIndexFilter
		setupMock      func(*MockBackingRepository, *MockTypeRepository)
		expectedError  error
		expectedResult []models.Backing
	}{
		{
			name:   "success export backings",
			search: "",
			filter: backingDto.ReqBackingIndexFilter{},
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				backings := []models.Backing{
					{
						ID:          uuid.New(),
						BackingCode: "01",
						Name:        "Backing 1",
						CreatedAt:   time.Now(),
					},
					{
						ID:          uuid.New(),
						BackingCode: "02",
						Name:        "Backing 2",
						CreatedAt:   time.Now(),
					},
				}
				m.On("GetAll", ctx, backingDto.ReqBackingIndexFilter{}).Return(backings, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Backing{
				{BackingCode: "01", Name: "Backing 1"},
				{BackingCode: "02", Name: "Backing 2"},
			},
		},
		{
			name:   "success export backings with search",
			search: "Backing 1",
			filter: backingDto.ReqBackingIndexFilter{Search: "Backing 1"},
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				backings := []models.Backing{
					{
						ID:          uuid.New(),
						BackingCode: "01",
						Name:        "Backing 1",
						CreatedAt:   time.Now(),
					},
				}
				m.On("GetAll", ctx, backingDto.ReqBackingIndexFilter{Search: "Backing 1"}).Return(backings, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Backing{
				{BackingCode: "01", Name: "Backing 1"},
			},
		},
		{
			name:   "error when repository fails",
			search: "",
			filter: backingDto.ReqBackingIndexFilter{},
			setupMock: func(m *MockBackingRepository, mt *MockTypeRepository) {
				m.On("GetAll", ctx, backingDto.ReqBackingIndexFilter{}).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockBackingRepository)
			mockTypeRepo := new(MockTypeRepository)
			tt.setupMock(mockRepo, mockTypeRepo)

			usecaseInstance := usecase.NewBackingUsecase(mockRepo, mockTypeRepo)

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
			mockTypeRepo.AssertExpectations(t)
		})
	}
}
