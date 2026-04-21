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
		return nil, err
	}

	return out, nil
}

func (r *UserRepoPG) GetUserByID(ctx context.Context, id int64) (*dto.GetUserByIDDTOOutput, error) {
	output := dto.GetUserByIDDTOOutput{}

	err := r.pool.QueryRow(
		ctx,
		`SELECT username, created_at, profile_picture_url
		FROM users
		WHERE id = $1`,
		id,
	).Scan(&output.Username, &output.CreatedAt, &output.ProfilePictureURL)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	return &output, nil
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

	return err
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

	return err
}
