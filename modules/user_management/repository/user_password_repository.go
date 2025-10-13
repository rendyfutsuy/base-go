package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
)

// IsUserPasswordCanUpdated checks if the password of a user with the given ID can be updated.
//
// Parameters:
// - id: The ID of the user.
//
// Returns:
// - user: The user object if found, nil otherwise.
// - err: An error if the user is not found, blocked, or inactive.
func (repo *userRepository) IsUserPasswordCanUpdated(id uuid.UUID) (bool, error) {
	// initialize user variable
	user := new(models.User)

	// fetch data from database by id that passed
	// assign return value to user variable
	user, err := repo.GetUserByID(id)

	if err != nil {
		return false, errors.New(constants.UserNotFound)
	}

	// assert user is not blocked
	if user.IsBlocked == true {
		return false, errors.New("User is Blocked, Please Unblock this user first")
	}

	// assert user is active
	if user.IsActive == false {
		return false, errors.New("User is Inactive, Please Activate this user first")
	}

	return true, nil
}
