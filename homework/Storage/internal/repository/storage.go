package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"hw3/storage/internal/domain"
	"hw3/storage/internal/repository/queries"
)

type StorageRepo interface {
	InsertMessage(ctx context.Context, messageInfo domain.MessageInfo) (domain.MessageInfo, error)
}

type storageRepo struct {
	*queries.Queries
	pool *pgxpool.Pool
}

func NewStorageRepo(pgxPool *pgxpool.Pool) StorageRepo {
	return &storageRepo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
	}
}
