package queue

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type RabbitMQHandler struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQHandler() *RabbitMQHandler {
	return &RabbitMQHandler{}
}

func (h *RabbitMQHandler) Driver() string {
	return "rabbitmq"
}

func (h *RabbitMQHandler) connect() error {
	if h.conn != nil && !h.conn.IsClosed() {
		return nil
	}

	host := utils.ConfigVars.String("rabbitmq.host")
	port := utils.ConfigVars.String("rabbitmq.port")
	username := utils.ConfigVars.String("rabbitmq.username")
	password := utils.ConfigVars.String("rabbitmq.password")

	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		zap.S().Errorf("Failed to connect to RabbitMQ: %v", err)
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		zap.S().Errorf("Failed to open RabbitMQ channel: %v", err)
		return err
	}

	h.conn = conn
	h.ch = ch

	zap.S().Info("Successfully connected to RabbitMQ")
	return nil
}

func (h *RabbitMQHandler) close() {
	if h.ch != nil {
		h.ch.Close()
	}
	if h.conn != nil {
		h.conn.Close()
	}
}

func (h *RabbitMQHandler) NewAsynqClient() (*asynq.Client, error) {
	return nil, fmt.Errorf("rabbitmq handler does not support asynq client - use redis handler instead")
}

func (h *RabbitMQHandler) NewAsynqServer() (*asynq.Server, error) {
	return nil, fmt.Errorf("rabbitmq handler does not support asynq server - use redis handler instead")
}

func (h *RabbitMQHandler) NewAsynqScheduler() (*asynq.Scheduler, error) {
	return nil, fmt.Errorf("rabbitmq handler does not support asynq scheduler - use redis handler instead")
}

// Send publishes a message to RabbitMQ
func (h *RabbitMQHandler) Send(queueName string, payload []byte) error {
	return h.PublishMessage(queueName, payload)
}

func (h *RabbitMQHandler) DeclareQueue(queueName string) (amqp.Queue, error) {
	if err := h.connect(); err != nil {
		return amqp.Queue{}, err
	}

	q, err := h.ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		zap.S().Errorf("Failed to declare RabbitMQ queue %s: %v", queueName, err)
		return q, err
	}

	zap.S().Infof("Successfully declared RabbitMQ queue: %s", queueName)
	return q, nil
}

func (h *RabbitMQHandler) PublishMessage(queueName string, body []byte) error {
	if err := h.connect(); err != nil {
		return err
	}

	q, err := h.DeclareQueue(queueName)
	if err != nil {
		return err
	}

	err = h.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		zap.S().Errorf("Failed to publish message to RabbitMQ queue %s: %v", queueName, err)
		return err
	}

	zap.S().Infof("Successfully published message to RabbitMQ queue: %s", queueName)
	return nil
}

func (h *RabbitMQHandler) ConsumeMessages(queueName string, handler func([]byte) error) error {
	if err := h.connect(); err != nil {
		return err
	}

	q, err := h.DeclareQueue(queueName)
	if err != nil {
		return err
	}

	msgs, err := h.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		zap.S().Errorf("Failed to register RabbitMQ consumer for queue %s: %v", queueName, err)
		return err
	}

	go func() {
		for d := range msgs {
			zap.S().Infof("Received message from RabbitMQ queue %s", queueName)
			if err := handler(d.Body); err != nil {
				zap.S().Errorf("Error processing RabbitMQ message: %v", err)
			}
		}
	}()

	zap.S().Infof("Successfully started consuming messages from RabbitMQ queue: %s", queueName)
	return nil
}
