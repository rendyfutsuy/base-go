package test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/user_management/usecase"
	"github.com/rendyfutsuy/base-go/utils/token_storage"
	"github.com/stretchr/testify/mock"
)

// setupTestLogger initializes a no-op logger for testing
// This prevents nil pointer panics when Logger is used in usecase code
func setupTestLogger() {
	usecase.SetupTestLogger()
}

// MockTokenStorage is a mock implementation of token_storage.TokenStorage
type MockTokenStorage struct {
	mock.Mock
}

func (m *MockTokenStorage) SaveSession(ctx context.Context, user models.User, accessToken, refreshToken, accessJTI, refreshJTI string, refreshTTL time.Duration) error {
	args := m.Called(ctx, user, accessToken, refreshToken, accessJTI, refreshJTI, refreshTTL)
	return args.Error(0)
}

func (m *MockTokenStorage) GetRefreshTokenMetadata(ctx context.Context, refreshJTI string) (token_storage.RefreshTokenMeta, error) {
	args := m.Called(ctx, refreshJTI)
	return args.Get(0).(token_storage.RefreshTokenMeta), args.Error(1)
}

func (m *MockTokenStorage) MarkRefreshTokenUsed(ctx context.Context, refreshJTI string) error {
	args := m.Called(ctx, refreshJTI)
	return args.Error(0)
}

func (m *MockTokenStorage) DestroySession(ctx context.Context, accessToken string) error {
	args := m.Called(ctx, accessToken)
	return args.Error(0)
}

func (m *MockTokenStorage) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockTokenStorage) ValidateAccessToken(ctx context.Context, accessToken string) (models.User, error) {
	args := m.Called(ctx, accessToken)
	return args.Get(0).(models.User), args.Error(1)
}
