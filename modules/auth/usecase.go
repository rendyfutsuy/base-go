package auth

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
)

// AuthenticateResult represents the result of authentication
type AuthenticateResult struct {
	AccessToken      string
	RefreshToken     string
	IsFirstTimeLogin bool
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenMeta struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
	Used      bool
	AccessJTI string
}

// Usecase represent the auth's usecases
type Usecase interface {
	// every new usecase on ..modules/auth/usecase/, please register it here
	Authenticate(ctx context.Context, login string, password string) (AuthenticateResult, error)
	SignOut(ctx context.Context, token string) error
	GetProfile(ctx context.Context, accessToken string) (user models.User, err error)
	UpdateProfile(ctx context.Context, profileChunks dto.ReqUpdateProfile, userId string) error
	UpdateMyPassword(ctx context.Context, passwordChunks dto.ReqUpdatePassword, userId string) error
	UpdateMyAvatar(ctx context.Context, user models.User, file *multipart.FileHeader) error
	IsUserPasswordExpired(ctx context.Context, login string) error

	// for reset password
	RequestResetPassword(ctx context.Context, email string) error
	ResetUserPassword(ctx context.Context, newPassword string, token string) error

	// for refresh token
	RefreshToken(ctx context.Context, refreshToken string) (RefreshResult, error)
}
