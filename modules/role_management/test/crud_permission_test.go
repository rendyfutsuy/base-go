package test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/stretchr/testify/assert"
)

func TestGetPermissionByID(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()
	expectedPermission := &models.Permission{
		ID:   validID,
		Name: "Test Permission",
	}

	tests := []struct {
		name           string
		id             string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name: "Positive case - successful get permission by ID",
			id:   validIDString,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionByID", ctx, validID).Return(expectedPermission, nil).Once()
			},
			expectedError: false,
			description:   "Valid UUID should return permission",
		},
		{
			name: "Negative case - invalid UUID format",
			id:   "invalid-uuid",
			setupMock: func() {
				// No repository call should be made
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "Invalid UUID format should return error",
		},
		{
			name: "Negative case - empty ID",
			id:   "",
			setupMock: func() {
				// No repository call should be made
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "Empty ID should return error",
		},
		{
			name: "Negative case - permission not found",
			id:   validIDString,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionByID", ctx, validID).Return(nil, errors.New("permission not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "permission not found",
			description:    "Non-existent permission should return error",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in ID",
			id:   "'; DROP TABLE permissions; --",
			setupMock: func() {
				// Should fail UUID parsing, preventing SQL injection
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "SQL injection attempt should fail UUID validation",
		},
		{
			name: "Negative case - database error",
			id:   validIDString,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionByID", ctx, validID).Return(nil, errors.New("database connection error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database connection error",
			description:    "Database error should be returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo.ExpectedCalls = nil
			mockRoleRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.GetPermissionByID(c.Request().Context(), tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, expectedPermission.ID, result.ID)
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}

func TestGetIndexPermission(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	validPageReq := request.PageRequest{
		Page:      1,
		PerPage:   10,
		Search:    "test",
		SortBy:    "name",
		SortOrder: "ASC",
	}

	expectedPermissions := []models.Permission{
		{ID: uuid.New(), Name: "Permission 1"},
		{ID: uuid.New(), Name: "Permission 2"},
	}

	tests := []struct {
		name           string
		req            request.PageRequest
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name: "Positive case - successful get index",
			req:  validPageReq,
			setupMock: func() {
				mockRoleRepo.On("GetIndexPermission", ctx, validPageReq).Return(expectedPermissions, 2, nil).Once()
			},
			expectedError: false,
			description:   "Valid page request should return permissions",
		},
		{
			name: "Negative case - database error",
			req:  validPageReq,
			setupMock: func() {
				mockRoleRepo.On("GetIndexPermission", ctx, validPageReq).Return(nil, 0, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in search",
			req: request.PageRequest{
				Page:      1,
				PerPage:   10,
				Search:    "'; DROP TABLE permissions; --",
				SortBy:    "name",
				SortOrder: "ASC",
			},
			setupMock: func() {
				// SQL injection should be treated as normal search string
				injectedReq := request.PageRequest{
					Page:      1,
					PerPage:   10,
					Search:    "'; DROP TABLE permissions; --",
					SortBy:    "name",
					SortOrder: "ASC",
				}
				mockRoleRepo.On("GetIndexPermission", ctx, injectedReq).Return(expectedPermissions, 2, nil).Once()
			},
			expectedError: false,
			description:   "SQL injection attempt in search should be treated as normal string",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in sort_by",
			req: request.PageRequest{
				Page:      1,
				PerPage:   10,
				Search:    "",
				SortBy:    "'; DROP TABLE permissions; --",
				SortOrder: "ASC",
			},
			setupMock: func() {
				// Should be sanitized by ValidateAndSanitizeSortColumn
				injectedReq := request.PageRequest{
					Page:      1,
					PerPage:   10,
					Search:    "",
					SortBy:    "'; DROP TABLE permissions; --",
					SortOrder: "ASC",
				}
				mockRoleRepo.On("GetIndexPermission", ctx, injectedReq).Return(expectedPermissions, 2, nil).Once()
			},
			expectedError: false,
			description:   "SQL injection attempt in sort_by should be sanitized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo.ExpectedCalls = nil
			mockRoleRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

			permissions, total, err := usecaseInstance.GetIndexPermission(c.Request().Context(), tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, permissions)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, permissions)
				assert.GreaterOrEqual(t, total, len(permissions))
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllPermission(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	expectedPermissions := []models.Permission{
		{ID: uuid.New(), Name: "Permission 1"},
		{ID: uuid.New(), Name: "Permission 2"},
	}

	tests := []struct {
		name           string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name: "Positive case - successful get all permissions",
			setupMock: func() {
				mockRoleRepo.On("GetAllPermission", ctx).Return(expectedPermissions, nil).Once()
			},
			expectedError: false,
			description:   "Should return all permissions successfully",
		},
		{
			name: "Negative case - database error",
			setupMock: func() {
				mockRoleRepo.On("GetAllPermission", ctx).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo.ExpectedCalls = nil
			mockRoleRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

			permissions, err := usecaseInstance.GetAllPermission(c.Request().Context())

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, permissions)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, permissions)
				assert.Equal(t, len(expectedPermissions), len(permissions))
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}
