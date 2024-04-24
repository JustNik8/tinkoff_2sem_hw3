package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"hw3/chat-service/internal/domain"
	"hw3/chat-service/internal/repo/queries"
)

type ChatRepo interface {
	GetLastMessages(ctx context.Context, count int) ([]domain.MessageInfo, error)
	InsertMessage(ctx context.Context, messageInfo domain.MessageInfo) (domain.MessageInfo, error)
}

type chatRepo struct {
	*queries.Queries
	pool *pgxpool.Pool
}

func NewRepo(pgxPool *pgxpool.Pool) ChatRepo {
	return &chatRepo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
	}
}
