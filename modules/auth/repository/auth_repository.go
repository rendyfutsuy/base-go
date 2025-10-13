package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
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
// - email: The email of the user to retrieve.
//
// Returns:
// - user: The retrieved user.
// - err:  An error if the retrieval fails.
func (repo *authRepository) FindByEmailOrUsername(login string) (user models.User, err error) {
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
	err = repo.Conn.QueryRow(query, login).Scan(&user.ID, &user.Email, &user.Password)

	// Handle the error.
	if err != nil {
		// Print an error message if scanning the row fails.
		fmt.Println("Error scanning row:", err)

		// Handle case where no row is found.
		if err == sql.ErrNoRows {
			log.Printf("No user found with email/username: %s", login)
			return user, errors.New("User Not Found, please check to Customer Services...")
		}

		// Log other errors for debugging.
		log.Printf("QueryRow scan error: %v", err)
		return user, err
	}

	return user, nil
}

// AssertPasswordRight checks if the provided password matches the hashed password in the database for the given user ID.
//
// Parameters:
// - password: The password to compare.
// - userId: The unique identifier of the user.
//
// Returns:
// - bool: True if the passwords match, false otherwise.
// - error: An error if the comparison fails or if there are database errors.
func (repo *authRepository) AssertPasswordRight(password string, userId uuid.UUID) (bool, error) {

	// get user from database by email
	query := `SELECT password FROM users WHERE id = $1 AND deleted_at IS NULL AND is_active = true`

	// Initialize an empty hashedPassword variable
	var hashedPassword string

	// Execute the query and scan the result into the user struct
	// Get user's registered Password
	// append to hashedPassword
	err := repo.Conn.QueryRow(query, userId).Scan(&hashedPassword)

	// Handle the error, such as not finding the user or database errors
	// if user not active and soft deleted, return error
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User not found")
			return false, errors.New("User Not Found, please check to Customer Services...")
		}
		fmt.Println("Error querying database:", err)
		return false, err
	}

	// Compare the provided password with the hashed password from the database
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		// password do not match, add counter on users table
		repo.Conn.QueryRow(
			`UPDATE users SET counter = counter + 1, updated_at = $1 WHERE id = $2 RETURNING id`,
			time.Now(),
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
func (repo *authRepository) AssertPasswordNeverUsesByUser(newPassword string, userId uuid.UUID) (bool, error) {

	// Query the password history
	query := "SELECT hashed_password FROM password_histories WHERE user_id = $1"

	rows, err := repo.Conn.Query(query, userId)

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
// - userId: The unique identifier of the user.
//
// Returns:
// - bool: True if the password has expired, false otherwise.
// - error: An error if the query to the database fails.
func (repo *authRepository) AssertPasswordExpiredIsPassed(userId uuid.UUID) (bool, error) {

	// get user from database by id
	query := `SELECT password_expired_at FROM users WHERE id = $1`

	// Initialize an empty hashedPassword variable
	var expirationDate time.Time

	// Execute the query and scan the result into the user struct
	// Get user's registered Password
	// append to expirationDate
	err := repo.Conn.QueryRow(query, userId).Scan(&expirationDate)

	// Handle the error, such as not finding the user or database errors
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User not found")
			return false, errors.New("User Not Found")
		}
		fmt.Println("Error querying database:", err)
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

func (repo *authRepository) AssertPasswordAttemptPassed(userId uuid.UUID) (bool, error) {

	// get user from database by id
	query := `SELECT counter FROM users WHERE id = $1`

	// Initialize an empty hashedPassword variable
	var attempt int

	// Execute the query and scan the result into the attempt variable
	err := repo.Conn.QueryRow(query, userId).Scan(&attempt)

	// Handle the error, such as not finding the user or database errors
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("User not found")
			return false, errors.New("User Not Found")
		}
		fmt.Println("Error querying database:", err)
		return false, err
	}

	// if attempt above or equals t0 3, return false
	if attempt >= 3 {
		// Password has expired
		return false, errors.New("Password Attempt is above 3, you're blocked. please contact admin")
	}

	return true, nil
}

// AddUserAccessToken inserts a new access token for a user into the database.
//
// Parameters:
// - accessToken: The access token to be inserted.
// - userId: The unique identifier of the user.
//
// Returns:
// - error: An error if the insertion fails.
func (repo *authRepository) AddUserAccessToken(accessToken string, userId uuid.UUID) error {

	// insert new access token record into database
	err := repo.Conn.QueryRow(
		`INSERT INTO jwt_tokens
			(access_token, user_id, created_at, updated_at) 
		VALUES 
			($1, $2, $3, $4) 
		RETURNING access_token, user_id`,
		accessToken,
		userId,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err.Err()
	}

	return nil
}

// AddUserAccessToken inserts a new access token for a user into the database.
//
// Parameters:
// - accessToken: The access token to be inserted.
// - userId: The unique identifier of the user.
//
// Returns:
// - error: An error if the insertion fails.
func (repo *authRepository) AddPasswordHistory(hashedPassword string, userId uuid.UUID) error {

	// insert new access token record into database
	err := repo.Conn.QueryRow(
		`INSERT INTO password_histories
			(hashed_password, user_id, created_at, updated_at) 
		VALUES 
			($1, $2, $3, $4) 
		RETURNING hashed_password, user_id`,
		hashedPassword,
		userId,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err.Err()
	}

	return nil
}

// GetUserByAccessToken retrieves a user from the database based on the provided access token.
//
// Parameters:
// - accessToken: The access token used to identify the user.
//
// Returns:
// - user: The retrieved user object.
// - errorMain: An error if the retrieval fails, or if the access token is not valid.
func (repo *authRepository) GetUserByAccessToken(accessToken string) (user models.User, errorMain error) {
	// SQL query to retrieve the user with the given email.
	query := `
		SELECT 
			usr.id as id, usr.email
		FROM 
			users usr
		JOIN
			jwt_tokens jwt
		on
			jwt.user_id = usr.id
		WHERE 
			jwt.access_token = $1
	`
	// Execute the query and scan the result into the user struct.
	err := repo.Conn.QueryRow(query, accessToken).Scan(&user.ID, &user.Email)

	// Handle the error.
	if err != nil {
		// Print an error message if scanning the row fails.
		fmt.Println("Error scanning row:", err)

		// Handle case where no row is found.
		if err == sql.ErrNoRows {
			log.Printf("No user found with this")
			return user, errors.New("User Not Found, the access token is not valid please re-login")
		}

		// Log other errors for debugging.
		log.Printf("QueryRow scan error: %v", err)
		return user, err
	}

	return user, nil
}

// DestroyToken deletes a JWT token from the database.
//
// Parameters:
// - accessToken: The access token to be deleted.
//
// Returns:
// - error: An error if the deletion fails, nil otherwise.
func (repo *authRepository) DestroyToken(accessToken string) error {
	// SQL query to delete the user with the given access token.
	query := `
		DELETE FROM jwt_tokens
		WHERE access_token = $1
	`
	// Execute the query and delete requested row.
	_, err := repo.Conn.Exec(query, accessToken)

	// Handle the error.
	if err != nil {
		// Print an error message if delete row fails.
		fmt.Println("Error scanning row:", err)
		return err
	}
	return nil
}

// FindByCurrentSession retrieves user profile based on the provided access token.
//
// Parameters:
// - accessToken: The access token used to identify the user session.
//
// Returns:
// - profile: User profile information.
// - err: An error if the retrieval fails, nil otherwise.
func (repo *authRepository) FindByCurrentSession(accessToken string) (profile dto.UserProfile, err error) {
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
	err = repo.Conn.QueryRow(query, accessToken).Scan(
		&profile.UserId,
		&profile.Email,
		&profile.Name,
		&profile.Role,
		&profile.Status,
		&profile.Gender)

	// Handle the error.
	if err != nil {
		// Print an error message if scanning the row fails.
		fmt.Println("Error scanning row:", err)

		// Handle case where no row is found.
		if err == sql.ErrNoRows {
			log.Printf("No user found")
			return profile, errors.New("User Not Found, please check to Customer Services...")
		}

		// Log other errors for debugging.
		log.Printf("QueryRow scan error: %v", err)
		return profile, err
	}
	return profile, nil
}

// UpdateProfileById updates the full name of a user profile by ID.
//
// Parameters:
// - profileChunks: The updated profile information.
// - userId: The unique identifier of the user.
//
// Returns:
// - bool: True if the profile was successfully updated, false otherwise.
// - error: An error if the update fails.
func (repo *authRepository) UpdateProfileById(profileChunks dto.ReqUpdateProfile, userId uuid.UUID) (bool, error) {
	// Update Profile
	// column updated: full_name
	query := `
		UPDATE users
		SET full_name = $2
		WHERE id = $1
	`

	// Execute the query and scan the result into the profile struct.
	err := repo.Conn.QueryRow(query, userId, profileChunks.Name)

	// Handle the error.
	if err != nil {
		return false, err.Err()
	}

	return true, nil
}

// UpdatePasswordById updates the password of a user identified by their userId.
//
// Parameters:
// - newPassword: The new password to be set.
// - userId: The unique identifier of the user.
//
// Returns:
// - bool: True if the password is successfully updated, false otherwise.
// - error: An error if the update operation fails.
func (repo *authRepository) UpdatePasswordById(newPassword string, userId uuid.UUID) (bool, error) {
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
	_, err = repo.Conn.Exec(query, userId, string(hashedPassword))
	if err != nil {
		return false, err
	}

	return true, nil
}

// DestroyAllToken deletes all tokens associated with a specific user ID.
//
// Parameters:
// - userId: The unique identifier of the user whose tokens are to be deleted.
// Returns:
// - error: An error if the deletion fails, nil otherwise.
func (repo *authRepository) DestroyAllToken(userId uuid.UUID) error {
	// SQL query to delete the user with the given user ID.
	query := `
		DELETE FROM jwt_tokens
		WHERE user_id = $1
	`
	// Execute the query and delete requested row.
	_, err := repo.Conn.Exec(query, userId)

	// Handle the error.
	if err != nil {
		// Print an error message if delete row fails.
		fmt.Println("Error scanning row:", err)
		return err
	}
	return nil
}
