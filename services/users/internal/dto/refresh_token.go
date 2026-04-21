package dto

import "time"

type CreateRefreshTokenDTOInput struct {
	OwnerID   int64
	TokenHash string
	ExpiresAt time.Time
}

type CreateRefreshTokenDTOOutput struct {
	ID        int64
	OwnerID   int64
	TokenHash string
	ExpiresAt time.Time
}

type GetRefreshTokenByOwnerIDDTOOutput struct {
	ID        int64
	TokenHash string
	ExpiresAt time.Time
}

type SetRefreshTokenDTOInput struct {
	ID        int64
	OwnerID   int64
	TokenHash string
	ExpiresAt time.Time
}
