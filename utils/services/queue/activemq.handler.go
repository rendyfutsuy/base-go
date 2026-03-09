package queue

import (
	"fmt"

	"github.com/go-stomp/stomp"
	"github.com/hibiken/asynq"
	"github.com/rendyfutsuy/base-go/utils"
	"go.uber.org/zap"
)

type ActiveMQHandler struct {
	conn *stomp.Conn
}

func NewActiveMQHandler() *ActiveMQHandler {
	return &ActiveMQHandler{}
}

func (h *ActiveMQHandler) Driver() string {
	return "activemq"
}

func (h *ActiveMQHandler) connect() error {
	if h.conn != nil {
		return nil
	}

	host := utils.ConfigVars.String("activemq.host")
	port := utils.ConfigVars.String("activemq.port")
	username := utils.ConfigVars.String("activemq.username")
	password := utils.ConfigVars.String("activemq.password")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		// default STOMP port
		port = "61613"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := stomp.Dial("tcp", addr, stomp.ConnOpt.Login(username, password))
	if err != nil {
		zap.S().Errorf("Failed to connect to ActiveMQ (STOMP %s): %v", addr, err)
		return err
	}
	h.conn = conn
	zap.S().Info("Successfully connected to ActiveMQ via STOMP")
	return nil
}

func (h *ActiveMQHandler) Close() error {
	if h.conn != nil {
		if err := h.conn.Disconnect(); err != nil {
			zap.S().Errorf("Failed to disconnect ActiveMQ STOMP connection: %v", err)
			return err
		}
		h.conn = nil
	}
	return nil
}

func (h *ActiveMQHandler) NewAsynqClient() (*asynq.Client, error) {
	return nil, fmt.Errorf("activemq handler does not support asynq client - use redis handler instead")
}

func (h *ActiveMQHandler) NewAsynqServer() (*asynq.Server, error) {
	return nil, fmt.Errorf("activemq handler does not support asynq server - use redis handler instead")
}

func (h *ActiveMQHandler) NewAsynqScheduler() (*asynq.Scheduler, error) {
	return nil, fmt.Errorf("activemq handler does not support asynq scheduler - use redis handler instead")
}

// Send publishes a message to ActiveMQ
func (h *ActiveMQHandler) Send(queueName string, payload []byte) error {
	return h.PublishMessage(queueName, payload)
}

// Run starts consuming provided workers from ActiveMQ and blocks
func (h *ActiveMQHandler) Run(workers map[string]func([]byte) error) error {
	for qname, handler := range workers {
		if err := h.ConsumeMessages(qname, handler); err != nil {
			return err
		}
	}
	select {}
}

func (h *ActiveMQHandler) PublishMessage(queueName string, body []byte) error {
	if err := h.connect(); err != nil {
		return err
	}
	dest := fmt.Sprintf("/queue/%s", queueName)
	if err := h.conn.Send(dest, "application/json", body, nil); err != nil {
		zap.S().Errorf("Failed to publish message to ActiveMQ queue %s: %v", queueName, err)
		return err
	}
	zap.S().Infof("Successfully published message to ActiveMQ queue: %s", queueName)
	return nil
}

func (h *ActiveMQHandler) ConsumeMessages(queueName string, handler func([]byte) error) error {
	if err := h.connect(); err != nil {
		return err
	}
	dest := fmt.Sprintf("/queue/%s", queueName)
	sub, err := h.conn.Subscribe(dest, stomp.AckAuto)
	if err != nil {
		zap.S().Errorf("Failed to subscribe to ActiveMQ queue %s: %v", queueName, err)
		return err
	}

	go func() {
		for {
			msg := sub.C
			if msg == nil {
				zap.S().Warnf("ActiveMQ subscription channel closed for queue: %s", queueName)
				return
			}
			m := <-msg
			if m.Err != nil {
				zap.S().Errorf("Error receiving ActiveMQ message: %v", m.Err)
				continue
			}
			if err := handler(m.Body); err != nil {
				zap.S().Errorf("Error processing ActiveMQ message: %v", err)
			}
		}
	}()

	zap.S().Infof("Successfully started consuming messages from ActiveMQ queue: %s", queueName)
	return nil
}
