package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/jwtGO/models"
)

type RefreshToken interface {
	Create(ctx context.Context, refresh *models.Refresh) error
	Get(ctx context.Context, tokenID string) (*models.Refresh, error)
	Update(ctx context.Context, oldTokenID string, newRefresh *models.Refresh) error
	Delete(ctx context.Context, tokenID string) error
}

type Repository struct {
	RefreshToken
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		RefreshToken: NewTokensRepo(db),
	}
}
