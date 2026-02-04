package token_storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
)

type RefreshTokenMeta struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
	Used      bool
	AccessJTI string
}

type TokenStorage interface {
	// SaveSession saves both access and refresh token info.
	SaveSession(ctx context.Context, user models.User, accessToken, refreshToken, accessJTI, refreshJTI string, refreshTTL time.Duration) error

	// GetRefreshTokenMetadata retrieves metadata by refresh token JTI.
	GetRefreshTokenMetadata(ctx context.Context, refreshJTI string) (RefreshTokenMeta, error)

	// MarkRefreshTokenUsed marks the refresh token as used.
	MarkRefreshTokenUsed(ctx context.Context, refreshJTI string) error

	// DestroySession destroys the session associated with the access token.
	DestroySession(ctx context.Context, accessToken string) error

	// RevokeAllUserSessions revokes all sessions for a user.
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
}
