package constants

import "errors"

var (
	ErrPasswordExpired  = errors.New("password has expired")
	ErrTokenRevoked     = errors.New("token is revoked please re-login from login form again..")
	TokenRevokedMessage = "token is revoked please re-login from login form again.."
)
