package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"hw3/chat-service/internal/config"
	"hw3/chat-service/internal/domain"
)

const (
	messagesKey = "messages"
)

type redisMessage struct {
	Nickname string    `json:"nickname"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
}

type StorageCache struct {
	client *redis.Client
}

func NewStorageCache(cfg config.RedisConfig) *StorageCache {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	opts := &redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	client := redis.NewClient(opts)

	return &StorageCache{
		client: client,
	}
}

func (s *StorageCache) GetLastMessages(ctx context.Context, count int64) ([]domain.MessageInfo, error) {
	cmd := s.client.LRange(ctx, messagesKey, 0, count-1)
	resp := make([]domain.MessageInfo, 0)

	messages, err := cmd.Result()
	log.Println(messages, err)
	if err != nil {
		return nil, err
	}

	log.Println("MESSAGES", messages)
	for i := range messages {

		var message redisMessage
		err := json.Unmarshal([]byte(messages[i]), &message)
		if err != nil {
			return nil, err
		}

		messageInfo := domain.MessageInfo{
			Nickname: message.Nickname,
			Message:  message.Message,
			Time:     message.Time,
		}
		resp = append(resp, messageInfo)
	}

	log.Println(resp)
	return resp, nil
}
