package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
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

var (
	refreshSecret = os.Getenv("REFRESH_SECRET")
)

const (
	refreshTTL    = 7 * 24 * time.Hour
)

type TokenRefreshClaims struct {
	jwt.StandardClaims
}

func (s *RefreshService) Create(ctx context.Context, td *TokensData) (string, error) {
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

	err = s.repo.Create(ctx, data)
	if err != nil {
		return "", err
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
			return "", ErrHasNotToken
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

	stored := storedToken.RefreshTokenHash
	if err := compareHash(stored, refreshToken); err != nil {
		return "", "", err
	}

	return storedToken.TokenID, storedToken.IP, nil
}

func compareHash(hashedToken string, inputToken string) error {
    sha256Hash := sha256.Sum256([]byte(inputToken))
    sha256Hex := hex.EncodeToString(sha256Hash[:])

    err := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(sha256Hex))
    if err != nil {
        return err
    }

    return nil
}

func hash(token string) (string, error) {
	// get fixed len <= 72, for use bcrypt
	sha256Hash := sha256.Sum256([]byte(token))
	sha256Hex := hex.EncodeToString(sha256Hash[:])

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(sha256Hex), bcrypt.DefaultCost)
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
	return token.SignedString([]byte(refreshSecret))
}
