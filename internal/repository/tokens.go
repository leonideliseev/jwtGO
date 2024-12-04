package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokensRepo struct {
	db *pgxpool.Pool
}

func NewTokensRepo(db *pgxpool.Pool) *TokensRepo {
	return &TokensRepo{
		db: db,
	}
}

func (r *TokensRepo) Create(ctx context.Context, hashedRefreshToken, userID string) error {
	_, err := r.db.Exec(ctx, "INSERT INTO refresh_tokens(user_id, token_hash) VALUES($1, $2)", userID, hashedRefreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (r *TokensRepo) Update(ctx context.Context, hashedRefreshToken, userID string) error {
	_, err := r.db.Exec(ctx, "UPDATE refresh_tokens SET token_hash=$1 WHERE user_id=$2", hashedRefreshToken, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *TokensRepo) Get(ctx context.Context, userID string) (string, error) {
	var hash string
	err := r.db.QueryRow(ctx, "SELECT token_hash FROM refresh_tokens WHERE user_id=$1", userID).Scan(&hash)
	if err != nil {
		return "", nil
	}

	return hash, err
}

func (r *TokensRepo) CheckUser(ctx context.Context, userID string) (bool, error) {
	return false, nil
}
