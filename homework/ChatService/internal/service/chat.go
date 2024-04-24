package service

import (
	"context"

	"github.com/google/uuid"
	"hw3/chat-service/internal/domain"
	"hw3/chat-service/internal/repo"
)

type ChatService interface {
	GetLastMessages(ctx context.Context, count int) ([]domain.MessageInfo, error)
	InsertMessage(ctx context.Context, messageInfo domain.MessageInfo) (domain.MessageInfo, error)
}

type chatService struct {
	repo repo.ChatRepo
}

func NewChatService(repo repo.ChatRepo) ChatService {
	return &chatService{repo: repo}
}

func (s *chatService) GetLastMessages(ctx context.Context, count int) ([]domain.MessageInfo, error) {
	return s.repo.GetLastMessages(ctx, count)
}

func (s *chatService) InsertMessage(ctx context.Context, messageInfo domain.MessageInfo) (domain.MessageInfo, error) {
	id := uuid.New()
	messageInfo.ID = id.String()

	return s.repo.InsertMessage(ctx, messageInfo)
}
