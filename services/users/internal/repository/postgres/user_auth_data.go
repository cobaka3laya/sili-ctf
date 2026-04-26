package repository_postgres

import (
	"context"
	"fmt"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/db/dbtx"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/db/txcontext"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/dto"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserAuthDataPG struct {
	pool *pgxpool.Pool
}

func NewUserAuthDataRepo(pool *pgxpool.Pool) repository.UserAuthDataRepository {
	return &UserAuthDataPG{pool: pool}
}

func (r *UserAuthDataPG) getExecutor(ctx context.Context) dbtx.DBTX {
	if tx, ok := txcontext.TxFromContext(ctx); ok {
		return tx
	}
	return r.pool
}

func (r *UserAuthDataPG) CreateUserAuthData(ctx context.Context, data dto.CreateUserAuthDataDTOInput) (*dto.CreateUserAuthDataDTOOutput, error) {
	executor := r.getExecutor(ctx)

	out := &dto.CreateUserAuthDataDTOOutput{}

	err := executor.QueryRow(
		ctx,
		`INSERT INTO
		users_auth_data (user_id, hashed_password)
	 	VALUES ($1, $2)
		RETURNING id`,
		data.UserID,
		data.HashedPassword,
	).Scan(&out.ID)

	if err != nil {
		return nil, fmt.Errorf("repo: create user auth data failed: query row scan: %w", err)
	}

	return out, nil
}

func (r *UserAuthDataPG) GetUserAuthDataByUserID(ctx context.Context, id int64) (*dto.GetUserAuthDataByUserIDDTOOutput, error) {
	out := &dto.GetUserAuthDataByUserIDDTOOutput{}

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, hashed_password
		FROM users_auth_data
		WHERE user_id = $1`,
		id,
	).Scan(&out.ID, &out.HashedPassword)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrUserAuthDataNotFound
		}
		return nil, fmt.Errorf("repo: get user auth data by user id for %d failed: query row scan: %w", id, err)
	}

	return out, nil
}

func (r *UserAuthDataPG) UpdateUserAuthDataByUserID(ctx context.Context, id int64, data dto.UpdateUserAuthDataDTOInput) error {
	executor := r.getExecutor(ctx)

	_, err := executor.Exec(
		ctx,
		`UPDATE users_auth_data
		SET hashed_password = $1
		WHERE user_id = $2
		`,
		data.HashedPassword,
		id,
	)

	if err != nil {
		return fmt.Errorf("repo: update user auth data by user id for %d failed: exec: %w", id, err)
	}

	return nil
}

func (r *UserAuthDataPG) DeleteUserAuthDataByUserID(ctx context.Context, id int64) error {
	executor := r.getExecutor(ctx)

	_, err := executor.Exec(
		ctx,
		`DELETE
	 	FROM users_auth_data
		WHERE user_id = $1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("repo: delete user auth data by user id for %d failed: exec: %w", id, err)
	}

	return nil
}
