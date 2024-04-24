package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"hw3/chat-service/internal/transport/dto"
)

type ChatHandler struct {
	producer sarama.SyncProducer
}

func NewChatHandler(addrs []string) (*ChatHandler, error) {
	producer, err := sarama.NewSyncProducer(addrs, nil)
	if err != nil {
		return nil, err
	}

	return &ChatHandler{
		producer: producer,
	}, nil
}

func (h *ChatHandler) ProduceMessage(topic string, message dto.MessageInfoRequest) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	resp := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(message.Nickname),
		Value: sarama.ByteEncoder(messageBytes),
	}

	_, _, err = h.producer.SendMessage(resp)
	return err
}

func (h *ChatHandler) Close() {
	err := h.producer.Close()
	if err != nil {
		log.Printf("Error while closing kafka producer: %v", err)
	}
}
