package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rendyfutsuy/base-go/constants"
	models "github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/auth/tasks"
	"github.com/rendyfutsuy/base-go/utils"
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
	// SQL query to retrieve the user with the given email.
	query := `
		SELECT 
			usr.id as id, usr.email
		FROM 
			users usr
		JOIN
			reset_password_tokens rst_pwd
		on
			rst_pwd.user_id = usr.id
		WHERE 
			rst_pwd.access_token = $1
	`
	// Execute the query and scan the result into the user struct.
	err := repo.Conn.QueryRowContext(ctx, query, token).Scan(&user.ID, &user.Email)

	// Handle the error.
	if err != nil {
		// Print an error message if scanning the row fails.
		fmt.Println(constants.SQLErrorScanRow, err)

		// Handle case where no row is found.
		if err == sql.ErrNoRows {
			log.Printf("No user found with this")
			return user, errors.New("User Not Found, the access token is not valid please restart the process...")
		}

		// Log other errors for debugging.
		log.Printf(constants.SQLErrorQueryRow, err)
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

	// insert new access token record into database
	_, err := repo.Conn.ExecContext(
		ctx,
		`INSERT INTO reset_password_tokens
				(access_token, user_id, created_at, updated_at) 
			VALUES 
				($1, $2, $3, $4) 
			RETURNING access_token, user_id`,
		token,
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

// DestroyResetPasswordToken deletes the reset password token from the database based on the provided token.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - token: The access token to identify the reset password token.
// Returns:
// - error: An error if the deletion fails, nil otherwise.
func (repo *authRepository) DestroyResetPasswordToken(ctx context.Context, token string) error {
	// SQL query to delete the user with the given access token.
	query := `
		DELETE FROM reset_password_tokens
		WHERE access_token = $1
	`
	// Execute the query and delete requested row.
	_, err := repo.Conn.ExecContext(ctx, query, token)

	// Handle the error.
	if err != nil {
		// Print an error message if delete row fails.
		fmt.Println(constants.SQLErrorScanRow, err)
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
	// SQL query to delete the user with the given access token.
	query := `
		DELETE FROM reset_password_tokens
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

func (repo *authRepository) IncreasePasswordExpiredAt(ctx context.Context, userId uuid.UUID) error {
	// Calculate the expiration date to be 3 months from today
	expiredAt := time.Now().AddDate(0, 3, 0)

	// SQL query to delete the user with the given access token.
	query := `
		UPDATE users
		SET password_expired_at = $2
		WHERE id = $1
	`
	// Execute the query and delete requested row.
	_, err := repo.Conn.ExecContext(ctx, query, userId, expiredAt)

	// Handle the error.
	if err != nil {
		// Print an error message if delete row fails.
		fmt.Println(constants.SQLErrorScanRow, err)
		return err
	}
	return nil
}
