package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/jwtGO/models"
)

type TokensRepo struct {
	db *pgxpool.Pool
}

func NewTokensRepo(db *pgxpool.Pool) *TokensRepo {
	return &TokensRepo{
		db: db,
	}
}

func (r *TokensRepo) Create(ctx context.Context, refresh *models.Refresh) error {
	_, err := r.db.Exec(ctx, "INSERT INTO refresh_tokens(token_id, ip, token_hash) VALUES($1, $2, $3)", 
		refresh.TokenID, refresh.IP, refresh.RefreshTokenHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *TokensRepo) Update(ctx context.Context, refresh *models.Refresh) error {
	commandTag, err := r.db.Exec(ctx, "UPDATE refresh_tokens SET token_hash=$1, ip=$2 WHERE token_id=$3",
		refresh.RefreshTokenHash, refresh.IP, refresh.TokenID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *TokensRepo) Get(ctx context.Context, tokenID string) (*models.Refresh, error) {
	var refresh models.Refresh
	err := r.db.QueryRow(ctx, "SELECT token_hash FROM refresh_tokens WHERE token_id=$1", tokenID).Scan(&refresh)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &refresh, err
}
