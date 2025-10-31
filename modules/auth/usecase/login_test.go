package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockAuthRepository is a mock implementation of auth.Repository
type MockAuthRepository struct {
	mock.Mock
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

func (m *MockAuthRepository) FindByCurrentSession(ctx context.Context, accessToken string) (dto.UserProfile, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).(dto.UserProfile), args.Error(1)
}

func (m *MockAuthRepository) UpdateProfileById(ctx context.Context, profileChunks dto.ReqUpdateProfile, userId uuid.UUID) (bool, error) {
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

func TestAuthenticate(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAuthRepository)
	signingKey := []byte("test-secret-key-for-jwt")
	hashSalt := "test-salt"
	timeout := 5 * time.Second

	usecase := &authUsecase{
		authRepo:       mockRepo,
		contextTimeout: timeout,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
	}

	testUserID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{
		ID:       testUserID,
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	tests := []struct {
		name           string
		login          string
		password       string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:     "Positive case - successful authentication",
			login:    "test@example.com",
			password: "password123",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordAttemptPassed", ctx, testUserID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, "password123", testUserID).Return(true, nil).Once()
				mockRepo.On("ResetPasswordAttempt", ctx, testUserID).Return(nil).Once()
				mockRepo.On("AssertPasswordExpiredIsPassed", ctx, testUserID).Return(false, nil).Once()
				mockRepo.On("AddUserAccessToken", ctx, mock.AnythingOfType("string"), testUserID).Return(nil).Once()
			},
			expectedError: false,
			description:   "Valid credentials should return access token",
		},
		{
			name:     "Negative case - user not found",
			login:    "nonexistent@example.com",
			password: "password123",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "nonexistent@example.com").Return(models.User{}, errors.New(constants.UserInvalid)).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.UserInvalid,
			description:    "Non-existent user should return error",
		},
		{
			name:     "Negative case - too many password attempts",
			login:    "test@example.com",
			password: "wrongpassword",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordAttemptPassed", ctx, testUserID).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "too many password attempts",
			description:    "Too many attempts should return error",
		},
		{
			name:     "Negative case - wrong password",
			login:    "test@example.com",
			password: "wrongpassword",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordAttemptPassed", ctx, testUserID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, "wrongpassword", testUserID).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "invalid credentials",
			description:    "Wrong password should return error",
		},
		{
			name:     "Negative case - password expired",
			login:    "test@example.com",
			password: "password123",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordAttemptPassed", ctx, testUserID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, "password123", testUserID).Return(true, nil).Once()
				mockRepo.On("ResetPasswordAttempt", ctx, testUserID).Return(nil).Once()
				mockRepo.On("AssertPasswordExpiredIsPassed", ctx, testUserID).Return(true, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.ErrPasswordExpired.Error(),
			description:    "Expired password should return error",
		},
		{
			name:     "Negative-Positive case - SQL injection attempt in login",
			login:    "test@example.com' OR '1'='1",
			password: "password123",
			setupMock: func() {
				// Should not find user even with SQL injection attempt due to parameterized query
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com' OR '1'='1").Return(models.User{}, errors.New(constants.UserInvalid)).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.UserInvalid,
			description:    "SQL injection attempt should be treated as invalid login, not executed",
		},
		{
			name:     "Negative-Positive case - SQL injection attempt in password",
			login:    "test@example.com",
			password: "password123'; DROP TABLE users; --",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordAttemptPassed", ctx, testUserID).Return(true, nil).Once()
				// Password should not match due to parameterized query, SQL injection should not execute
				mockRepo.On("AssertPasswordRight", ctx, "password123'; DROP TABLE users; --", testUserID).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "invalid credentials",
			description:    "SQL injection in password should be treated as wrong password, not executed",
		},
		{
			name:     "Negative case - empty login",
			login:    "",
			password: "password123",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "").Return(models.User{}, errors.New(constants.UserInvalid)).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.UserInvalid,
			description:    "Empty login should return error",
		},
		{
			name:     "Negative case - empty password",
			login:    "test@example.com",
			password: "",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordAttemptPassed", ctx, testUserID).Return(true, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, "", testUserID).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "invalid credentials",
			description:    "Empty password should return error",
		},
		{
			name:     "Negative case - database error on FindByEmailOrUsername",
			login:    "test@example.com",
			password: "password123",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").Return(models.User{}, errors.New("database connection error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database connection error",
			description:    "Database error should be returned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.setupMock()

			token, err := usecase.Authenticate(ctx, tt.login, tt.password)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestIsUserPasswordExpired(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockAuthRepository)
	signingKey := []byte("test-secret-key")
	hashSalt := "test-salt"
	timeout := 5 * time.Second

	usecase := &authUsecase{
		authRepo:       mockRepo,
		contextTimeout: timeout,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
	}

	testUserID := uuid.New()
	testUser := models.User{
		ID:    testUserID,
		Email: "test@example.com",
	}

	tests := []struct {
		name           string
		login          string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:  "Positive case - password not expired",
			login: "test@example.com",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordExpiredIsPassed", ctx, testUserID).Return(false, nil).Once()
			},
			expectedError: false,
			description:   "Valid user with non-expired password should not return error",
		},
		{
			name:  "Negative case - password expired",
			login: "test@example.com",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordExpiredIsPassed", ctx, testUserID).Return(true, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.ErrPasswordExpired.Error(),
			description:    "Expired password should return error",
		},
		{
			name:  "Negative case - user not found",
			login: "nonexistent@example.com",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "nonexistent@example.com").Return(models.User{}, errors.New(constants.UserInvalid)).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.UserInvalid,
			description:    "Non-existent user should return error",
		},
		{
			name:  "Negative-Positive case - SQL injection attempt",
			login: "test@example.com' OR '1'='1",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com' OR '1'='1").Return(models.User{}, errors.New(constants.UserInvalid)).Once()
			},
			expectedError:  true,
			expectedErrMsg: constants.UserInvalid,
			description:    "SQL injection attempt should not be executed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.setupMock()

			err := usecase.IsUserPasswordExpired(ctx, tt.login)

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
