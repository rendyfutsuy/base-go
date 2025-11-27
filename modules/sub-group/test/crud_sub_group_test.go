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
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	groupDto "github.com/rendyfutsuy/base-go/modules/group/dto"
	subGroupDto "github.com/rendyfutsuy/base-go/modules/sub-group/dto"
	"github.com/rendyfutsuy/base-go/modules/sub-group/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func (m *MockSubGroupRepository) ExistsInTypes(ctx context.Context, subGroupID uuid.UUID) (bool, error) {
	args := m.Called(ctx, subGroupID)
	return args.Bool(0), args.Error(1)
}

// MockGroupRepository is a mock implementation of group.Repository
type MockGroupRepository struct {
	mock.Mock
}

func (m *MockGroupRepository) Create(ctx context.Context, name string, createdBy string) (*models.GoodsGroup, error) {
	args := m.Called(ctx, name, createdBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GoodsGroup), args.Error(1)
}

func (m *MockGroupRepository) Update(ctx context.Context, id uuid.UUID, name string, updatedBy string) (*models.GoodsGroup, error) {
	args := m.Called(ctx, id, name, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GoodsGroup), args.Error(1)
}

func (m *MockGroupRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *MockGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.GoodsGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GoodsGroup), args.Error(1)
}

func (m *MockGroupRepository) GetIndex(ctx context.Context, req request.PageRequest, filter groupDto.ReqGroupIndexFilter) ([]models.GoodsGroup, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.GoodsGroup), args.Int(1), args.Error(2)
}

func (m *MockGroupRepository) GetAll(ctx context.Context, filter groupDto.ReqGroupIndexFilter) ([]models.GoodsGroup, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.GoodsGroup), args.Error(1)
}

func (m *MockGroupRepository) ExistsByName(ctx context.Context, name string, excludeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockGroupRepository) ExistsInSubGroups(ctx context.Context, groupID uuid.UUID) (bool, error) {
	args := m.Called(ctx, groupID)
	return args.Bool(0), args.Error(1)
}

func TestCreateSubGroup(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	goodsGroupID := uuid.New()
	testUserID := uuid.New()

	tests := []struct {
		name           string
		req            *subGroupDto.ReqCreateSubGroup
		authId         string
		setupMock      func(*MockSubGroupRepository, *MockGroupRepository)
		setupContext   func(*echo.Context)
		expectedError  error
		expectedResult *models.SubGroup
	}{
		{
			name: "success create sub-group",
			req: &subGroupDto.ReqCreateSubGroup{
				GoodsGroupID: goodsGroupID,
				Name:         "Test Sub-Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				mg.On("GetByID", ctx, goodsGroupID).Return(&models.GoodsGroup{
					ID:   goodsGroupID,
					Name: "Test Group",
				}, nil).Once()
				m.On("ExistsByName", ctx, goodsGroupID, "Test Sub-Group", uuid.Nil).Return(false, nil).Once()
				m.On("Create", ctx, goodsGroupID, "Test Sub-Group", testUserID.String()).Return(&models.SubGroup{
					ID:           uuid.New(),
					GoodsGroupID: goodsGroupID,
					SubgroupCode: "01",
					Name:         "Test Sub-Group",
					CreatedAt:    time.Now(),
					CreatedBy:    testUserID.String(),
					UpdatedAt:    time.Now(),
					UpdatedBy:    testUserID.String(),
				}, nil).Once()
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError: nil,
			expectedResult: &models.SubGroup{
				GoodsGroupID: goodsGroupID,
				SubgroupCode: "01",
				Name:         "Test Sub-Group",
			},
		},
		{
			name: "error when name already exists in group",
			req: &subGroupDto.ReqCreateSubGroup{
				GoodsGroupID: goodsGroupID,
				Name:         "Existing Sub-Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				mg.On("GetByID", ctx, goodsGroupID).Return(&models.GoodsGroup{
					ID:   goodsGroupID,
					Name: "Test Group",
				}, nil).Once()
				m.On("ExistsByName", ctx, goodsGroupID, "Existing Sub-Group", uuid.Nil).Return(true, nil).Once()
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError:  errors.New("Sub-group name already exists in this group"),
			expectedResult: nil,
		},
		{
			name: "error when repository check exists fails",
			req: &subGroupDto.ReqCreateSubGroup{
				GoodsGroupID: goodsGroupID,
				Name:         "Test Sub-Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				mg.On("GetByID", ctx, goodsGroupID).Return(&models.GoodsGroup{
					ID:   goodsGroupID,
					Name: "Test Group",
				}, nil).Once()
				m.On("ExistsByName", ctx, goodsGroupID, "Test Sub-Group", uuid.Nil).Return(false, errors.New("database error")).Once()
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
		{
			name: "error when repository create fails",
			req: &subGroupDto.ReqCreateSubGroup{
				GoodsGroupID: goodsGroupID,
				Name:         "Test Sub-Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				mg.On("GetByID", ctx, goodsGroupID).Return(&models.GoodsGroup{
					ID:   goodsGroupID,
					Name: "Test Group",
				}, nil).Once()
				m.On("ExistsByName", ctx, goodsGroupID, "Test Sub-Group", uuid.Nil).Return(false, nil).Once()
				m.On("Create", ctx, goodsGroupID, "Test Sub-Group", testUserID.String()).Return(nil, errors.New("create failed")).Once()
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError:  errors.New("create failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSubGroupRepository)
			mockGroupRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo, mockGroupRepo)

			usecaseInstance := usecase.NewSubGroupUsecase(mockRepo, mockGroupRepo)

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))
			if tt.setupContext != nil {
				tt.setupContext(&c)
			}

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
					assert.Equal(t, tt.expectedResult.GoodsGroupID, result.GoodsGroupID)
					assert.Equal(t, tt.expectedResult.SubgroupCode, result.SubgroupCode)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateSubGroup(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	goodsGroupID := uuid.New()
	testUserID := uuid.New()

	tests := []struct {
		name           string
		id             string
		req            *subGroupDto.ReqUpdateSubGroup
		authId         string
		setupMock      func(*MockSubGroupRepository, *MockGroupRepository)
		setupContext   func(*echo.Context)
		expectedError  error
		expectedResult *models.SubGroup
	}{
		{
			name: "success update sub-group",
			id:   validID.String(),
			req: &subGroupDto.ReqUpdateSubGroup{
				GoodsGroupID: goodsGroupID,
				Name:         "Updated Sub-Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				mg.On("GetByID", ctx, goodsGroupID).Return(&models.GoodsGroup{
					ID:   goodsGroupID,
					Name: "Test Group",
				}, nil).Once()
				m.On("ExistsByName", ctx, goodsGroupID, "Updated Sub-Group", validID).Return(false, nil).Once()
				m.On("Update", ctx, validID, goodsGroupID, "Updated Sub-Group", testUserID.String()).Return(&models.SubGroup{
					ID:           validID,
					GoodsGroupID: goodsGroupID,
					SubgroupCode: "01",
					Name:         "Updated Sub-Group",
					UpdatedAt:    time.Now(),
					UpdatedBy:    testUserID.String(),
				}, nil).Once()
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError: nil,
			expectedResult: &models.SubGroup{
				ID:           validID,
				GoodsGroupID: goodsGroupID,
				SubgroupCode: "01",
				Name:         "Updated Sub-Group",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			req: &subGroupDto.ReqUpdateSubGroup{
				GoodsGroupID: goodsGroupID,
				Name:         "Updated Sub-Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				// No mock calls expected
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when name already exists in group",
			id:   validID.String(),
			req: &subGroupDto.ReqUpdateSubGroup{
				GoodsGroupID: goodsGroupID,
				Name:         "Existing Sub-Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				mg.On("GetByID", ctx, goodsGroupID).Return(&models.GoodsGroup{
					ID:   goodsGroupID,
					Name: "Test Group",
				}, nil).Once()
				m.On("ExistsByName", ctx, goodsGroupID, "Existing Sub-Group", validID).Return(true, nil).Once()
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError:  errors.New("Sub-group name already exists in this group"),
			expectedResult: nil,
		},
		{
			name: "error when repository update fails",
			id:   validID.String(),
			req: &subGroupDto.ReqUpdateSubGroup{
				GoodsGroupID: goodsGroupID,
				Name:         "Updated Sub-Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				mg.On("GetByID", ctx, goodsGroupID).Return(&models.GoodsGroup{
					ID:   goodsGroupID,
					Name: "Test Group",
				}, nil).Once()
				m.On("ExistsByName", ctx, goodsGroupID, "Updated Sub-Group", validID).Return(false, nil).Once()
				m.On("Update", ctx, validID, goodsGroupID, "Updated Sub-Group", testUserID.String()).Return(nil, errors.New("update failed")).Once()
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError:  errors.New("update failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSubGroupRepository)
			mockGroupRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo, mockGroupRepo)

			usecaseInstance := usecase.NewSubGroupUsecase(mockRepo, mockGroupRepo)

			req := httptest.NewRequest(http.MethodPut, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))
			if tt.setupContext != nil {
				tt.setupContext(&c)
			}

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
					assert.Equal(t, tt.expectedResult.GoodsGroupID, result.GoodsGroupID)
					assert.Equal(t, tt.expectedResult.SubgroupCode, result.SubgroupCode)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteSubGroup(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	testUserID := uuid.New()

	tests := []struct {
		name          string
		id            string
		authId        string
		setupMock     func(*MockSubGroupRepository, *MockGroupRepository)
		setupContext  func(*echo.Context)
		expectedError error
	}{
		{
			name:   "success delete sub-group",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				m.On("ExistsInTypes", ctx, validID).Return(false, nil).Once()
				m.On("Delete", ctx, validID, testUserID.String()).Return(nil).Once()
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError: nil,
		},
		{
			name:   "error when invalid UUID",
			id:     "invalid-uuid",
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				// No mock calls expected
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError: errors.New("requested param is string"),
		},
		{
			name:   "error when sub-group still used in types",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				m.On("ExistsInTypes", ctx, validID).Return(true, nil).Once()
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError: errors.New(constants.SubGroupStillUsedInTypes),
		},
		{
			name:   "error when repository delete fails",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				m.On("ExistsInTypes", ctx, validID).Return(false, nil).Once()
				m.On("Delete", ctx, validID, testUserID.String()).Return(errors.New("delete failed")).Once()
			},
			setupContext: func(c *echo.Context) {
				(*c).Set("user", models.User{ID: testUserID, Username: "testuser"})
			},
			expectedError: errors.New("delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSubGroupRepository)
			mockGroupRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo, mockGroupRepo)

			usecaseInstance := usecase.NewSubGroupUsecase(mockRepo, mockGroupRepo)

			req := httptest.NewRequest(http.MethodDelete, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))
			if tt.setupContext != nil {
				tt.setupContext(&c)
			}

			err := usecaseInstance.Delete(c, tt.id, tt.authId)

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

func TestGetSubGroupByID(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	goodsGroupID := uuid.New()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*MockSubGroupRepository, *MockGroupRepository)
		expectedError  error
		expectedResult *models.SubGroup
	}{
		{
			name: "success get sub-group by id",
			id:   validID.String(),
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				m.On("GetByID", ctx, validID).Return(&models.SubGroup{
					ID:             validID,
					GoodsGroupID:   goodsGroupID,
					GoodsGroupName: "Test Goods Group",
					SubgroupCode:   "01",
					Name:           "Test Sub-Group",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.SubGroup{
				ID:             validID,
				GoodsGroupID:   goodsGroupID,
				GoodsGroupName: "Test Goods Group",
				SubgroupCode:   "01",
				Name:           "Test Sub-Group",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when sub-group not found",
			id:   validID.String(),
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				m.On("GetByID", ctx, validID).Return(nil, errors.New("record not found")).Once()
			},
			expectedError:  errors.New("record not found"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSubGroupRepository)
			mockGroupRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo, mockGroupRepo)

			usecaseInstance := usecase.NewSubGroupUsecase(mockRepo, mockGroupRepo)

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
					assert.Equal(t, tt.expectedResult.GoodsGroupID, result.GoodsGroupID)
					assert.Equal(t, tt.expectedResult.SubgroupCode, result.SubgroupCode)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetIndexSubGroup(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	goodsGroupID := uuid.New()

	tests := []struct {
		name           string
		req            request.PageRequest
		filter         subGroupDto.ReqSubGroupIndexFilter
		setupMock      func(*MockSubGroupRepository, *MockGroupRepository)
		expectedError  error
		expectedResult []models.SubGroup
		expectedTotal  int
	}{
		{
			name: "success get index with pagination",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: subGroupDto.ReqSubGroupIndexFilter{},
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				subGroups := []models.SubGroup{
					{
						ID:           uuid.New(),
						GoodsGroupID: goodsGroupID,
						SubgroupCode: "01",
						Name:         "Sub-Group 1",
						CreatedAt:    time.Now(),
					},
					{
						ID:           uuid.New(),
						GoodsGroupID: goodsGroupID,
						SubgroupCode: "02",
						Name:         "Sub-Group 2",
						CreatedAt:    time.Now(),
					},
				}
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, subGroupDto.ReqSubGroupIndexFilter{}).Return(subGroups, 2, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.SubGroup{
				{GoodsGroupID: goodsGroupID, SubgroupCode: "01", Name: "Sub-Group 1"},
				{GoodsGroupID: goodsGroupID, SubgroupCode: "02", Name: "Sub-Group 2"},
			},
			expectedTotal: 2,
		},
		{
			name: "success get index with filter",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: subGroupDto.ReqSubGroupIndexFilter{
				GoodsGroupIDs: []uuid.UUID{goodsGroupID},
			},
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				subGroups := []models.SubGroup{
					{
						ID:           uuid.New(),
						GoodsGroupID: goodsGroupID,
						SubgroupCode: "01",
						Name:         "Sub-Group 1",
						CreatedAt:    time.Now(),
					},
				}
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, subGroupDto.ReqSubGroupIndexFilter{
					GoodsGroupIDs: []uuid.UUID{goodsGroupID},
				}).Return(subGroups, 1, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.SubGroup{
				{GoodsGroupID: goodsGroupID, SubgroupCode: "01", Name: "Sub-Group 1"},
			},
			expectedTotal: 1,
		},
		{
			name: "error when repository fails",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: subGroupDto.ReqSubGroupIndexFilter{},
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, subGroupDto.ReqSubGroupIndexFilter{}).Return(nil, 0, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
			expectedTotal:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSubGroupRepository)
			mockGroupRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo, mockGroupRepo)

			usecaseInstance := usecase.NewSubGroupUsecase(mockRepo, mockGroupRepo)

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
					assert.Equal(t, tt.expectedResult[0].GoodsGroupID, result[0].GoodsGroupID)
					assert.Equal(t, tt.expectedResult[0].SubgroupCode, result[0].SubgroupCode)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExportSubGroup(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	goodsGroupID := uuid.New()

	tests := []struct {
		name           string
		search         string
		filter         subGroupDto.ReqSubGroupIndexFilter
		setupMock      func(*MockSubGroupRepository, *MockGroupRepository)
		expectedError  error
		expectedResult []models.SubGroup
	}{
		{
			name:   "success export sub-groups",
			search: "",
			filter: subGroupDto.ReqSubGroupIndexFilter{},
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				subGroups := []models.SubGroup{
					{
						ID:           uuid.New(),
						GoodsGroupID: goodsGroupID,
						SubgroupCode: "01",
						Name:         "Sub-Group 1",
						CreatedAt:    time.Now(),
					},
					{
						ID:           uuid.New(),
						GoodsGroupID: goodsGroupID,
						SubgroupCode: "02",
						Name:         "Sub-Group 2",
						CreatedAt:    time.Now(),
					},
				}
				m.On("GetAll", ctx, subGroupDto.ReqSubGroupIndexFilter{}).Return(subGroups, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.SubGroup{
				{GoodsGroupID: goodsGroupID, SubgroupCode: "01", Name: "Sub-Group 1"},
				{GoodsGroupID: goodsGroupID, SubgroupCode: "02", Name: "Sub-Group 2"},
			},
		},
		{
			name:   "success export sub-groups with search",
			search: "Sub-Group 1",
			filter: subGroupDto.ReqSubGroupIndexFilter{Search: "Sub-Group 1"},
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				subGroups := []models.SubGroup{
					{
						ID:           uuid.New(),
						GoodsGroupID: goodsGroupID,
						SubgroupCode: "01",
						Name:         "Sub-Group 1",
						CreatedAt:    time.Now(),
					},
				}
				m.On("GetAll", ctx, subGroupDto.ReqSubGroupIndexFilter{Search: "Sub-Group 1"}).Return(subGroups, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.SubGroup{
				{GoodsGroupID: goodsGroupID, SubgroupCode: "01", Name: "Sub-Group 1"},
			},
		},
		{
			name:   "error when repository fails",
			search: "",
			filter: subGroupDto.ReqSubGroupIndexFilter{},
			setupMock: func(m *MockSubGroupRepository, mg *MockGroupRepository) {
				m.On("GetAll", ctx, subGroupDto.ReqSubGroupIndexFilter{}).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSubGroupRepository)
			mockGroupRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo, mockGroupRepo)

			usecaseInstance := usecase.NewSubGroupUsecase(mockRepo, mockGroupRepo)

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
		})
	}
}
