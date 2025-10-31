package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthUsecase is a mock implementation of auth.Usecase
type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Authenticate(ctx context.Context, login string, password string) (string, error) {
	args := m.Called(ctx, login, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthUsecase) SignOut(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthUsecase) GetProfile(ctx context.Context, accessToken string) (dto.UserProfile, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).(dto.UserProfile), args.Error(1)
}

func (m *MockAuthUsecase) UpdateProfile(ctx context.Context, profileChunks dto.ReqUpdateProfile, userId string) error {
	args := m.Called(ctx, profileChunks, userId)
	return args.Error(0)
}

func (m *MockAuthUsecase) UpdateMyPassword(ctx context.Context, passwordChunks dto.ReqUpdatePassword, userId string) error {
	args := m.Called(ctx, passwordChunks, userId)
	return args.Error(0)
}

func (m *MockAuthUsecase) IsUserPasswordExpired(ctx context.Context, login string) error {
	args := m.Called(ctx, login)
	return args.Error(0)
}

func (m *MockAuthUsecase) RequestResetPassword(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthUsecase) ResetUserPassword(ctx context.Context, newPassword string, token string) error {
	args := m.Called(ctx, newPassword, token)
	return args.Error(0)
}

// setupEchoContext is a helper function (not currently used but kept for reference)
// func setupEchoContext() (echo.Context, *httptest.ResponseRecorder) {
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/", nil)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)
// 	return c, rec
// }

func TestAuthenticate(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := &AuthHandler{
		AuthUseCase: mockUsecase,
		validator:   validator.New(),
	}

	validToken := "valid-access-token"

	tests := []struct {
		name           string
		requestBody    dto.ReqAuthUser
		setupMock      func()
		expectedStatus int
		expectedError  bool
		expectedMsg    string
		description    string
	}{
		{
			name: "Positive case - successful authentication",
			requestBody: dto.ReqAuthUser{
				Login:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockUsecase.On("IsUserPasswordExpired", mock.Anything, "test@example.com").Return(nil).Once()
				mockUsecase.On("Authenticate", mock.Anything, "test@example.com", "password123").Return(validToken, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			description:    "Valid credentials should return access token",
		},
		{
			name: "Negative case - invalid request body",
			requestBody: dto.ReqAuthUser{
				Login:    "",
				Password: "",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "Empty login and password should return validation error",
		},
		{
			name: "Negative case - invalid credentials",
			requestBody: dto.ReqAuthUser{
				Login:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func() {
				mockUsecase.On("IsUserPasswordExpired", mock.Anything, "test@example.com").Return(nil).Once()
				mockUsecase.On("Authenticate", mock.Anything, "test@example.com", "wrongpassword").Return("", errors.New("invalid credentials")).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			expectedMsg:    "invalid credentials",
			description:    "Wrong password should return error",
		},
		{
			name: "Negative case - password expired",
			requestBody: dto.ReqAuthUser{
				Login:    "test@example.com",
				Password: "password123",
			},
			setupMock: func() {
				mockUsecase.On("IsUserPasswordExpired", mock.Anything, "test@example.com").Return(constants.ErrPasswordExpired).Once()
			},
			expectedStatus: 419, // StatusCode 419 for password expired
			expectedError:  true,
			expectedMsg:    constants.ErrPasswordExpired.Error(),
			description:    "Expired password should return 419 status",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in login",
			requestBody: dto.ReqAuthUser{
				Login:    "test@example.com' OR '1'='1",
				Password: "password123",
			},
			setupMock: func() {
				mockUsecase.On("IsUserPasswordExpired", mock.Anything, "test@example.com' OR '1'='1").Return(nil).Once()
				mockUsecase.On("Authenticate", mock.Anything, "test@example.com' OR '1'='1", "password123").Return("", errors.New("User Not Found")).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "SQL injection attempt should be treated as invalid login",
		},
		{
			name: "Negative-Positive case - SQL injection attempt in password",
			requestBody: dto.ReqAuthUser{
				Login:    "test@example.com",
				Password: "password123'; DROP TABLE users; --",
			},
			setupMock: func() {
				mockUsecase.On("IsUserPasswordExpired", mock.Anything, "test@example.com").Return(nil).Once()
				mockUsecase.On("Authenticate", mock.Anything, "test@example.com", "password123'; DROP TABLE users; --").Return("", errors.New("invalid credentials")).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "SQL injection in password should be treated as wrong password",
		},
		{
			name: "Negative case - missing login field",
			requestBody: dto.ReqAuthUser{
				Password: "password123",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "Missing login should return validation error",
		},
		{
			name: "Negative case - missing password field",
			requestBody: dto.ReqAuthUser{
				Login: "test@example.com",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "Missing password should return validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.ExpectedCalls = nil
			mockUsecase.Calls = nil
			tt.setupMock()

			e := echo.New()
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Authenticate(c)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)

				var response ResponseAuth
				json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NotEmpty(t, response.AccessToken)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSignOut(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := &AuthHandler{
		AuthUseCase: mockUsecase,
		validator:   validator.New(),
	}

	validToken := "valid-access-token"

	tests := []struct {
		name           string
		token          string
		setupMock      func()
		expectedStatus int
		expectedError  bool
		description    string
	}{
		{
			name:  "Positive case - successful logout",
			token: validToken,
			setupMock: func() {
				mockUsecase.On("SignOut", mock.Anything, validToken).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			description:    "Valid token should be destroyed successfully",
		},
		{
			name:  "Negative case - invalid token",
			token: "invalid-token",
			setupMock: func() {
				mockUsecase.On("SignOut", mock.Anything, "invalid-token").Return(errors.New("token not found")).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "Invalid token should return error",
		},
		{
			name:  "Negative-Positive case - SQL injection attempt in token",
			token: "token'; DROP TABLE jwt_tokens; --",
			setupMock: func() {
				mockUsecase.On("SignOut", mock.Anything, "token'; DROP TABLE jwt_tokens; --").Return(errors.New("token not found")).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "SQL injection attempt should be treated as invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.ExpectedCalls = nil
			mockUsecase.Calls = nil
			tt.setupMock()

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", nil)
			req.Header.Set("Authorization", "Bearer "+tt.token)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("token", tt.token)

			err := handler.SignOut(c)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetProfile(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := &AuthHandler{
		AuthUseCase: mockUsecase,
		validator:   validator.New(),
	}

	validToken := "valid-access-token"
	expectedProfile := dto.UserProfile{
		UserId: "user-id",
		Name:   "Test User",
		Email:  "test@example.com",
	}

	tests := []struct {
		name           string
		token          string
		setupMock      func()
		expectedStatus int
		expectedError  bool
		description    string
	}{
		{
			name:  "Positive case - successful get profile",
			token: validToken,
			setupMock: func() {
				mockUsecase.On("GetProfile", mock.Anything, validToken).Return(expectedProfile, nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			description:    "Valid token should return user profile",
		},
		{
			name:  "Negative case - invalid token",
			token: "invalid-token",
			setupMock: func() {
				mockUsecase.On("GetProfile", mock.Anything, "invalid-token").Return(dto.UserProfile{}, errors.New("token not found")).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "Invalid token should return error",
		},
		{
			name:  "Negative-Positive case - SQL injection attempt in token",
			token: "token'; DROP TABLE jwt_tokens; --",
			setupMock: func() {
				mockUsecase.On("GetProfile", mock.Anything, "token'; DROP TABLE jwt_tokens; --").Return(dto.UserProfile{}, errors.New("token not found")).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "SQL injection attempt should be treated as invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.ExpectedCalls = nil
			mockUsecase.Calls = nil
			tt.setupMock()

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/v1/auth/profile", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("token", tt.token)

			err := handler.GetProfile(c)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := &AuthHandler{
		AuthUseCase: mockUsecase,
		validator:   validator.New(),
	}

	validUserId := "user-id"

	tests := []struct {
		name           string
		userId         string
		requestBody    dto.ReqUpdateProfile
		setupMock      func()
		expectedStatus int
		expectedError  bool
		description    string
	}{
		{
			name:   "Positive case - successful update profile",
			userId: validUserId,
			requestBody: dto.ReqUpdateProfile{
				Name: "Updated Name",
			},
			setupMock: func() {
				mockUsecase.On("UpdateProfile", mock.Anything, dto.ReqUpdateProfile{Name: "Updated Name"}, validUserId).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			description:    "Valid profile update should succeed",
		},
		{
			name:   "Negative case - empty name",
			userId: validUserId,
			requestBody: dto.ReqUpdateProfile{
				Name: "",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "Empty name should return validation error",
		},
		{
			name:   "Negative-Positive case - SQL injection attempt in name",
			userId: validUserId,
			requestBody: dto.ReqUpdateProfile{
				Name: "'; DROP TABLE users; --",
			},
			setupMock: func() {
				// Handler validates required, so injection string should pass validation
				// but usecase/repository should treat it as literal
				mockUsecase.On("UpdateProfile", mock.Anything, dto.ReqUpdateProfile{Name: "'; DROP TABLE users; --"}, validUserId).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			description:    "SQL injection attempt should be treated as literal string",
		},
		{
			name:   "Negative case - missing name field",
			userId: validUserId,
			requestBody: dto.ReqUpdateProfile{
				Name: "",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "Missing name should return validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.ExpectedCalls = nil
			mockUsecase.Calls = nil
			tt.setupMock()

			e := echo.New()
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/profile", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("userId", tt.userId)

			err := handler.UpdateProfile(c)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUpdateMyPassword(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	handler := &AuthHandler{
		AuthUseCase: mockUsecase,
		validator:   validator.New(),
	}

	validUserId := "user-id"

	tests := []struct {
		name           string
		userId         string
		requestBody    dto.ReqUpdatePassword
		setupMock      func()
		expectedStatus int
		expectedError  bool
		description    string
	}{
		{
			name:   "Positive case - successful password update",
			userId: validUserId,
			requestBody: dto.ReqUpdatePassword{
				OldPassword:          "oldpassword123",
				NewPassword:          "newpassword123",
				PasswordConfirmation: "newpassword123",
			},
			setupMock: func() {
				mockUsecase.On("UpdateMyPassword", mock.Anything, mock.AnythingOfType("dto.ReqUpdatePassword"), validUserId).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
			description:    "Valid password update should succeed",
		},
		{
			name:   "Negative case - password confirmation mismatch",
			userId: validUserId,
			requestBody: dto.ReqUpdatePassword{
				OldPassword:          "oldpassword123",
				NewPassword:          "newpassword123",
				PasswordConfirmation: "differentpassword",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "Password confirmation mismatch should return validation error",
		},
		{
			name:   "Negative case - empty old password",
			userId: validUserId,
			requestBody: dto.ReqUpdatePassword{
				OldPassword:          "",
				NewPassword:          "newpassword123",
				PasswordConfirmation: "newpassword123",
			},
			setupMock:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "Empty old password should return validation error",
		},
		{
			name:   "Negative-Positive case - SQL injection attempt in passwords",
			userId: validUserId,
			requestBody: dto.ReqUpdatePassword{
				OldPassword:          "'; DROP TABLE users; --",
				NewPassword:          "newpass'; DROP TABLE users; --",
				PasswordConfirmation: "newpass'; DROP TABLE users; --",
			},
			setupMock: func() {
				// Should fail at old password check, but injection should not execute
				mockUsecase.On("UpdateMyPassword", mock.Anything, mock.AnythingOfType("dto.ReqUpdatePassword"), validUserId).Return(errors.New("Old Password not Match")).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
			description:    "SQL injection in passwords should be treated as literal strings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsecase.ExpectedCalls = nil
			mockUsecase.Calls = nil
			tt.setupMock()

			e := echo.New()
			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/password", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("userId", tt.userId)

			err := handler.UpdateMyPassword(c)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}
