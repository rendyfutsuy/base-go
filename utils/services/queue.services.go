package services

import (
	"sync"

	"github.com/rendyfutsuy/base-go/utils"
	"github.com/rendyfutsuy/base-go/utils/services/queue"
	"go.uber.org/zap"
)

const (
	RedisDriver    = "redis"
	RabbitMQDriver = "rabbitmq"
)

var (
	queueOnce sync.Once
	queueInst queue.QueueService
	queueErr  error
)

// GetQueue returns a queue service instance based on the driver name.
// It supports redis and rabbitmq drivers.
func GetQueue(driver string) (queue.QueueService, error) {
	switch driver {
	case RedisDriver:
		return queue.NewRedisHandler(), nil
	case RabbitMQDriver:
		return queue.NewRabbitMQHandler(), nil
	default:
		zap.S().Errorf("unsupported queue driver: %s", driver)
		return nil, nil
	}
}

// NewQueueService creates a new queue service instance based on configuration.
// It defaults to redis if no driver is specified.
func NewQueueService() queue.QueueService {
	driver := ""
	if utils.ConfigVars.Exists("queue.driver") {
		driver = utils.ConfigVars.String("queue.driver")
	}
	if driver == "" {
		driver = RedisDriver
	}

	q, err := GetQueue(driver)
	if err != nil {
		zap.S().Errorf("failed to create queue service: %v", err)
		return queue.NewRedisHandler() // fallback to redis
	}
	return q
}
