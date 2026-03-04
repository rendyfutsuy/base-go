package queue

import (
	"errors"
	"github.com/hibiken/asynq"
)

type RabbitMQHandler struct{}

func NewRabbitMQHandler() *RabbitMQHandler {
	return &RabbitMQHandler{}
}

func (h *RabbitMQHandler) Driver() string {
	return "rabbitmq"
}

func (h *RabbitMQHandler) NewAsynqClient() (*asynq.Client, error) {
	return nil, errors.New("rabbitmq not supported for asynq client")
}

func (h *RabbitMQHandler) NewAsynqServer() (*asynq.Server, error) {
	return nil, errors.New("rabbitmq not supported for asynq server")
}

func (h *RabbitMQHandler) NewAsynqScheduler() (*asynq.Scheduler, error) {
	return nil, errors.New("rabbitmq not supported for asynq scheduler")
}