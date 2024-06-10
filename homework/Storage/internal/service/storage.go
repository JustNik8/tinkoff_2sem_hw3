package service

import (
	"context"

	"github.com/google/uuid"
	"hw3/storage/internal/domain"
	"hw3/storage/internal/repository"
	"hw3/storage/internal/repository/redis"
)

type StorageService interface {
	InsertMessage(ctx context.Context, message domain.MessageInfo) (domain.MessageInfo, error)
}

type storageService struct {
	repo  repository.StorageRepo
	cache *redis.StorageCache
}

func NewStorageService(repo repository.StorageRepo, cache *redis.StorageCache) StorageService {
	return &storageService{
		repo:  repo,
		cache: cache,
	}
}

func (s *storageService) InsertMessage(ctx context.Context, message domain.MessageInfo) (domain.MessageInfo, error) {
	message.ID = uuid.New().String()

	message, err := s.repo.InsertMessage(ctx, message)
	if err != nil {
		return domain.MessageInfo{}, err
	}

	err = s.cache.PushMessage(ctx, message)
	if err != nil {
		return domain.MessageInfo{}, err
	}

	return message, nil
}
