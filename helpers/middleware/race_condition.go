package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/utils"
	"go.uber.org/zap"
)

type IRaceConditionMiddleware interface {
	PreventRaceCondition(keyPrefix string) echo.MiddlewareFunc
}

type RaceConditionMiddleware struct {
	redisClient *redis.Client
	lockTTL     time.Duration
}

func NewRaceConditionMiddleware(redisClient *redis.Client) IRaceConditionMiddleware {
	return &RaceConditionMiddleware{
		redisClient: redisClient,
		lockTTL:     30 * time.Second, // Default lock TTL
	}
}

// PreventRaceCondition creates a middleware that prevents race conditions
// by implementing distributed locking using Redis
func (rc *RaceConditionMiddleware) PreventRaceCondition(keyPrefix string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Generate unique lock key based on request context
			lockKey := rc.generateLockKey(c, keyPrefix)
			lockID := uuid.New().String()

			// Try to acquire lock
			acquired, err := rc.acquireLock(c.Request().Context(), lockKey, lockID)
			if err != nil {
				utils.Logger.Error("Failed to acquire lock for race condition prevention",
					zap.String("lock_key", lockKey),
					zap.Error(err),
				)
				return c.JSON(http.StatusInternalServerError, response.SetErrorResponse(
					http.StatusInternalServerError,
					"Failed to process request due to system error",
				))
			}

			if !acquired {
				return c.JSON(http.StatusTooManyRequests, response.SetErrorResponse(
					http.StatusTooManyRequests,
					"Request is being processed, please try again",
				))
			}

			// Ensure lock is released after processing
			defer rc.releaseLock(c.Request().Context(), lockKey, lockID)

			// Process request
			return next(c)
		}
	}
}

// generateLockKey creates a unique lock key based on request parameters
func (rc *RaceConditionMiddleware) generateLockKey(c echo.Context, keyPrefix string) string {
	var keyBuilder strings.Builder
	keyBuilder.WriteString("race_condition:")
	keyBuilder.WriteString(keyPrefix)
	keyBuilder.WriteString(":")

	// Include user ID if available
	if userID := c.Get("userId"); userID != nil {
		keyBuilder.WriteString(fmt.Sprintf("user:%v:", userID))
	}

	// Include relevant request parameters
	keyBuilder.WriteString(fmt.Sprintf("method:%s:", c.Request().Method))
	keyBuilder.WriteString(fmt.Sprintf("path:%s:", c.Path()))

	// Include specific parameters that could cause race conditions
	// For example, resource IDs, action types, etc.
	if id := c.Param("id"); id != "" {
		keyBuilder.WriteString(fmt.Sprintf("id:%s:", id))
	}

	// Create hash of the key to ensure consistent length
	hash := sha256.Sum256([]byte(keyBuilder.String()))
	return hex.EncodeToString(hash[:])
}

// acquireLock attempts to acquire a distributed lock
func (rc *RaceConditionMiddleware) acquireLock(ctx context.Context, lockKey, lockID string) (bool, error) {
	if rc.redisClient == nil {
		// If Redis is not available, skip race condition prevention
		return true, nil
	}

	// Use SET with NX (only set if not exists) and EX (expire time)
	result := rc.redisClient.SetNX(ctx, lockKey, lockID, rc.lockTTL)
	return result.Result()
}

// releaseLock releases the distributed lock
func (rc *RaceConditionMiddleware) releaseLock(ctx context.Context, lockKey, lockID string) error {
	if rc.redisClient == nil {
		return nil
	}

	// Use Lua script to ensure only the lock owner can release it
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result := rc.redisClient.Eval(ctx, script, []string{lockKey}, lockID)
	_, err := result.Result()
	return err
}
