package repository_postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/db/dbtx"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/db/txcontext"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/dto"
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

func (r *UserRepoPG) CreateUser(ctx context.Context, data dto.CreateUserDTOInput) (*dto.CreateUserDTOOutput, error) {
	executor := r.getExecutor(ctx)

	now := time.Now()
	out := &dto.CreateUserDTOOutput{CreatedAt: now, Username: data.Username}

	err := executor.QueryRow(
		ctx,
		`INSERT INTO
		users (username, created_at, profile_picture_url)
	 	VALUES ($1, $2, $3)
		RETURNING id`,
		data.Username,
		now,
		"",
	).Scan(&out.ID)

	if err != nil {
		return nil, fmt.Errorf("repo: create user failed: query row scan: %w", err)
	}

	return out, nil
}

func (r *UserRepoPG) GetUserByID(ctx context.Context, id int64) (*dto.GetUserDTOOutput, error) {
	out := &dto.GetUserDTOOutput{ID: id}

	err := r.pool.QueryRow(
		ctx,
		`SELECT username, created_at, profile_picture_url
		FROM users
		WHERE id = $1`,
		id,
	).Scan(out.Username, out.CreatedAt, out.ProfilePictureURL)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("repo: get user by id for %d failed: query row scan: %w", id, err)
	}

	return out, nil
}

func (r *UserRepoPG) GetUserByUsername(ctx context.Context, username string) (*dto.GetUserDTOOutput, error) {
	out := &dto.GetUserDTOOutput{Username: username}

	err := r.pool.QueryRow(
		ctx,
		`SELECT id, created_at, profile_picture_url
		FROM users
		WHERE username = $1`,
		username,
	).Scan(out.ID, out.CreatedAt, out.ProfilePictureURL)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("repo: get user by username for %s failed: query row scan: %w", username, err)
	}

	return out, nil
}

func (r *UserRepoPG) UpdateUserByID(ctx context.Context, id int64, data dto.UpdateUserDTOInput) error {
	executor := r.getExecutor(ctx)

	updateArgs := []string{}
	queryArgs := []any{}
	i := 1

	if data.Username != nil {
		updateArgs = append(updateArgs, fmt.Sprintf("username = $%d", i))
		queryArgs = append(queryArgs, *data.Username)
		i++
	}

	if data.ProfilePictureURL != nil {
		updateArgs = append(updateArgs, fmt.Sprintf("profile_picture_url = $%d", i))
		queryArgs = append(queryArgs, *data.ProfilePictureURL)
		i++
	}

	queryArgs = append(queryArgs, id)

	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d",
		strings.Join(updateArgs, ", "),
		i,
	)

	_, err := executor.Exec(ctx, query, queryArgs...)

	if err != nil {
		return fmt.Errorf("repo: update user by id for %d failed: exec: %w", id, err)
	}

	return nil
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
		return fmt.Errorf("repo: delete user by id for %d failed: exec: %w", id, err)
	}

	return nil
}
