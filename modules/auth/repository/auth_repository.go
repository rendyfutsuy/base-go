package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rendyfutsuy/base-go/constants"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authRepository struct {
	DB           *gorm.DB
	EmailService *services.EmailService
	QueueClient  *asynq.Client
}

func NewAuthRepository(DB *gorm.DB, EmailService *services.EmailService) auth.Repository {

	QueueClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     utils.ConfigVars.String("redis.address"),
		Password: utils.ConfigVars.String("redis.password"),
		DB:       utils.ConfigVars.Int("redis.db"),
	})

	return &authRepository{
		DB,
		EmailService,
		QueueClient,
	}
}

// FindByEmail retrieves a user from the database by email.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - email: The email of the user to retrieve.
//
// Returns:
// - user: The retrieved user.
// - err:  An error if the retrieval fails.
func (repo *authRepository) FindByEmailOrUsername(ctx context.Context, login string) (user models.User, err error) {
	err = repo.DB.WithContext(ctx).
		Select("id, email, password").
		Where("(email = ? OR username = ?) AND deleted_at IS NULL AND is_active = ?", login, login, true).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No user found with email/username: %s", login)
			return user, errors.New(constants.UserInvalid)
		}
		log.Printf("Error querying user: %v", err)
		return user, err
	}

	return user, nil
}

// AssertPasswordRight checks if the provided password matches the hashed password in the database for the given user ID.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - password: The password to compare.
// - userId: The unique identifier of the user.
//
// Returns:
// - bool: True if the passwords match, false otherwise.
// - error: An error if the comparison fails or if there are database errors.
func (repo *authRepository) AssertPasswordRight(ctx context.Context, password string, userId uuid.UUID) (bool, error) {
	var user models.User
	err := repo.DB.WithContext(ctx).
		Select("password").
		Where("id = ? AND deleted_at IS NULL AND is_active = ?", userId, true).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println(constants.UserNotFound)
			return false, errors.New(constants.UserInvalid)
		}
		fmt.Println("Error querying database:", err)
		return false, err
	}

	// Compare the provided password with the hashed password from the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		// password do not match, add counter on users table
		repo.DB.WithContext(ctx).
			Model(&models.User{}).
			Where("id = ?", userId).
			Updates(map[string]interface{}{
				"counter":    gorm.Expr("counter + ?", 1),
				"updated_at": time.Now().UTC(),
			})

		// Passwords do not match, return error
		return false, errors.New("Password Not Match")
	}

	return true, nil
}

// AssertPasswordNeverUsesByUser checks if the new password has been used before by the user.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - newPassword: The new password to check.
// - userId: The unique identifier of the user.
//
// Returns:
// - bool: True if the new password has not been used before, false otherwise.
// - error: An error if there are database query errors or if the new password matches an old password.
// case:
// - if new password matches an password in password history, return error
// - if there are database query errors, return error
// - if new password has not been used before and no present on password history, return true
func (repo *authRepository) AssertPasswordNeverUsesByUser(ctx context.Context, newPassword string, userId uuid.UUID) (bool, error) {
	var histories []models.PasswordHistory
	err := repo.DB.WithContext(ctx).
		Select("hashed_password").
		Where("user_id = ?", userId).
		Find(&histories).Error

	if err != nil {
		log.Fatal(err)
		return false, err
	}

	for _, history := range histories {
		// Compare the new password with each old hashed password
		err = bcrypt.CompareHashAndPassword([]byte(history.HashedPassword), []byte(newPassword))
		if err == nil {
			// Password matches an old password
			return false, fmt.Errorf("Youre already used this password, please try another one..")
		} else if err != bcrypt.ErrMismatchedHashAndPassword {
			// Unknown error
			return false, fmt.Errorf("error comparing hashed password: %w", err)
		}
	}

	return true, nil
}

// AssertPasswordExpiredIsPassed checks if the password expiration date of a user has passed.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - userId: The unique identifier of the user.
//
// Returns:
// - bool: True if the password has expired, false otherwise.
// - error: An error if the query to the database fails.
func (repo *authRepository) AssertPasswordExpiredIsPassed(ctx context.Context, userId uuid.UUID) (bool, error) {
	var user models.User
	err := repo.DB.WithContext(ctx).
		Select("password_expired_at").
		Where("id = ?", userId).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println(constants.UserNotFound)
			return false, errors.New(constants.UserNotFound)
		}
		fmt.Println("Error querying database:", err)
		return false, err
	}

	// Get the current time
	currentTime := time.Now()

	// Compare expirationDate with currentTime to check if it has passed
	if user.PasswordExpiredAt.Before(currentTime) {
		// Password has expired
		return true, errors.New("password has expired, please change your password now")
	}

	// if password not expired, return false
	return false, nil
}

func (repo *authRepository) AssertPasswordAttemptPassed(ctx context.Context, userId uuid.UUID) (bool, error) {
	var user models.User
	err := repo.DB.WithContext(ctx).
		Select("counter").
		Where("id = ?", userId).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println(constants.UserNotFound)
			return false, errors.New(constants.UserNotFound)
		}
		fmt.Println("Error querying database:", err)
		return false, err
	}

	// if attempt above or equals to 3, return false
	if user.Counter >= 3 {
		// Password has expired
		return false, errors.New("Password Attempt is above 3, you're blocked. please contact admin")
	}

	return true, nil
}

func (repo *authRepository) ResetPasswordAttempt(ctx context.Context, userId uuid.UUID) error {
	// reset attempt to 0
	return repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userId).
		Updates(map[string]interface{}{
			"counter":    0,
			"updated_at": time.Now().UTC(),
		}).Error
}

// AddUserAccessToken inserts a new access token for a user into the database.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - accessToken: The access token to be inserted.
// - userId: The unique identifier of the user.
//
// Returns:
// - error: An error if the insertion fails.
func (repo *authRepository) AddUserAccessToken(ctx context.Context, accessToken string, userId uuid.UUID) error {
	now := time.Now().UTC()
	jwtToken := models.JWTToken{
		AccessToken: accessToken,
		UserId:      userId,
		CreatedAt:   now,
		UpdatedAt:   &now,
	}

	err := repo.DB.WithContext(ctx).Create(&jwtToken).Error
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}

// AddPasswordHistory inserts a new password history for a user into the database.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - hashedPassword: The hashed password to be inserted.
// - userId: The unique identifier of the user.
//
// Returns:
// - error: An error if the insertion fails.
func (repo *authRepository) AddPasswordHistory(ctx context.Context, hashedPassword string, userId uuid.UUID) error {
	now := time.Now().UTC()
	passwordHistory := models.PasswordHistory{
		HashedPassword: hashedPassword,
		UserId:         userId,
		CreatedAt:      now,
		UpdatedAt:      &now,
	}

	err := repo.DB.WithContext(ctx).Create(&passwordHistory).Error
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return nil
}

// GetUserByAccessToken retrieves a user from the database based on the provided access token.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - accessToken: The access token used to identify the user.
//
// Returns:
// - user: The retrieved user object.
// - errorMain: An error if the retrieval fails, or if the access token is not valid.
func (repo *authRepository) GetUserByAccessToken(ctx context.Context, accessToken string) (user models.User, errorMain error) {
	err := repo.DB.WithContext(ctx).
		Table("users usr").
		Select("usr.id as id, usr.full_name as full_name, usr.email, usr.role_id, roles.name as role_name").
		Joins("JOIN jwt_tokens jwt ON jwt.user_id = usr.id").
		Joins("JOIN roles ON roles.id = usr.role_id").
		Where("jwt.access_token = ?", accessToken).
		Scan(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No user found with this access token")
			return user, errors.New("User Not Found, the access token is not valid please re-login")
		}
		log.Printf("Error querying user: %v", err)
		return user, err
	}

	return user, nil
}

// DestroyToken deletes a JWT token from the database.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - accessToken: The access token to be deleted.
//
// Returns:
// - error: An error if the deletion fails, nil otherwise.
func (repo *authRepository) DestroyToken(ctx context.Context, accessToken string) error {
	err := repo.DB.WithContext(ctx).
		Where("access_token = ?", accessToken).
		Delete(&models.JWTToken{}).Error

	if err != nil {
		fmt.Println("Error deleting token:", err)
		return err
	}
	return nil
}

// FindByCurrentSession retrieves user profile based on the provided access token.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - accessToken: The access token used to identify the user session.
//
// Returns:
// - profile: User profile information.
// - err: An error if the retrieval fails, nil otherwise.
func (repo *authRepository) FindByCurrentSession(ctx context.Context, accessToken string) (profile dto.UserProfile, err error) {
	err = repo.DB.WithContext(ctx).
		Table("users usr").
		Select(`
			usr.id AS user_id,
			usr.email,
			usr.full_name AS name,
			rls.name AS role,
			CASE 
				WHEN usr.is_active THEN 'Active' 
				ELSE 'In Active' 
			END AS is_active,
			usr.gender
		`).
		Joins("JOIN roles rls ON rls.id = usr.role_id").
		Joins("JOIN jwt_tokens jwt ON jwt.user_id = usr.id").
		Where("jwt.access_token = ? AND usr.deleted_at IS NULL", accessToken).
		Scan(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No user found")
			return profile, errors.New(constants.UserInvalid)
		}
		log.Printf("Error querying user profile: %v", err)
		return profile, err
	}
	return profile, nil
}

// UpdateProfileById updates the full name of a user profile by ID.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - profileChunks: The updated profile information.
// - userId: The unique identifier of the user.
//
// Returns:
// - bool: True if the profile was successfully updated, false otherwise.
// - error: An error if the update fails.
func (repo *authRepository) UpdateProfileById(ctx context.Context, profileChunks dto.ReqUpdateProfile, userId uuid.UUID) (bool, error) {
	err := repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userId).
		Update("full_name", profileChunks.Name).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdatePasswordById updates the password of a user identified by their userId.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - newPassword: The new password to be set.
// - userId: The unique identifier of the user.
//
// Returns:
// - bool: True if the password is successfully updated, false otherwise.
// - error: An error if the update operation fails.
func (repo *authRepository) UpdatePasswordById(ctx context.Context, newPassword string, userId uuid.UUID) (bool, error) {
	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}

	// Update password
	err = repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userId).
		Update("password", string(hashedPassword)).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

// DestroyAllToken deletes all tokens associated with a specific user ID.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - userId: The unique identifier of the user whose tokens are to be deleted.
// Returns:
// - error: An error if the deletion fails, nil otherwise.
func (repo *authRepository) DestroyAllToken(ctx context.Context, userId uuid.UUID) error {
	err := repo.DB.WithContext(ctx).
		Where("user_id = ?", userId).
		Delete(&models.JWTToken{}).Error

	if err != nil {
		fmt.Println("Error deleting all tokens:", err)
		return err
	}
	return nil
}
