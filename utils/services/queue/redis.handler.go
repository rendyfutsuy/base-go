package queue

import (
	"github.com/hibiken/asynq"
	"github.com/rendyfutsuy/base-go/utils"
)

type RedisHandler struct{}

func NewRedisHandler() *RedisHandler {
	return &RedisHandler{}
}

func (h *RedisHandler) Driver() string {
	return "redis"
}

func (h *RedisHandler) redisOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     utils.ConfigVars.String("redis.address"),
		Password: utils.ConfigVars.String("redis.password"),
		DB:       utils.ConfigVars.Int("redis.db"),
	}
}

func (h *RedisHandler) serverConfig() asynq.Config {
	return asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	}
}

func (h *RedisHandler) NewAsynqClient() (*asynq.Client, error) {
	return asynq.NewClient(h.redisOpt()), nil
}

func (h *RedisHandler) NewAsynqServer() (*asynq.Server, error) {
	return asynq.NewServer(h.redisOpt(), h.serverConfig()), nil
}

func (h *RedisHandler) NewAsynqScheduler() (*asynq.Scheduler, error) {
	return asynq.NewScheduler(h.redisOpt(), nil), nil
}