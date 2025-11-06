package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/utils"
)

func (u *authUsecase) Authenticate(ctx context.Context, login string, password string) (string, error) {
	// get user by email
	user, err := u.authRepo.FindByEmailOrUsername(ctx, login)

	// if fail to get user return error
	if err != nil {
		return "", err

	}

	// assert login attempt is not above 3
	isAttemptPassed, err := u.authRepo.AssertPasswordAttemptPassed(ctx, user.ID)
	if err != nil {
		return "", err // Return error from the check itself.
	}
	if !isAttemptPassed {
		// This should return a specific "too many attempts" error.
		// For now, we'll assume the repo returns it in the 'err' variable.
		return "", errors.New(constants.AuthTooManyPasswordAttempts)
	}

	// assert password given is same with saved password
	isPasswordRight, err := u.authRepo.AssertPasswordRight(ctx, password, user.ID)

	if err != nil {
		return "", err // Return error from the check itself.
	}
	if !isPasswordRight {
		// This should return a specific "invalid credentials" error.
		return "", errors.New(constants.AuthInvalidCredentials)
	}

	// Reset password attempt counter to 0 since login was successful.
	if err := u.authRepo.ResetPasswordAttempt(ctx, user.ID); err != nil {
		return "", err // If fail to reset, return error.
	}

	// assert if password expiration passed
	isPasswordExpired, err := u.authRepo.AssertPasswordExpiredIsPassed(ctx, user.ID)
	if err != nil {
		return "", err // Return error from the check itself.
	}
	// When the password has expired, `isPasswordExpired` is true.
	// You must return the specific `ErrPasswordExpired` variable, not the `err` variable
	// from the line above (which is nil in this case).
	if isPasswordExpired {
		return "", constants.ErrPasswordExpired
	}

	// --- JWT Generation Logic ---
	// Get Redis TTL from config (for session storage)
	ttlSeconds := utils.ConfigVars.Int("auth.access_token_ttl_seconds")
	if ttlSeconds <= 0 {
		ttlSeconds = 2 * 24 * 60 * 60 // Default 2 days
	}

	// Get JWT expiration time from config (for token expiration)
	jwtExpiresAtSeconds := utils.ConfigVars.Int("auth.jwt_expires_at_seconds")
	if jwtExpiresAtSeconds <= 0 {
		jwtExpiresAtSeconds = 30 * 60 // Default 30 minutes
	}

	now := time.Now().UTC()
	issuedAt := now
	expireTime := now.Add(time.Duration(jwtExpiresAtSeconds) * time.Second)

	// Generate JWT token
	// JWT expires based on jwt_expires_at_seconds (30 minutes)
	// Redis session TTL is based on access_token_ttl_seconds (2 days)
	claims := AuthClaims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(u.signingKey)
	if err != nil {
		return "", err // If fail to sign, return error.
	}

	// Record access token to the database.
	if err := u.authRepo.AddUserAccessToken(ctx, accessToken, user.ID); err != nil {
		return "", err // If fail to record, return error.
	}

	return accessToken, nil
}

func (u *authUsecase) IsUserPasswordExpired(ctx context.Context, login string) error {
	// get user by email
	user, err := u.authRepo.FindByEmailOrUsername(ctx, login)

	if err != nil {
		return err
	}

	// assert if password expiration passed
	isPasswordExpired, err := u.authRepo.AssertPasswordExpiredIsPassed(ctx, user.ID)

	if err != nil {
		return err
	}

	// if password expired return error
	if isPasswordExpired {
		return constants.ErrPasswordExpired
	}

	return nil
}
