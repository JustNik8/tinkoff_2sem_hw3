package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"hw3/chat-service/internal/config"
	"hw3/chat-service/internal/domain"
)

const (
	messagesKey = "messages"
)

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
	resp := make([]domain.MessageInfo, count)

	messages := cmd.Args()
	for i, m := range messages {
		strMsg := messages[i].(string)

		var message domain.MessageInfo
		err := json.Unmarshal([]byte(strMsg), &message)
		if err != nil {
			return nil, err
		}

		resp[i] = m.(domain.MessageInfo)
	}

	return resp, nil
}
