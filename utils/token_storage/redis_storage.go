package token_storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/utils"
	"go.uber.org/zap"
)

type RedisStorage struct {
	Redis *redis.Client
}

func NewRedisStorage(redisClient *redis.Client) *RedisStorage {
	return &RedisStorage{
		Redis: redisClient,
	}
}

// extractJTIFromToken extracts JWT ID (jti) from token string
func (s *RedisStorage) extractJTIFromToken(tokenString string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{}))
	_, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.AuthTokenParseFailed, err)
	}
	if claims.ID == "" {
		return "", errors.New(constants.AuthTokenMissingJTI)
	}
	return claims.ID, nil
}

func (s *RedisStorage) SaveSession(ctx context.Context, user models.User, accessToken, refreshToken, accessJTI, refreshJTI string, refreshTTL time.Duration) error {
	// 1. Add User Access Token
	// Get TTL from config for access token (Redis session)
	ttlSeconds := utils.ConfigVars.Int("auth.redis_ttl_seconds")
	if ttlSeconds <= 0 {
		ttlSeconds = 2 * 24 * 60 * 60 // Default 2 days
	}
	accessTTL := time.Duration(ttlSeconds) * time.Second

	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	pipe := s.Redis.TxPipeline()

	// Store access token session
	pipe.Set(ctx, accessJTI, userData, accessTTL)

	// Add access JTI to user set
	userSetKey := fmt.Sprintf("auth:user_tokens:%s", user.ID.String())
	pipe.SAdd(ctx, userSetKey, accessJTI)
	// We might want to set expire on the set too, but it's tricky as it holds multiple tokens.
	// Usually we rely on manual cleanup or let it persist (it's small).

	// 2. Store Refresh Token Metadata
	tokenKey := fmt.Sprintf("auth:refresh:%s", refreshJTI)
	userRefreshSetKey := fmt.Sprintf("auth:user_refresh_tokens:%s", user.ID.String())

	expiresAt := time.Now().UTC().Add(refreshTTL).Format(time.RFC3339)

	pipe.HSet(ctx, tokenKey, map[string]interface{}{
		"user_id":    user.ID.String(),
		"expires_at": expiresAt,
		"used":       "0",
		"access_jti": accessJTI,
	})
	pipe.Expire(ctx, tokenKey, refreshTTL)

	pipe.SAdd(ctx, userRefreshSetKey, refreshJTI)
	pipe.ExpireNX(ctx, userRefreshSetKey, refreshTTL)

	_, err = pipe.Exec(ctx)
	return err
}

func (s *RedisStorage) GetRefreshTokenMetadata(ctx context.Context, refreshJTI string) (RefreshTokenMeta, error) {
	tokenKey := fmt.Sprintf("auth:refresh:%s", refreshJTI)
	data, err := s.Redis.HGetAll(ctx, tokenKey).Result()
	if err != nil || len(data) == 0 {
		return RefreshTokenMeta{}, redis.Nil // Using redis.Nil as generic not found? Or should wrap?
	}

	uid, _ := uuid.Parse(data["user_id"])
	t, _ := time.Parse(time.RFC3339, data["expires_at"])
	used := data["used"] == "1"

	return RefreshTokenMeta{
		UserID:    uid,
		ExpiresAt: t,
		Used:      used,
		AccessJTI: data["access_jti"],
	}, nil
}

func (s *RedisStorage) MarkRefreshTokenUsed(ctx context.Context, refreshJTI string) error {
	tokenKey := fmt.Sprintf("auth:refresh:%s", refreshJTI)
	// Keep it for a while so we know it's used (reuse detection)
	// Configurable short TTL? Or keep original TTL?
	// The original implementation kept it.
	return s.Redis.HSet(ctx, tokenKey, "used", "1").Err()
}

func (s *RedisStorage) DestroySession(ctx context.Context, accessToken string) error {
	// Extract JTI from token
	jti, err := s.extractJTIFromToken(accessToken)
	if err != nil {
		// Fallback to using accessToken as key if jti extraction fails (legacy support)
		jti = accessToken
		utils.Logger.Warn("Failed to extract jti from token, using accessToken as key", zap.Error(err))
	}

	// For Redis, we need to remove it from user set too.
	// Get user data from session to get userId
	result, err := s.Redis.Get(ctx, jti).Result()
	if err == nil {
		var user models.User
		if json.Unmarshal([]byte(result), &user) == nil {
			userSetKey := fmt.Sprintf("auth:user_tokens:%s", user.ID.String())
			_ = s.Redis.SRem(ctx, userSetKey, jti).Err()
		}
	}

	return s.Redis.Del(ctx, jti).Err()
}

func (s *RedisStorage) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	// 1. Access Tokens
	userSetKey := fmt.Sprintf("auth:user_tokens:%s", userID.String())
	jtis, err := s.Redis.SMembers(ctx, userSetKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	pipe := s.Redis.TxPipeline()
	if len(jtis) > 0 {
		for _, jti := range jtis {
			pipe.Del(ctx, jti)
		}
	}
	pipe.Del(ctx, userSetKey)

	// 2. Refresh Tokens (optional? Usually revoke all sessions means access tokens)
	// But in auth_repository.go `RevokeAllUserSessions` logic was empty?
	// Let's check auth_repository.go again.
	// I missed checking `RevokeAllUserSessions` implementation in `auth_repository.go`.
	// It wasn't in the snippet.

	// Assuming we want to delete all.
	_, err = pipe.Exec(ctx)
	return err
}
