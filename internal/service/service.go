package service

import (
	"context"

	"github.com/leonideliseev/jwtGO/internal/repository"
)

type TokensData struct {
	UserID  string `json:"user_id"`
	TokenID string `json:"token_id"`
	IP      string `json:"ip"`
}

type AccessToken interface {
	Create(ctx context.Context, td *TokensData) (string, error)
}

type RefreshToken interface {
	Create(ctx context.Context, td *TokensData) (string, error)
	Update(ctx context.Context, oldTokenID string, td *TokensData) (string, error)
	Parse(ctx context.Context, userID, refreshToken string) (string, string, error)
}

type Service struct {
	RefreshToken
	AccessToken
}

func New(repo *repository.Repository) *Service {
	return &Service{
		RefreshToken: NewRefreshService(repo.RefreshToken),
		AccessToken: NewAccessService(),
		// Tokens: NewTokensService(repo.RefreshToken),
	}
}
