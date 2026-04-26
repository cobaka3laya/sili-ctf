package repository_redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/dto"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/repository"
	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepoCacheRedis struct {
	client *redis.Client
}

func NewRefreshTokenCacheRepoRedis(client *redis.Client) repository.RefreshTokenCacheRepository {
	return &RefreshTokenRepoCacheRedis{client: client}
}

func (r *RefreshTokenRepoCacheRedis) key(id int64) string {
	return fmt.Sprintf("services:users:cache:refresh_tokens:owner_id:%d", id)
}

func (r *RefreshTokenRepoCacheRedis) SetRefreshToken(ctx context.Context, data dto.SetRefreshTokenDTOInput) error {
	pipe := r.client.Pipeline()
	currentKey := r.key(data.OwnerID)

	values := map[string]any{
		"id":         data.ID,
		"token_hash": data.TokenHash,
		"expires_at": data.ExpiresAt.Format(time.RFC3339),
	}

	pipe.HSet(ctx, currentKey, values)
	pipe.Expire(ctx, currentKey, time.Until(data.ExpiresAt))
	_, err := pipe.Exec(ctx)

	if err != nil {
		return fmt.Errorf("cache: set refresh token for id %d failed: pipe exec: %w", data.ID, err)
	}

	return err
}

func (r *RefreshTokenRepoCacheRedis) GetRefreshTokenByOwnerID(ctx context.Context, id int64) (*dto.GetRefreshTokenByOwnerIDDTOOutput, error) {
	out := &dto.GetRefreshTokenByOwnerIDDTOOutput{}
	values, err := r.client.HGetAll(ctx, r.key(id)).Result()

	if err != nil {
		return nil, fmt.Errorf("cache: get refresh token by owner id for %d failed: HGetAll: %w", id, err)
	}

	if len(values) == 0 {
		return nil, repository.ErrRefreshTokenNotFound
	}

	if tokenID, ok := values["id"]; ok {
		convertedTokenID, err := strconv.ParseInt(tokenID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cache: get refresh token by owner id for %d failed: ParseInt tokenID for %s: %w", id, tokenID, err)
		}
		if convertedTokenID <= 0 {
			return nil, repository.ErrInvalidRefreshTokenFound
		}
		out.ID = convertedTokenID
	} else {
		return nil, repository.ErrInvalidRefreshTokenFound
	}

	if tokenHash, ok := values["token_hash"]; ok && tokenHash != "" {
		out.TokenHash = tokenHash
	} else {
		return nil, repository.ErrInvalidRefreshTokenFound
	}

	if expiresAt, ok := values["expires_at"]; ok && expiresAt != "" {
		t, err := time.Parse(time.RFC3339, expiresAt)
		if err != nil {
			return nil, fmt.Errorf("cache: get refresh token by owner id for %d failed: time.Parse expiresAt for %s: %w", id, expiresAt, err)
		}
		out.ExpiresAt = t
	} else {
		return nil, repository.ErrInvalidRefreshTokenFound
	}

	return out, nil
}

func (r *RefreshTokenRepoCacheRedis) DeleteRefreshTokenByOwnerID(ctx context.Context, id int64) error {
	err := r.client.Del(ctx, r.key(id)).Err()

	if err != nil {
		return fmt.Errorf("cache: delete refresh token by owner id for %d failed: Del: %w", id, err)
	}

	return nil
}
