package queue

import (
	"github.com/hibiken/asynq"
)

// QueueService defines the interface that any queue provider (redis, rabbitmq) must implement.
type QueueService interface {
	Driver() string
	NewAsynqClient() (*asynq.Client, error)
	NewAsynqServer() (*asynq.Server, error)
	NewAsynqScheduler() (*asynq.Scheduler, error)
}
