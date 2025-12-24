package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rendyfutsuy/base-go/modules/auth/usecase"
	"github.com/stretchr/testify/assert"
)

func TestSignOut(t *testing.T) {
	// Setup test logger to prevent nil pointer panics
	setupTestLogger()

	ctx := context.Background()
	mockRepo := new(MockAuthRepository)
	signingKey := []byte("test-secret-key")
	refreshSigningKey := []byte("test-secret-refresh-key")
	hashSalt := "test-salt"
	timeout := 5 * time.Second

	usecaseInstance := usecase.NewTestAuthUsecase(mockRepo, nil, timeout, hashSalt, signingKey, refreshSigningKey, 24*time.Hour)

	validToken := "valid-access-token"

	tests := []struct {
		name           string
		token          string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:  "Positive case - successful logout",
			token: validToken,
			setupMock: func() {
				mockRepo.On("DestroyToken", ctx, validToken).Return(nil).Once()
			},
			expectedError: false,
			description:   "Valid token should be destroyed successfully",
		},
		{
			name:  "Negative case - invalid token",
			token: "invalid-token",
			setupMock: func() {
				mockRepo.On("DestroyToken", ctx, "invalid-token").Return(errors.New("token not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "token not found",
			description:    "Invalid token should return error",
		},
		{
			name:  "Negative case - empty token",
			token: "",
			setupMock: func() {
				mockRepo.On("DestroyToken", ctx, "").Return(errors.New("token is required")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "token is required",
			description:    "Empty token should return error",
		},
		{
			name:  "Negative case - database error",
			token: validToken,
			setupMock: func() {
				mockRepo.On("DestroyToken", ctx, validToken).Return(errors.New("database connection error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "database connection error",
			description:    "Database error should be returned",
		},
		{
			name:  "Negative-Positive case - SQL injection attempt",
			token: "token'; DROP TABLE jwt_tokens; --",
			setupMock: func() {
				// Should try to delete with the literal string, not execute SQL
				mockRepo.On("DestroyToken", ctx, "token'; DROP TABLE jwt_tokens; --").Return(errors.New("token not found")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "token not found",
			description:    "SQL injection attempt should be treated as literal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.setupMock()

			err := usecaseInstance.SignOut(ctx, tt.token)

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
