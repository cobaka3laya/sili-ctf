package repository_postgres

import (
	"context"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/db/dbtx"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/db/txcontext"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/domain"
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

func (r *UserAuthDataPG) CreateUserAuthData(ctx context.Context, userID int64, hashedPassword string) (*domain.UserAuthData, error) {
	executor := r.getExecutor(ctx)

	createdUserAuthData, err := domain.NewUserAuthData(userID, hashedPassword)

	if err != nil {
		return nil, err
	}

	err = executor.QueryRow(
		ctx,
		`INSERT INTO
		users_auth_data (user_id, hashed_password)
	 	VALUES ($1, $2)
		RETURNING ID`,
		userID,
		hashedPassword,
	).Scan(&createdUserAuthData.ID)

	if err != nil {
		return nil, err
	}

	return createdUserAuthData, nil
}

func (r *UserAuthDataPG) GetUserAuthDataByUserID(ctx context.Context, id int64) (*domain.UserAuthData, error) {
	userAuthData := &domain.UserAuthData{}

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, user_id, hashed_password
		FROM users_auth_data
		WHERE user_id = $1`,
		id,
	).Scan(&userAuthData.ID, &userAuthData.UserID, &userAuthData.HashedPassword)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrUserAuthDataNotFound
		}
		return nil, err
	}

	return userAuthData, nil
}

func (r *UserAuthDataPG) UpdateUserAuthDataByUserID(ctx context.Context, id int64, userAuthData domain.UserAuthData) (*domain.UserAuthData, error) {
	executor := r.getExecutor(ctx)

	_, err := executor.Exec(
		ctx,
		`UPDATE users_auth_data
		SET hashed_password = $1
		WHERE user_id = $2
		`,
		userAuthData.HashedPassword,
		id,
	)

	if err != nil {
		return nil, err
	}

	updatedUserAuthData := userAuthData

	return &updatedUserAuthData, nil
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
		return err
	}

	return nil
}
