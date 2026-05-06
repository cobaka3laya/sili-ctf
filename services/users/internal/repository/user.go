package repository

import (
	"context"
	"errors"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/dto"
)

type UserRepository interface {
	CreateUser(ctx context.Context, data dto.CreateUserDTOInput) (*dto.CreateUserDTOOutput, error)

	GetUserByID(ctx context.Context, id int64) (*dto.GetUserDTOOutput, error)

	GetUserByUsername(ctx context.Context, username string) (*dto.GetUserDTOOutput, error)

	UpdateUserByID(ctx context.Context, id int64, data dto.UpdateUserDTOInput) error

	DeleteUserByID(ctx context.Context, id int64) error
}

type UserCacheRepository interface {
	SetUser(ctx context.Context, data dto.SetUserDTOInput) error

	GetUserByID(ctx context.Context, id int64) (*dto.GetUserDTOOutput, error)

	GetUserIDByUsername(ctx context.Context, username string) (int64, error)

	DeleteUserByID(ctx context.Context, id int64) error
}

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidUserFound = errors.New("invalid user found")
