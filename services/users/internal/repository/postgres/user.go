package repository_postgres

import (
	"context"
	"time"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/db/dbtx"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/db/txcontext"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/domain"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepoPG struct {
	pool *pgxpool.Pool
}

func NewUserRepoPG(pool *pgxpool.Pool) repository.UserRepository {
	return &UserRepoPG{pool: pool}
}

func (r *UserRepoPG) getExecutor(ctx context.Context) dbtx.DBTX {
	if tx, ok := txcontext.TxFromContext(ctx); ok {
		return tx
	}
	return r.pool
}

func (r *UserRepoPG) CreateUser(ctx context.Context, username string) (*domain.User, error) {
	executor := r.getExecutor(ctx)

	now := time.Now()
	createdUser, err := domain.NewUser(username, now)

	if err != nil {
		return nil, err
	}

	err = executor.QueryRow(
		ctx,
		`INSERT INTO
		users (username, created_at, profile_picture_url)
	 	VALUES ($1, $2, $3)
		RETURNING ID`,
		username,
		now,
		"",
	).Scan(&createdUser.ID)

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (r *UserRepoPG) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	user := &domain.User{}

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, username, created_at, profile_picture_url
		FROM users
		WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.ProfilePictureUrl)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepoPG) UpdateUserByID(ctx context.Context, id int64, user domain.User) (*domain.User, error) {
	executor := r.getExecutor(ctx)

	_, err := executor.Exec(
		ctx,
		`UPDATE users
		SET username=$1, created_at=$2, profile_picture_url=$3
		WHERE id = $4
		`,
		user.Username,
		user.CreatedAt,
		user.ProfilePictureUrl,
		id,
	)

	if err != nil {
		return nil, err
	}

	updatedUser := user
	updatedUser.ID = id

	return &updatedUser, nil
}

func (r *UserRepoPG) DeleteUserByID(ctx context.Context, id int64) error {
	executor := r.getExecutor(ctx)

	_, err := executor.Exec(
		ctx,
		`DELETE
	 	FROM users
		WHERE id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
