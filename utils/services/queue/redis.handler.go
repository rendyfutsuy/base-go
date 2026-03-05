package queue

import (
	"context"

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

// Send enqueues a task to Redis using Asynq
func (h *RedisHandler) Send(queueName string, payload []byte) error {
	client, err := h.NewAsynqClient()
	if err != nil {
		return err
	}
	defer client.Close()
	_, err = client.Enqueue(asynq.NewTask(queueName, payload), asynq.MaxRetry(5))
	return err
}

// Run starts Asynq server and scheduler, wiring provided workers into ServeMux
func (h *RedisHandler) Run(workers map[string]func([]byte) error) error {
	srv, err := h.NewAsynqServer()
	if err != nil {
		return err
	}
	mux := asynq.NewServeMux()
	for qname, handler := range workers {
		h := handler
		mux.HandleFunc(qname, func(ctx context.Context, t *asynq.Task) error {
			return h(t.Payload())
		})
	}
	scheduler, err := h.NewAsynqScheduler()
	if err != nil {
		return err
	}
	go func() {
		_ = srv.Run(mux)
	}()
	return scheduler.Run()
}
