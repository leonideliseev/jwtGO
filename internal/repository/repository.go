package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshToken interface {
	Create(ctx context.Context, token, userID string) error
	Get(ctx context.Context, userID string) (string, error)
	Update(ctx context.Context, token, userID string) error
	// CheckUser(ctx context.Context, userID string) (bool, error)
}

type Repository struct {
	RefreshToken
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		RefreshToken: NewTokensRepo(db),
	}
}
