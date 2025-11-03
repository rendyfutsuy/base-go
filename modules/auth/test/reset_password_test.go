package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRequestResetPassword(t *testing.T) {
	// Setup test logger to prevent nil pointer panics
	setupTestLogger()

	ctx := context.Background()
	mockRepo := new(MockAuthRepository)
	signingKey := []byte("test-secret-key")
	hashSalt := "test-salt"
	timeout := 5 * time.Second

	usecaseInstance := usecase.NewTestAuthUsecase(mockRepo, timeout, hashSalt, signingKey, 24*time.Hour)

	testUser := models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	tests := []struct {
		name           string
		email          string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:  "Positive case - successful request reset password",
			email: "test@example.com",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").Return(testUser, nil).Once()
				mockRepo.On("RequestResetPassword", ctx, "test@example.com").Return(nil).Once()
			},
			expectedError: false,
			description:   "Valid email should trigger reset password request",
		},
		{
			name:  "Negative case - email not found",
			email: "nonexistent@example.com",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "nonexistent@example.com").Return(models.User{}, errors.New("User Not Found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Email not found",
			description:    "Non-existent email should return error",
		},
		{
			name:  "Negative case - empty email",
			email: "",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "").Return(models.User{}, errors.New("User Not Found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Email not found",
			description:    "Empty email should return error",
		},
		{
			name:  "Negative case - database error on RequestResetPassword",
			email: "test@example.com",
			setupMock: func() {
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com").Return(testUser, nil).Once()
				mockRepo.On("RequestResetPassword", ctx, "test@example.com").Return(errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database error",
			description:    "Database error should be returned",
		},
		{
			name:  "Negative-Positive case - SQL injection attempt in email",
			email: "test@example.com' OR '1'='1",
			setupMock: func() {
				// Should not find user even with SQL injection attempt
				mockRepo.On("FindByEmailOrUsername", ctx, "test@example.com' OR '1'='1").Return(models.User{}, errors.New("User Not Found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Email not found",
			description:    "SQL injection attempt should not be executed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.setupMock()

			err := usecaseInstance.RequestResetPassword(ctx, tt.email)

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

func TestResetUserPassword(t *testing.T) {
	// Setup test logger to prevent nil pointer panics
	setupTestLogger()

	ctx := context.Background()
	mockRepo := new(MockAuthRepository)
	signingKey := []byte("test-secret-key")
	hashSalt := "test-salt"
	timeout := 5 * time.Second

	usecaseInstance := usecase.NewTestAuthUsecase(mockRepo, timeout, hashSalt, signingKey, 24*time.Hour)

	testUserID := uuid.New()
	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword123"), bcrypt.DefaultCost)
	testUser := models.User{
		ID:       testUserID,
		Email:    "test@example.com",
		Password: string(hashedOldPassword),
	}
	resetToken := "valid-reset-token"
	newPassword := "newpassword123"

	tests := []struct {
		name           string
		newPassword    string
		token          string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:        "Positive case - successful password reset",
			newPassword: newPassword,
			token:       resetToken,
			setupMock: func() {
				mockRepo.On("GetUserByResetPasswordToken", ctx, resetToken).Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, newPassword, testUserID).Return(false, errors.New("Password Not Match")).Once()
				mockRepo.On("AssertPasswordNeverUsesByUser", ctx, newPassword, testUserID).Return(true, nil).Once()
				mockRepo.On("UpdatePasswordById", ctx, newPassword, testUserID).Return(true, nil).Once()
				mockRepo.On("IncreasePasswordExpiredAt", ctx, testUserID).Return(nil).Once()
				mockRepo.On("AddPasswordHistory", ctx, mock.AnythingOfType("string"), testUserID).Return(nil).Once()
				mockRepo.On("ResetPasswordAttempt", ctx, testUserID).Return(nil).Once()
				mockRepo.On("DestroyAllResetPasswordToken", ctx, testUserID).Return(nil).Once()
				mockRepo.On("DestroyAllToken", ctx, testUserID).Return(nil).Once()
			},
			expectedError: false,
			description:   "Valid token and password should reset password successfully",
		},
		{
			name:        "Negative case - invalid token",
			newPassword: newPassword,
			token:       "invalid-token",
			setupMock: func() {
				mockRepo.On("GetUserByResetPasswordToken", ctx, "invalid-token").Return(models.User{}, errors.New("token not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "token not found",
			description:    "Invalid token should return error",
		},
		{
			name:        "Negative case - new password same as current password",
			newPassword: "oldpassword123",
			token:       resetToken,
			setupMock: func() {
				mockRepo.On("GetUserByResetPasswordToken", ctx, resetToken).Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, "oldpassword123", testUserID).Return(true, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "New Password should not be same with Current Password",
			description:    "New password same as current should return error",
		},
		{
			name:        "Negative case - new password already used",
			newPassword: newPassword,
			token:       resetToken,
			setupMock: func() {
				mockRepo.On("GetUserByResetPasswordToken", ctx, resetToken).Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, newPassword, testUserID).Return(false, errors.New("Password Not Match")).Once()
				mockRepo.On("AssertPasswordNeverUsesByUser", ctx, newPassword, testUserID).Return(false, errors.New("Youre already used this password")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Youre already used this password",
			description:    "Password already used should return error",
		},
		{
			name:        "Negative case - empty token",
			newPassword: newPassword,
			token:       "",
			setupMock: func() {
				mockRepo.On("GetUserByResetPasswordToken", ctx, "").Return(models.User{}, errors.New("token not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "token not found",
			description:    "Empty token should return error",
		},
		{
			name:        "Negative case - empty new password",
			newPassword: "",
			token:       resetToken,
			setupMock: func() {
				mockRepo.On("GetUserByResetPasswordToken", ctx, resetToken).Return(testUser, nil).Once()
				mockRepo.On("AssertPasswordRight", ctx, "", testUserID).Return(false, errors.New("Password Not Match")).Once()
				mockRepo.On("AssertPasswordNeverUsesByUser", ctx, "", testUserID).Return(true, nil).Once()
				mockRepo.On("UpdatePasswordById", ctx, "", testUserID).Return(true, nil).Once()
				mockRepo.On("IncreasePasswordExpiredAt", ctx, testUserID).Return(nil).Once()
				mockRepo.On("AddPasswordHistory", ctx, mock.AnythingOfType("string"), testUserID).Return(nil).Once()
				mockRepo.On("ResetPasswordAttempt", ctx, testUserID).Return(nil).Once()
				mockRepo.On("DestroyAllResetPasswordToken", ctx, testUserID).Return(nil).Once()
				mockRepo.On("DestroyAllToken", ctx, testUserID).Return(nil).Once()
			},
			expectedError: false,
			description:   "Empty password is technically valid but should be validated at handler level",
		},
		{
			name:        "Negative-Positive case - SQL injection attempt in token",
			newPassword: newPassword,
			token:       "token'; DROP TABLE reset_password_tokens; --",
			setupMock: func() {
				// Should not find user even with SQL injection attempt
				mockRepo.On("GetUserByResetPasswordToken", ctx, "token'; DROP TABLE reset_password_tokens; --").Return(models.User{}, errors.New("token not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "token not found",
			description:    "SQL injection attempt should not be executed",
		},
		{
			name:        "Negative-Positive case - SQL injection attempt in new password",
			newPassword: "'; DROP TABLE users; --",
			token:       resetToken,
			setupMock: func() {
				mockRepo.On("GetUserByResetPasswordToken", ctx, resetToken).Return(testUser, nil).Once()
				// Password should not match due to parameterized query
				mockRepo.On("AssertPasswordRight", ctx, "'; DROP TABLE users; --", testUserID).Return(false, errors.New("Password Not Match")).Once()
				// Should check password history with the literal string
				mockRepo.On("AssertPasswordNeverUsesByUser", ctx, "'; DROP TABLE users; --", testUserID).Return(true, nil).Once()
				mockRepo.On("UpdatePasswordById", ctx, "'; DROP TABLE users; --", testUserID).Return(true, nil).Once()
				mockRepo.On("IncreasePasswordExpiredAt", ctx, testUserID).Return(nil).Once()
				mockRepo.On("AddPasswordHistory", ctx, mock.AnythingOfType("string"), testUserID).Return(nil).Once()
				mockRepo.On("ResetPasswordAttempt", ctx, testUserID).Return(nil).Once()
				mockRepo.On("DestroyAllResetPasswordToken", ctx, testUserID).Return(nil).Once()
				mockRepo.On("DestroyAllToken", ctx, testUserID).Return(nil).Once()
			},
			expectedError: false,
			description:   "SQL injection in new password should be treated as literal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.setupMock()

			err := usecaseInstance.ResetUserPassword(ctx, tt.newPassword, tt.token)

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
