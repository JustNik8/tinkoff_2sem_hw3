package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"hw3/storage/internal/converter"
	"hw3/storage/internal/service"
	"hw3/storage/internal/transport/dto"
)

type StorageKafka struct {
	consumer         sarama.Consumer
	storageService   service.StorageService
	messageConverter converter.MessageConverter
}

func NewStorageKafka(addrs []string, storageService service.StorageService, messageConverter converter.MessageConverter) (*StorageKafka, error) {
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	return &StorageKafka{
		consumer:         consumer,
		storageService:   storageService,
		messageConverter: messageConverter,
	}, nil
}

func (s *StorageKafka) ConsumeMessages(topic string) error {
	partConsumer, err := s.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	for {
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				return fmt.Errorf("channel closed")
			}

			// Десериализация входящего сообщения из JSON
			var messageDTO dto.MessageInfoDTO
			err := json.Unmarshal(msg.Value, &messageDTO)

			if err != nil {
				log.Printf("Error unmarshaling JSON: %v\n", err)
				continue
			}

			message := s.messageConverter.MapDtoToDomain(messageDTO)
			log.Printf("Received message: %+v\n", messageDTO)

			ctx := context.Background()
			savedMessage, err := s.storageService.InsertMessage(ctx, message)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(savedMessage)
		}
	}
}
