package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
)

type AccessService struct {
}

func NewAccessService() *AccessService {
	return &AccessService{}
}

const (
	accessSecret  = "your_secret_key"
)

type TokenAccessClaims struct {
	IP string `json:"ip"`
	jwt.StandardClaims
}

func (s *AccessService) Create(ctx context.Context, td *TokensData) (string, error) {
	claims := &TokenAccessClaims{
		IP: td.IP,
		StandardClaims: jwt.StandardClaims{
			Subject:   td.UserID,
			ExpiresAt: time.Now().Add(accessTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "tokens_service",
			Id:        td.TokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(accessSecret))
}
