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
	groupDto "github.com/rendyfutsuy/base-go/modules/group/dto"
	"github.com/rendyfutsuy/base-go/modules/group/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGroupRepository is a mock implementation of group.Repository
type MockGroupRepository struct {
	mock.Mock
}

func (m *MockGroupRepository) Create(ctx context.Context, name string) (*models.GoodsGroup, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GoodsGroup), args.Error(1)
}

func (m *MockGroupRepository) Update(ctx context.Context, id uuid.UUID, name string) (*models.GoodsGroup, error) {
	args := m.Called(ctx, id, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GoodsGroup), args.Error(1)
}

func (m *MockGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
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

func TestCreateGroup(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		req            *groupDto.ReqCreateGroup
		authId         string
		setupMock      func(*MockGroupRepository)
		expectedError  error
		expectedResult *models.GoodsGroup
	}{
		{
			name: "success create group",
			req: &groupDto.ReqCreateGroup{
				Name: "Test Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockGroupRepository) {
				m.On("ExistsByName", ctx, "Test Group", uuid.Nil).Return(false, nil).Once()
				m.On("Create", ctx, "Test Group").Return(&models.GoodsGroup{
					ID:        uuid.New(),
					GroupCode: "01",
					Name:      "Test Group",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.GoodsGroup{
				GroupCode: "01",
				Name:      "Test Group",
			},
		},
		{
			name: "error when name already exists",
			req: &groupDto.ReqCreateGroup{
				Name: "Existing Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockGroupRepository) {
				m.On("ExistsByName", ctx, "Existing Group", uuid.Nil).Return(true, nil).Once()
			},
			expectedError:  errors.New("Group name already exists"),
			expectedResult: nil,
		},
		{
			name: "error when repository check exists fails",
			req: &groupDto.ReqCreateGroup{
				Name: "Test Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockGroupRepository) {
				m.On("ExistsByName", ctx, "Test Group", uuid.Nil).Return(false, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
		{
			name: "error when repository create fails",
			req: &groupDto.ReqCreateGroup{
				Name: "Test Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockGroupRepository) {
				m.On("ExistsByName", ctx, "Test Group", uuid.Nil).Return(false, nil).Once()
				m.On("Create", ctx, "Test Group").Return(nil, errors.New("create failed")).Once()
			},
			expectedError:  errors.New("create failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewGroupUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))

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
					assert.Equal(t, tt.expectedResult.GroupCode, result.GroupCode)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateGroup(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name           string
		id             string
		req            *groupDto.ReqUpdateGroup
		authId         string
		setupMock      func(*MockGroupRepository)
		expectedError  error
		expectedResult *models.GoodsGroup
	}{
		{
			name: "success update group",
			id:   validID.String(),
			req: &groupDto.ReqUpdateGroup{
				Name: "Updated Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockGroupRepository) {
				m.On("ExistsByName", ctx, "Updated Group", validID).Return(false, nil).Once()
				m.On("Update", ctx, validID, "Updated Group").Return(&models.GoodsGroup{
					ID:        validID,
					GroupCode: "02",
					Name:      "Updated Group",
					UpdatedAt: time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.GoodsGroup{
				ID:        validID,
				GroupCode: "02",
				Name:      "Updated Group",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			req: &groupDto.ReqUpdateGroup{
				Name: "Updated Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockGroupRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when name already exists",
			id:   validID.String(),
			req: &groupDto.ReqUpdateGroup{
				Name: "Existing Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockGroupRepository) {
				m.On("ExistsByName", ctx, "Existing Group", validID).Return(true, nil).Once()
			},
			expectedError:  errors.New("Group name already exists"),
			expectedResult: nil,
		},
		{
			name: "error when repository update fails",
			id:   validID.String(),
			req: &groupDto.ReqUpdateGroup{
				Name: "Updated Group",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockGroupRepository) {
				m.On("ExistsByName", ctx, "Updated Group", validID).Return(false, nil).Once()
				m.On("Update", ctx, validID, "Updated Group").Return(nil, errors.New("update failed")).Once()
			},
			expectedError:  errors.New("update failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewGroupUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPut, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

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
					assert.Equal(t, tt.expectedResult.GroupCode, result.GroupCode)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteGroup(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name          string
		id            string
		authId        string
		setupMock     func(*MockGroupRepository)
		expectedError error
	}{
		{
			name:   "success delete group",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockGroupRepository) {
				m.On("Delete", ctx, validID).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:   "error when invalid UUID",
			id:     "invalid-uuid",
			authId: "test-auth-id",
			setupMock: func(m *MockGroupRepository) {
				// No mock calls expected
			},
			expectedError: errors.New("requested param is string"),
		},
		{
			name:   "error when repository delete fails",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockGroupRepository) {
				m.On("Delete", ctx, validID).Return(errors.New("delete failed")).Once()
			},
			expectedError: errors.New("delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewGroupUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodDelete, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))

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

func TestGetGroupByID(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*MockGroupRepository)
		expectedError  error
		expectedResult *models.GoodsGroup
	}{
		{
			name: "success get group by id",
			id:   validID.String(),
			setupMock: func(m *MockGroupRepository) {
				m.On("GetByID", ctx, validID).Return(&models.GoodsGroup{
					ID:        validID,
					GroupCode: "01",
					Name:      "Test Group",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.GoodsGroup{
				ID:        validID,
				GroupCode: "01",
				Name:      "Test Group",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			setupMock: func(m *MockGroupRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when group not found",
			id:   validID.String(),
			setupMock: func(m *MockGroupRepository) {
				m.On("GetByID", ctx, validID).Return(nil, errors.New("record not found")).Once()
			},
			expectedError:  errors.New("record not found"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewGroupUsecase(mockRepo)

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
					assert.Equal(t, tt.expectedResult.GroupCode, result.GroupCode)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetIndexGroup(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		req            request.PageRequest
		filter         groupDto.ReqGroupIndexFilter
		setupMock      func(*MockGroupRepository)
		expectedError  error
		expectedResult []models.GoodsGroup
		expectedTotal  int
	}{
		{
			name: "success get index with pagination",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: groupDto.ReqGroupIndexFilter{},
			setupMock: func(m *MockGroupRepository) {
				groups := []models.GoodsGroup{
					{
						ID:        uuid.New(),
						GroupCode: "01",
						Name:      "Group 1",
						CreatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						GroupCode: "02",
						Name:      "Group 2",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, groupDto.ReqGroupIndexFilter{}).Return(groups, 2, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.GoodsGroup{
				{GroupCode: "01", Name: "Group 1"},
				{GroupCode: "02", Name: "Group 2"},
			},
			expectedTotal: 2,
		},
		{
			name: "success get index with search",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: groupDto.ReqGroupIndexFilter{},
			setupMock: func(m *MockGroupRepository) {
				groups := []models.GoodsGroup{
					{
						ID:        uuid.New(),
						GroupCode: "01",
						Name:      "Group 1",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, groupDto.ReqGroupIndexFilter{}).Return(groups, 1, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.GoodsGroup{
				{GroupCode: "01", Name: "Group 1"},
			},
			expectedTotal: 1,
		},
		{
			name: "error when repository fails",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: groupDto.ReqGroupIndexFilter{},
			setupMock: func(m *MockGroupRepository) {
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, groupDto.ReqGroupIndexFilter{}).Return(nil, 0, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
			expectedTotal:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewGroupUsecase(mockRepo)

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
					assert.Equal(t, tt.expectedResult[0].GroupCode, result[0].GroupCode)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExportGroup(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		search         string
		filter         groupDto.ReqGroupIndexFilter
		setupMock      func(*MockGroupRepository)
		expectedError  error
		expectedResult []models.GoodsGroup
	}{
		{
			name:   "success export groups",
			search: "",
			filter: groupDto.ReqGroupIndexFilter{},
			setupMock: func(m *MockGroupRepository) {
				groups := []models.GoodsGroup{
					{
						ID:        uuid.New(),
						GroupCode: "01",
						Name:      "Group 1",
						CreatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						GroupCode: "02",
						Name:      "Group 2",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetAll", ctx, groupDto.ReqGroupIndexFilter{}).Return(groups, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.GoodsGroup{
				{GroupCode: "01", Name: "Group 1"},
				{GroupCode: "02", Name: "Group 2"},
			},
		},
		{
			name:   "success export groups with search",
			search: "Group 1",
			filter: groupDto.ReqGroupIndexFilter{Search: "Group 1"},
			setupMock: func(m *MockGroupRepository) {
				groups := []models.GoodsGroup{
					{
						ID:        uuid.New(),
						GroupCode: "01",
						Name:      "Group 1",
						CreatedAt: time.Now(),
					},
				}
				m.On("GetAll", ctx, groupDto.ReqGroupIndexFilter{Search: "Group 1"}).Return(groups, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.GoodsGroup{
				{GroupCode: "01", Name: "Group 1"},
			},
		},
		{
			name:   "error when repository fails",
			search: "",
			filter: groupDto.ReqGroupIndexFilter{},
			setupMock: func(m *MockGroupRepository) {
				m.On("GetAll", ctx, groupDto.ReqGroupIndexFilter{}).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockGroupRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewGroupUsecase(mockRepo)

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
