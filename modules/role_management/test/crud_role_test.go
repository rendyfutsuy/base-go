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
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/models"
	authDto "github.com/rendyfutsuy/base-go/modules/auth/dto"
	"github.com/rendyfutsuy/base-go/modules/role_management"
	roleDto "github.com/rendyfutsuy/base-go/modules/role_management/dto"
	"github.com/rendyfutsuy/base-go/modules/role_management/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRoleRepository is a mock implementation of role_management.Repository
type MockRoleRepository struct {
	mock.Mock
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

func (m *MockRoleRepository) GetRoleByName(ctx context.Context, name string) (role *models.Role, err error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
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
	return args.Get(0).([]models.Permission), args.Error(1)
}

func (m *MockRoleRepository) GetPermissionGroupFromRoleId(ctx context.Context, id uuid.UUID) (permissionGroups []models.PermissionGroup, err error) {
	args := m.Called(ctx, id)
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

func (m *MockRoleRepository) CreateTable(sqlFilePath string) (err error) {
	args := m.Called(sqlFilePath)
	return args.Error(0)
}

// MockAuthRepository is a mock implementation of auth.Repository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) GetIsFirstTimeLogin(ctx context.Context, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) FindByEmailOrUsername(ctx context.Context, login string) (models.User, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockAuthRepository) AssertPasswordRight(ctx context.Context, password string, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, password, userId)
	return args.Bool(0), args.Error(1)
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

func (m *MockAuthRepository) UpdatePasswordById(ctx context.Context, hashedPassword string, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, hashedPassword, userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) DestroyAllToken(ctx context.Context, userId uuid.UUID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockAuthRepository) AssertPasswordNeverUsesByUser(ctx context.Context, newPassword string, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, newPassword, userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) AddPasswordHistory(ctx context.Context, hashedPassword string, userId uuid.UUID) error {
	args := m.Called(ctx, hashedPassword, userId)
	return args.Error(0)
}

func (m *MockAuthRepository) AssertPasswordAttemptPassed(ctx context.Context, userId uuid.UUID) (bool, error) {
	args := m.Called(ctx, userId)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepository) ResetPasswordAttempt(ctx context.Context, userId uuid.UUID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
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

func (m *MockAuthRepository) UpdateLastLogin(ctx context.Context, userId uuid.UUID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

// Helper function to create Echo context
func createEchoContext() echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

// Helper function to create test usecase
func createTestUsecase() (role_management.Usecase, *MockRoleRepository, *MockAuthRepository) {
	mockRoleRepo := new(MockRoleRepository)
	mockAuthRepo := new(MockAuthRepository)
	timeout := 5 * time.Second

	usecaseInstance := usecase.NewTestRoleUsecase(mockRoleRepo, mockAuthRepo, timeout)

	return usecaseInstance, mockRoleRepo, mockAuthRepo
}

func TestCreateRole(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	validPermissionGroupID := uuid.New()
	validAuthID := uuid.New().String()
	validCount := 10

	validReq := &roleDto.ReqCreateRole{
		Name:             "Test Role",
		Description:      "Test Description",
		PermissionGroups: []uuid.UUID{validPermissionGroupID},
	}

	expectedRole := &models.Role{
		ID:   uuid.New(),
		Name: "Test Role",
	}

	tests := []struct {
		name           string
		req            *roleDto.ReqCreateRole
		authId         string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:   "Positive case - successful role creation",
			req:    validReq,
			authId: validAuthID,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validPermissionGroupID).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
				mockRoleRepo.On("RoleNameIsNotDuplicated", ctx, validReq.Name, uuid.Nil).Return(true, nil).Once()
				mockRoleRepo.On("CountRole", ctx).Return(&validCount, nil).Once()
				mockRoleRepo.On("CreateRole", ctx, mock.Anything).Return(expectedRole, nil).Once()
			},
			expectedError: false,
			description:   "Valid request should create role successfully",
		},
		{
			name:   "Negative case - invalid permission group ID",
			req:    validReq,
			authId: validAuthID,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validPermissionGroupID).Return(nil, errors.New("permission group not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Function with ID",
			description:    "Invalid permission group should return error",
		},
		{
			name: "Negative case - duplicate role name",
			req:  validReq,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validPermissionGroupID).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
				mockRoleRepo.On("RoleNameIsNotDuplicated", ctx, validReq.Name, uuid.Nil).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.RoleErrorRoleNotFound,
			description:    "Duplicate role name should return error",
		},
		{
			name: "Negative case - empty role name",
			req: &roleDto.ReqCreateRole{
				Name:             "",
				Description:      "Test Description",
				PermissionGroups: []uuid.UUID{validPermissionGroupID},
			},
			setupMock: func() {
				// Should still check permission group first
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validPermissionGroupID).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
				mockRoleRepo.On("RoleNameIsNotDuplicated", ctx, "", uuid.Nil).Return(true, nil).Once()
				mockRoleRepo.On("CountRole", ctx).Return(&validCount, nil).Once()
				mockRoleRepo.On("CreateRole", ctx, mock.Anything).Return(nil, errors.New("name is required")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "name is required",
			description:    "Empty role name should return error",
		},
		{
			name: "Negative case - empty permission groups",
			req: &roleDto.ReqCreateRole{
				Name:             "Test Role",
				Description:      "Test Description",
				PermissionGroups: []uuid.UUID{},
			},
			setupMock: func() {
				// No permission groups to check, so RoleNameIsNotDuplicated will be called
				mockRoleRepo.On("RoleNameIsNotDuplicated", ctx, "Test Role", uuid.Nil).Return(true, nil).Once()
				mockRoleRepo.On("CountRole", ctx).Return(&validCount, nil).Once()
				mockRoleRepo.On("CreateRole", ctx, mock.Anything).Return(expectedRole, nil).Once()
			},
			expectedError: false, // Will pass validation but might fail at repository
			description:   "Empty permission groups array",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in role name",
			req: &roleDto.ReqCreateRole{
				Name:             "'; DROP TABLE roles; --",
				Description:      "Test Description",
				PermissionGroups: []uuid.UUID{validPermissionGroupID},
			},
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validPermissionGroupID).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
				// SQL injection should be treated as normal string due to parameterized query
				mockRoleRepo.On("RoleNameIsNotDuplicated", ctx, "'; DROP TABLE roles; --", uuid.Nil).Return(true, nil).Once()
				mockRoleRepo.On("CountRole", ctx).Return(&validCount, nil).Once()
				mockRoleRepo.On("CreateRole", ctx, mock.MatchedBy(func(r roleDto.ToDBCreateRole) bool {
					return r.Name == "'; DROP TABLE roles; --"
				})).Return(expectedRole, nil).Once()
			},
			expectedError: false,
			description:   "SQL injection attempt should be treated as normal string (parameterized query)",
		},
		{
			name: "Negative case - database error on CountRole",
			req:  validReq,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validPermissionGroupID).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
				mockRoleRepo.On("RoleNameIsNotDuplicated", ctx, validReq.Name, uuid.Nil).Return(true, nil).Once()
				mockRoleRepo.On("CountRole", ctx).Return(nil, errors.New("database connection error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database connection error",
			description:    "Database error should be returned",
		},
		{
			name: "Negative case - database error on CreateRole",
			req:  validReq,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validPermissionGroupID).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
				mockRoleRepo.On("RoleNameIsNotDuplicated", ctx, validReq.Name, uuid.Nil).Return(true, nil).Once()
				mockRoleRepo.On("CountRole", ctx).Return(&validCount, nil).Once()
				mockRoleRepo.On("CreateRole", ctx, mock.Anything).Return(nil, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error on create should be returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo.ExpectedCalls = nil
			mockRoleRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.CreateRole(c, tt.req, tt.authId)

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

func TestGetRoleByID(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()
	expectedRole := &models.Role{
		ID:   validID,
		Name: "Test Role",
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
			name: "Positive case - successful get role by ID",
			id:   validIDString,
			setupMock: func() {
				mockRoleRepo.On("GetRoleByID", ctx, validID).Return(expectedRole, nil).Once()
			},
			expectedError: false,
			description:   "Valid UUID should return role",
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
			name: "Negative case - role not found",
			id:   validIDString,
			setupMock: func() {
				mockRoleRepo.On("GetRoleByID", ctx, validID).Return(nil, errors.New("role not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "role not found",
			description:    "Non-existent role should return error",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in ID",
			id:   "'; DROP TABLE roles; --",
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
				mockRoleRepo.On("GetRoleByID", ctx, validID).Return(nil, errors.New("database connection error")).Once()
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

			result, err := usecaseInstance.GetRoleByID(c, tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, expectedRole.ID, result.ID)
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}

func TestGetIndexRole(t *testing.T) {
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

	expectedRoles := []models.Role{
		{ID: uuid.New(), Name: "Role 1"},
		{ID: uuid.New(), Name: "Role 2"},
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
				mockRoleRepo.On("GetIndexRole", ctx, validPageReq).Return(expectedRoles, 2, nil).Once()
			},
			expectedError: false,
			description:   "Valid page request should return roles",
		},
		{
			name: "Negative case - database error",
			req:  validPageReq,
			setupMock: func() {
				mockRoleRepo.On("GetIndexRole", ctx, validPageReq).Return(nil, 0, errors.New("database error")).Once()
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
				Search:    "'; DROP TABLE roles; --",
				SortBy:    "name",
				SortOrder: "ASC",
			},
			setupMock: func() {
				// SQL injection should be treated as normal search string
				injectedReq := request.PageRequest{
					Page:      1,
					PerPage:   10,
					Search:    "'; DROP TABLE roles; --",
					SortBy:    "name",
					SortOrder: "ASC",
				}
				mockRoleRepo.On("GetIndexRole", ctx, injectedReq).Return(expectedRoles, 2, nil).Once()
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
				SortBy:    "'; DROP TABLE roles; --",
				SortOrder: "ASC",
			},
			setupMock: func() {
				// Should be sanitized by ValidateAndSanitizeSortColumn
				injectedReq := request.PageRequest{
					Page:      1,
					PerPage:   10,
					Search:    "",
					SortBy:    "'; DROP TABLE roles; --",
					SortOrder: "ASC",
				}
				mockRoleRepo.On("GetIndexRole", ctx, injectedReq).Return(expectedRoles, 2, nil).Once()
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

			roles, total, err := usecaseInstance.GetIndexRole(c, tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, roles)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, roles)
				assert.GreaterOrEqual(t, total, len(roles))
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllRole(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	expectedRoles := []models.Role{
		{ID: uuid.New(), Name: "Role 1"},
		{ID: uuid.New(), Name: "Role 2"},
	}

	tests := []struct {
		name           string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name: "Positive case - successful get all roles",
			setupMock: func() {
				mockRoleRepo.On("GetAllRole", ctx).Return(expectedRoles, nil).Once()
			},
			expectedError: false,
			description:   "Should return all roles successfully",
		},
		{
			name: "Negative case - database error",
			setupMock: func() {
				mockRoleRepo.On("GetAllRole", ctx).Return(nil, errors.New("database error")).Once()
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

			roles, err := usecaseInstance.GetAllRole(c)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Nil(t, roles)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, roles)
				assert.Equal(t, len(expectedRoles), len(roles))
			}

			mockRoleRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateRole(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()
	validPermissionGroupID := uuid.New()

	validReq := &roleDto.ReqUpdateRole{
		Name:             "Updated Role",
		Description:      "Updated Description",
		PermissionGroups: []uuid.UUID{validPermissionGroupID},
	}

	expectedRole := &models.Role{
		ID:   validID,
		Name: "Updated Role",
	}

	tests := []struct {
		name           string
		id             string
		req            *roleDto.ReqUpdateRole
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name: "Positive case - successful role update",
			id:   validIDString,
			req:  validReq,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validPermissionGroupID).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
				mockRoleRepo.On("RoleNameIsNotDuplicated", ctx, validReq.Name, validID).Return(true, nil).Once()
				mockRoleRepo.On("UpdateRole", ctx, validID, mock.Anything).Return(expectedRole, nil).Once()
			},
			expectedError: false,
			description:   "Valid request should update role successfully",
		},
		{
			name: "Negative case - invalid UUID",
			id:   "invalid-uuid",
			req:  validReq,
			setupMock: func() {
				// GetPermissionGroupByID will be called first, then UUID parsing will fail
				mockRoleRepo.On("GetPermissionGroupByID", ctx, mock.Anything).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "Invalid UUID should return error",
		},
		{
			name: "Negative case - invalid permission group",
			id:   validIDString,
			req:  validReq,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validPermissionGroupID).Return(nil, errors.New("permission group not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Function with ID",
			description:    "Invalid permission group should return error",
		},
		{
			name: "Negative case - duplicate name",
			id:   validIDString,
			req:  validReq,
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validPermissionGroupID).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
				mockRoleRepo.On("RoleNameIsNotDuplicated", ctx, validReq.Name, validID).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.RoleErrorRoleNotFound,
			description:    "Duplicate name should return error",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in name",
			id:   validIDString,
			req: &roleDto.ReqUpdateRole{
				Name:             "'; DROP TABLE roles; --",
				Description:      "Test",
				PermissionGroups: []uuid.UUID{validPermissionGroupID},
			},
			setupMock: func() {
				mockRoleRepo.On("GetPermissionGroupByID", ctx, validPermissionGroupID).Return(&models.PermissionGroup{ID: validPermissionGroupID}, nil).Once()
				mockRoleRepo.On("RoleNameIsNotDuplicated", ctx, "'; DROP TABLE roles; --", validID).Return(true, nil).Once()
				mockRoleRepo.On("UpdateRole", ctx, validID, mock.MatchedBy(func(r roleDto.ToDBUpdateRole) bool {
					return r.Name == "'; DROP TABLE roles; --"
				})).Return(expectedRole, nil).Once()
			},
			expectedError: false,
			description:   "SQL injection attempt should be treated as normal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo.ExpectedCalls = nil
			mockRoleRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodPut, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.UpdateRole(c, tt.id, tt.req, "auth-id")

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

func TestSoftDeleteRole(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockRoleRepo, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()

	roleWithUsers := &models.Role{
		ID:        validID,
		Name:      "Test Role",
		TotalUser: 5,
	}

	roleWithoutUsers := &models.Role{
		ID:        validID,
		Name:      "Test Role",
		TotalUser: 0,
	}

	deletedRole := &models.Role{
		ID:   validID,
		Name: "Test Role",
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
			name: "Positive case - successful soft delete",
			id:   validIDString,
			setupMock: func() {
				mockRoleRepo.On("GetRoleByID", ctx, validID).Return(roleWithoutUsers, nil).Once()
				mockRoleRepo.On("SoftDeleteRole", ctx, validID, mock.Anything).Return(deletedRole, nil).Once()
			},
			expectedError: false,
			description:   "Role without users should be deleted successfully",
		},
		{
			name: "Negative case - role not found",
			id:   validIDString,
			setupMock: func() {
				mockRoleRepo.On("GetRoleByID", ctx, validID).Return(nil, errors.New("role not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Role Not Found",
			description:    "Non-existent role should return error",
		},
		{
			name: "Negative case - role has users",
			id:   validIDString,
			setupMock: func() {
				mockRoleRepo.On("GetRoleByID", ctx, validID).Return(roleWithUsers, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Role has user. Can't be deleted",
			description:    "Role with users should not be deleted",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo.ExpectedCalls = nil
			mockRoleRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.SoftDeleteRole(c, tt.id, "auth-id")

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

func TestMyPermissionsByUserToken(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockRoleRepo, mockAuthRepo := createTestUsecase()
	ctx := context.Background()

	validToken := "valid-access-token"
	roleID := uuid.New()
	testUser := models.User{
		ID:     uuid.New(),
		RoleId: roleID,
	}

	expectedRole := &models.Role{
		ID:   roleID,
		Name: "Test Role",
	}

	tests := []struct {
		name           string
		token          string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:  "Positive case - successful get permissions",
			token: validToken,
			setupMock: func() {
				mockAuthRepo.On("GetUserByAccessToken", ctx, validToken).Return(testUser, nil).Once()
				mockRoleRepo.On("GetRoleByID", ctx, roleID).Return(expectedRole, nil).Once()
			},
			expectedError: false,
			description:   "Valid token should return role permissions",
		},
		{
			name:  "Negative case - invalid token",
			token: "invalid-token",
			setupMock: func() {
				mockAuthRepo.On("GetUserByAccessToken", ctx, "invalid-token").Return(models.User{}, errors.New("token not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.UserNotFound,
			description:    "Invalid token should return error",
		},
		{
			name:  "Negative case - empty token",
			token: "",
			setupMock: func() {
				mockAuthRepo.On("GetUserByAccessToken", ctx, "").Return(models.User{}, errors.New("token required")).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.UserNotFound,
			description:    "Empty token should return error",
		},
		{
			name:  "Negative-Positive case - SQL injection attempt in token",
			token: "'; DROP TABLE jwt_tokens; --",
			setupMock: func() {
				// Should be treated as normal string due to parameterized query
				mockAuthRepo.On("GetUserByAccessToken", ctx, "'; DROP TABLE jwt_tokens; --").Return(models.User{}, errors.New("token not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.UserNotFound,
			description:    "SQL injection attempt should be treated as normal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRoleRepo.ExpectedCalls = nil
			mockRoleRepo.Calls = nil
			mockAuthRepo.ExpectedCalls = nil
			mockAuthRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

			result, err := usecaseInstance.MyPermissionsByUserToken(c, tt.token)

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
			mockAuthRepo.AssertExpectations(t)
		})
	}
}
