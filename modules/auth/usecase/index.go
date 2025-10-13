package usecase

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rendyfutsuy/base-go.git/modules/auth"
)

type AuthClaims struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}

type authUsecase struct {
	authRepo       auth.Repository
	contextTimeout time.Duration
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUsecase(r auth.Repository, timeout time.Duration, hashSalt string, signingKey []byte) auth.Usecase {
	// Expire Time Calculation BEGIN
	// Determine the current time in UTC+7 (Asia/Bangkok timezone)
	loc := time.FixedZone("UTC+7", 7*60*60) // UTC+7 is 7 hours ahead of UTC
	now := time.Now().In(loc)

	// Calculate the next 03:00 AM UTC+7 time
	next03AM := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, loc)
	if now.After(next03AM) {
		next03AM = next03AM.Add(24 * time.Hour) // Add 24 hours if current time is after 03:00 AM
	}

	// Calculate expireDuration based on the difference between now and next 03:00 AM UTC+7
	expireDuration := next03AM.Sub(now)
	// Expire Time Calculation END

	return &authUsecase{
		authRepo:       r,
		contextTimeout: timeout,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: expireDuration,
	}
}
