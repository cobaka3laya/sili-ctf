package repository

import (
	"context"
	"errors"
	"time"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/domain"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, ownerID int64, tokenHash string, expiresAt time.Time) (*domain.RefreshToken, error)

	GetRefreshTokenByOwnerID(ctx context.Context, id int64) (*domain.RefreshToken, error)

	DeleteRefreshTokenByOwnerID(ctx context.Context, id int64) error
}

type RefreshTokenCacheRepository interface {
	SetRefreshToken(ctx context.Context, refreshToken domain.RefreshToken) error

	GetRefreshTokenByOwnerID(ctx context.Context, id int64) (*domain.RefreshToken, error)

	DeleteRefreshTokenByOwnerID(ctx context.Context, id int64) error
}

var ErrRefreshTokenNotFound = errors.New("refresh token not found")
