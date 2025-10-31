// file: modules/auth/repository/redis_repository.go
package repository

import (
	"context"
	"encoding/json"
	"time"

	models "github.com/rendyfutsuy/base-go/models"
	// Updated import path
)

func (r *authRepository) CreateSession(ctx context.Context, jti string, user *models.User, ttl time.Duration) error {
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.Redis.Set(ctx, jti, userData, ttl).Err()
}

func (r *authRepository) GetSessionData(ctx context.Context, jti string) (*models.User, error) {
	result, err := r.Redis.Get(ctx, jti).Result()
	if err != nil {
		// This includes redis.Nil, meaning the session is not on the allowlist.
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal([]byte(result), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) DeleteSession(ctx context.Context, jti string) error {
	return r.Redis.Del(ctx, jti).Err()
}
