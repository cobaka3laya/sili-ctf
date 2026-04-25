package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/cobaka3laya/sili-ctf/services/users/config"
)

type RefreshToken struct {
	ID        int64     `json:"id"`
	OwnerID   int64     `json:"owner_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
}

var ErrInvalidOwnerID = errors.New("owner ID is invalid")
var ErrTokenHashEmpty = errors.New("token hash is empty")
var ErrAlreadyExpired = errors.New("refresh token is already expired")
var ErrTTLTooShort = errors.New("token ttl is to short")

func NewRefreshToken(ownerID int64, tokenHash string, expiresAt time.Time, cfg *config.RefreshTokenConfig) (*RefreshToken, error) {
	if ownerID <= 0 {
		return nil, ErrInvalidOwnerID
	}

	if strings.TrimSpace(tokenHash) == "" {
		return nil, ErrTokenHashEmpty
	}

	if expiresAt.Before(time.Now()) {
		return nil, ErrAlreadyExpired
	}

	if time.Until(expiresAt) < time.Minute*time.Duration(cfg.MinTTL) {
		return nil, ErrTTLTooShort
	}

	return &RefreshToken{
		OwnerID:   ownerID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}, nil
}
