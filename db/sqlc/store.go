package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg CreateTransferParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParam) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParam) (VerifyEmailTxResult, error)
}

type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{connPool: pool, Queries: New(pool)}
}
