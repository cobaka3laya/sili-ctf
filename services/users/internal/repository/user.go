package repository

import (
	"context"
	"errors"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, username string) (*domain.User, error)

	GetUserByID(ctx context.Context, id int64) (*domain.User, error)

	UpdateUserByID(ctx context.Context, id int64, user domain.User) (*domain.User, error)

	DeleteUserByID(ctx context.Context, id int64) error
}

type UserCacheRepository interface {
	SetUser(ctx context.Context, user domain.User) error

	GetUserByID(ctx context.Context, id int64) (*domain.User, error)

	DeleteUserByID(ctx context.Context, id int64) error
}

var ErrUserNotFound = errors.New("user not found")
