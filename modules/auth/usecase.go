package auth

import (
	"context"

	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
)

// Usecase represent the auth's usecases
type Usecase interface {
	// every new usecase on ..modules/auth/usecase/, please register it here
	Authenticate(ctx context.Context, login string, password string) (string, error)
	SignOut(ctx context.Context, token string) error
	GetProfile(ctx context.Context, accessToken string) (user models.User, err error)
	UpdateProfile(ctx context.Context, profileChunks dto.ReqUpdateProfile, userId string) error
	UpdateMyPassword(ctx context.Context, passwordChunks dto.ReqUpdatePassword, userId string) error
	IsUserPasswordExpired(ctx context.Context, login string) error

	// for reset password
	RequestResetPassword(ctx context.Context, email string) error
	ResetUserPassword(ctx context.Context, newPassword string, token string) error

	// for refresh token
	RefreshToken(ctx context.Context, accessToken string) (string, error)
}
