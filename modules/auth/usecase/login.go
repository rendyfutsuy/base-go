package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth"
	"github.com/rendyfutsuy/base-go/utils"
	"go.uber.org/zap"
)

// Authenticate: returns Access + Refresh tokens
func (u *authUsecase) Authenticate(ctx context.Context, login string, password string) (auth.AuthenticateResult, error) {
	// 1) load user
	user, err := u.authRepo.FindByEmailOrUsername(ctx, login)
	if err != nil {
		// keep generic message for auth failures
		return auth.AuthenticateResult{}, errors.New(constants.AuthUsernamePasswordNotFound)
	}

	// 2) check password attempt limit
	isAttemptPassed, err := u.authRepo.AssertPasswordAttemptPassed(ctx, user.ID)
	if err != nil {
		// treat as exceeded for security
		return auth.AuthenticateResult{}, errors.New(constants.AuthPasswordAttemptExceeded)
	}
	if !isAttemptPassed {
		return auth.AuthenticateResult{}, errors.New(constants.AuthPasswordAttemptExceeded)
	}

	// 3) check password correctness
	isPasswordRight, err := u.authRepo.AssertPasswordRight(ctx, password, user.ID)
	if err != nil || !isPasswordRight {
		// keep generic message
		return auth.AuthenticateResult{}, errors.New(constants.AuthUsernamePasswordNotFound)
	}

	// 4) reset attempts
	if err := u.authRepo.ResetPasswordAttempt(ctx, user.ID); err != nil {
		utils.Logger.Warn("failed to reset password attempt counter", zap.Error(err))
		// non-fatal for authentication success — but you may choose to return error
	}

	// 5) check password expiry
	isPasswordExpired, err := u.authRepo.AssertPasswordExpiredIsPassed(ctx, user.ID)
	if err != nil {
		return auth.AuthenticateResult{}, err
	}
	if isPasswordExpired {
		return auth.AuthenticateResult{}, constants.ErrPasswordExpired
	}

	// 6) get first time login flag
	isFirstTimeLogin, err := u.authRepo.GetIsFirstTimeLogin(ctx, user.ID)
	if err != nil {
		return auth.AuthenticateResult{}, err
	}

	// 7) create access token
	accessToken, accessJTI, err := u.createAccessToken(user)
	if err != nil {
		return auth.AuthenticateResult{}, err
	}

	// 8) store access token session
	// 	save to redis
	if err := u.authRepo.AddUserAccessToken(ctx, accessToken, user.ID); err != nil {
		return auth.AuthenticateResult{}, fmt.Errorf("failed to store access token: %w", err)
	}

	// 9) create refresh token + store metadata
	refreshToken, refreshJTI, refreshTTL, err := u.createRefreshToken(user)
	if err != nil {
		return auth.AuthenticateResult{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	// 10) store store metadata for token

	// save to redis
	if err := u.authRepo.StoreRefreshToken(ctx, refreshJTI, user.ID, accessJTI, refreshTTL); err != nil {
		// best effort: if store fails, revoke created access token & return error
		_ = u.authRepo.DestroyToken(ctx, accessToken)
		return auth.AuthenticateResult{}, fmt.Errorf("failed to persist refresh token: %w", err)
	}

	return auth.AuthenticateResult{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		IsFirstTimeLogin: isFirstTimeLogin,
	}, nil
}

// RefreshToken: accepts a refresh token string, rotates tokens and returns new pair.
func (u *authUsecase) RefreshToken(ctx context.Context, refreshTokenString string) (auth.RefreshResult, error) {
	// 1) Parse refresh token to get JTI
	claims := &jwt.RegisteredClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(refreshTokenString, claims)
	if err != nil || claims.ID == "" {
		return auth.RefreshResult{}, constants.ErrTokenRevoked
	}
	refreshJTI := claims.ID

	// 2) Load refresh token metadata
	meta, err := u.authRepo.GetRefreshTokenMetadata(ctx, refreshJTI)
	if err != nil {
		return auth.RefreshResult{}, constants.ErrTokenRevoked
	}

	// 3) If already used → token theft detected
	if meta.Used {
		_ = u.authRepo.RevokeAllUserSessions(ctx, meta.UserID)
		return auth.RefreshResult{}, constants.ErrTokenRevoked
	}

	// 4) Expired refresh token?
	if time.Now().UTC().After(meta.ExpiresAt) {
		_ = u.authRepo.MarkRefreshTokenUsed(ctx, refreshJTI)
		return auth.RefreshResult{}, constants.ErrTokenRevoked
	}

	// 5) Mark refresh token as used
	_ = u.authRepo.MarkRefreshTokenUsed(ctx, refreshJTI)

	// 6) REVOKE OLD ACCESS TOKEN IF EXISTS
	if meta.AccessJTI != "" {
		_ = u.authRepo.DestroyToken(ctx, meta.AccessJTI)
	}

	// 7) Issue NEW access token
	user := models.User{ID: meta.UserID}
	newAccessToken, newAccessJTI, err := u.createAccessToken(user)
	if err != nil {
		return auth.RefreshResult{}, err
	}

	_ = u.authRepo.AddUserAccessToken(ctx, newAccessToken, meta.UserID)

	// 8) Create NEW refresh token (bind to new accessJTI)
	newRefreshToken, newRefreshJTI, newTTL, err := u.createRefreshToken(user)
	if err != nil {
		_ = u.authRepo.DestroyToken(ctx, newAccessJTI)
		return auth.RefreshResult{}, err
	}

	err = u.authRepo.StoreRefreshToken(ctx, newRefreshJTI, meta.UserID, newAccessJTI, newTTL)
	if err != nil {
		_ = u.authRepo.DestroyToken(ctx, newAccessJTI)
		return auth.RefreshResult{}, err
	}

	return auth.RefreshResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// --- Helper token generation functions ---

// createAccessToken creates a signed JWT access token (short-lived)
func (u *authUsecase) createAccessToken(user models.User) (tokenString string, accessJTI string, err error) {
	expires := utils.ConfigVars.Int("auth.access_token_ttl_seconds")
	if expires <= 0 {
		expires = 1800
	}

	now := time.Now().UTC()
	jti := uuid.NewString() // NEW — keep JTI returned
	claims := AuthClaims{
		UserID: user.ID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expires) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(u.signingKey)
	if err != nil {
		return "", "", err
	}
	return signed, jti, nil
}

// createRefreshToken returns (tokenString, jti, ttl, error)
func (u *authUsecase) createRefreshToken(user models.User) (string, string, time.Duration, error) {
	// refresh token TTL
	refreshTTLSeconds := utils.ConfigVars.Int("auth.refresh_token_ttl_seconds")
	if refreshTTLSeconds <= 0 {
		// default e.g. 14 days
		refreshTTLSeconds = 14 * 24 * 60 * 60
	}
	refreshTTL := time.Duration(refreshTTLSeconds) * time.Second

	now := time.Now().UTC()
	jti := uuid.NewString()

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(refreshTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        jti,
		// optionally set Subject = user.ID.String()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(u.refreshSigningKey)
	if err != nil {
		return "", "", 0, err
	}

	return signed, jti, refreshTTL, nil
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
