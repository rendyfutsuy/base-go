package test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	authDto "github.com/rendyfutsuy/base-go/modules/auth/dto"
	roleDto "github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/rendyfutsuy/base-go/modules/user_management"
	userDto "github.com/rendyfutsuy/base-go/modules/user_management/dto"
	"github.com/rendyfutsuy/base-go/modules/user_management/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of user_management.Repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, userReq userDto.ToDBCreateUser) (userRes *models.User, err error) {
	args := m.Called(ctx, userReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (user *models.User, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAllUser(ctx context.Context) (users []models.User, err error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetIndexUser(ctx context.Context, req request.PageRequest, filter userDto.ReqUserIndexFilter) (users []models.User, total int, err error) {
	args := m.Called(ctx, req, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.User), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, id uuid.UUID, userReq userDto.ToDBUpdateUser) (userRes *models.User, err error) {
	args := m.Called(ctx, id, userReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) SoftDeleteUser(ctx context.Context, id uuid.UUID, userReq userDto.ToDBDeleteUser) (userRes *models.User, err error) {
	args := m.Called(ctx, id, userReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UserNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetDuplicatedUser(ctx context.Context, name string, excludedId uuid.UUID) (user *models.User, err error) {
	args := m.Called(ctx, name, excludedId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UserNameIsNotDuplicatedOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) GetDuplicatedUserOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (user *models.User, err error) {
	args := m.Called(ctx, name, excludedId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) BlockUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UnBlockUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) ActivateUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) DisActivateUser(ctx context.Context, id uuid.UUID) (userRes *models.User, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) EmailIsNotDuplicated(ctx context.Context, email string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, email, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UsernameIsNotDuplicated(ctx context.Context, username string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, username, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) NikIsNotDuplicated(ctx context.Context, nik string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, nik, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CheckBatchDuplication(ctx context.Context, emails, usernames, niks []string) (duplicatedEmails, duplicatedUsernames, duplicatedNiks map[string]bool, err error) {
	args := m.Called(ctx, emails, usernames, niks)
	if args.Get(0) == nil {
		return nil, nil, nil, args.Error(3)
	}
	return args.Get(0).(map[string]bool), args.Get(1).(map[string]bool), args.Get(2).(map[string]bool), args.Error(3)
}

func (m *MockUserRepository) BulkCreateUsers(ctx context.Context, usersReq []userDto.ToDBCreateUser) (err error) {
	args := m.Called(ctx, usersReq)
	return args.Error(0)
}

func (m *MockUserRepository) CountUser(ctx context.Context) (count *int, err error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockUserRepository) IsUserPasswordCanUpdated(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) CreateTable(sqlFilePath string) (err error) {
	args := m.Called(sqlFilePath)
	return args.Error(0)
}

// MockAuthRepository is a mock implementation of auth.Repository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) UpdatePasswordById(ctx context.Context, hashedPassword string, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, hashedPassword, userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) DestroyAllToken(ctx context.Context, userId uuid.UUID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockAuthRepository) AssertPasswordRight(ctx context.Context, password string, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, password, userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) AssertPasswordNeverUsesByUser(ctx context.Context, newPassword string, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, newPassword, userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) AddPasswordHistory(ctx context.Context, hashedPassword string, userId uuid.UUID) error {
	args := m.Called(ctx, hashedPassword, userId)
	return args.Error(0)
}

func (m *MockAuthRepository) ResetPasswordAttempt(ctx context.Context, userId uuid.UUID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockAuthRepository) FindByEmailOrUsername(ctx context.Context, login string) (models.User, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockAuthRepository) AssertPasswordExpiredIsPassed(ctx context.Context, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) AddUserAccessToken(ctx context.Context, accessToken string, userId uuid.UUID) error {
	args := m.Called(ctx, accessToken, userId)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserByAccessToken(ctx context.Context, accessToken string) (models.User, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockAuthRepository) DestroyToken(ctx context.Context, accessToken string) error {
	args := m.Called(ctx, accessToken)
	return args.Error(0)
}

func (m *MockAuthRepository) FindByCurrentSession(ctx context.Context, accessToken string) (models.User, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockAuthRepository) UpdateProfileById(ctx context.Context, profileChunks authDto.ReqUpdateProfile, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, profileChunks, userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) AssertPasswordAttemptPassed(ctx context.Context, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) RequestResetPassword(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthRepository) GetUserByResetPasswordToken(ctx context.Context, token string) (models.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockAuthRepository) DestroyResetPasswordToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthRepository) IncreasePasswordExpiredAt(ctx context.Context, userId uuid.UUID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockAuthRepository) DestroyAllResetPasswordToken(ctx context.Context, userId uuid.UUID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

// MockRoleRepository is a mock implementation of roleManagement.Repository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetRoleByName(ctx context.Context, name string) (role *models.Role, err error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) CreateTable(sqlFilePath string) (err error) {
	args := m.Called(sqlFilePath)
	return args.Error(0)
}

func (m *MockRoleRepository) CreateRole(ctx context.Context, roleReq roleDto.ToDBCreateRole) (roleRes *models.Role, err error) {
	args := m.Called(ctx, roleReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (role *models.Role, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) GetAllRole(ctx context.Context) (roles []models.Role, err error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRoleRepository) GetIndexRole(ctx context.Context, req request.PageRequest) (roles []models.Role, total int, err error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Role), args.Int(1), args.Error(2)
}

func (m *MockRoleRepository) UpdateRole(ctx context.Context, id uuid.UUID, roleReq roleDto.ToDBUpdateRole) (roleRes *models.Role, err error) {
	args := m.Called(ctx, id, roleReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) SoftDeleteRole(ctx context.Context, id uuid.UUID, roleReq roleDto.ToDBDeleteRole) (roleRes *models.Role, err error) {
	args := m.Called(ctx, id, roleReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) RoleNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleRepository) GetDuplicatedRole(ctx context.Context, name string, excludedId uuid.UUID) (role *models.Role, err error) {
	args := m.Called(ctx, name, excludedId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) RoleNameIsNotDuplicatedOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleRepository) GetDuplicatedRoleOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (role *models.Role, err error) {
	args := m.Called(ctx, name, excludedId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) CountRole(ctx context.Context) (count *int, err error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockRoleRepository) ReAssignPermissionGroup(ctx context.Context, id uuid.UUID, permissionGroupReq roleDto.ToDBUpdatePermissionGroupAssignmentToRole) error {
	args := m.Called(ctx, id, permissionGroupReq)
	return args.Error(0)
}

func (m *MockRoleRepository) GetTotalUser(ctx context.Context, id uuid.UUID) (total int, err error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockRoleRepository) GetPermissionFromRoleId(ctx context.Context, id uuid.UUID) (permissions []models.Permission, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockRoleRepository) GetPermissionGroupFromRoleId(ctx context.Context, id uuid.UUID) (permissionGroups []models.PermissionGroup, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PermissionGroup), args.Error(1)
}

func (m *MockRoleRepository) AssignUsers(ctx context.Context, roleId uuid.UUID, userReq []uuid.UUID) error {
	args := m.Called(ctx, roleId, userReq)
	return args.Error(0)
}

func (m *MockRoleRepository) ReAssignPermissionsToPermissionGroup(ctx context.Context, id uuid.UUID, permissions []uuid.UUID) error {
	args := m.Called(ctx, id, permissions)
	return args.Error(0)
}

func (m *MockRoleRepository) GetUserByID(ctx context.Context, id uuid.UUID) (user *models.User, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRoleRepository) GetPermissionByID(ctx context.Context, id uuid.UUID) (permission *models.Permission, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockRoleRepository) GetAllPermission(ctx context.Context) (permissions []models.Permission, err error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockRoleRepository) GetIndexPermission(ctx context.Context, req request.PageRequest) (permissions []models.Permission, total int, err error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.Permission), args.Int(1), args.Error(2)
}

func (m *MockRoleRepository) PermissionNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleRepository) GetDuplicatedPermission(ctx context.Context, name string, excludedId uuid.UUID) (permission *models.Permission, err error) {
	args := m.Called(ctx, name, excludedId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockRoleRepository) CountPermission(ctx context.Context) (count *int, err error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockRoleRepository) GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (permissionGroup *models.PermissionGroup, err error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PermissionGroup), args.Error(1)
}

func (m *MockRoleRepository) GetAllPermissionGroup(ctx context.Context) (permissionGroups []models.PermissionGroup, err error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PermissionGroup), args.Error(1)
}

func (m *MockRoleRepository) GetIndexPermissionGroup(ctx context.Context, req request.PageRequest) (permissionGroups []models.PermissionGroup, total int, err error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]models.PermissionGroup), args.Int(1), args.Error(2)
}

func (m *MockRoleRepository) PermissionGroupNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleRepository) GetDuplicatedPermissionGroup(ctx context.Context, name string, excludedId uuid.UUID) (permissionGroup *models.PermissionGroup, err error) {
	args := m.Called(ctx, name, excludedId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PermissionGroup), args.Error(1)
}

func (m *MockRoleRepository) CountPermissionGroup(ctx context.Context) (count *int, err error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func createTestUsecase() (user_management.Usecase, *MockUserRepository, *MockAuthRepository, *MockRoleRepository) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	mockRoleRepo := new(MockRoleRepository)
	timeout := 5 * time.Second

	usecaseInstance := usecase.NewTestUserUsecase(mockUserRepo, mockRoleRepo, mockAuthRepo, timeout)

	return usecaseInstance, mockUserRepo, mockAuthRepo, mockRoleRepo
}

func TestCreateUser(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockUserRepo, mockAuthRepo, _ := createTestUsecase()
	ctx := context.Background()

	validRoleID := uuid.New()
	validAuthID := uuid.New().String()
	validCount := 10

	validReq := &userDto.ReqCreateUser{
		FullName:             "Test User",
		RoleId:               validRoleID,
		Email:                "test@example.com",
		IsActive:             true,
		Gender:               "male",
		Password:             "password123",
		PasswordConfirmation: "password123",
	}

	expectedUser := &models.User{
		ID:       uuid.New(),
		FullName: "Test User",
		Email:    "test@example.com",
	}

	tests := []struct {
		name           string
		req            *userDto.ReqCreateUser
		authId         string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:   "Positive case - successful user creation",
			req:    validReq,
			authId: validAuthID,
			setupMock: func() {
				mockUserRepo.On("EmailIsNotDuplicated", ctx, validReq.Email, uuid.Nil).Return(true, nil).Once()
				mockUserRepo.On("CountUser", ctx).Return(&validCount, nil).Once()
				mockUserRepo.On("CreateUser", ctx, mock.Anything).Return(expectedUser, nil).Once()
				mockAuthRepo.On("UpdatePasswordById", ctx, validReq.Password, expectedUser.ID).Return(true, nil).Once()
			},
			expectedError: false,
			description:   "Valid request should create user successfully",
		},
		{
			name:   "Negative case - duplicate email",
			req:    validReq,
			authId: validAuthID,
			setupMock: func() {
				mockUserRepo.On("EmailIsNotDuplicated", ctx, validReq.Email, uuid.Nil).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "User Email already exists",
			description:    "Duplicate email should return error",
		},
		{
			name:   "Negative case - database error on email check",
			req:    validReq,
			authId: validAuthID,
			setupMock: func() {
				mockUserRepo.On("EmailIsNotDuplicated", ctx, validReq.Email, uuid.Nil).Return(false, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in email",
			req: &userDto.ReqCreateUser{
				FullName:             "Test User",
				RoleId:               validRoleID,
				Email:                "'; DROP TABLE users; --",
				IsActive:             true,
				Gender:               "male",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			authId: validAuthID,
			setupMock: func() {
				// Email should be validated and treated as normal string
				mockUserRepo.On("EmailIsNotDuplicated", ctx, "'; DROP TABLE users; --", uuid.Nil).Return(true, nil).Once()
				mockUserRepo.On("CountUser", ctx).Return(&validCount, nil).Once()
				mockUserRepo.On("CreateUser", ctx, mock.Anything).Return(expectedUser, nil).Once()
				mockAuthRepo.On("UpdatePasswordById", ctx, "password123", expectedUser.ID).Return(true, nil).Once()
			},
			expectedError: false,
			description:   "SQL injection attempt in email should be treated as normal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil
			mockAuthRepo.ExpectedCalls = nil
			mockAuthRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.CreateUser(c, tt.req, tt.authId)

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

			mockUserRepo.AssertExpectations(t)
			mockAuthRepo.AssertExpectations(t)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockUserRepo, _, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()
	expectedUser := &models.User{
		ID:       validID,
		FullName: "Test User",
		Email:    "test@example.com",
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
			name: "Positive case - successful get user by ID",
			id:   validIDString,
			setupMock: func() {
				mockUserRepo.On("GetUserByID", ctx, validID).Return(expectedUser, nil).Once()
			},
			expectedError: false,
			description:   "Valid UUID should return user",
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
			name: "Negative case - user not found",
			id:   validIDString,
			setupMock: func() {
				mockUserRepo.On("GetUserByID", ctx, validID).Return(nil, errors.New("user not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "user not found",
			description:    "Non-existent user should return error",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in ID",
			id:   "'; DROP TABLE users; --",
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
				mockUserRepo.On("GetUserByID", ctx, validID).Return(nil, errors.New("database connection error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database connection error",
			description:    "Database error should be returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.GetUserByID(c, tt.id)

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

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestGetIndexUser(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockUserRepo, _, _ := createTestUsecase()
	ctx := context.Background()

	validPageReq := request.PageRequest{
		Page:      1,
		PerPage:   10,
		Search:    "test",
		SortBy:    "name",
		SortOrder: "ASC",
	}
	validFilter := userDto.ReqUserIndexFilter{}

	expectedUsers := []models.User{
		{ID: uuid.New(), FullName: "User 1"},
		{ID: uuid.New(), FullName: "User 2"},
	}

	tests := []struct {
		name           string
		req            request.PageRequest
		filter         userDto.ReqUserIndexFilter
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:   "Positive case - successful get index",
			req:    validPageReq,
			filter: validFilter,
			setupMock: func() {
				mockUserRepo.On("GetIndexUser", ctx, validPageReq, validFilter).Return(expectedUsers, 2, nil).Once()
			},
			expectedError: false,
			description:   "Valid page request should return users",
		},
		{
			name:   "Negative case - database error",
			req:    validPageReq,
			filter: validFilter,
			setupMock: func() {
				mockUserRepo.On("GetIndexUser", ctx, validPageReq, validFilter).Return([]models.User(nil), 0, errors.New("database error")).Once()
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
				Search:    "'; DROP TABLE users; --",
				SortBy:    "name",
				SortOrder: "ASC",
			},
			filter: validFilter,
			setupMock: func() {
				injectedReq := request.PageRequest{
					Page:      1,
					PerPage:   10,
					Search:    "'; DROP TABLE users; --",
					SortBy:    "name",
					SortOrder: "ASC",
				}
				mockUserRepo.On("GetIndexUser", ctx, injectedReq, validFilter).Return(expectedUsers, 2, nil).Once()
			},
			expectedError: false,
			description:   "SQL injection attempt in search should be treated as normal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

			users, total, err := usecaseInstance.GetIndexUser(c, tt.req, tt.filter)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, users)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, users)
				assert.Greater(t, total, 0)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllUser(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockUserRepo, _, _ := createTestUsecase()
	ctx := context.Background()

	expectedUsers := []models.User{
		{ID: uuid.New(), FullName: "User 1"},
		{ID: uuid.New(), FullName: "User 2"},
	}

	tests := []struct {
		name           string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name: "Positive case - successful get all users",
			setupMock: func() {
				mockUserRepo.On("GetAllUser", ctx).Return(expectedUsers, nil).Once()
			},
			expectedError: false,
			description:   "Should return all users successfully",
		},
		{
			name: "Negative case - database error",
			setupMock: func() {
				mockUserRepo.On("GetAllUser", ctx).Return([]models.User(nil), errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

			users, err := usecaseInstance.GetAllUser(c)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, users)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockUserRepo, _, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()
	validRoleID := uuid.New()
	validAuthID := uuid.New().String()

	validReq := &userDto.ReqUpdateUser{
		FullName: "Updated User",
		RoleId:   validRoleID,
		Email:    "updated@example.com",
		IsActive: true,
		Gender:   "female",
	}

	expectedUser := &models.User{
		ID:       validID,
		FullName: "Updated User",
		Email:    "updated@example.com",
	}

	tests := []struct {
		name           string
		id             string
		req            *userDto.ReqUpdateUser
		authId         string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:   "Positive case - successful user update",
			id:     validIDString,
			req:    validReq,
			authId: validAuthID,
			setupMock: func() {
				mockUserRepo.On("EmailIsNotDuplicated", ctx, validReq.Email, validID).Return(true, nil).Once()
				mockUserRepo.On("UpdateUser", ctx, validID, mock.Anything).Return(expectedUser, nil).Once()
			},
			expectedError: false,
			description:   "Valid request should update user successfully",
		},
		{
			name:   "Negative case - invalid UUID",
			id:     "invalid-uuid",
			req:    validReq,
			authId: validAuthID,
			setupMock: func() {
				// No repository call should be made
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "Invalid UUID should return error",
		},
		{
			name:   "Negative case - duplicate email",
			id:     validIDString,
			req:    validReq,
			authId: validAuthID,
			setupMock: func() {
				mockUserRepo.On("EmailIsNotDuplicated", ctx, validReq.Email, validID).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "User Email already exists",
			description:    "Duplicate email should return error",
		},
		{
			name:   "Negative case - database error",
			id:     validIDString,
			req:    validReq,
			authId: validAuthID,
			setupMock: func() {
				mockUserRepo.On("EmailIsNotDuplicated", ctx, validReq.Email, validID).Return(true, nil).Once()
				mockUserRepo.On("UpdateUser", ctx, validID, mock.Anything).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in email",
			id:   validIDString,
			req: &userDto.ReqUpdateUser{
				FullName: "Updated User",
				RoleId:   validRoleID,
				Email:    "'; DROP TABLE users; --",
				IsActive: true,
				Gender:   "female",
			},
			authId: validAuthID,
			setupMock: func() {
				mockUserRepo.On("EmailIsNotDuplicated", ctx, "'; DROP TABLE users; --", validID).Return(true, nil).Once()
				mockUserRepo.On("UpdateUser", ctx, validID, mock.Anything).Return(expectedUser, nil).Once()
			},
			expectedError: false,
			description:   "SQL injection attempt in email should be treated as normal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodPut, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.UpdateUser(c, tt.id, tt.req, tt.authId)

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

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestSoftDeleteUser(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockUserRepo, _, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()
	validAuthID := uuid.New().String()

	existingUser := &models.User{
		ID:       validID,
		FullName: "Test User",
		Email:    "test@example.com",
	}

	deletedUser := &models.User{
		ID:       validID,
		FullName: "Test User",
		Email:    "test@example.com",
	}

	tests := []struct {
		name           string
		id             string
		authId         string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:   "Positive case - successful soft delete",
			id:     validIDString,
			authId: validAuthID,
			setupMock: func() {
				mockUserRepo.On("GetUserByID", ctx, validID).Return(existingUser, nil).Once()
				mockUserRepo.On("SoftDeleteUser", ctx, validID, mock.Anything).Return(deletedUser, nil).Once()
			},
			expectedError: false,
			description:   "User should be deleted successfully",
		},
		{
			name:   "Negative case - user not found",
			id:     validIDString,
			authId: validAuthID,
			setupMock: func() {
				mockUserRepo.On("GetUserByID", ctx, validID).Return(nil, errors.New("user not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "User Not Found",
			description:    "Non-existent user should return error",
		},
		{
			name:   "Negative case - database error on delete",
			id:     validIDString,
			authId: validAuthID,
			setupMock: func() {
				mockUserRepo.On("GetUserByID", ctx, validID).Return(existingUser, nil).Once()
				mockUserRepo.On("SoftDeleteUser", ctx, validID, mock.Anything).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.SoftDeleteUser(c, tt.id, tt.authId)

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

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestBlockUser(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockUserRepo, mockAuthRepo, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()

	validBlockReq := &userDto.ReqBlockUser{
		IsBlock: true,
	}

	validUnblockReq := &userDto.ReqBlockUser{
		IsBlock: false,
	}

	expectedUser := &models.User{
		ID:       validID,
		FullName: "Test User",
		Email:    "test@example.com",
	}

	tests := []struct {
		name           string
		id             string
		req            *userDto.ReqBlockUser
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name: "Positive case - successful block user",
			id:   validIDString,
			req:  validBlockReq,
			setupMock: func() {
				mockUserRepo.On("BlockUser", ctx, validID).Return(expectedUser, nil).Once()
				mockAuthRepo.On("DestroyAllToken", ctx, validID).Return(nil).Once()
				mockUserRepo.On("GetUserByID", ctx, validID).Return(expectedUser, nil).Once()
			},
			expectedError: false,
			description:   "User should be blocked successfully",
		},
		{
			name: "Positive case - successful unblock user",
			id:   validIDString,
			req:  validUnblockReq,
			setupMock: func() {
				mockUserRepo.On("UnBlockUser", ctx, validID).Return(expectedUser, nil).Once()
				mockAuthRepo.On("DestroyAllToken", ctx, validID).Return(nil).Once()
				mockUserRepo.On("GetUserByID", ctx, validID).Return(expectedUser, nil).Once()
			},
			expectedError: false,
			description:   "User should be unblocked successfully",
		},
		{
			name: "Negative case - invalid UUID",
			id:   "invalid-uuid",
			req:  validBlockReq,
			setupMock: func() {
				// No repository call should be made
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "Invalid UUID should return error",
		},
		{
			name: "Negative case - database error on block",
			id:   validIDString,
			req:  validBlockReq,
			setupMock: func() {
				mockUserRepo.On("BlockUser", ctx, validID).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in ID",
			id:   "'; DROP TABLE users; --",
			req:  validBlockReq,
			setupMock: func() {
				// Should fail UUID parsing, preventing SQL injection
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "SQL injection attempt should fail UUID validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil
			mockAuthRepo.ExpectedCalls = nil
			mockAuthRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodPut, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.BlockUser(c, tt.id, tt.req)

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

			mockUserRepo.AssertExpectations(t)
			mockAuthRepo.AssertExpectations(t)
		})
	}
}

func TestActivateUser(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockUserRepo, mockAuthRepo, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()

	validActivateReq := &userDto.ReqActivateUser{
		IsActive: true,
	}

	validDeactivateReq := &userDto.ReqActivateUser{
		IsActive: false,
	}

	expectedUser := &models.User{
		ID:       validID,
		FullName: "Test User",
		Email:    "test@example.com",
	}

	tests := []struct {
		name           string
		id             string
		req            *userDto.ReqActivateUser
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name: "Positive case - successful activate user",
			id:   validIDString,
			req:  validActivateReq,
			setupMock: func() {
				mockUserRepo.On("ActivateUser", ctx, validID).Return(expectedUser, nil).Once()
				mockAuthRepo.On("DestroyAllToken", ctx, validID).Return(nil).Once()
				mockUserRepo.On("GetUserByID", ctx, validID).Return(expectedUser, nil).Once()
			},
			expectedError: false,
			description:   "User should be activated successfully",
		},
		{
			name: "Positive case - successful deactivate user",
			id:   validIDString,
			req:  validDeactivateReq,
			setupMock: func() {
				mockUserRepo.On("DisActivateUser", ctx, validID).Return(expectedUser, nil).Once()
				mockAuthRepo.On("DestroyAllToken", ctx, validID).Return(nil).Once()
				mockUserRepo.On("GetUserByID", ctx, validID).Return(expectedUser, nil).Once()
			},
			expectedError: false,
			description:   "User should be deactivated successfully",
		},
		{
			name: "Negative case - invalid UUID",
			id:   "invalid-uuid",
			req:  validActivateReq,
			setupMock: func() {
				// No repository call should be made
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "Invalid UUID should return error",
		},
		{
			name: "Negative case - database error on activate",
			id:   validIDString,
			req:  validActivateReq,
			setupMock: func() {
				mockUserRepo.On("ActivateUser", ctx, validID).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in ID",
			id:   "'; DROP TABLE users; --",
			req:  validActivateReq,
			setupMock: func() {
				// Should fail UUID parsing, preventing SQL injection
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "SQL injection attempt should fail UUID validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil
			mockAuthRepo.ExpectedCalls = nil
			mockAuthRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodPut, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.ActivateUser(c, tt.id, tt.req)

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

			mockUserRepo.AssertExpectations(t)
			mockAuthRepo.AssertExpectations(t)
		})
	}
}

func TestUserNameIsNotDuplicated(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockUserRepo, _, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validName := "Test User"

	expectedUser := &models.User{
		ID:       validID,
		FullName: validName,
	}

	tests := []struct {
		name           string
		fullName       string
		id             uuid.UUID
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:     "Positive case - user name not duplicated",
			fullName: validName,
			id:       uuid.Nil,
			setupMock: func() {
				mockUserRepo.On("GetDuplicatedUser", ctx, validName, uuid.Nil).Return(nil, errors.New("not found")).Once()
			},
			expectedError: true,
			description:   "Should return error when user not found (name not duplicated)",
		},
		{
			name:     "Negative case - user name duplicated",
			fullName: validName,
			id:       uuid.Nil,
			setupMock: func() {
				mockUserRepo.On("GetDuplicatedUser", ctx, validName, uuid.Nil).Return(expectedUser, nil).Once()
			},
			expectedError: false,
			description:   "Should return user when name is duplicated",
		},
		{
			name:     "Negative case - database error",
			fullName: validName,
			id:       uuid.Nil,
			setupMock: func() {
				mockUserRepo.On("GetDuplicatedUser", ctx, validName, uuid.Nil).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
		{
			name:     "Negative-Positive case - SQL injection attempt in name",
			fullName: "'; DROP TABLE users; --",
			id:       uuid.Nil,
			setupMock: func() {
				mockUserRepo.On("GetDuplicatedUser", ctx, "'; DROP TABLE users; --", uuid.Nil).Return(nil, errors.New("not found")).Once()
			},
			expectedError: true,
			description:   "SQL injection attempt in name should be treated as normal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.ExpectedCalls = nil
			mockUserRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.UserNameIsNotDuplicated(c, tt.fullName, tt.id)

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

			mockUserRepo.AssertExpectations(t)
		})
	}
}
