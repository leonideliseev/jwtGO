package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/leonideliseev/jwtGO/internal/repository"
	"github.com/leonideliseev/jwtGO/models"
	"golang.org/x/crypto/bcrypt"
)

type RefreshService struct {
	repo repository.RefreshToken
}

func NewRefreshService(repo repository.RefreshToken) *RefreshService {
	return &RefreshService{
		repo: repo,
	}
}

const (
	refreshSecret = "your_refresh_secret"
	accessTTL     = time.Hour
	refreshTTL    = 7 * 24 * time.Hour
)

type TokenRefreshClaims struct {
	jwt.StandardClaims
}

func (s *RefreshService) Create(ctx context.Context, td *TokensData) (string, error) {
	refreshToken, err := generateRefreshToken(td)
	if err != nil {
		return "", nil
	}

	hashedRefreshToken, err := hash(refreshToken)
	if err != nil {
		return "", err
	}

	data := &models.Refresh{
		TokenID:          td.TokenID,
		IP:               td.IP,
		RefreshTokenHash: hashedRefreshToken,
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		return "", nil
	}

	return refreshToken, nil
}

func (s *RefreshService) Update(ctx context.Context, oldTokenID string, td *TokensData) (string, error) {
	refreshToken, err := generateRefreshToken(td)
	if err != nil {
		return "", err
	}

	hashedRefreshToken, err := hash(refreshToken)
	if err != nil {
		return "", err
	}

	data := &models.Refresh{
		TokenID:          td.TokenID,
		IP:               td.IP,
		RefreshTokenHash: hashedRefreshToken,
	}

	err = s.repo.Update(ctx, oldTokenID, data)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrConcurency
		}

		return "", err
	}

	return refreshToken, nil
}

func (s *RefreshService) Parse(ctx context.Context, userID, refreshToken string) (string, string, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &TokenRefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(refreshSecret), nil
	})
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*TokenRefreshClaims)
	if !ok || !token.Valid {
		return "", "", errors.New("invalid token")
	}

	if claims.Subject != userID {
		return "", "", errors.New("user ID does not match")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return "", "", errors.New("refresh token expired")
	}

	storedToken, err := s.repo.Get(ctx, claims.Id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", ErrHasNotToken
		}

		return "", "", ErrInternal
	}

	requestTokenHask, err := hash(refreshToken)
	if err != nil {
		return "", "", ErrInternal
	}

	if storedToken.RefreshTokenHash != requestTokenHask {
		return "", "", errors.New("wrong token")
	}

	return storedToken.TokenID, storedToken.IP, nil
}

func hash(token string) (string, error) {
	if len(token) > 72 {
		token = token[:72]
	}

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedToken), nil
}

func generateRefreshToken(td *TokensData) (string, error) {
	claims := &TokenRefreshClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   td.UserID,
			ExpiresAt: time.Now().Add(refreshTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "tokens_service",
			Id:        td.TokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	refreshToken, err := token.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
