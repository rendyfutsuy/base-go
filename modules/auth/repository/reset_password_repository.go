package repository

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth/tasks"
	"github.com/rendyfutsuy/base-go/utils"
	"gorm.io/gorm"
)

// RequestResetPassword generates a random password reset session and sends it to the user's email.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - email: The user's email.
//
// It takes the user's email as a parameter and returns an error if any.
func (repo *authRepository) RequestResetPassword(ctx context.Context, email string) error {
	// get user by email
	user, err := repo.FindByEmailOrUsername(ctx, email)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// generate a random string of length 64 (because Base64 encoding increases length)
	token := utils.GenerateRandomString(64)

	// encode the random string in Base64 to get a 16-character string
	session := base64.StdEncoding.EncodeToString([]byte(token))

	// add token to Database
	err = repo.AddResetPasswordToken(ctx, token, user.ID)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	// enqueue task
	task, err := tasks.NewEmailResetPasswordRequestTask(user.ID, user.Email, session)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	_, err = repo.QueueClient.Enqueue(task)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}

// GetUserByResetPasswordToken retrieves a user from the database based on the provided reset password token.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - token: The reset password token used to identify the user.
//
// Returns:
// - user: The retrieved user object.
// - errorMain: An error if the retrieval fails, or if the reset password token is not valid.
func (repo *authRepository) GetUserByResetPasswordToken(ctx context.Context, token string) (user models.User, errorMain error) {
	err := repo.DB.WithContext(ctx).
		Table("users usr").
		Select("usr.id as id, usr.email").
		Joins("JOIN reset_password_tokens rst_pwd ON rst_pwd.user_id = usr.id").
		Where("rst_pwd.access_token = ?", token).
		Scan(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No user found with this reset token")
			return user, errors.New(constants.AuthTokenInvalidRestart)
		}
		log.Printf("Error querying user: %v", err)
		return user, err
	}

	return user, nil
}

// AddResetPasswordToken inserts a new reset password token for a user into the database.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - token: The reset password token to be inserted.
// - userId: The unique identifier of the user.
//
// Returns:
// - error: An error if the insertion fails.
func (repo *authRepository) AddResetPasswordToken(ctx context.Context, token string, userId uuid.UUID) error {
	now := time.Now().UTC()
	resetToken := models.ResetPasswordToken{
		AccessToken: token,
		UserId:      userId,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	err := repo.DB.WithContext(ctx).Create(&resetToken).Error
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}

// DestroyResetPasswordToken deletes the reset password token from the database based on the provided token.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - token: The access token to identify the reset password token.
// Returns:
// - error: An error if the deletion fails, nil otherwise.
func (repo *authRepository) DestroyResetPasswordToken(ctx context.Context, token string) error {
	err := repo.DB.WithContext(ctx).
		Where("access_token = ?", token).
		Delete(&models.ResetPasswordToken{}).Error

	if err != nil {
		fmt.Println("Error deleting reset password token:", err)
		return err
	}
	return nil
}

// DestroyAllResetPasswordToken deletes all reset password tokens associated with a specific user ID.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - userId: The unique identifier of the user whose reset password tokens are to be deleted.
// Returns:
// - error: An error if the deletion fails, nil otherwise.
func (repo *authRepository) DestroyAllResetPasswordToken(ctx context.Context, userId uuid.UUID) error {
	err := repo.DB.WithContext(ctx).
		Where("user_id = ?", userId).
		Delete(&models.ResetPasswordToken{}).Error

	if err != nil {
		fmt.Println("Error deleting all reset password tokens:", err)
		return err
	}
	return nil
}

func (repo *authRepository) IncreasePasswordExpiredAt(ctx context.Context, userId uuid.UUID) error {
	// Calculate the expiration date to be 3 months from today
	expiredAt := time.Now().AddDate(0, 3, 0)

	err := repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userId).
		Update("password_expired_at", expiredAt).Error

	if err != nil {
		fmt.Println("Error updating password expiration:", err)
		return err
	}
	return nil
}
