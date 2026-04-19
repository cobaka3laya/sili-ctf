package repository

import (
	"context"
	"errors"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/domain"
)

type UserAuthDataRepository interface {
	CreateUserAuthData(ctx context.Context, userID int64, hashedPassword string) (*domain.UserAuthData, error)

	GetUserAuthDataByUserID(ctx context.Context, id int64) (*domain.UserAuthData, error)

	UpdateUserAuthDataByUserID(ctx context.Context, id int64, userAuthData domain.UserAuthData) (*domain.UserAuthData, error)

	DeleteUserAuthDataByUserID(ctx context.Context, id int64) error
}

type UserAuthDataCacheRepository interface {
	SetUserAuthData(ctx context.Context, userAuthData domain.UserAuthData) error

	GetUserAuthDataByUserID(ctx context.Context, id int64) (*domain.UserAuthData, error)

	DeleteUserAuthDataByUserID(ctx context.Context, id int64) error
}

var ErrUserAuthDataNotFound = errors.New("user auth data not found")
