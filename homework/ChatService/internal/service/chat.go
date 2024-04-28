package service

import (
	"context"

	"hw3/chat-service/internal/domain"
	"hw3/chat-service/internal/repo/redis"
)

type ChatService interface {
	GetLastMessages(ctx context.Context, count int64) ([]domain.MessageInfo, error)
}

type chatService struct {
	storageCache *redis.StorageCache
}

func NewChatService(storageCache *redis.StorageCache) ChatService {
	return &chatService{
		storageCache: storageCache,
	}
}

func (s *chatService) GetLastMessages(ctx context.Context, count int64) ([]domain.MessageInfo, error) {
	return s.storageCache.GetLastMessages(ctx, count)
}
