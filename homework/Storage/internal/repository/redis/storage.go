package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"hw3/storage/internal/config"
	"hw3/storage/internal/domain"
)

const (
	messagesKey = "messages"
	startIdx    = 0
	endIdx      = 9
)

type redisMessage struct {
	ID       string    `json:"id"`
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

func (s *StorageCache) PushMessage(ctx context.Context, message domain.MessageInfo) error {
	cacheMessage := redisMessage{
		ID:       message.ID,
		Nickname: message.Nickname,
		Message:  message.Message,
		Time:     message.Time,
	}

	messageBytes, err := json.Marshal(cacheMessage)
	if err != nil {
		return err
	}

	s.client.LPush(ctx, messagesKey, messageBytes)
	s.client.LTrim(ctx, messagesKey, startIdx, endIdx)

	return nil
}
