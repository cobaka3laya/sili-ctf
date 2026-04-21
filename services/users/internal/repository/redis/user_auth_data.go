package repository_redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cobaka3laya/sili-ctf/services/users/config"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/dto"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/repository"
	"github.com/redis/go-redis/v9"
)

type UserAuthDataRepoCacheRedis struct {
	client *redis.Client
}

func NewUserAuthDataCacheRepoRedis(client *redis.Client) repository.UserAuthDataCacheRepository {
	return &UserAuthDataRepoCacheRedis{client: client}
}

func (r *UserAuthDataRepoCacheRedis) key(id int64) string {
	return fmt.Sprintf("services:users:cache:users_auth_data:user_id:%d", id)
}

func (r *UserAuthDataRepoCacheRedis) SetUserAuthData(ctx context.Context, data dto.SetUserAuthDataByUserIDDTOInput) error {
	pipe := r.client.Pipeline()
	currentKey := r.key(data.UserID)

	values := map[string]any{}

	if data.ID != nil {
		values["id"] = *data.ID
	}

	if data.HashedPassword != nil {
		values["hashed_password"] = *data.HashedPassword
	}

	if len(values) == 0 {
		return repository.ErrUserAuthDataNotFound
	}

	pipe.HSet(ctx, currentKey, values)
	pipe.Expire(ctx, currentKey, time.Duration(config.MustLoad().Data.UserAuthData.Cache.Expires)*time.Minute)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *UserAuthDataRepoCacheRedis) GetUserAuthDataByUserID(ctx context.Context, id int64) (*dto.GetUserAuthDataByUserIDDTOOutput, error) {
	out := &dto.GetUserAuthDataByUserIDDTOOutput{}
	values, err := r.client.HGetAll(ctx, r.key(id)).Result()

	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, repository.ErrInvalidUserAuthDataFound
	}

	if userAuthDataID, ok := values["id"]; ok {
		convertedID, err := strconv.ParseInt(userAuthDataID, 10, 64)
		if err != nil {
			return nil, err
		}
		out.ID = convertedID
	} else {
		return nil, repository.ErrInvalidUserAuthDataFound
	}

	if hashedPassword, ok := values["hashed_password"]; ok && hashedPassword != "" {
		out.HashedPassword = hashedPassword
	} else {
		return nil, repository.ErrInvalidUserAuthDataFound
	}

	return out, nil
}

func (r *UserAuthDataRepoCacheRedis) DeleteUserAuthDataByUserID(ctx context.Context, id int64) error {
	return r.client.Del(ctx, r.key(id)).Err()
}
