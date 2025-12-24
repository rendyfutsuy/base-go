package test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	userDto "github.com/rendyfutsuy/base-go/modules/user_management/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateUserPassword(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, mockUserRepo, mockAuthRepo, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()

	validReq := &userDto.ReqUpdateUserPassword{
		OldPassword:          "oldpassword123",
		NewPassword:          "newpassword123",
		PasswordConfirmation: "newpassword123",
	}

	tests := []struct {
		name           string
		id             string
		req            *userDto.ReqUpdateUserPassword
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name: "Positive case - successful password update",
			id:   validIDString,
			req:  validReq,
			setupMock: func() {
				mockUserRepo.On("IsUserPasswordCanUpdated", ctx, validID).Return(true, nil).Once()
				mockAuthRepo.On("AssertPasswordRight", ctx, validReq.OldPassword, validID).Return(true, nil).Once()
				mockAuthRepo.On("AssertPasswordRight", ctx, validReq.NewPassword, validID).Return(false, nil).Once()
				mockAuthRepo.On("AssertPasswordNeverUsesByUser", ctx, validReq.NewPassword, validID).Return(true, nil).Once()
				mockAuthRepo.On("AddPasswordHistory", ctx, mock.Anything, validID).Return(nil).Once()
				mockAuthRepo.On("ResetPasswordAttempt", ctx, validID).Return(nil).Once()
				mockAuthRepo.On("UpdatePasswordById", ctx, validReq.NewPassword, validID).Return(true, nil).Once()
				mockAuthRepo.On("DestroyAllToken", ctx, validID).Return(nil).Once()
			},
			expectedError: false,
			description:   "Valid request should update password successfully",
		},
		{
			name: "Negative case - invalid UUID",
			id:   "invalid-uuid",
			req:  validReq,
			setupMock: func() {
				// No repository call should be made
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "Invalid UUID should return error",
		},
		{
			name: "Negative case - password cannot be updated",
			id:   validIDString,
			req:  validReq,
			setupMock: func() {
				mockUserRepo.On("IsUserPasswordCanUpdated", ctx, validID).Return(false, errors.New("password cannot be updated")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "password cannot be updated",
			description:    "Password update restriction should return error",
		},
		{
			name: "Negative case - old password not match",
			id:   validIDString,
			req:  validReq,
			setupMock: func() {
				mockUserRepo.On("IsUserPasswordCanUpdated", ctx, validID).Return(true, nil).Once()
				mockAuthRepo.On("AssertPasswordRight", ctx, validReq.OldPassword, validID).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Old Password not Match",
			description:    "Wrong old password should return error",
		},
		{
			name: "Negative case - new password same as current",
			id:   validIDString,
			req:  validReq,
			setupMock: func() {
				mockUserRepo.On("IsUserPasswordCanUpdated", ctx, validID).Return(true, nil).Once()
				mockAuthRepo.On("AssertPasswordRight", ctx, validReq.OldPassword, validID).Return(true, nil).Once()
				mockAuthRepo.On("AssertPasswordRight", ctx, validReq.NewPassword, validID).Return(true, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "New Password should not be same with Current Password",
			description:    "New password same as current should return error",
		},
		{
			name: "Negative case - new password used before",
			id:   validIDString,
			req:  validReq,
			setupMock: func() {
				mockUserRepo.On("IsUserPasswordCanUpdated", ctx, validID).Return(true, nil).Once()
				mockAuthRepo.On("AssertPasswordRight", ctx, validReq.OldPassword, validID).Return(true, nil).Once()
				mockAuthRepo.On("AssertPasswordRight", ctx, validReq.NewPassword, validID).Return(false, nil).Once()
				mockAuthRepo.On("AssertPasswordNeverUsesByUser", ctx, validReq.NewPassword, validID).Return(false, errors.New("password used before")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "password used before",
			description:    "New password used before should return error",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in ID",
			id:   "'; DROP TABLE users; --",
			req:  validReq,
			setupMock: func() {
				// Should fail UUID parsing, preventing SQL injection
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "SQL injection attempt should fail UUID validation",
		},
		{
			name: "Negative case - database error on update",
			id:   validIDString,
			req:  validReq,
			setupMock: func() {
				mockUserRepo.On("IsUserPasswordCanUpdated", ctx, validID).Return(true, nil).Once()
				mockAuthRepo.On("AssertPasswordRight", ctx, validReq.OldPassword, validID).Return(true, nil).Once()
				mockAuthRepo.On("AssertPasswordRight", ctx, validReq.NewPassword, validID).Return(false, nil).Once()
				mockAuthRepo.On("AssertPasswordNeverUsesByUser", ctx, validReq.NewPassword, validID).Return(true, nil).Once()
				mockAuthRepo.On("AddPasswordHistory", ctx, mock.Anything, validID).Return(nil).Once()
				mockAuthRepo.On("ResetPasswordAttempt", ctx, validID).Return(nil).Once()
				mockAuthRepo.On("UpdatePasswordById", ctx, validReq.NewPassword, validID).Return(false, errors.New("database error")).Once()
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
			mockAuthRepo.ExpectedCalls = nil
			mockAuthRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodPut, "/", nil), httptest.NewRecorder())

			err := usecaseInstance.UpdateUserPassword(c.Request().Context(), tt.id, tt.req)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
			mockAuthRepo.AssertExpectations(t)
		})
	}
}

func TestAssertCurrentUserPassword(t *testing.T) {
	setupTestLogger()

	e := echo.New()
	usecaseInstance, _, mockAuthRepo, _ := createTestUsecase()
	ctx := context.Background()

	validID := uuid.New()
	validIDString := validID.String()
	validPassword := "password123"

	tests := []struct {
		name           string
		id             string
		password       string
		setupMock      func()
		expectedError  bool
		expectedErrMsg string
		description    string
	}{
		{
			name:     "Positive case - password matches",
			id:       validIDString,
			password: validPassword,
			setupMock: func() {
				mockAuthRepo.On("AssertPasswordRight", ctx, validPassword, validID).Return(true, nil).Once()
			},
			expectedError: false,
			description:   "Correct password should pass validation",
		},
		{
			name:     "Negative case - invalid UUID",
			id:       "invalid-uuid",
			password: validPassword,
			setupMock: func() {
				// No repository call should be made
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "Invalid UUID should return error",
		},
		{
			name:     "Negative case - password does not match",
			id:       validIDString,
			password: "wrongpassword",
			setupMock: func() {
				mockAuthRepo.On("AssertPasswordRight", ctx, "wrongpassword", validID).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Given Password not Match with Current Password",
			description:    "Wrong password should return error",
		},
		{
			name:     "Negative case - database error",
			id:       validIDString,
			password: validPassword,
			setupMock: func() {
				mockAuthRepo.On("AssertPasswordRight", ctx, validPassword, validID).Return(false, errors.New("database error")).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Given Password not Match with Current Password",
			description:    "Database error should return error",
		},
		{
			name:     "Negative-Positive case - SQL injection attempt in ID",
			id:       "'; DROP TABLE users; --",
			password: validPassword,
			setupMock: func() {
				// Should fail UUID parsing, preventing SQL injection
			},
			expectedError:  true,
			expectedErrMsg: "requested param is string",
			description:    "SQL injection attempt should fail UUID validation",
		},
		{
			name:     "Negative-Positive case - SQL injection attempt in password",
			id:       validIDString,
			password: "'; DROP TABLE users; --",
			setupMock: func() {
				// Password should be treated as normal string and hashed before comparison
				mockAuthRepo.On("AssertPasswordRight", ctx, "'; DROP TABLE users; --", validID).Return(false, nil).Once()
			},
			expectedError:  true,
			expectedErrMsg: "Given Password not Match with Current Password",
			description:    "SQL injection attempt in password should be treated as normal string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthRepo.ExpectedCalls = nil
			mockAuthRepo.Calls = nil
			tt.setupMock()

			c := e.NewContext(httptest.NewRequest(http.MethodPost, "/", nil), httptest.NewRecorder())

			err := usecaseInstance.AssertCurrentUserPassword(c.Request().Context(), tt.id, tt.password)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			mockAuthRepo.AssertExpectations(t)
		})
	}
}
