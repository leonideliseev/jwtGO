package repository

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leonideliseev/jwtGO/models"
)

type TokensRepo struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewTokensRepo(db *pgxpool.Pool) *TokensRepo {
	return &TokensRepo{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

const (
	refreshTokensTable = "refresh_tokens"
)

const (
	token_id_F   = "token_id"
	ip_F         = "ip"
	token_hash_F = "token_hash"
)

func (r *TokensRepo) Create(ctx context.Context, refresh *models.Refresh) error {
	q, args, err := r.builder.
		Insert(refreshTokensTable).
		Columns(token_id_F, ip_F, token_hash_F).
		Values(refresh.TokenID, refresh.IP, refresh.RefreshTokenHash).
		ToSql()

	_, err = r.db.Exec(ctx, q, args)
	if err != nil {
		return err
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

func (r *TokensRepo) Update(ctx context.Context, oldTokenID string, newRefresh *models.Refresh) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = r.Delete(ctx, oldTokenID)
	if err != nil {
		return err
	}

	err = r.Create(ctx, newRefresh)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *TokensRepo) Delete(ctx context.Context, tokenID string) error {
	q, args, err := r.builder.
		Delete(refreshTokensTable).
		Where(squirrel.Eq{token_id_F: tokenID}).
		ToSql()
	if err != nil {
		return err
	}

	commandTag, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
