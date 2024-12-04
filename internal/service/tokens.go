package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/leonideliseev/jwtGO/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	secretKey     = "your_secret_key"
	refreshSecret = "your_refresh_secret"
	accessTTL     = time.Hour
	refreshTTL    = 7 * 24 * time.Hour
)

type TokensData struct {
	UserID  string `json:"user_id"`
	TokenID string `json:"token_id"`
	IP      string `json:"ip"`
}

type TokensClaims struct {
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

func (s *TokensService) GenerateAccessToken(ctx context.Context, td *TokensData) (string, error) {
	claims := &TokensClaims{
		IP:     td.IP,
		StandardClaims: jwt.StandardClaims{
			Subject: td.UserID,
			ExpiresAt: time.Now().Add(accessTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "tokens_service",
			Id: td.TokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secretKey))
}

func (s *TokensService) GenerateRefreshToken(ctx context.Context, td *TokensData) (string, error) {
	var err error
	refreshToken, err := generateRefreshToken(td)
	if err != nil {
		return "", nil
	}

	_, err = hashRefreshToken(refreshToken) //TODO: hashedRefreshToken
	if err != nil {
		return "", err
	}

	ok, err := s.repo.Get(ctx, td.UserID) // TODO: возможно стоит добавить ошибку NotFound и использовать проверку через неё
	if err != nil {
		return "", err
	}

	// TODO: заменить на хэш
	if ok == "" {
		err = s.repo.Update(ctx, refreshToken, td.UserID)
	} else {
		err = s.repo.Create(ctx, refreshToken, td.UserID)
	}
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (s *TokensService) UpdateRefreshToken(ctx context.Context, td *TokensData) (string, error) {
	token, err := generateRefreshToken(td)
	if err != nil {
		return "", err
	}

	hashedRefreshToken, err := hashRefreshToken(token)
	if err != nil {
		return "", err
	}

	err = s.repo.Update(ctx, string(hashedRefreshToken), td.UserID)
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

	err = verifyRefreshToken(refreshToken, storedHash) // TODO: изменить проверку
	if err != nil {
		return fmt.Errorf("Incorrect refresh token")
	}

	return nil
}

func verifyRefreshToken(refreshToken string, storedHash string) error {
	signedToken := fmt.Sprintf("%s.%s", refreshToken, refreshSecret)

	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(signedToken))
}

func hashRefreshToken(token string) (string, error) {
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
	claims := &TokensClaims{
		IP:      td.IP,
		StandardClaims: jwt.StandardClaims{
			Subject: td.UserID,
			ExpiresAt: time.Now().Add(refreshTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "tokens_service",
			Id: td.TokenID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	refreshToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}
