package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/leonideliseev/jwtGO/internal/repository"
	"github.com/leonideliseev/jwtGO/models"
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

type TokenAccessClaims struct {
	IP      string `json:"ip"`
	jwt.StandardClaims
}

type TokenRefreshClaims struct {
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
	claims := &TokenAccessClaims{
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

	hashedRefreshToken, err := hashRefreshToken(refreshToken) //TODO: проверить хэш-функцию
	if err != nil {
		return "", err
	}

	data := &models.Refresh{
		UserID: td.UserID,
		IP: td.IP,
		RefreshTokenHash: hashedRefreshToken,
	}

	_, err = s.repo.Get(ctx, td.UserID)
	var errRepo error
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errRepo = s.repo.Create(ctx, data)
	case err == nil:
		errRepo = s.repo.Update(ctx, data)
	default:
		return "", nil
	}
	if errRepo != nil {
		if errors.Is(errRepo, repository.ErrNotFound) {
			return "", ErrConcurency
		}

		return "", nil
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

	data := &models.Refresh{
		UserID: td.UserID,
		IP: td.IP,
		RefreshTokenHash: hashedRefreshToken,
	}

	err = s.repo.Update(ctx, data)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", HasNotToken
		}

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

// TODO: проверка должна быть на равенство с хэшем и изменился ли ip
func verifyRefreshToken(refreshToken string, storedHash *models.Refresh) error {
	return nil
	/*signedToken := fmt.Sprintf("%s.%s", refreshToken, refreshSecret)

	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(signedToken))*/
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
	claims := &TokenRefreshClaims{
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
