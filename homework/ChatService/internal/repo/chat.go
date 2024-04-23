package repo

import (
	"context"

	"2sem/hw1/homework/internal/domain"
	"2sem/hw1/homework/internal/repo/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepo interface {
	GetLastMessages(ctx context.Context, count int) ([]*domain.MessageInfo, error)
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
