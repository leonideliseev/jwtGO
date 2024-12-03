package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/leonideliseev/jwtGO/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	secretKey      = "your_secret_key"
	refreshSecret  = "your_refresh_secret"
	accessTTL      = time.Hour
)

type Claims struct {
	UserID  string `json:"user_id"`
	IP      string `json:"ip"`
	jwt.StandardClaims
}

type TokensService struct {
	repo repository.RefreshToken
}

func NewTokensService(repo repository.RefreshToken) *TokensService {
	return &TokensService{
		repo: repo,
	}
}

func (s *TokensService) GenerateAccessToken(ctx context.Context, userID, ip string) (string, error) {
	claims := &Claims{
		UserID: userID,
		IP:     ip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTTL).Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer:    "tokens_service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secretKey))
}

func (s *TokensService) GenerateRefreshToken(ctx context.Context, userID string) (string, error) {
	token, hashedToken, err := generateRefreshToken()
	if err != nil {
		return "", err
	}

	err = s.repo.Create(ctx, string(hashedToken), userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *TokensService) UpdateRefreshToken(ctx context.Context, userID string) (string, error) {
	token, hashedToken, err := generateRefreshToken()
	if err != nil {
		return "", err
	}

	err = s.repo.Update(ctx, string(hashedToken), userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *TokensService) CheckRefreshToken(ctx context.Context, userID, refreshToken string) error {
	storedHash, err := s.repo.Get(ctx, userID)
	if err != nil {
		return err
	}

	err = verifyRefreshToken(refreshToken, storedHash)
	if err != nil {
		return fmt.Errorf("Incorrect refresh token")
	}

	return nil
}

func verifyRefreshToken(refreshToken string, storedHash string) error {
	signedToken := fmt.Sprintf("%s.%s", refreshToken, refreshSecret)

	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(signedToken))
}

func generateRefreshToken() (string, string, error) {
	token := make([]byte, 64)
	_, err := rand.Read(token)
	if err != nil {
		return "", "", err
	}

	signedToken := fmt.Sprintf("%s.%s", base64.URLEncoding.EncodeToString(token), refreshSecret)

	if len(signedToken) > 72 {
		signedToken = signedToken[:72]
	}

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(signedToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return base64.URLEncoding.EncodeToString(token), string(hashedToken), nil
}
