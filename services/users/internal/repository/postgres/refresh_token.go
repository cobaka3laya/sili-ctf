package repository_postgres

import (
	"context"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/db/dbtx"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/db/txcontext"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/dto"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepoPG struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepoPG(pool *pgxpool.Pool) repository.RefreshTokenRepository {
	return &RefreshTokenRepoPG{pool: pool}
}

func (r *RefreshTokenRepoPG) getExecutor(ctx context.Context) dbtx.DBTX {
	if tx, ok := txcontext.TxFromContext(ctx); ok {
		return tx
	}
	return r.pool
}

func (r *RefreshTokenRepoPG) CreateRefreshToken(ctx context.Context, data dto.CreateRefreshTokenDTOInput) (*dto.CreateRefreshTokenDTOOutput, error) {
	executor := r.getExecutor(ctx)

	createdRefreshToken := &dto.CreateRefreshTokenDTOOutput{OwnerID: data.OwnerID, TokenHash: data.TokenHash, ExpiresAt: data.ExpiresAt}
	err := executor.QueryRow(
		ctx,
		`INSERT INTO
		refresh_tokens (owner_id, token_hash, expires_at)
	 	VALUES ($1, $2, $3)
		RETURNING id`,
		data.OwnerID,
		data.TokenHash,
		data.ExpiresAt,
		"",
	).Scan(&createdRefreshToken.ID)

	if err != nil {
		return nil, err
	}

	return createdRefreshToken, nil
}

func (r *RefreshTokenRepoPG) GetRefreshTokenByOwnerID(ctx context.Context, id int64) (*dto.GetRefreshTokenByOwnerIDDTOOutput, error) {
	refreshToken := &dto.GetRefreshTokenByOwnerIDDTOOutput{}

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, token_hash, expires_at
		FROM refresh_tokens
		WHERE owner_id = $1`,
		id,
	).Scan(&refreshToken.ID, &refreshToken.TokenHash, &refreshToken.ExpiresAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepoPG) DeleteRefreshTokenByOwnerID(ctx context.Context, id int64) error {
	executor := r.getExecutor(ctx)

	_, err := executor.Exec(
		ctx,
		`DELETE
	 	FROM refresh_tokens
		WHERE owner_id = $1`,
		id,
	)

	return err
}
