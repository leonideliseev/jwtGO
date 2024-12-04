package service

import (
	"context"

	"github.com/leonideliseev/jwtGO/internal/repository"
)

type Tokens interface {
	GenerateAccessToken(ctx context.Context, td *TokensData) (string, error)
	GenerateRefreshToken(ctx context.Context, td *TokensData) (string, error)
	UpdateRefreshToken(ctx context.Context, td *TokensData) (string, error)
	CheckRefreshToken(ctx context.Context, userID, refreshToken string) error
}

type Service struct {
	Tokens
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Tokens: NewTokensService(repo.RefreshToken),
	}
}
