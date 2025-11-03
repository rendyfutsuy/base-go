package usecase

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/utils"
)

// RefreshToken generates a new access token based on the provided access token.
// If the old token is revoked (not found in Redis), it returns an error.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - accessToken: The existing access token to be refreshed.
//
// Returns:
// - string: The new access token.
// - error: An error if the token refresh fails or if the token is revoked.
func (u *authUsecase) RefreshToken(ctx context.Context, accessToken string) (string, error) {
	// Check if token exists in Redis (not revoked)
	user, err := u.authRepo.GetUserByAccessToken(ctx, accessToken)
	if err != nil {
		if err == redis.Nil {
			return "", constants.ErrTokenRevoked
		}
		return "", constants.ErrTokenRevoked
	}

	// If user not found, token is revoked
	if user.ID == uuid.Nil {
		return "", constants.ErrTokenRevoked
	}

	// Destroy old token
	if err := u.authRepo.DestroyToken(ctx, accessToken); err != nil {
		// Log error but continue with token generation
	}

	// Generate new JWT token with same logic as Authenticate
	// Get TTL from config (same as AddUserAccessToken)
	ttlSeconds := utils.ConfigVars.Int("auth.access_token_ttl_seconds")
	if ttlSeconds <= 0 {
		ttlSeconds = 24 * 60 * 60 // Default 24 hours
	}

	now := time.Now().UTC()
	issuedAt := now
	expireTime := now.Add(time.Duration(ttlSeconds) * time.Second)

	// Generate new JWT token
	// access token expires based on TTL from config (same as Redis session)
	claims := AuthClaims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newAccessToken, err := token.SignedString(u.signingKey)
	if err != nil {
		return "", err // If fail to sign, return error.
	}

	// Record new access token to Redis
	if err := u.authRepo.AddUserAccessToken(ctx, newAccessToken, user.ID); err != nil {
		return "", err // If fail to record, return error.
	}

	return newAccessToken, nil
}
