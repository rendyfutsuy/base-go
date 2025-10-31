package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
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
			return "", errors.New("token is revoked please re-login from login form again..")
		}
		return "", errors.New("token is revoked please re-login from login form again..")
	}

	// If user not found, token is revoked
	if user.ID == uuid.Nil {
		return "", errors.New("token is revoked please re-login from login form again..")
	}

	// Destroy old token
	if err := u.authRepo.DestroyToken(ctx, accessToken); err != nil {
		// Log error but continue with token generation
	}

	// Generate new JWT token with same logic as Authenticate
	now := time.Now().UTC()
	location, _ := time.LoadLocation("Asia/Jakarta") // WIB is UTC+7
	nowInJakarta := now.In(location)

	nextDay := nowInJakarta.AddDate(0, 0, 1)
	expireTime := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 3, 0, 0, 0, location)

	// Generate new JWT token
	// access token would always expire in next day on 03:00 AM WIB (UTC+7)
	claims := AuthClaims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().In(location)),
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
