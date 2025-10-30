package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rendyfutsuy/base-go/constants"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth"
	"github.com/rendyfutsuy/base-go/modules/auth/dto"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/services"
	"golang.org/x/crypto/bcrypt"
)

type authRepository struct {
	Conn         *sql.DB
	EmailService *services.EmailService
	QueueClient  *asynq.Client
}

func NewAuthRepository(Conn *sql.DB, EmailService *services.EmailService) auth.Repository {

	QueueClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     utils.ConfigVars.String("redis.address"),
		Password: utils.ConfigVars.String("redis.password"),
		DB:       utils.ConfigVars.Int("redis.db"),
	})

	return &authRepository{
		Conn,
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
	// SQL query to retrieve the user with the given email.
	query := `
		SELECT 
			id, email, password 
		FROM 
			users 
		WHERE 
			(email = $1 OR username = $1)
		AND deleted_at IS NULL 
		AND is_active = true
	`

	// Execute the query and scan the result into the user struct.
	err = repo.Conn.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Email, &user.Password)

	// Handle the error.
	if err != nil {
		// Print an error message if scanning the row fails.
		fmt.Println(constants.SQLErrorScanRow, err)

		// Handle case where no row is found.
		if err == sql.ErrNoRows {
			log.Printf("No user found with email/username: %s", login)
			return user, errors.New(constants.UserInvalid)
		}

		// Log other errors for debugging.
		log.Printf(constants.SQLErrorQueryRow, err)
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

	// get user from database by email
	query := `SELECT password FROM users WHERE id = $1 AND deleted_at IS NULL AND is_active = true`

	// Initialize an empty hashedPassword variable
	var hashedPassword string

	// Execute the query and scan the result into the user struct
	// Get user's registered Password
	// append to hashedPassword
	err := repo.Conn.QueryRowContext(ctx, query, userId).Scan(&hashedPassword)

	// Handle the error, such as not finding the user or database errors
	// if user not active and soft deleted, return error
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(constants.UserNotFound)
			return false, errors.New(constants.UserInvalid)
		}
		fmt.Println(constants.SQLErrorQueryDatabase, err)
		return false, err
	}

	// Compare the provided password with the hashed password from the database
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		// password do not match, add counter on users table
		repo.Conn.ExecContext(
			ctx,
			`UPDATE users SET counter = counter + 1, updated_at = $1 WHERE id = $2 RETURNING id`,
			time.Now().UTC(),
			userId,
		)

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

	// Query the password history
	query := "SELECT hashed_password FROM password_histories WHERE user_id = $1"

	rows, err := repo.Conn.QueryContext(ctx, query, userId)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var oldHashedPassword string
		if err := rows.Scan(&oldHashedPassword); err != nil {
			log.Fatal(err)
		}

		// Compare the new password with each old hashed password
		err = bcrypt.CompareHashAndPassword([]byte(oldHashedPassword), []byte(newPassword))
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

	// get user from database by id
	query := `SELECT password_expired_at FROM users WHERE id = $1`

	// Initialize an empty hashedPassword variable
	var expirationDate time.Time

	// Execute the query and scan the result into the user struct
	// Get user's registered Password
	// append to expirationDate
	err := repo.Conn.QueryRowContext(ctx, query, userId).Scan(&expirationDate)

	// Handle the error, such as not finding the user or database errors
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(constants.UserNotFound)
			return false, errors.New(constants.UserNotFound)
		}
		fmt.Println(constants.SQLErrorQueryDatabase, err)
		return false, err
	}

	// Get the current time
	currentTime := time.Now()

	// Compare expirationDate with currentTime to check if it has passed
	if expirationDate.Before(currentTime) {
		// Password has expired
		return true, errors.New("password has expired, please change your password now")
	}

	// if password not expired, return false
	return false, nil
}

func (repo *authRepository) AssertPasswordAttemptPassed(ctx context.Context, userId uuid.UUID) (bool, error) {

	// get user from database by id
	query := `SELECT counter FROM users WHERE id = $1`

	// Initialize an empty hashedPassword variable
	var attempt int

	// Execute the query and scan the result into the attempt variable
	err := repo.Conn.QueryRowContext(ctx, query, userId).Scan(&attempt)

	// Handle the error, such as not finding the user or database errors
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println(constants.UserNotFound)
			return false, errors.New(constants.UserNotFound)
		}
		fmt.Println(constants.SQLErrorQueryDatabase, err)
		return false, err
	}

	// if attempt above or equals t0 3, return false
	if attempt >= 3 {
		// Password has expired
		return false, errors.New("Password Attempt is above 3, you're blocked. please contact admin")
	}

	return true, nil
}

func (repo *authRepository) ResetPasswordAttempt(ctx context.Context, userId uuid.UUID) error {

	// reset attempt to 0
	repo.Conn.ExecContext(
		ctx,
		`UPDATE users SET counter = 0, updated_at = $1 WHERE id = $2 RETURNING id`,
		time.Now().UTC(),
		userId,
	)

	return nil
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

	// insert new access token record into database
	_, err := repo.Conn.ExecContext(
		ctx,
		`INSERT INTO jwt_tokens
			(access_token, user_id, created_at, updated_at) 
		VALUES 
			($1, $2, $3, $4) 
		RETURNING access_token, user_id`,
		accessToken,
		userId,
		time.Now().UTC(),
		time.Now().UTC(),
	)

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

	// insert new access token record into database
	_, err := repo.Conn.ExecContext(
		ctx,
		`INSERT INTO password_histories
			(hashed_password, user_id, created_at, updated_at) 
		VALUES 
			($1, $2, $3, $4) 
		RETURNING hashed_password, user_id`,
		hashedPassword,
		userId,
		time.Now().UTC(),
		time.Now().UTC(),
	)

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
	// SQL query to retrieve the user with the given email.
	query := `
		SELECT 
			usr.id as id,
			usr.full_name as full_name,
			usr.email,
			usr.role_id,
			roles.name
		FROM 
			users usr
		JOIN
			jwt_tokens jwt
		on
			jwt.user_id = usr.id
		JOIN
			roles
		on
			roles.id = usr.role_id
		WHERE 
			jwt.access_token = $1
	`
	// Execute the query and scan the result into the user struct.
	err := repo.Conn.QueryRowContext(ctx, query, accessToken).Scan(&user.ID, &user.FullName, &user.Email, &user.RoleId, &user.RoleName)

	// Handle the error.
	if err != nil {
		// Print an error message if scanning the row fails.
		fmt.Println(constants.SQLErrorScanRow, err)

		// Handle case where no row is found.
		if err == sql.ErrNoRows {
			log.Printf("No user found with this")
			return user, errors.New("User Not Found, the access token is not valid please re-login")
		}

		// Log other errors for debugging.
		log.Printf(constants.SQLErrorQueryRow, err)
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
	// SQL query to delete the user with the given access token.
	query := `
		DELETE FROM jwt_tokens
		WHERE access_token = $1
	`
	// Execute the query and delete requested row.
	_, err := repo.Conn.ExecContext(ctx, query, accessToken)

	// Handle the error.
	if err != nil {
		// Print an error message if delete row fails.
		fmt.Println(constants.SQLErrorScanRow, err)
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
	// SQL query to retrieve the user with the given access token.
	// cast is_active as Active if true
	// cast is_active as In Active if false
	query := `
		SELECT 
			usr.id AS user_id,
			usr.email,
			usr.full_name AS name,
			rls.name AS role,
			CASE 
				WHEN usr.is_active THEN 'Active' 
				ELSE 'In Active' 
			END AS is_active,
			usr.gender
		FROM 
			users usr
		JOIN
			roles rls
		ON
			rls.id = usr.role_id
		JOIN
			jwt_tokens jwt
		on
			jwt.user_id = usr.id
		WHERE 
			jwt.access_token = $1
		AND usr.deleted_at IS NULL 
	`

	// Execute the query and scan the result into the profile struct.
	err = repo.Conn.QueryRowContext(ctx, query, accessToken).Scan(
		&profile.UserId,
		&profile.Email,
		&profile.Name,
		&profile.Role,
		&profile.Status,
		&profile.Gender)

	// Handle the error.
	if err != nil {
		// Print an error message if scanning the row fails.
		fmt.Println(constants.SQLErrorScanRow, err)

		// Handle case where no row is found.
		if err == sql.ErrNoRows {
			log.Printf("No user found")
			return profile, errors.New(constants.UserInvalid)
		}

		// Log other errors for debugging.
		log.Printf(constants.SQLErrorQueryRow, err)
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
	// Update Profile
	// column updated: full_name
	query := `
		UPDATE users
		SET full_name = $2
		WHERE id = $1
	`

	// Execute the query and scan the result into the profile struct.
	_, err := repo.Conn.ExecContext(ctx, query, userId, profileChunks.Name)

	// Handle the error.
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

	// Update Profile
	// column updated: password
	query := `
		UPDATE users
		SET password = $2
		WHERE id = $1
	`

	// Execute the query with the hashed password
	_, err = repo.Conn.ExecContext(ctx, query, userId, string(hashedPassword))
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
	// SQL query to delete the user with the given user ID.
	query := `
		DELETE FROM jwt_tokens
		WHERE user_id = $1
	`
	// Execute the query and delete requested row.
	_, err := repo.Conn.ExecContext(ctx, query, userId)

	// Handle the error.
	if err != nil {
		// Print an error message if delete row fails.
		fmt.Println(constants.SQLErrorScanRow, err)
		return err
	}
	return nil
}
