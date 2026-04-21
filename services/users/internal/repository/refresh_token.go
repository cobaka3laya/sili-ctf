package repository

import (
	"context"
	"errors"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/dto"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, data dto.CreateRefreshTokenDTOInput) (*dto.CreateRefreshTokenDTOOutput, error)

	GetRefreshTokenByOwnerID(ctx context.Context, id int64) (*dto.GetRefreshTokenByOwnerIDDTOOutput, error)

	DeleteRefreshTokenByOwnerID(ctx context.Context, id int64) error
}

type RefreshTokenCacheRepository interface {
	SetRefreshToken(ctx context.Context, data dto.SetRefreshTokenDTOInput) error

	GetRefreshTokenByOwnerID(ctx context.Context, id int64) (*dto.GetRefreshTokenByOwnerIDDTOOutput, error)

	DeleteRefreshTokenByOwnerID(ctx context.Context, id int64) error
}

var ErrRefreshTokenNotFound = errors.New("refresh token not found")
var ErrInvalidRefreshTokenFound = errors.New("invalid refresh token found")
