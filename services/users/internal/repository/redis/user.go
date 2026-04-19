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

type UserRepoCacheRedis struct {
	client *redis.Client
}

func NewUserCacheRepoRedis(client *redis.Client) repository.UserCacheRepository {
	return &UserRepoCacheRedis{client: client}
}

func (r *UserRepoCacheRedis) SetUser(ctx context.Context, user domain.User) error {
	userJson, err := json.Marshal(user)

	if err != nil {
		return err
	}

	cfg := config.MustLoad()

	err = r.client.Set(ctx, fmt.Sprintf("services:users:cache:users:%d", user.ID), userJson, time.Minute*time.Duration(cfg.Data.User.Cache.Expires)).Err()

	return err
}

func (r *UserRepoCacheRedis) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	user := domain.User{}
	val, err := r.client.Get(ctx, fmt.Sprintf("services:users:cache:users:%d", id)).Result()

	if err != nil {
		if err == redis.Nil {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepoCacheRedis) DeleteUserByID(ctx context.Context, id int64) error {
	return r.client.Del(ctx, fmt.Sprintf("services:users:cache:users:%d", id)).Err()
}
