package usecase

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/models"
	roleDto "github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReAssignPermissionByGroup(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecase, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	validRoleID := uuid.New()
	validRoleIDString := validRoleID.String()
	validPermissionGroupID := uuid.New()

	validReq := &roleDto.ReqUpdatePermissionGroupAssignmentToRole{
		PermissionGroupIds: []uuid.UUID{validPermissionGroupID},
	}

	expectedRole := &models.Role{
		ID:   validRoleID,
		Name: "Test Role",
	}

	tests := []struct {
		name           string
		roleId         string
		req            *roleDto.ReqUpdatePermissionGroupAssignmentToRole
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:   "Positive case - successful reassign permission group",
			roleId: validRoleIDString,
			req:    validReq,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, mock.Anything).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
				mockRoleRepo.On("ReAssignPermissionGroup", ctx, validRoleID, mock.Anything).Return(nil).Once()
				mockRoleRepo.On("GetRoleByID", ctx, validRoleID).Return(expectedRole, nil).Once()
			},
			expectedError: false,
			description:   "Valid request should reassign permission group successfully",
		},
		{
			name:   "Negative case - invalid role UUID",
			roleId: "invalid-uuid",
			req:    validReq,
			setupMock: func() {
				// GetPermissionGroupByID will be called first, then UUID parsing will fail
				mockRoleRepo.On("GetPermissionGroupByID", ctx, mock.Anything).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "Invalid UUID should return error",
		},
		{
			name:   "Negative case - invalid permission group",
			roleId: validRoleIDString,
			req:    validReq,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, mock.Anything).Return(nil, errors.New("permission group not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Permission Group with ID",
			description:    "Invalid permission group should return error",
		},
		{
			name:   "Negative case - database error on ReAssignPermissionGroup",
			roleId: validRoleIDString,
			req:    validReq,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, mock.Anything).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
				mockRoleRepo.On("ReAssignPermissionGroup", ctx, validRoleID, mock.Anything).Return(errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
		{
			name:   "Negative-Positive case - SQL injection attempt in role ID",
			roleId: "'; DROP TABLE roles; --",
			req:    validReq,
			setupMock: func() {
				// GetPermissionGroupByID will be called first, then UUID parsing will fail
				mockRoleRepo.On("GetPermissionGroupByID", ctx, mock.Anything).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "SQL injection attempt should fail UUID validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo.ExpectedCalls = nil
			mockRoleRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodPut, "/", nil), httptest.NewRecorder())

			result, err := usecase.ReAssignPermissionByGroup(c, tt.roleId, tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}

func TestAssignUsersToRole(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecase, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	validRoleID := uuid.New()
	validRoleIDString := validRoleID.String()
	validUserID := uuid.New()

	validReq := &roleDto.ReqUpdateAssignUsersToRole{
		UserIds: []uuid.UUID{validUserID},
	}

	expectedRole := &models.Role{
		ID:   validRoleID,
		Name: "Test Role",
	}

	testUser := &models.User{
		ID:       validUserID,
		FullName: "Test User",
	}

	tests := []struct {
		name           string
		roleId         string
		req            *roleDto.ReqUpdateAssignUsersToRole
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:   "Positive case - successful assign users to role",
			roleId: validRoleIDString,
			req:    validReq,
			setupMock: func() {
				mockRoleRepo.On("GetUserByID", ctx, mock.Anything).Return(testUser, nil).Once()
				mockRoleRepo.On("GetRoleByID", ctx, validRoleID).Return(expectedRole, nil).Once()
				mockRoleRepo.On("AssignUsers", ctx, validRoleID, mock.Anything).Return(nil).Once()
				mockRoleRepo.On("GetRoleByID", ctx, validRoleID).Return(expectedRole, nil).Once()
			},
			expectedError: false,
			description:   "Valid request should assign users to role successfully",
		},
		{
			name:   "Negative case - invalid role UUID",
			roleId: "invalid-uuid",
			req:    validReq,
			setupMock: func() {
				// GetUserByID will be called first, then UUID parsing will fail
				mockRoleRepo.On("GetUserByID", ctx, mock.Anything).Return(testUser, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "Invalid UUID should return error",
		},
		{
			name:   "Negative case - invalid user",
			roleId: validRoleIDString,
			req:    validReq,
			setupMock: func() {
				mockRoleRepo.On("GetUserByID", ctx, mock.Anything).Return(nil, errors.New("user not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "User with ID",
			description:    "Invalid user should return error",
		},
		{
			name:   "Negative case - role not found",
			roleId: validRoleIDString,
			req:    validReq,
			setupMock: func() {
				mockRoleRepo.On("GetUserByID", ctx, mock.Anything).Return(testUser, nil).Once()
				mockRoleRepo.On("GetRoleByID", ctx, validRoleID).Return(nil, errors.New("role not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Role with ID",
			description:    "Non-existent role should return error",
		},
		{
			name:   "Negative case - database error on AssignUsers",
			roleId: validRoleIDString,
			req:    validReq,
			setupMock: func() {
				mockRoleRepo.On("GetUserByID", ctx, mock.Anything).Return(testUser, nil).Once()
				mockRoleRepo.On("GetRoleByID", ctx, validRoleID).Return(expectedRole, nil).Once()
				mockRoleRepo.On("AssignUsers", ctx, validRoleID, mock.Anything).Return(errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Something went wrong when assigning users to role",
			description:    "Database error should be returned",
		},
		{
			name:   "Negative-Positive case - SQL injection attempt in role ID",
			roleId: "'; DROP TABLE roles; --",
			req:    validReq,
			setupMock: func() {
				// GetUserByID will be called first, then UUID parsing will fail
				mockRoleRepo.On("GetUserByID", ctx, mock.Anything).Return(testUser, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "SQL injection attempt should fail UUID validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo.ExpectedCalls = nil
			mockRoleRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodPut, "/", nil), httptest.NewRecorder())

			result, err := usecase.AssignUsersToRole(c, tt.roleId, tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}
