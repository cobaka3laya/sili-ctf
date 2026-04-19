package repository_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cobaka3laya/sili-ctf/services/users/config"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/domain"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/repository"
	"github.com/redis/go-redis/v9"
)

type UserAuthDataRepoCacheRedis struct {
	client *redis.Client
}

func NewUserAuthDataCacheRepoRedis(client *redis.Client) repository.UserAuthDataCacheRepository {
	return &UserAuthDataRepoCacheRedis{client: client}
}

func (r *UserAuthDataRepoCacheRedis) SetUserAuthData(ctx context.Context, userAuthData domain.UserAuthData) error {
	userAuthDataJson, err := json.Marshal(userAuthData)

	if err != nil {
		return err
	}

	cfg := config.MustLoad()

	err = r.client.Set(ctx, fmt.Sprintf("services:users:cache:users_auth_data:%d", userAuthData.UserID), userAuthDataJson, time.Minute*time.Duration(cfg.Data.UserAuthData.Cache.Expires)).Err()

	return err
}

func (r *UserAuthDataRepoCacheRedis) GetUserAuthDataByUserID(ctx context.Context, id int64) (*domain.UserAuthData, error) {
	userAuthData := domain.UserAuthData{}
	val, err := r.client.Get(ctx, fmt.Sprintf("services:users:cache:users_auth_data:%d", id)).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &userAuthData)

	if err != nil {
		return nil, err
	}

	return &userAuthData, nil
}

func (r *UserAuthDataRepoCacheRedis) DeleteUserAuthDataByUserID(ctx context.Context, id int64) error {
	return r.client.Del(ctx, fmt.Sprintf("services:users:cache:users_auth_data:%d", id)).Err()
}
