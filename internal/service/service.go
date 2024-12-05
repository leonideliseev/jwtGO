package service

import (
	"context"

	"github.com/leonideliseev/jwtGO/internal/repository"
)

type Tokens interface {
	CreateAccessToken(ctx context.Context, td *TokensData) (string, error)
	CreateRefreshToken(ctx context.Context, td *TokensData) (string, error)
	UpdateRefreshToken(ctx context.Context, oldTokenID string, td *TokensData) (string, error)
	ParseRefreshToken(ctx context.Context, userID, refreshToken string) (string, string, error)
}

type Service struct {
	Tokens
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Tokens: NewTokensService(repo.RefreshToken),
	}
}
