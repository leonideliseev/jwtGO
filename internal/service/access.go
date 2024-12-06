package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/leonideliseev/jwtGO/config"
)

type AccessService struct {
	accessSecret string
	accessTTL    time.Duration
}

func NewAccessService(cfg config.JWT) *AccessService {
	return &AccessService{
		accessSecret: cfg.AccessSignKey,
		accessTTL:    cfg.AccessTokenTTL,
	}
}

type TokenAccessClaims struct {
	IP string `json:"ip"`
	jwt.StandardClaims
}

func (s *AccessService) Create(ctx context.Context, td *TokensData) (string, error) {
	claims := &TokenAccessClaims{
		IP: td.IP,
		StandardClaims: jwt.StandardClaims{
			Subject:   td.UserID,
			ExpiresAt: time.Now().Add(s.accessTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "tokens_service",
			Id:        td.TokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(s.accessSecret))
}
