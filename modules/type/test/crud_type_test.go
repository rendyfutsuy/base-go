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
	subGroupDto "github.com/rendyfutsuy/base-go/modules/sub-group/dto"
	typeDto "github.com/rendyfutsuy/base-go/modules/type/dto"
	"github.com/rendyfutsuy/base-go/modules/type/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTypeRepository is a mock implementation of type_module.Repository
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

// MockSubGroupRepository is a mock implementation of sub_group.Repository
type MockSubGroupRepository struct {
	mock.Mock
}

func (m *MockSubGroupRepository) Create(ctx context.Context, goodsGroupID uuid.UUID, name string, createdBy string) (*models.SubGroup, error) {
	args := m.Called(ctx, goodsGroupID, name, createdBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubGroup), args.Error(1)
}

func (m *MockSubGroupRepository) Update(ctx context.Context, id uuid.UUID, goodsGroupID uuid.UUID, name string, updatedBy string) (*models.SubGroup, error) {
	args := m.Called(ctx, id, goodsGroupID, name, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubGroup), args.Error(1)
}

func (m *MockSubGroupRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *MockSubGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.SubGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SubGroup), args.Error(1)
}

func (m *MockSubGroupRepository) GetIndex(ctx context.Context, req request.PageRequest, filter subGroupDto.ReqSubGroupIndexFilter) ([]models.SubGroup, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.SubGroup), args.Int(1), args.Error(2)
}

func (m *MockSubGroupRepository) GetAll(ctx context.Context, filter subGroupDto.ReqSubGroupIndexFilter) ([]models.SubGroup, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SubGroup), args.Error(1)
}

func (m *MockSubGroupRepository) ExistsByName(ctx context.Context, goodsGroupID uuid.UUID, name string, excludeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, goodsGroupID, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func TestCreateType(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	testUserID := uuid.New()
	subgroupID := uuid.New()

	setupContext := func(c echo.Context) {
		c.Set("user", models.User{ID: testUserID, Username: "testuser"})
	}

	tests := []struct {
		name           string
		req            *typeDto.ReqCreateType
		authId         string
		setupMock      func(*MockTypeRepository, *MockSubGroupRepository)
		expectedError  error
		expectedResult *models.Type
	}{
		{
			name: "success create type",
			req: &typeDto.ReqCreateType{
				SubgroupID: subgroupID,
				Name:       "Test Type",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				msg.On("GetByID", ctx, subgroupID).Return(&models.SubGroup{
					ID:   subgroupID,
					Name: "Test SubGroup",
				}, nil).Once()
				m.On("ExistsByNameInSubgroup", ctx, subgroupID, "Test Type", uuid.Nil).Return(false, nil).Once()
				m.On("Create", ctx, subgroupID, "Test Type", testUserID.String()).Return(&models.Type{
					ID:         uuid.New(),
					SubgroupID: subgroupID,
					TypeCode:   "1",
					Name:       "Test Type",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					CreatedBy:  testUserID.String(),
					UpdatedBy:  testUserID.String(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Type{
				SubgroupID: subgroupID,
				TypeCode:   "1",
				Name:       "Test Type",
			},
		},
		{
			name: "error when name already exists in subgroup",
			req: &typeDto.ReqCreateType{
				SubgroupID: subgroupID,
				Name:       "Existing Type",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				msg.On("GetByID", ctx, subgroupID).Return(&models.SubGroup{
					ID:   subgroupID,
					Name: "Test SubGroup",
				}, nil).Once()
				m.On("ExistsByNameInSubgroup", ctx, subgroupID, "Existing Type", uuid.Nil).Return(true, nil).Once()
			},
			expectedError:  errors.New("Type name already exists in this subgroup"),
			expectedResult: nil,
		},
		{
			name: "error when repository check exists fails",
			req: &typeDto.ReqCreateType{
				SubgroupID: subgroupID,
				Name:       "Test Type",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				msg.On("GetByID", ctx, subgroupID).Return(&models.SubGroup{
					ID:   subgroupID,
					Name: "Test SubGroup",
				}, nil).Once()
				m.On("ExistsByNameInSubgroup", ctx, subgroupID, "Test Type", uuid.Nil).Return(false, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
		{
			name: "error when repository create fails",
			req: &typeDto.ReqCreateType{
				SubgroupID: subgroupID,
				Name:       "Test Type",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				msg.On("GetByID", ctx, subgroupID).Return(&models.SubGroup{
					ID:   subgroupID,
					Name: "Test SubGroup",
				}, nil).Once()
				m.On("ExistsByNameInSubgroup", ctx, subgroupID, "Test Type", uuid.Nil).Return(false, nil).Once()
				m.On("Create", ctx, subgroupID, "Test Type", testUserID.String()).Return(nil, errors.New("create failed")).Once()
			},
			expectedError:  errors.New("create failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTypeRepository)
			mockSubGroupRepo := new(MockSubGroupRepository)
			tt.setupMock(mockRepo, mockSubGroupRepo)

			usecaseInstance := usecase.NewTypeUsecase(mockRepo, mockSubGroupRepo)

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))
			setupContext(c)

			result, err := usecaseInstance.Create(c, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Name, result.Name)
					assert.Equal(t, tt.expectedResult.TypeCode, result.TypeCode)
					assert.Equal(t, tt.expectedResult.SubgroupID, result.SubgroupID)
				}
			}

			mockRepo.AssertExpectations(t)
			mockSubGroupRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateType(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	testUserID := uuid.New()
	validID := uuid.New()
	subgroupID := uuid.New()

	setupContext := func(c echo.Context) {
		c.Set("user", models.User{ID: testUserID, Username: "testuser"})
	}

	tests := []struct {
		name           string
		id             string
		req            *typeDto.ReqUpdateType
		authId         string
		setupMock      func(*MockTypeRepository, *MockSubGroupRepository)
		expectedError  error
		expectedResult *models.Type
	}{
		{
			name: "success update type",
			id:   validID.String(),
			req: &typeDto.ReqUpdateType{
				SubgroupID: subgroupID,
				Name:       "Updated Type",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				msg.On("GetByID", ctx, subgroupID).Return(&models.SubGroup{
					ID:   subgroupID,
					Name: "Test SubGroup",
				}, nil).Once()
				m.On("ExistsByNameInSubgroup", ctx, subgroupID, "Updated Type", validID).Return(false, nil).Once()
				m.On("Update", ctx, validID, subgroupID, "Updated Type", testUserID.String()).Return(&models.Type{
					ID:         validID,
					SubgroupID: subgroupID,
					TypeCode:   "2",
					Name:       "Updated Type",
					UpdatedAt:  time.Now(),
					UpdatedBy:  testUserID.String(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Type{
				ID:         validID,
				SubgroupID: subgroupID,
				TypeCode:   "2",
				Name:       "Updated Type",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			req: &typeDto.ReqUpdateType{
				SubgroupID: subgroupID,
				Name:       "Updated Type",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when name already exists in subgroup",
			id:   validID.String(),
			req: &typeDto.ReqUpdateType{
				SubgroupID: subgroupID,
				Name:       "Existing Type",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				msg.On("GetByID", ctx, subgroupID).Return(&models.SubGroup{
					ID:   subgroupID,
					Name: "Test SubGroup",
				}, nil).Once()
				m.On("ExistsByNameInSubgroup", ctx, subgroupID, "Existing Type", validID).Return(true, nil).Once()
			},
			expectedError:  errors.New("Type name already exists in this subgroup"),
			expectedResult: nil,
		},
		{
			name: "error when repository update fails",
			id:   validID.String(),
			req: &typeDto.ReqUpdateType{
				SubgroupID: subgroupID,
				Name:       "Updated Type",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				msg.On("GetByID", ctx, subgroupID).Return(&models.SubGroup{
					ID:   subgroupID,
					Name: "Test SubGroup",
				}, nil).Once()
				m.On("ExistsByNameInSubgroup", ctx, subgroupID, "Updated Type", validID).Return(false, nil).Once()
				m.On("Update", ctx, validID, subgroupID, "Updated Type", testUserID.String()).Return(nil, errors.New("update failed")).Once()
			},
			expectedError:  errors.New("update failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTypeRepository)
			mockSubGroupRepo := new(MockSubGroupRepository)
			tt.setupMock(mockRepo, mockSubGroupRepo)

			usecaseInstance := usecase.NewTypeUsecase(mockRepo, mockSubGroupRepo)

			req := httptest.NewRequest(http.MethodPut, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))
			setupContext(c)

			result, err := usecaseInstance.Update(c, tt.id, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.Name, result.Name)
					assert.Equal(t, tt.expectedResult.TypeCode, result.TypeCode)
					assert.Equal(t, tt.expectedResult.SubgroupID, result.SubgroupID)
				}
			}

			mockRepo.AssertExpectations(t)
			mockSubGroupRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteType(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	testUserID := uuid.New()
	validID := uuid.New()

	setupContext := func(c echo.Context) {
		c.Set("user", models.User{ID: testUserID, Username: "testuser"})
	}

	tests := []struct {
		name          string
		id            string
		authId        string
		setupMock     func(*MockTypeRepository, *MockSubGroupRepository)
		expectedError error
	}{
		{
			name:   "success delete type",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				m.On("Delete", ctx, validID, testUserID.String()).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:   "error when invalid UUID",
			id:     "invalid-uuid",
			authId: "test-auth-id",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				// No mock calls expected
			},
			expectedError: errors.New("requested param is string"),
		},
		{
			name:   "error when repository delete fails",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				m.On("Delete", ctx, validID, testUserID.String()).Return(errors.New("delete failed")).Once()
			},
			expectedError: errors.New("delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTypeRepository)
			mockSubGroupRepo := new(MockSubGroupRepository)
			tt.setupMock(mockRepo, mockSubGroupRepo)

			usecaseInstance := usecase.NewTypeUsecase(mockRepo, mockSubGroupRepo)

			req := httptest.NewRequest(http.MethodDelete, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))
			setupContext(c)

			err := usecaseInstance.Delete(c, tt.id, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockSubGroupRepo.AssertExpectations(t)
		})
	}
}

func TestGetTypeByID(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	subgroupID := uuid.New()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*MockTypeRepository, *MockSubGroupRepository)
		expectedError  error
		expectedResult *models.Type
	}{
		{
			name: "success get type by id",
			id:   validID.String(),
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				m.On("GetByID", ctx, validID).Return(&models.Type{
					ID:           validID,
					SubgroupID:   subgroupID,
					SubgroupName: "Test Sub-Group",
					TypeCode:     "1",
					Name:         "Test Type",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Type{
				ID:           validID,
				SubgroupID:   subgroupID,
				SubgroupName: "Test Sub-Group",
				TypeCode:     "1",
				Name:         "Test Type",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when type not found",
			id:   validID.String(),
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				m.On("GetByID", ctx, validID).Return(nil, errors.New("record not found")).Once()
			},
			expectedError:  errors.New("record not found"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTypeRepository)
			mockSubGroupRepo := new(MockSubGroupRepository)
			tt.setupMock(mockRepo, mockSubGroupRepo)

			usecaseInstance := usecase.NewTypeUsecase(mockRepo, mockSubGroupRepo)

			req := httptest.NewRequest(http.MethodGet, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

			result, err := usecaseInstance.GetByID(c, tt.id)

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
					assert.Equal(t, tt.expectedResult.TypeCode, result.TypeCode)
					assert.Equal(t, tt.expectedResult.SubgroupName, result.SubgroupName)
				}
			}

			mockRepo.AssertExpectations(t)
			mockSubGroupRepo.AssertExpectations(t)
		})
	}
}

func TestGetIndexType(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	subgroupID := uuid.New()

	tests := []struct {
		name           string
		req            request.PageRequest
		filter         typeDto.ReqTypeIndexFilter
		setupMock      func(*MockTypeRepository, *MockSubGroupRepository)
		expectedError  error
		expectedResult []models.Type
		expectedTotal  int
	}{
		{
			name: "success get index with pagination",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: typeDto.ReqTypeIndexFilter{},
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				types := []models.Type{
					{
						ID:           uuid.New(),
						SubgroupID:   subgroupID,
						SubgroupName: "Sub-Group 1",
						TypeCode:     "1",
						Name:         "Type 1",
						CreatedAt:    time.Now(),
					},
					{
						ID:           uuid.New(),
						SubgroupID:   subgroupID,
						SubgroupName: "Sub-Group 1",
						TypeCode:     "2",
						Name:         "Type 2",
						CreatedAt:    time.Now(),
					},
				}
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, typeDto.ReqTypeIndexFilter{}).Return(types, 2, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Type{
				{TypeCode: "1", Name: "Type 1"},
				{TypeCode: "2", Name: "Type 2"},
			},
			expectedTotal: 2,
		},
		{
			name: "success get index with filters",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: typeDto.ReqTypeIndexFilter{
				TypeCodes:   []string{"1", "2"},
				SubgroupIDs: []string{subgroupID.String()},
				Names:       []string{"Type 1"},
			},
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				types := []models.Type{
					{
						ID:           uuid.New(),
						SubgroupID:   subgroupID,
						SubgroupName: "Sub-Group 1",
						TypeCode:     "1",
						Name:         "Type 1",
						CreatedAt:    time.Now(),
					},
				}
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, typeDto.ReqTypeIndexFilter{
					TypeCodes:   []string{"1", "2"},
					SubgroupIDs: []string{subgroupID.String()},
					Names:       []string{"Type 1"},
				}).Return(types, 1, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Type{
				{TypeCode: "1", Name: "Type 1"},
			},
			expectedTotal: 1,
		},
		{
			name: "error when repository fails",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: typeDto.ReqTypeIndexFilter{},
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, typeDto.ReqTypeIndexFilter{}).Return(nil, 0, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
			expectedTotal:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTypeRepository)
			mockSubGroupRepo := new(MockSubGroupRepository)
			tt.setupMock(mockRepo, mockSubGroupRepo)

			usecaseInstance := usecase.NewTypeUsecase(mockRepo, mockSubGroupRepo)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

			result, total, err := usecaseInstance.GetIndex(c, tt.req, tt.filter)

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
					assert.Equal(t, tt.expectedResult[0].TypeCode, result[0].TypeCode)
				}
			}

			mockRepo.AssertExpectations(t)
			mockSubGroupRepo.AssertExpectations(t)
		})
	}
}

func TestExportType(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	subgroupID := uuid.New()

	tests := []struct {
		name           string
		search         string
		filter         typeDto.ReqTypeIndexFilter
		setupMock      func(*MockTypeRepository, *MockSubGroupRepository)
		expectedError  error
		expectedResult []models.Type
	}{
		{
			name:   "success export types",
			search: "",
			filter: typeDto.ReqTypeIndexFilter{},
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				types := []models.Type{
					{
						ID:           uuid.New(),
						SubgroupID:   subgroupID,
						SubgroupName: "Sub-Group 1",
						TypeCode:     "1",
						Name:         "Type 1",
						CreatedAt:    time.Now(),
					},
					{
						ID:           uuid.New(),
						SubgroupID:   subgroupID,
						SubgroupName: "Sub-Group 1",
						TypeCode:     "2",
						Name:         "Type 2",
						CreatedAt:    time.Now(),
					},
				}
				m.On("GetAll", ctx, typeDto.ReqTypeIndexFilter{}).Return(types, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Type{
				{TypeCode: "1", Name: "Type 1"},
				{TypeCode: "2", Name: "Type 2"},
			},
		},
		{
			name:   "success export types with search",
			search: "Type 1",
			filter: typeDto.ReqTypeIndexFilter{Search: "Type 1"},
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				types := []models.Type{
					{
						ID:           uuid.New(),
						SubgroupID:   subgroupID,
						SubgroupName: "Sub-Group 1",
						TypeCode:     "1",
						Name:         "Type 1",
						CreatedAt:    time.Now(),
					},
				}
				m.On("GetAll", ctx, typeDto.ReqTypeIndexFilter{Search: "Type 1"}).Return(types, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Type{
				{TypeCode: "1", Name: "Type 1"},
			},
		},
		{
			name:   "error when repository fails",
			search: "",
			filter: typeDto.ReqTypeIndexFilter{},
			setupMock: func(m *MockTypeRepository, msg *MockSubGroupRepository) {
				m.On("GetAll", ctx, typeDto.ReqTypeIndexFilter{}).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTypeRepository)
			mockSubGroupRepo := new(MockSubGroupRepository)
			tt.setupMock(mockRepo, mockSubGroupRepo)

			usecaseInstance := usecase.NewTypeUsecase(mockRepo, mockSubGroupRepo)

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

			result, err := usecaseInstance.Export(c, tt.filter)

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
			mockSubGroupRepo.AssertExpectations(t)
		})
	}
}
