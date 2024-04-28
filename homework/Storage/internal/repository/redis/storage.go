package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"hw3/storage/internal/config"
	"hw3/storage/internal/transport/dto"
)

const (
	messagesKey = "messages"
	startIdx    = 0
	endIdx      = 9
	messagesCnt = 10
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

func (s *StorageCache) PushMessage(ctx context.Context, message dto.MessageInfoResponse) {
	s.client.LPush(ctx, messagesKey, message)
	s.client.LTrim(ctx, messagesKey, startIdx, endIdx)
}
