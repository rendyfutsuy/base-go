// file: database/redis.go
package database

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rendyfutsuy/base-go/utils"
	"go.uber.org/zap"
)

// ConnectToRedis initializes and returns a Redis client
func ConnectToRedis() *redis.Client {
	address := utils.ConfigVars.String("redis.address")
	if address == "" {
		address = utils.ConfigVars.String("redis.addr")
	}
	password := utils.ConfigVars.String("redis.password")
	db := utils.ConfigVars.Int("redis.db")

	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	// Ping the server to check the connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		utils.Logger.Error("Could not connect to Redis", zap.Error(err))
		return nil
	}

	utils.Logger.Info("Successfully connected to Redis!")
	return rdb
}
