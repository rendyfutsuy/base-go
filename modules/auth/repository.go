package auth

import (
	// "database/sql"

	"github.com/google/uuid"
	models "github.com/rendyfutsuy/base-go.git/models"
	"github.com/rendyfutsuy/base-go.git/modules/auth/dto"
)

// Repository represent the auth's repository contract
type Repository interface {
	// every new method on ..modules/auth/repository/, please register it here
	FindByEmailOrUsername(login string) (user models.User, err error)
	AssertPasswordRight(password string, userId uuid.UUID) (bool, error)
	AssertPasswordExpiredIsPassed(userId uuid.UUID) (bool, error)
	AddUserAccessToken(accessToken string, userId uuid.UUID) error
	GetUserByAccessToken(accessToken string) (user models.User, errorMain error)
	DestroyToken(accessToken string) error
	FindByCurrentSession(accessToken string) (profile dto.UserProfile, err error)
	UpdateProfileById(profileChunks dto.ReqUpdateProfile, userId uuid.UUID) (bool, error)
	UpdatePasswordById(hashedPassword string, userId uuid.UUID) (bool, error)
	DestroyAllToken(userId uuid.UUID) error
	AssertPasswordNeverUsesByUser(newPassword string, userId uuid.UUID) (bool, error)
	AddPasswordHistory(hashedPassword string, userId uuid.UUID) error
	AssertPasswordAttemptPassed(userId uuid.UUID) (bool, error)

	// for reset password
	RequestResetPassword(email string) error
	GetUserByResetPasswordToken(token string) (user models.User, errorMain error)
	DestroyResetPasswordToken(token string) error
	IncreasePasswordExpiredAt(userId uuid.UUID) error
	DestroyAllResetPasswordToken(userId uuid.UUID) error
}
