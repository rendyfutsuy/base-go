package constants

import "errors"

var (
	ErrPasswordExpired  = errors.New("password has expired")
	ErrTokenRevoked     = errors.New("token is revoked please re-login from login form again..")
	TokenRevokedMessage = "token is revoked please re-login from login form again.."
)

const (
	// Authentication errors
	AuthPasswordNotMatch         = "Password Not Match"
	AuthPasswordAlreadyUsed      = "Youre already used this password, please try another one.."
	AuthPasswordExpiredMessage   = "password has expired, please change your password now"
	AuthPasswordAttemptExceeded  = "Password Attempt is above 3, you're blocked. please contact admin"
	AuthTokenParseFailed         = "failed to parse token"
	AuthTokenMissingJTI          = "token does not contain jti claim"
	AuthTokenInvalid             = "User Not Found, the access token is not valid please re-login"
	AuthTokenInvalidRestart      = "User Not Found, the access token is not valid please restart the process..."
	AuthTooManyPasswordAttempts  = "too many password attempts"
	AuthInvalidCredentials       = "invalid credentials"
	AuthUsernamePasswordNotFound = "Username/Password not found..."
	AuthOldPasswordNotMatch      = "Old Password not Match"
	AuthNewPasswordSameAsOld     = "New Password should not be same with Current Password"
	AuthEmailNotFound            = "Email not found, please be sure no typo on email input..."
	AuthInvalidToken             = "Invalid token"

	// Success messages
	AuthResetEmailSent         = "Successfully Send Reset Email Request"
	AuthPasswordResetSuccess   = "Successfully Reset Password"
	AuthLogoutSuccess          = "Successfully Logged Out"
	AuthProfileUpdated         = "Successfully Updated Profile"
	AuthPasswordUpdated        = "Successfully Updated My Password"
	FirstTimeLoginErrorMessage = "User has not been activate, please change your password now..."
	AuthPasswordAlreadyChanged = "User has change, if you want to change your password again. please contact admin"

	// define role
	AuthRoleSuperAdmin = "Super Admin"
)
