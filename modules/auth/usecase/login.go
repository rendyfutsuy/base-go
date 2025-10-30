package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
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
		return "", errors.New("too many password attempts")
	}

	// assert password given is same with saved password
	isPasswordRight, err := u.authRepo.AssertPasswordRight(ctx, password, user.ID)

	if err != nil {
		return "", err // Return error from the check itself.
	}
	if !isPasswordRight {
		// This should return a specific "invalid credentials" error.
		return "", errors.New("invalid credentials")
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
	now := time.Now().UTC()
	location, _ := time.LoadLocation("Asia/Jakarta") // WIB is UTC+7
	nowInJakarta := now.In(location)

	nextDay := nowInJakarta.AddDate(0, 0, 1)
	expireTime := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 3, 0, 0, 0, location)

	// Generate JWT token
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
