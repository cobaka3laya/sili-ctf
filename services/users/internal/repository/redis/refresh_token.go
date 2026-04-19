package repository_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/domain"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/repository"
	"github.com/redis/go-redis/v9"
)

type RefreshTokenRepoCacheRedis struct {
	client *redis.Client
}

func NewRefreshTokenCacheRepoRedis(client *redis.Client) repository.RefreshTokenCacheRepository {
	return &RefreshTokenRepoCacheRedis{client: client}
}

func (r *RefreshTokenRepoCacheRedis) SetRefreshToken(ctx context.Context, refreshToken domain.RefreshToken) error {
	refreshTokenJson, err := json.Marshal(refreshToken)

	if err != nil {
		return err
	}

	err = r.client.Set(ctx, fmt.Sprintf("services:users:cache:refresh_tokens:%d", refreshToken.ID), refreshTokenJson, time.Until(refreshToken.ExpiresAt)).Err()

	return err
}

func (r *RefreshTokenRepoCacheRedis) GetRefreshTokenByOwnerID(ctx context.Context, id int64) (*domain.RefreshToken, error) {
	refreshToken := domain.RefreshToken{}
	val, err := r.client.Get(ctx, fmt.Sprintf("services:users:cache:refresh_tokens:%d", id)).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &refreshToken)

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepoCacheRedis) DeleteRefreshTokenByOwnerID(ctx context.Context, id int64) error {
	return r.client.Del(ctx, fmt.Sprintf("services:users:cache:refresh_tokens:%d", id)).Err()
}
