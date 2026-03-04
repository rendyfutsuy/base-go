package queue

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"time"
)

type KafkaHandler struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaHandler() *KafkaHandler {
	return &KafkaHandler{}
}

func (h *KafkaHandler) Driver() string {
	return "kafka"
}

func (h *KafkaHandler) getKafkaBrokers() []string {
	broker := utils.ConfigVars.String("kafka.broker")
	if broker == "" {
		broker = "localhost:9092"
	}
	return []string{broker}
}

func (h *KafkaHandler) NewAsynqClient() (*asynq.Client, error) {
	return nil, fmt.Errorf("kafka handler does not support asynq client - use redis handler instead")
}

func (h *KafkaHandler) NewAsynqServer() (*asynq.Server, error) {
	return nil, fmt.Errorf("kafka handler does not support asynq server - use redis handler instead")
}

func (h *KafkaHandler) NewAsynqScheduler() (*asynq.Scheduler, error) {
	return nil, fmt.Errorf("kafka handler does not support asynq scheduler - use redis handler instead")
}

func (h *KafkaHandler) CreateTopic(topic string) error {
	conn, err := kafka.Dial("tcp", h.getKafkaBrokers()[0])
	if err != nil {
		zap.S().Errorf("Failed to connect to Kafka broker: %v", err)
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		zap.S().Errorf("Failed to get Kafka controller: %v", err)
		return err
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		zap.S().Errorf("Failed to connect to Kafka controller: %v", err)
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		zap.S().Errorf("Failed to create Kafka topic %s: %v", topic, err)
		return err
	}

	zap.S().Infof("Successfully created Kafka topic: %s", topic)
	return nil
}

func (h *KafkaHandler) PublishMessage(topic string, key string, value []byte) error {
	if h.writer == nil {
		h.writer = &kafka.Writer{
			Addr:     kafka.TCP(h.getKafkaBrokers()...),
			Topic:    topic,
			Balancer: &kafka.Hash{},
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	})
	if err != nil {
		zap.S().Errorf("Failed to publish message to Kafka topic %s: %v", topic, err)
		return err
	}

	zap.S().Infof("Successfully published message to Kafka topic: %s", topic)
	return nil
}

func (h *KafkaHandler) ConsumeMessages(topic string, groupID string, handler func([]byte) error) error {
	if h.reader == nil {
		h.reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers: h.getKafkaBrokers(),
			Topic:   topic,
			GroupID: groupID,
			MinBytes: 1,
			MaxBytes: 10e6, // 10MB
		})
	}

	go func() {
		for {
			ctx := context.Background()
			msg, err := h.reader.ReadMessage(ctx)
			if err != nil {
				zap.S().Errorf("Failed to read message from Kafka topic %s: %v", topic, err)
				continue
			}

			zap.S().Infof("Received message from Kafka topic %s", topic)
			if err := handler(msg.Value); err != nil {
				zap.S().Errorf("Error processing Kafka message: %v", err)
			}
		}
	}()

	zap.S().Infof("Successfully started consuming messages from Kafka topic: %s", topic)
	return nil
}

func (h *KafkaHandler) Close() error {
	if h.writer != nil {
		if err := h.writer.Close(); err != nil {
			zap.S().Errorf("Failed to close Kafka writer: %v", err)
			return err
		}
	}
	if h.reader != nil {
		if err := h.reader.Close(); err != nil {
			zap.S().Errorf("Failed to close Kafka reader: %v", err)
			return err
		}
	}
	return nil
}