package token_storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rendyfutsuy/base-go/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	TokenStorageLocal = "local"
	TokenStorageRedis = "redis"
)

var (
	tokenStorageOnce    sync.Once
	defaultTokenStorage TokenStorage
)

func GetTokenStorage(driver string, db *gorm.DB, redisClient *redis.Client) (TokenStorage, error) {
	switch driver {
	case TokenStorageLocal:
		return NewLocalStorage(db), nil
	case TokenStorageRedis:
		return NewRedisStorage(redisClient, db), nil
	default:
		zap.S().Errorf("unsupported token storage driver: %s", driver)
		return nil, fmt.Errorf("unsupported token storage driver: %s", driver)
	}
}

func InitTokenStorage(driver string, db *gorm.DB, redisClient *redis.Client) error {
	var err error
	tokenStorageOnce.Do(func() {
		defaultTokenStorage, err = GetTokenStorage(driver, db, redisClient)
	})
	return err
}

func GetTokenStorageInstance() (TokenStorage, error) {
	if defaultTokenStorage == nil {
		return nil, fmt.Errorf("token storage not initialized")
	}
	return defaultTokenStorage, nil
}

func SetTokenStorage(storage TokenStorage) {
	defaultTokenStorage = storage
}

// Wrapper functions

func SaveSession(ctx context.Context, user models.User, accessToken, refreshToken, accessJTI, refreshJTI string, refreshTTL time.Duration) error {
	s, err := GetTokenStorageInstance()
	if err != nil {
		return err
	}
	return s.SaveSession(ctx, user, accessToken, refreshToken, accessJTI, refreshJTI, refreshTTL)
}

func GetRefreshTokenMetadata(ctx context.Context, refreshJTI string) (RefreshTokenMeta, error) {
	s, err := GetTokenStorageInstance()
	if err != nil {
		return RefreshTokenMeta{}, err
	}
	return s.GetRefreshTokenMetadata(ctx, refreshJTI)
}

func MarkRefreshTokenUsed(ctx context.Context, refreshJTI string) error {
	s, err := GetTokenStorageInstance()
	if err != nil {
		return err
	}
	return s.MarkRefreshTokenUsed(ctx, refreshJTI)
}

func DestroySession(ctx context.Context, accessToken string) error {
	s, err := GetTokenStorageInstance()
	if err != nil {
		return err
	}
	return s.DestroySession(ctx, accessToken)
}

func RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	s, err := GetTokenStorageInstance()
	if err != nil {
		return err
	}
	return s.RevokeAllUserSessions(ctx, userID)
}

func ValidateAccessToken(ctx context.Context, accessToken string) (models.User, error) {
	s, err := GetTokenStorageInstance()
	if err != nil {
		return models.User{}, err
	}
	return s.ValidateAccessToken(ctx, accessToken)
}
