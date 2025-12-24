package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rendyfutsuy/base-go/constants"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/services"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authRepository struct {
	DB           *gorm.DB
	EmailService *services.EmailService
	QueueClient  *asynq.Client
	Redis        *redis.Client
}

func NewAuthRepository(DB *gorm.DB, EmailService *services.EmailService, RedisClient *redis.Client) auth.Repository {

	QueueClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     utils.ConfigVars.String("redis.address"),
		Password: utils.ConfigVars.String("redis.password"),
		DB:       utils.ConfigVars.Int("redis.db"),
	})

	return &authRepository{
		DB:           DB,
		EmailService: EmailService,
		QueueClient:  QueueClient,
		Redis:        RedisClient,
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
func (repo *authRepository) FindByEmailOrUsername(ctx context.Context, login string) (models.User, error) {
	var dbUser models.User

	// Query actual table "users"
	err := repo.DB.WithContext(ctx).
		Where("(email = ? OR username = ?) AND deleted_at IS NULL AND is_active = ?",
			login, login, true).
		First(&dbUser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New(constants.UserInvalid)
		}
		return models.User{}, fmt.Errorf("failed querying user: %w", err)
	}

	// Map DB struct → models.User for usecase
	return models.User{
		ID:       dbUser.ID,
		Email:    dbUser.Email,
		Username: dbUser.Username,
	}, nil
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
		return false, errors.New(constants.AuthPasswordNotMatch)
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
			return false, fmt.Errorf(constants.AuthPasswordAlreadyUsed)
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
		return true, errors.New(constants.AuthPasswordExpiredMessage)
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
		return false, errors.New(constants.AuthPasswordAttemptExceeded)
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

// extractJTIFromToken extracts JWT ID (jti) from token string
func (repo *authRepository) extractJTIFromToken(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// We don't validate here, just extract claims
		return nil, nil
	})
	if err != nil {
		// If parsing fails, try to extract jti directly without validation
		// Parse without validation to get claims
		parser := jwt.NewParser(jwt.WithValidMethods([]string{}))
		_, _, err := parser.ParseUnverified(tokenString, claims)
		if err != nil {
			return "", fmt.Errorf("%s: %w", constants.AuthTokenParseFailed, err)
		}
	}
	if claims.ID == "" {
		return "", errors.New(constants.AuthTokenMissingJTI)
	}
	return claims.ID, nil
}

// getFullUserData retrieves full user data including role information
func (repo *authRepository) getFullUserData(ctx context.Context, userId uuid.UUID) (*models.User, error) {
	var user models.User
	err := repo.DB.WithContext(ctx).
		Table("users usr").
		Select("usr.id, usr.full_name, usr.email, usr.username, usr.is_active, usr.gender, usr.role_id, roles.name as role_name").
		Joins("LEFT JOIN roles ON roles.id = usr.role_id AND roles.deleted_at IS NULL").
		Where("usr.id = ? AND usr.deleted_at IS NULL", userId).
		Scan(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// AddUserAccessToken inserts a new access token for a user into Redis using redis_repository helper.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - accessToken: The access token to be inserted.
// - userId: The unique identifier of the user.
//
// Returns:
// - error: An error if the insertion fails.
func (repo *authRepository) AddUserAccessToken(ctx context.Context, accessToken string, userId uuid.UUID) error {
	// Extract JTI from token
	jti, err := repo.extractJTIFromToken(accessToken)
	if err != nil {
		// Fallback to using accessToken as key if jti extraction fails
		jti = accessToken
		utils.Logger.Warn("Failed to extract jti from token, using accessToken as key", zap.Error(err))
	}

	// Get full user data
	user, err := repo.getFullUserData(ctx, userId)
	if err != nil {
		utils.Logger.Error("Failed to get user data", zap.Error(err))
		return err
	}

	// Get TTL from config
	ttlSeconds := utils.ConfigVars.Int("auth.redis_ttl_seconds")
	if ttlSeconds <= 0 {
		ttlSeconds = 2 * 24 * 60 * 60 // Default 2 days
	}
	ttl := time.Duration(ttlSeconds) * time.Second

	// Use CreateSession helper from redis_repository.go
	if err := repo.CreateSession(ctx, jti, user, ttl); err != nil {
		return err
	}

	// Maintain a set of jtis per user for efficient DestroyAllToken
	userSetKey := fmt.Sprintf("auth:user_tokens:%s", userId.String())
	if err := repo.Redis.SAdd(ctx, userSetKey, jti).Err(); err != nil {
		utils.Logger.Warn("Failed to add jti to user token set", zap.Error(err))
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

// GetUserByAccessToken retrieves a user from Redis using redis_repository helper.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - accessToken: The access token used to identify the user.
//
// Returns:
// - user: The retrieved user object.
// - errorMain: An error if the retrieval fails, or if the access token is not valid.
func (repo *authRepository) GetUserByAccessToken(ctx context.Context, accessToken string) (user models.User, errorMain error) {
	// Extract JTI from token
	jti, err := repo.extractJTIFromToken(accessToken)
	if err != nil {
		// Fallback to using accessToken as key if jti extraction fails
		jti = accessToken
		utils.Logger.Warn("Failed to extract jti from token, using accessToken as key", zap.Error(err))
	}

	// Use GetSessionData helper from redis_repository.go
	userData, rerr := repo.GetSessionData(ctx, jti)
	if rerr != nil {
		if rerr == redis.Nil {
			log.Printf("No session found")
			return user, errors.New(constants.UserInvalid)
		}
		log.Printf("Error querying redis session: %v", rerr)
		return user, rerr
	}

	if userData == nil {
		log.Printf("No user data found in session")
		return user, errors.New(constants.UserInvalid)
	}

	// Fetch user
	err = repo.DB.WithContext(ctx).
		Table("users").
		Select("id", "full_name", "email", "username", "is_active", "gender", "role_id", "is_first_time_login").
		Where("username = ? AND deleted_at IS NULL", userData.Username).
		First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

// DestroyToken deletes a JWT token from Redis using redis_repository helper.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - accessToken: The access token to be deleted.
//
// Returns:
// - error: An error if the deletion fails, nil otherwise.
func (repo *authRepository) DestroyToken(ctx context.Context, accessToken string) error {
	// Extract JTI from token
	jti, err := repo.extractJTIFromToken(accessToken)
	if err != nil {
		// Fallback to using accessToken as key if jti extraction fails
		jti = accessToken
		utils.Logger.Warn("Failed to extract jti from token, using accessToken as key", zap.Error(err))
	}

	// Get user data from session to get userId for set cleanup
	userData, err := repo.GetSessionData(ctx, jti)
	if err == nil && userData != nil {
		// Remove jti from user's token set
		userSetKey := fmt.Sprintf("auth:user_tokens:%s", userData.ID.String())
		_ = repo.Redis.SRem(ctx, userSetKey, jti).Err()
	}

	// Use DeleteSession helper from redis_repository.go
	return repo.DeleteSession(ctx, jti)
}

// FindByCurrentSession retrieves user based on the provided access token using redis_repository helper.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - accessToken: The access token used to identify the user session.
//
// Returns:
// - user: User information.
// - err: An error if the retrieval fails, nil otherwise.
func (repo *authRepository) FindByCurrentSession(ctx context.Context, accessToken string) (user models.User, err error) {
	// Extract JTI from token
	jti, rerr := repo.extractJTIFromToken(accessToken)
	if rerr != nil {
		// Fallback to using accessToken as key if jti extraction fails
		jti = accessToken
		utils.Logger.Warn("Failed to extract jti from token, using accessToken as key", zap.Error(rerr))
	}

	// Use GetSessionData helper from redis_repository.go
	userData, rerr := repo.GetSessionData(ctx, jti)
	if rerr != nil {
		if rerr == redis.Nil {
			log.Printf("No session found")
			return user, errors.New(constants.UserInvalid)
		}
		log.Printf("Error querying redis session: %v", rerr)
		return user, rerr
	}

	if userData == nil {
		log.Printf("No user data found in session")
		return user, errors.New(constants.UserInvalid)
	}

	// Fetch user
	err = repo.DB.WithContext(ctx).
		Table("users").
		Select("id", "full_name", "email", "username", "is_active", "gender", "role_id", "is_first_time_login").
		Where("username = ? AND deleted_at IS NULL", userData.Username).
		First(&user).Error

	if err != nil {
		return user, err
	}

	// If RoleName is empty in session, fetch from database
	if user.RoleName == "" && user.RoleId != uuid.Nil {
		// Fetch role name from database
		var role models.Role
		err := repo.DB.WithContext(ctx).
			Table("roles").
			Select("name").
			Where("id = ? AND deleted_at IS NULL", user.RoleId).
			First(&role).Error
		if err == nil {
			user.RoleName = role.Name
		}
	}

	return user, nil
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
	updates := map[string]interface{}{
		"password":            profileChunks.Name,
		"is_first_time_login": false,
	}
	err := repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userId).
		Updates(updates).Error

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

	// Update password and set is_first_time_login to false
	updates := map[string]interface{}{
		"password":            string(hashedPassword),
		"is_first_time_login": false,
	}
	err = repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userId).
		Updates(updates).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

// GetIsFirstTimeLogin retrieves the is_first_time_login status of a user by ID.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - userId: The unique identifier of the user.
//
// Returns:
// - bool: The is_first_time_login status of the user.
// - error: An error if the retrieval fails.
func (repo *authRepository) GetIsFirstTimeLogin(ctx context.Context, userId uuid.UUID) (bool, error) {
	var user models.User
	err := repo.DB.WithContext(ctx).
		Select("is_first_time_login").
		Where("id = ? AND deleted_at IS NULL", userId).
		First(&user).Error

	if err != nil {
		return false, err
	}

	return user.IsFirstTimeLogin, nil
}

// DestroyAllToken deletes all tokens associated with a specific user ID.
// Uses the user token set maintained in AddUserAccessToken for efficient deletion.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - userId: The unique identifier of the user whose tokens are to be deleted.
// Returns:
// - error: An error if the deletion fails, nil otherwise.
func (repo *authRepository) DestroyAllToken(ctx context.Context, userId uuid.UUID) error {
	// Get all jtis for this user from the set
	userSetKey := fmt.Sprintf("auth:user_tokens:%s", userId.String())
	jtis, err := repo.Redis.SMembers(ctx, userSetKey).Result()
	if err != nil && err != redis.Nil {
		utils.Logger.Error("Error fetching user token set", zap.Error(err))
		return err
	}

	// Delete all session keys (jtis)
	if len(jtis) > 0 {
		for _, jti := range jtis {
			if err := repo.DeleteSession(ctx, jti); err != nil && err != redis.Nil {
				utils.Logger.Warn("Error deleting session", zap.String("jti", jti), zap.Error(err))
				// Continue deleting other sessions even if one fails
			}
		}
	}

	// Delete the user token set itself
	_ = repo.Redis.Del(ctx, userSetKey).Err()

	return nil
}

func (repo *authRepository) StoreRefreshToken(
	ctx context.Context,
	refreshJTI string,
	userID uuid.UUID,
	accessJTI string,
	ttl time.Duration,
) error {

	tokenKey := fmt.Sprintf("auth:refresh:%s", refreshJTI)
	userSetKey := fmt.Sprintf("auth:user_refresh_tokens:%s", userID.String())

	expiresAt := time.Now().UTC().Add(ttl).Format(time.RFC3339)

	pipe := repo.Redis.TxPipeline()

	pipe.HSet(ctx, tokenKey, map[string]interface{}{
		"user_id":    userID.String(),
		"expires_at": expiresAt,
		"used":       "0",
		"access_jti": accessJTI, // NEW
	})
	pipe.Expire(ctx, tokenKey, ttl)

	pipe.SAdd(ctx, userSetKey, refreshJTI)
	pipe.ExpireNX(ctx, userSetKey, ttl)

	_, err := pipe.Exec(ctx)
	return err
}

func (repo *authRepository) GetRefreshTokenMetadata(
	ctx context.Context,
	jti string,
) (auth.RefreshTokenMeta, error) {

	tokenKey := fmt.Sprintf("auth:refresh:%s", jti)
	data, err := repo.Redis.HGetAll(ctx, tokenKey).Result()
	if err != nil || len(data) == 0 {
		return auth.RefreshTokenMeta{}, redis.Nil
	}

	uid, _ := uuid.Parse(data["user_id"])
	expAt, _ := time.Parse(time.RFC3339, data["expires_at"])
	used := data["used"] == "1"
	accessJTI := data["access_jti"]

	return auth.RefreshTokenMeta{
		UserID:    uid,
		ExpiresAt: expAt,
		Used:      used,
		AccessJTI: accessJTI,
	}, nil
}

func (repo *authRepository) MarkRefreshTokenUsed(ctx context.Context, jti string) error {
	tokenKey := fmt.Sprintf("auth:refresh:%s", jti)

	pipe := repo.Redis.TxPipeline()

	// Mark token as used
	pipe.HSet(ctx, tokenKey, "used", "1")

	// Optional: shorten TTL so used refresh tokens disappear faster
	// Set to e.g. 24 hours (or keep original TTL — your choice)
	pipe.Expire(ctx, tokenKey, 24*time.Hour)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to mark refresh token used: %w", err)
	}

	return nil
}

func (repo *authRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {

	accessSetKey := fmt.Sprintf("auth:user_tokens:%s", userID.String())
	refreshSetKey := fmt.Sprintf("auth:user_refresh_tokens:%s", userID.String())

	// 1) Get all access token JTIs
	accessJTIs, err := repo.Redis.SMembers(ctx, accessSetKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("failed to read user access tokens: %w", err)
	}

	// 2) Get all refresh token JTIs
	refreshJTIs, err := repo.Redis.SMembers(ctx, refreshSetKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("failed to read user refresh tokens: %w", err)
	}

	pipe := repo.Redis.TxPipeline()

	// --- Delete all access token sessions ---
	for _, jti := range accessJTIs {
		pipe.Del(ctx, fmt.Sprintf("auth:session:%s", jti))
	}

	// --- Delete all refresh token metadata ---
	for _, jti := range refreshJTIs {
		pipe.Del(ctx, fmt.Sprintf("auth:refresh:%s", jti))
	}

	// --- Clear user token sets ---
	pipe.Del(ctx, accessSetKey)
	pipe.Del(ctx, refreshSetKey)

	// Execute pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to revoke all user sessions: %w", err)
	}

	return nil
}
