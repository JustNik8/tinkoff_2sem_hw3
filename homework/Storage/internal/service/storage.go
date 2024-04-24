package service

import (
	"context"

	"github.com/google/uuid"
	"hw3/storage/internal/domain"
	"hw3/storage/internal/repository"
)

type StorageService interface {
	InsertMessage(ctx context.Context, message domain.MessageInfo) (domain.MessageInfo, error)
}

type storageService struct {
	repo repository.StorageRepo
}

func NewStorageService(repo repository.StorageRepo) StorageService {
	return &storageService{repo: repo}
}

func (s *storageService) InsertMessage(ctx context.Context, message domain.MessageInfo) (domain.MessageInfo, error) {
	message.ID = uuid.New().String()

	return s.repo.InsertMessage(ctx, message)
}
