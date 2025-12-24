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
	expeditionMod "github.com/rendyfutsuy/base-go/modules/expedition"
	expeditionDto "github.com/rendyfutsuy/base-go/modules/expedition/dto"
	"github.com/rendyfutsuy/base-go/modules/expedition/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockExpeditionRepository is a mock implementation of expedition.Repository
type MockExpeditionRepository struct {
	mock.Mock
}

func (m *MockExpeditionRepository) Create(ctx context.Context, params expeditionMod.CreateExpeditionParams) (*models.Expedition, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result := args.Get(0).(*models.Expedition)
	if result.ID == uuid.Nil {
		result.ID = uuid.New()
	}
	return result, args.Error(1)
}

func (m *MockExpeditionRepository) Update(ctx context.Context, id uuid.UUID, params expeditionMod.UpdateExpeditionParams) (*models.Expedition, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expedition), args.Error(1)
}

func (m *MockExpeditionRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}

func (m *MockExpeditionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Expedition, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expedition), args.Error(1)
}

func (m *MockExpeditionRepository) GetIndex(ctx context.Context, req request.PageRequest, filter expeditionDto.ReqExpeditionIndexFilter) ([]models.Expedition, int, error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Expedition), args.Int(1), args.Error(2)
}

func (m *MockExpeditionRepository) GetAll(ctx context.Context, filter expeditionDto.ReqExpeditionIndexFilter) ([]models.Expedition, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Expedition), args.Error(1)
}

func (m *MockExpeditionRepository) GetAllForExport(ctx context.Context, filter expeditionDto.ReqExpeditionIndexFilter) ([]expeditionDto.ExpeditionExport, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]expeditionDto.ExpeditionExport), args.Error(1)
}

func (m *MockExpeditionRepository) ExistsByExpeditionName(ctx context.Context, expeditionName string, excludeID uuid.UUID) (bool, error) {
	args := m.Called(ctx, expeditionName, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockExpeditionRepository) GetContactsByExpeditionID(ctx context.Context, expeditionID uuid.UUID) ([]models.ExpeditionContact, error) {
	args := m.Called(ctx, expeditionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.ExpeditionContact), args.Error(1)
}

func TestCreateExpedition(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	telpNumber := "021-1234567"
	phoneNumber := "081234567890"
	notes := "Test notes"
	testUserID := uuid.New()

	tests := []struct {
		name           string
		req            *expeditionDto.ReqCreateExpedition
		authId         string
		setupMock      func(*MockExpeditionRepository)
		expectedError  error
		expectedResult *models.Expedition
	}{
		{
			name: "success create expedition",
			req: &expeditionDto.ReqCreateExpedition{
				ExpeditionName: "JNE",
				Address:        "Jl. Ahmad Yani No. 123",
				TelpNumbers:    []expeditionDto.TelpNumberItem{{PhoneNumber: telpNumber, AreaCode: nil}},
				PhoneNumbers:   []string{phoneNumber},
				Notes:          &notes,
			},
			authId: "test-auth-id",
			setupMock: func(m *MockExpeditionRepository) {
				m.On("ExistsByExpeditionName", mock.Anything, "JNE", uuid.Nil).Return(false, nil).Once()
				createdExp := &models.Expedition{
					ID:             uuid.New(),
					ExpeditionCode: "01",
					ExpeditionName: "JNE",
					Address:        "Jl. Ahmad Yani No. 123",
					Notes:          &notes,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				m.On("Create", mock.Anything, expeditionMod.CreateExpeditionParams{
					ExpeditionName: "JNE",
					Address:        "Jl. Ahmad Yani No. 123",
					TelpNumbers:    []expeditionDto.TelpNumberItem{{PhoneNumber: telpNumber, AreaCode: nil}},
					PhoneNumbers:   []string{phoneNumber},
					Notes:          &notes,
					CreatedBy:      "test-auth-id",
				}).Return(createdExp, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Expedition{
				ExpeditionCode: "01",
				ExpeditionName: "JNE",
				Address:        "Jl. Ahmad Yani No. 123",
			},
		},
		{
			name: "error when expedition name already exists",
			req: &expeditionDto.ReqCreateExpedition{
				ExpeditionName: "JNE",
				Address:        "Jl. Ahmad Yani No. 123",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockExpeditionRepository) {
				m.On("ExistsByExpeditionName", mock.Anything, "JNE", uuid.Nil).Return(true, nil).Once()
			},
			expectedError:  errors.New("Expedition name already exists"),
			expectedResult: nil,
		},
		{
			name: "error when repository check exists fails",
			req: &expeditionDto.ReqCreateExpedition{
				ExpeditionName: "JNE",
				Address:        "Jl. Ahmad Yani No. 123",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockExpeditionRepository) {
				m.On("ExistsByExpeditionName", mock.Anything, "JNE", uuid.Nil).Return(false, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
		{
			name: "error when repository create fails",
			req: &expeditionDto.ReqCreateExpedition{
				ExpeditionName: "JNE",
				Address:        "Jl. Ahmad Yani No. 123",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockExpeditionRepository) {
				m.On("ExistsByExpeditionName", mock.Anything, "JNE", uuid.Nil).Return(false, nil).Once()
				m.On("Create", mock.Anything, mock.MatchedBy(func(params expeditionMod.CreateExpeditionParams) bool {
					return params.ExpeditionName == "JNE" &&
						params.Address == "Jl. Ahmad Yani No. 123" &&
						len(params.TelpNumbers) == 0 &&
						len(params.PhoneNumbers) == 0 &&
						params.Notes == nil &&
						params.CreatedBy == "test-auth-id"
				})).Return(nil, errors.New("create failed")).Once()
			},
			expectedError:  errors.New("create failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockExpeditionRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewExpeditionUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetRequest(req.WithContext(ctx))
			// Set user context
			c.Set("user", models.User{ID: testUserID, Username: "testuser"})

			result, err := usecaseInstance.Create(ctx, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.ExpeditionName, result.ExpeditionName)
					assert.Equal(t, tt.expectedResult.Address, result.Address)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateExpedition(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	testUserID := uuid.New()
	telpNumber := "021-1234567"
	phoneNumber := "081234567890"
	notes := "Updated notes"

	tests := []struct {
		name           string
		id             string
		req            *expeditionDto.ReqUpdateExpedition
		authId         string
		setupMock      func(*MockExpeditionRepository)
		expectedError  error
		expectedResult *models.Expedition
	}{
		{
			name: "success update expedition",
			id:   validID.String(),
			req: &expeditionDto.ReqUpdateExpedition{
				ExpeditionName: "JNE Updated",
				Address:        "Jl. Ahmad Yani No. 456",
				TelpNumbers:    []expeditionDto.TelpNumberItem{{PhoneNumber: telpNumber, AreaCode: nil}},
				PhoneNumbers:   []string{phoneNumber},
				Notes:          &notes,
			},
			authId: "test-auth-id",
			setupMock: func(m *MockExpeditionRepository) {
				m.On("ExistsByExpeditionName", mock.Anything, "JNE Updated", validID).Return(false, nil).Once()
				m.On("Update", mock.Anything, validID, expeditionMod.UpdateExpeditionParams{
					ExpeditionName: "JNE Updated",
					Address:        "Jl. Ahmad Yani No. 456",
					TelpNumbers:    []expeditionDto.TelpNumberItem{{PhoneNumber: telpNumber, AreaCode: nil}},
					PhoneNumbers:   []string{phoneNumber},
					Notes:          &notes,
					UpdatedBy:      "test-auth-id",
				}).Return(&models.Expedition{
					ID:             validID,
					ExpeditionCode: "01",
					ExpeditionName: "JNE Updated",
					Address:        "Jl. Ahmad Yani No. 456",
					Notes:          &notes,
					UpdatedAt:      time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Expedition{
				ID:             validID,
				ExpeditionCode: "01",
				ExpeditionName: "JNE Updated",
				Address:        "Jl. Ahmad Yani No. 456",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			req: &expeditionDto.ReqUpdateExpedition{
				ExpeditionName: "JNE Updated",
				Address:        "Jl. Ahmad Yani No. 456",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockExpeditionRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when expedition name already exists",
			id:   validID.String(),
			req: &expeditionDto.ReqUpdateExpedition{
				ExpeditionName: "JNE",
				Address:        "Jl. Ahmad Yani No. 123",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockExpeditionRepository) {
				m.On("ExistsByExpeditionName", mock.Anything, "JNE", validID).Return(true, nil).Once()
			},
			expectedError:  errors.New("Expedition name already exists"),
			expectedResult: nil,
		},
		{
			name: "error when repository update fails",
			id:   validID.String(),
			req: &expeditionDto.ReqUpdateExpedition{
				ExpeditionName: "JNE Updated",
				Address:        "Jl. Ahmad Yani No. 456",
			},
			authId: "test-auth-id",
			setupMock: func(m *MockExpeditionRepository) {
				m.On("ExistsByExpeditionName", mock.Anything, "JNE Updated", validID).Return(false, nil).Once()
				m.On("Update", mock.Anything, validID, mock.MatchedBy(func(params expeditionMod.UpdateExpeditionParams) bool {
					return params.ExpeditionName == "JNE Updated" &&
						params.Address == "Jl. Ahmad Yani No. 456" &&
						len(params.TelpNumbers) == 0 &&
						len(params.PhoneNumbers) == 0 &&
						params.Notes == nil &&
						params.UpdatedBy == "test-auth-id"
				})).Return(nil, errors.New("update failed")).Once()
			},
			expectedError:  errors.New("update failed"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockExpeditionRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewExpeditionUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodPut, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))
			// Set user context
			c.Set("user", models.User{ID: testUserID, Username: "testuser"})

			result, err := usecaseInstance.Update(ctx, tt.id, tt.req, tt.authId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectedResult != nil {
					assert.Equal(t, tt.expectedResult.ExpeditionName, result.ExpeditionName)
					assert.Equal(t, tt.expectedResult.Address, result.Address)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteExpedition(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()
	testUserID := uuid.New()

	tests := []struct {
		name          string
		id            string
		authId        string
		setupMock     func(*MockExpeditionRepository)
		expectedError error
	}{
		{
			name:   "success delete expedition",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockExpeditionRepository) {
				m.On("Delete", ctx, validID, "test-auth-id").Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:   "error when invalid UUID",
			id:     "invalid-uuid",
			authId: "test-auth-id",
			setupMock: func(m *MockExpeditionRepository) {
				// No mock calls expected
			},
			expectedError: errors.New("requested param is string"),
		},
		{
			name:   "error when repository delete fails",
			id:     validID.String(),
			authId: "test-auth-id",
			setupMock: func(m *MockExpeditionRepository) {
				m.On("Delete", ctx, validID, "test-auth-id").Return(errors.New("delete failed")).Once()
			},
			expectedError: errors.New("delete failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockExpeditionRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewExpeditionUsecase(mockRepo)

			req := httptest.NewRequest(http.MethodDelete, "/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id)
			c.SetRequest(req.WithContext(ctx))
			// Set user context
			c.Set("user", models.User{ID: testUserID, Username: "testuser"})

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

func TestGetExpeditionByID(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	validID := uuid.New()

	tests := []struct {
		name           string
		id             string
		setupMock      func(*MockExpeditionRepository)
		expectedError  error
		expectedResult *models.Expedition
	}{
		{
			name: "success get expedition by id",
			id:   validID.String(),
			setupMock: func(m *MockExpeditionRepository) {
				m.On("GetByID", ctx, validID).Return(&models.Expedition{
					ID:             validID,
					ExpeditionCode: "01",
					ExpeditionName: "JNE",
					Address:        "Jl. Ahmad Yani No. 123",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}, nil).Once()
			},
			expectedError: nil,
			expectedResult: &models.Expedition{
				ID:             validID,
				ExpeditionCode: "01",
				ExpeditionName: "JNE",
				Address:        "Jl. Ahmad Yani No. 123",
			},
		},
		{
			name: "error when invalid UUID",
			id:   "invalid-uuid",
			setupMock: func(m *MockExpeditionRepository) {
				// No mock calls expected
			},
			expectedError:  errors.New("requested param is string"),
			expectedResult: nil,
		},
		{
			name: "error when expedition not found",
			id:   validID.String(),
			setupMock: func(m *MockExpeditionRepository) {
				m.On("GetByID", ctx, validID).Return(nil, errors.New("record not found")).Once()
			},
			expectedError:  errors.New("record not found"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockExpeditionRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewExpeditionUsecase(mockRepo)

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
					assert.Equal(t, tt.expectedResult.ExpeditionName, result.ExpeditionName)
					assert.Equal(t, tt.expectedResult.ExpeditionCode, result.ExpeditionCode)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetIndexExpedition(t *testing.T) {
	e := echo.New()
	ctx := context.Background()
	telpNumber := "021-1234567"
	phoneNumber := "081234567890"

	tests := []struct {
		name           string
		req            request.PageRequest
		filter         expeditionDto.ReqExpeditionIndexFilter
		setupMock      func(*MockExpeditionRepository)
		expectedError  error
		expectedResult []models.Expedition
		expectedTotal  int
	}{
		{
			name: "success get index with pagination",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: expeditionDto.ReqExpeditionIndexFilter{},
			setupMock: func(m *MockExpeditionRepository) {
				expeditions := []models.Expedition{
					{
						ID:                 uuid.New(),
						ExpeditionCode:     "01",
						ExpeditionName:     "JNE",
						Address:            "Jl. Ahmad Yani No. 123",
						PrimaryTelpNumber:  &telpNumber,
						PrimaryPhoneNumber: &phoneNumber,
						CreatedAt:          time.Now(),
					},
					{
						ID:                 uuid.New(),
						ExpeditionCode:     "02",
						ExpeditionName:     "TIKI",
						Address:            "Jl. Sudirman No. 456",
						PrimaryTelpNumber:  &telpNumber,
						PrimaryPhoneNumber: &phoneNumber,
						CreatedAt:          time.Now(),
					},
				}
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, expeditionDto.ReqExpeditionIndexFilter{}).Return(expeditions, 2, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Expedition{
				{ExpeditionCode: "01", ExpeditionName: "JNE"},
				{ExpeditionCode: "02", ExpeditionName: "TIKI"},
			},
			expectedTotal: 2,
		},
		{
			name: "error when repository fails",
			req: request.PageRequest{
				Page:    1,
				PerPage: 10,
			},
			filter: expeditionDto.ReqExpeditionIndexFilter{},
			setupMock: func(m *MockExpeditionRepository) {
				m.On("GetIndex", ctx, request.PageRequest{Page: 1, PerPage: 10}, expeditionDto.ReqExpeditionIndexFilter{}).Return(nil, 0, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
			expectedTotal:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockExpeditionRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewExpeditionUsecase(mockRepo)

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
					assert.Equal(t, tt.expectedResult[0].ExpeditionName, result[0].ExpeditionName)
					assert.Equal(t, tt.expectedResult[0].ExpeditionCode, result[0].ExpeditionCode)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExportExpedition(t *testing.T) {
	e := echo.New()
	ctx := context.Background()

	tests := []struct {
		name           string
		search         string
		filter         expeditionDto.ReqExpeditionIndexFilter
		setupMock      func(*MockExpeditionRepository)
		expectedError  error
		expectedResult []models.Expedition
	}{
		{
			name:   "success export expeditions",
			search: "",
			filter: expeditionDto.ReqExpeditionIndexFilter{},
			setupMock: func(m *MockExpeditionRepository) {
				expeditions := []expeditionDto.ExpeditionExport{
					{
						ExpeditionCode: "01",
						ExpeditionName: "JNE",
						Address:        "Jl. Ahmad Yani No. 123",
						PhoneNumbers:   []string{"081234567890", "081234567891"},
						TelpNumbers:    []string{"021-1234567"},
						UpdatedAt:      time.Now(),
					},
					{
						ExpeditionCode: "02",
						ExpeditionName: "TIKI",
						Address:        "Jl. Sudirman No. 456",
						PhoneNumbers:   []string{"081234567892"},
						TelpNumbers:    []string{"021-1234568"},
						UpdatedAt:      time.Now(),
					},
				}
				m.On("GetAllForExport", ctx, expeditionDto.ReqExpeditionIndexFilter{}).Return(expeditions, nil).Once()
			},
			expectedError: nil,
			expectedResult: []models.Expedition{
				{ExpeditionCode: "01", ExpeditionName: "JNE"},
				{ExpeditionCode: "02", ExpeditionName: "TIKI"},
			},
		},
		{
			name:   "error when repository fails",
			search: "",
			filter: expeditionDto.ReqExpeditionIndexFilter{},
			setupMock: func(m *MockExpeditionRepository) {
				m.On("GetAllForExport", ctx, expeditionDto.ReqExpeditionIndexFilter{}).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  errors.New("database error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockExpeditionRepository)
			tt.setupMock(mockRepo)

			usecaseInstance := usecase.NewExpeditionUsecase(mockRepo)

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
