package repository

import (
	"context"
	"errors"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/dto"
)

type UserAuthDataRepository interface {
	CreateUserAuthData(ctx context.Context, data dto.CreateUserAuthDataDTOInput) (*dto.CreateUserAuthDataDTOOutput, error)

	GetUserAuthDataByUserID(ctx context.Context, id int64) (*dto.GetUserAuthDataByUserIDDTOOutput, error)

	UpdateUserAuthDataByUserID(ctx context.Context, id int64, data dto.UpdateUserAuthDataDTOInput) error

	DeleteUserAuthDataByUserID(ctx context.Context, id int64) error
}

type UserAuthDataCacheRepository interface {
	SetUserAuthData(ctx context.Context, data dto.SetUserAuthDataByUserIDDTOInput) error

	GetUserAuthDataByUserID(ctx context.Context, id int64) (*dto.GetUserAuthDataByUserIDDTOOutput, error)

	DeleteUserAuthDataByUserID(ctx context.Context, id int64) error
}

var ErrUserAuthDataNotFound = errors.New("user auth data not found")
var ErrInvalidUserAuthDataFound = errors.New("invalid user auth data found")
