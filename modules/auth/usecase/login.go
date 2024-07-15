package usecase

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func (u *authUsecase) Authenticate(c echo.Context, login string, password string) (string, error) {
	// get user by email
	user, err := u.authRepo.FindByEmailOrUsername(login)

	// if fail to get user return error
	if err != nil {
		return "", err

	}

	// assert login attempt is not above 3
	isAttemptPassed, err := u.authRepo.AssertPasswordAttemptPassed(user.ID)
	if !isAttemptPassed {
		return "", err
	}

	// assert password given is same with saved password
	isPasswordRight, err := u.authRepo.AssertPasswordRight(password, user.ID)

	// if password fail to match return error
	if !isPasswordRight {
		return "", err
	}

	// assert if password expiration passed
	isPasswordExpired, err := u.authRepo.AssertPasswordExpiredIsPassed(user.ID)

	// if password expired return error
	if isPasswordExpired {
		return "", err
	}

	// Generate JWT token
	// access token would always expire in next day on 03:00 AM WIB (UTC+7)
	claims := AuthClaims{
		UserID: user.ID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(u.expireDuration).Unix(),
		},
	}

	// append access token to return
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// assign signing key
	accessToken, err := token.SignedString(u.signingKey)

	// if fail to assign return error
	if err != nil {
		return "", err
	}

	// record access token to database
	err = u.authRepo.AddUserAccessToken(accessToken, user.ID)

	// if fail to assign return error
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
