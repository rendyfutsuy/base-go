package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
	"github.com/rendyfutsuy/base-go/modules/auth/usecase"
	roleManagementDto "github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRoleManagementRepository is a mock implementation of roleManagement.Repository
type MockRoleManagementRepository struct {
	mock.Mock
}

func (m *MockRoleManagementRepository) CreateTable(sqlFilePath string) error {
	args := m.Called(sqlFilePath)
	return args.Error(0)
}

func (m *MockRoleManagementRepository) CreateRole(ctx context.Context, roleReq roleManagementDto.ToDBCreateRole) (*models.Role, error) {
	args := m.Called(ctx, roleReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleManagementRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleManagementRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleManagementRepository) GetAllRole(ctx context.Context) ([]models.Role, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Role), args.Error(1)
}

func (m *MockRoleManagementRepository) GetIndexRole(ctx context.Context, req request.PageRequest) ([]models.Role, int, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.Role), args.Int(1), args.Error(2)
}

func (m *MockRoleManagementRepository) UpdateRole(ctx context.Context, id uuid.UUID, roleReq roleManagementDto.ToDBUpdateRole) (*models.Role, error) {
	args := m.Called(ctx, id, roleReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleManagementRepository) SoftDeleteRole(ctx context.Context, id uuid.UUID, roleReq roleManagementDto.ToDBDeleteRole) (*models.Role, error) {
	args := m.Called(ctx, id, roleReq)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleManagementRepository) RoleNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleManagementRepository) GetDuplicatedRole(ctx context.Context, name string, excludedId uuid.UUID) (*models.Role, error) {
	args := m.Called(ctx, name, excludedId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleManagementRepository) RoleNameIsNotDuplicatedOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleManagementRepository) GetDuplicatedRoleOnSoftDeleted(ctx context.Context, name string, excludedId uuid.UUID) (*models.Role, error) {
	args := m.Called(ctx, name, excludedId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleManagementRepository) CountRole(ctx context.Context) (*int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockRoleManagementRepository) ReAssignPermissionGroup(ctx context.Context, id uuid.UUID, permissionGroupReq roleManagementDto.ToDBUpdatePermissionGroupAssignmentToRole) error {
	args := m.Called(ctx, id, permissionGroupReq)
	return args.Error(0)
}

func (m *MockRoleManagementRepository) GetTotalUser(ctx context.Context, id uuid.UUID) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockRoleManagementRepository) GetPermissionFromRoleId(ctx context.Context, id uuid.UUID) ([]models.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return []models.Permission{}, args.Error(1)
	}
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockRoleManagementRepository) GetPermissionGroupFromRoleId(ctx context.Context, id uuid.UUID) ([]models.PermissionGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return []models.PermissionGroup{}, args.Error(1)
	}
	return args.Get(0).([]models.PermissionGroup), args.Error(1)
}

func (m *MockRoleManagementRepository) AssignUsers(ctx context.Context, roleId uuid.UUID, userReq []uuid.UUID) error {
	args := m.Called(ctx, roleId, userReq)
	return args.Error(0)
}

func (m *MockRoleManagementRepository) ReAssignPermissionsToPermissionGroup(ctx context.Context, id uuid.UUID, permissions []uuid.UUID) error {
	args := m.Called(ctx, id, permissions)
	return args.Error(0)
}

func (m *MockRoleManagementRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRoleManagementRepository) GetPermissionByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockRoleManagementRepository) GetAllPermission(ctx context.Context) ([]models.Permission, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockRoleManagementRepository) GetIndexPermission(ctx context.Context, req request.PageRequest) ([]models.Permission, int, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.Permission), args.Int(1), args.Error(2)
}

func (m *MockRoleManagementRepository) PermissionNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleManagementRepository) GetDuplicatedPermission(ctx context.Context, name string, excludedId uuid.UUID) (*models.Permission, error) {
	args := m.Called(ctx, name, excludedId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Permission), args.Error(1)
}

func (m *MockRoleManagementRepository) CountPermission(ctx context.Context) (*int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func (m *MockRoleManagementRepository) GetPermissionGroupByID(ctx context.Context, id uuid.UUID) (*models.PermissionGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PermissionGroup), args.Error(1)
}

func (m *MockRoleManagementRepository) GetAllPermissionGroup(ctx context.Context) ([]models.PermissionGroup, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.PermissionGroup), args.Error(1)
}

func (m *MockRoleManagementRepository) GetIndexPermissionGroup(ctx context.Context, req request.PageRequest) ([]models.PermissionGroup, int, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.PermissionGroup), args.Int(1), args.Error(2)
}

func (m *MockRoleManagementRepository) PermissionGroupNameIsNotDuplicated(ctx context.Context, name string, excludedId uuid.UUID) (bool, error) {
	args := m.Called(ctx, name, excludedId)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleManagementRepository) GetDuplicatedPermissionGroup(ctx context.Context, name string, excludedId uuid.UUID) (*models.PermissionGroup, error) {
	args := m.Called(ctx, name, excludedId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PermissionGroup), args.Error(1)
}

func (m *MockRoleManagementRepository) CountPermissionGroup(ctx context.Context) (*int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*int), args.Error(1)
}

func TestGetProfile(t *testing.T) {
	// Setup test logger to prevent nil pointer panics
	setupTestLogger()

	ctx := context.Background()
	mockRepo := new(MockAuthRepository)
	mockRoleManagementRepo := new(MockRoleManagementRepository)
	signingKey := []byte("test-secret-key")
	hashSalt := "test-salt"
	timeout := 5 * time.Second

	usecaseInstance := usecase.NewTestAuthUsecase(mockRepo, mockRoleManagementRepo, timeout, hashSalt, signingKey, 24*time.Hour)

	accessToken := "valid-access-token"
	testRoleId := uuid.New()
	expectedUser := models.User{
		ID:       uuid.New(),
		RoleId:   testRoleId,
		FullName: "Test User",
		Email:    "test@example.com",
		RoleName: "Admin",
		IsActive: true,
		Gender:   "Male",
	}

	tests := []struct {
		name           string
		accessToken    string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:        "Positive case - successful get profile",
			accessToken: accessToken,
			setupMock: func() {
				mockRepo.On("FindByCurrentSession", ctx, accessToken).Return(expectedUser, nil).Once()
				mockRoleManagementRepo.On("GetPermissionFromRoleId", ctx, testRoleId).Return([]models.Permission{}, nil).Once()
				mockRoleManagementRepo.On("GetPermissionGroupFromRoleId", ctx, testRoleId).Return([]models.PermissionGroup{}, nil).Once()
			},
			expectedError: false,
			description:   "Valid token should return user profile",
		},
		{
			name:        "Negative case - invalid token",
			accessToken: "invalid-token",
			setupMock: func() {
				mockRepo.On("FindByCurrentSession", ctx, "invalid-token").Return(models.User{}, errors.New(constants.UserInvalid)).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.UserInvalid,
			description:    "Invalid token should return error",
		},
		{
			name:        "Negative case - empty token",
			accessToken: "",
			setupMock: func() {
				mockRepo.On("FindByCurrentSession", ctx, "").Return(models.User{}, errors.New(constants.UserInvalid)).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.UserInvalid,
			description:    "Empty token should return error",
		},
		{
			name:        "Negative-Positive case - SQL injection attempt in token",
			accessToken: "token'; DROP TABLE jwt_tokens; --",
			setupMock: func() {
				// Should not find user even with SQL injection attempt due to parameterized query
				mockRepo.On("FindByCurrentSession", ctx, "token'; DROP TABLE jwt_tokens; --").Return(models.User{}, errors.New(constants.UserInvalid)).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.UserInvalid,
			description:    "SQL injection attempt should not be executed",
		},
		{
			name:        "Negative case - database error",
			accessToken: accessToken,
			setupMock: func() {
				mockRepo.On("FindByCurrentSession", ctx, accessToken).Return(models.User{}, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			mockRoleManagementRepo.ExpectedCalls = nil
			mockRoleManagementRepo.Calls = nil
			tt.setupMock()

			user, err := usecaseInstance.GetProfile(ctx, tt.accessToken)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Equal(t, uuid.Nil, user.ID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedUser.ID, user.ID)
				assert.Equal(t, expectedUser.FullName, user.FullName)
				assert.Equal(t, expectedUser.Email, user.Email)
				assert.Equal(t, expectedUser.RoleName, user.RoleName)
			}

			mockRepo.AssertExpectations(t)
			mockRoleManagementRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	// Setup test logger to prevent nil pointer panics
	setupTestLogger()

	ctx := context.Background()
	mockRepo := new(MockAuthRepository)
	signingKey := []byte("test-secret-key")
	hashSalt := "test-salt"
	timeout := 5 * time.Second

	usecaseInstance := usecase.NewTestAuthUsecase(mockRepo, nil, timeout, hashSalt, signingKey, 24*time.Hour)

	validUserId := uuid.New().String()

	tests := []struct {
		name           string
		userId         string
		profileChunks  dto.ReqUpdateProfile
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:   "Positive case - successful update profile",
			userId: validUserId,
			profileChunks: dto.ReqUpdateProfile{
				Name: "Updated Name",
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("UpdateProfileById", ctx, dto.ReqUpdateProfile{Name: "Updated Name"}, parsedUUID).Return(true, nil).Once()
			},
			expectedError: false,
			description:   "Valid profile update should succeed",
		},
		{
			name:   "Negative case - invalid UUID",
			userId: "invalid-uuid",
			profileChunks: dto.ReqUpdateProfile{
				Name: "Updated Name",
			},
			setupMock:      func() {},
			expectedError:  true,
			expectedErrMsg: "invalid UUID",
			description:    "Invalid UUID should return error",
		},
		{
			name:   "Negative case - empty name",
			userId: validUserId,
			profileChunks: dto.ReqUpdateProfile{
				Name: "",
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("UpdateProfileById", ctx, dto.ReqUpdateProfile{Name: ""}, parsedUUID).Return(true, nil).Once()
			},
			expectedError: false,
			description:   "Empty name is allowed (validation should be at handler level)",
		},
		{
			name:   "Negative case - database error",
			userId: validUserId,
			profileChunks: dto.ReqUpdateProfile{
				Name: "Updated Name",
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("UpdateProfileById", ctx, dto.ReqUpdateProfile{Name: "Updated Name"}, parsedUUID).Return(false, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
		{
			name:   "Negative-Positive case - SQL injection attempt in name",
			userId: validUserId,
			profileChunks: dto.ReqUpdateProfile{
				Name: "'; DROP TABLE users; --",
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				// Should update with the literal string, not execute SQL
				mockRepo.On("UpdateProfileById", ctx, dto.ReqUpdateProfile{Name: "'; DROP TABLE users; --"}, parsedUUID).Return(true, nil).Once()
			},
			expectedError: false,
			description:   "SQL injection attempt should be treated as literal string, not executed",
		},
		{
			name:   "Negative case - very long name",
			userId: validUserId,
			profileChunks: dto.ReqUpdateProfile{
				Name: string(make([]byte, 1000)), // Very long name
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("UpdateProfileById", ctx, mock.AnythingOfType("dto.ReqUpdateProfile"), parsedUUID).Return(false, errors.New("value too long")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "value too long",
			description:    "Very long name should be handled by database constraint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.setupMock()

			err := usecaseInstance.UpdateProfile(ctx, tt.profileChunks, tt.userId)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateMyPassword(t *testing.T) {
	// Setup test logger to prevent nil pointer panics
	setupTestLogger()

	ctx := context.Background()
	mockRepo := new(MockAuthRepository)
	signingKey := []byte("test-secret-key")
	hashSalt := "test-salt"
	timeout := 5 * time.Second

	usecaseInstance := usecase.NewTestAuthUsecase(mockRepo, nil, timeout, hashSalt, signingKey, 24*time.Hour)

	validUserId := uuid.New().String()
	oldPassword := "oldpassword123"
	newPassword := "newpassword123"

	tests := []struct {
		name           string
		userId         string
		passwordChunks dto.ReqUpdatePassword
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:   "Positive case - successful password update",
			userId: validUserId,
			passwordChunks: dto.ReqUpdatePassword{
				OldPassword:          oldPassword,
				NewPassword:          newPassword,
				PasswordConfirmation: newPassword,
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("GetIsFirstTimeLogin", ctx, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, oldPassword, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, newPassword, parsedUUID).Return(false, errors.New("Password Not Match")).Once()
				mockRepo.On("AssertPasswordNeverUsesByUser", ctx, newPassword, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AddPasswordHistory", ctx, mock.AnythingOfType("string"), parsedUUID).Return(nil).Once()
				mockRepo.On("ResetPasswordAttempt", ctx, parsedUUID).Return(nil).Once()
				mockRepo.On("UpdatePasswordById", ctx, newPassword, parsedUUID).Return(true, nil).Once()
				mockRepo.On("DestroyAllToken", ctx, parsedUUID).Return(nil).Once()
			},
			expectedError: false,
			description:   "Valid password update should succeed",
		},
		{
			name:   "Negative case - user already changed password",
			userId: validUserId,
			passwordChunks: dto.ReqUpdatePassword{
				OldPassword:          oldPassword,
				NewPassword:          newPassword,
				PasswordConfirmation: newPassword,
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("GetIsFirstTimeLogin", ctx, parsedUUID).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.AuthPasswordAlreadyChanged,
			description:    "User already changed password should return error",
		},
		{
			name:   "Negative case - invalid UUID",
			userId: "invalid-uuid",
			passwordChunks: dto.ReqUpdatePassword{
				OldPassword:          oldPassword,
				NewPassword:          newPassword,
				PasswordConfirmation: newPassword,
			},
			setupMock:      func() {},
			expectedError:  true,
			expectedErrMsg: "invalid UUID",
			description:    "Invalid UUID should return error",
		},
		{
			name:   "Negative case - old password not match",
			userId: validUserId,
			passwordChunks: dto.ReqUpdatePassword{
				OldPassword:          "wrongoldpassword",
				NewPassword:          newPassword,
				PasswordConfirmation: newPassword,
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("GetIsFirstTimeLogin", ctx, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, "wrongoldpassword", parsedUUID).Return(false, errors.New("Password Not Match")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Old Password not Match",
			description:    "Wrong old password should return error",
		},
		{
			name:   "Negative case - new password same as old password",
			userId: validUserId,
			passwordChunks: dto.ReqUpdatePassword{
				OldPassword:          oldPassword,
				NewPassword:          oldPassword,
				PasswordConfirmation: oldPassword,
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("GetIsFirstTimeLogin", ctx, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, oldPassword, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, oldPassword, parsedUUID).Return(true, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "New Password should not be same with Current Password",
			description:    "New password same as old should return error",
		},
		{
			name:   "Negative case - new password already used",
			userId: validUserId,
			passwordChunks: dto.ReqUpdatePassword{
				OldPassword:          oldPassword,
				NewPassword:          newPassword,
				PasswordConfirmation: newPassword,
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("GetIsFirstTimeLogin", ctx, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, oldPassword, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, newPassword, parsedUUID).Return(false, errors.New("Password Not Match")).Once()
				mockRepo.On("AssertPasswordNeverUsesByUser", ctx, newPassword, parsedUUID).Return(false, errors.New("Youre already used this password")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Youre already used this password",
			description:    "Password already used should return error",
		},
		{
			name:   "Negative case - empty old password",
			userId: validUserId,
			passwordChunks: dto.ReqUpdatePassword{
				OldPassword:          "",
				NewPassword:          newPassword,
				PasswordConfirmation: newPassword,
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("GetIsFirstTimeLogin", ctx, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, "", parsedUUID).Return(false, errors.New("Password Not Match")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Old Password not Match",
			description:    "Empty old password should return error",
		},
		{
			name:   "Negative case - empty new password",
			userId: validUserId,
			passwordChunks: dto.ReqUpdatePassword{
				OldPassword:          oldPassword,
				NewPassword:          "",
				PasswordConfirmation: "",
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("GetIsFirstTimeLogin", ctx, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, oldPassword, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, "", parsedUUID).Return(false, errors.New("Password Not Match")).Once()
				mockRepo.On("AssertPasswordNeverUsesByUser", ctx, "", parsedUUID).Return(true, nil).Once()
				mockRepo.On("AddPasswordHistory", ctx, mock.AnythingOfType("string"), parsedUUID).Return(nil).Once()
				mockRepo.On("ResetPasswordAttempt", ctx, parsedUUID).Return(nil).Once()
				mockRepo.On("UpdatePasswordById", ctx, "", parsedUUID).Return(true, nil).Once()
				mockRepo.On("DestroyAllToken", ctx, parsedUUID).Return(nil).Once()
			},
			expectedError: false,
			description:   "Empty new password is technically valid but should be validated at handler level",
		},
		{
			name:   "Negative-Positive case - SQL injection attempt in old password",
			userId: validUserId,
			passwordChunks: dto.ReqUpdatePassword{
				OldPassword:          "'; DROP TABLE users; --",
				NewPassword:          newPassword,
				PasswordConfirmation: newPassword,
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("GetIsFirstTimeLogin", ctx, parsedUUID).Return(true, nil).Once()
				// Should not match, SQL injection should not execute
				mockRepo.On("AssertPasswordRight", ctx, "'; DROP TABLE users; --", parsedUUID).Return(false, errors.New("Password Not Match")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Old Password not Match",
			description:    "SQL injection in old password should be treated as wrong password",
		},
		{
			name:   "Negative-Positive case - SQL injection attempt in new password",
			userId: validUserId,
			passwordChunks: dto.ReqUpdatePassword{
				OldPassword:          oldPassword,
				NewPassword:          "'; DROP TABLE users; --",
				PasswordConfirmation: "'; DROP TABLE users; --",
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("GetIsFirstTimeLogin", ctx, parsedUUID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, oldPassword, parsedUUID).Return(true, nil).Once()
				// New password should not match old password
				mockRepo.On("AssertPasswordRight", ctx, "'; DROP TABLE users; --", parsedUUID).Return(false, errors.New("Password Not Match")).Once()
				// Should check password history with the literal string
				mockRepo.On("AssertPasswordNeverUsesByUser", ctx, "'; DROP TABLE users; --", parsedUUID).Return(true, nil).Once()
				mockRepo.On("AddPasswordHistory", ctx, mock.AnythingOfType("string"), parsedUUID).Return(nil).Once()
				mockRepo.On("ResetPasswordAttempt", ctx, parsedUUID).Return(nil).Once()
				mockRepo.On("UpdatePasswordById", ctx, "'; DROP TABLE users; --", parsedUUID).Return(true, nil).Once()
				mockRepo.On("DestroyAllToken", ctx, parsedUUID).Return(nil).Once()
			},
			expectedError: false,
			description:   "SQL injection in new password should be treated as literal string",
		},
		{
			name:   "Negative case - GetIsFirstTimeLogin returns error",
			userId: validUserId,
			passwordChunks: dto.ReqUpdatePassword{
				OldPassword:          oldPassword,
				NewPassword:          newPassword,
				PasswordConfirmation: newPassword,
			},
			setupMock: func() {
				parsedUUID, _ := uuid.Parse(validUserId)
				mockRepo.On("GetIsFirstTimeLogin", ctx, parsedUUID).Return(false, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error on GetIsFirstTimeLogin should return error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.setupMock()

			err := usecaseInstance.UpdateMyPassword(ctx, tt.passwordChunks, tt.userId)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
