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

func TestGetPermissionGroupByID(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()
	expectedPermissionGroup := &models.PermissionGroup{
		ID:   validID,
		Name: "Test Permission Group",
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
			name: "Positive case - successful get permission group by ID",
			id:   validIDString,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validID).Return(expectedPermissionGroup, nil).Once()
			},
			expectedError: false,
			description:   "Valid UUID should return permission group",
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
			name: "Negative case - permission group not found",
			id:   validIDString,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validID).Return(nil, errors.New("permission group not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "permission group not found",
			description:    "Non-existent permission group should return error",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in ID",
			id:   "'; DROP TABLE permission_groups; --",
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
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validID).Return(nil, errors.New("database connection error")).Once()
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

			result, err := usecaseInstance.GetPermissionGroupByID(c.Request().Context(), tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, expectedPermissionGroup.ID, result.ID)
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}

func TestGetIndexPermissionGroup(t *testing.T) {
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

	expectedPermissionGroups := []models.PermissionGroup{
		{ID: uuid.New(), Name: "Permission Group 1"},
		{ID: uuid.New(), Name: "Permission Group 2"},
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
				mockRoleRepo.On("GetIndexPermissionGroup", ctx, validPageReq).Return(expectedPermissionGroups, 2, nil).Once()
			},
			expectedError: false,
			description:   "Valid page request should return permission groups",
		},
		{
			name: "Negative case - database error",
			req:  validPageReq,
			setupMock: func() {
				mockRoleRepo.On("GetIndexPermissionGroup", ctx, validPageReq).Return([]models.PermissionGroup(nil), 0, errors.New("database error")).Once()
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
				Search:    "'; DROP TABLE permission_groups; --",
				SortBy:    "name",
				SortOrder: "ASC",
			},
			setupMock: func() {
				// SQL injection should be treated as normal search string
				injectedReq := request.PageRequest{
					Page:      1,
					PerPage:   10,
					Search:    "'; DROP TABLE permission_groups; --",
					SortBy:    "name",
					SortOrder: "ASC",
				}
				mockRoleRepo.On("GetIndexPermissionGroup", ctx, injectedReq).Return(expectedPermissionGroups, 2, nil).Once()
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
				SortBy:    "'; DROP TABLE permission_groups; --",
				SortOrder: "ASC",
			},
			setupMock: func() {
				// Should be sanitized by ValidateAndSanitizeSortColumn
				injectedReq := request.PageRequest{
					Page:      1,
					PerPage:   10,
					Search:    "",
					SortBy:    "'; DROP TABLE permission_groups; --",
					SortOrder: "ASC",
				}
				mockRoleRepo.On("GetIndexPermissionGroup", ctx, injectedReq).Return(expectedPermissionGroups, 2, nil).Once()
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

			permissionGroups, total, err := usecaseInstance.GetIndexPermissionGroup(c.Request().Context(), tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, permissionGroups)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, permissionGroups)
				assert.GreaterOrEqual(t, total, len(permissionGroups))
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllPermissionGroup(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	expectedPermissionGroups := []models.PermissionGroup{
		{ID: uuid.New(), Name: "Permission Group 1"},
		{ID: uuid.New(), Name: "Permission Group 2"},
	}

	tests := []struct {
		name           string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name: "Positive case - successful get all permission groups",
			setupMock: func() {
				mockRoleRepo.On("GetAllPermissionGroup", ctx).Return(expectedPermissionGroups, nil).Once()
			},
			expectedError: false,
			description:   "Should return all permission groups successfully",
		},
		{
			name: "Negative case - database error",
			setupMock: func() {
				mockRoleRepo.On("GetAllPermissionGroup", ctx).Return(nil, errors.New("database error")).Once()
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

			permissionGroups, err := usecaseInstance.GetAllPermissionGroup(c.Request().Context())

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, permissionGroups)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, permissionGroups)
				assert.Equal(t, len(expectedPermissionGroups), len(permissionGroups))
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}
