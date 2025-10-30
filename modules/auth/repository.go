package auth

import (
	"context"

	"github.com/google/uuid"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
)

// Repository represent the auth's repository contract
type Repository interface {
	// every new method on ..modules/auth/repository/, please register it here
	FindByEmailOrUsername(ctx context.Context, login string) (user models.User, err error)
	AssertPasswordRight(ctx context.Context, password string, userId uuid.UUID) (bool, error)
	AssertPasswordExpiredIsPassed(ctx context.Context, userId uuid.UUID) (bool, error)
	AddUserAccessToken(ctx context.Context, accessToken string, userId uuid.UUID) error
	GetUserByAccessToken(ctx context.Context, accessToken string) (user models.User, errorMain error)
	DestroyToken(ctx context.Context, accessToken string) error
	FindByCurrentSession(ctx context.Context, accessToken string) (profile dto.UserProfile, err error)
	UpdateProfileById(ctx context.Context, profileChunks dto.ReqUpdateProfile, userId uuid.UUID) (bool, error)
	UpdatePasswordById(ctx context.Context, hashedPassword string, userId uuid.UUID) (bool, error)
	DestroyAllToken(ctx context.Context, userId uuid.UUID) error
	AssertPasswordNeverUsesByUser(ctx context.Context, newPassword string, userId uuid.UUID) (bool, error)
	AddPasswordHistory(ctx context.Context, hashedPassword string, userId uuid.UUID) error
	AssertPasswordAttemptPassed(ctx context.Context, userId uuid.UUID) (bool, error)
	ResetPasswordAttempt(ctx context.Context, userId uuid.UUID) error

	// for reset password
	RequestResetPassword(ctx context.Context, email string) error
	GetUserByResetPasswordToken(ctx context.Context, token string) (user models.User, errorMain error)
	DestroyResetPasswordToken(ctx context.Context, token string) error
	IncreasePasswordExpiredAt(ctx context.Context, userId uuid.UUID) error
	DestroyAllResetPasswordToken(ctx context.Context, userId uuid.UUID) error
}
