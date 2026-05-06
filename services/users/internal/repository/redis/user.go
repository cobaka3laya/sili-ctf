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

type UserRepoCacheRedis struct {
	client *redis.Client
	cfg    *config.UserCacheConfig
}

func NewUserCacheRepoRedis(client *redis.Client) repository.UserCacheRepository {
	return &UserRepoCacheRedis{client: client}
}

func (r *UserRepoCacheRedis) idKey(id int64) string {
	return fmt.Sprintf("services:users:cache:users:id:%d", id)
}

func (r *UserRepoCacheRedis) usernameKey(username string) string {
	return fmt.Sprintf("services:users:cache:users:username:%s", username)
}

func (r *UserRepoCacheRedis) SetUser(ctx context.Context, data dto.SetUserDTOInput) error {
	pipe := r.client.Pipeline()

	idKey := r.idKey(data.ID)
	var usernameKey string

	if data.Username != nil {
		usernameKey = r.usernameKey(*data.Username)
	}

	values := map[string]any{}

	if data.Username != nil {
		values["username"] = *data.Username
	}

	if data.CreatedAt != nil {
		values["created_at"] = data.CreatedAt.Format(time.RFC3339)
	}

	if data.ProfilePictureURL != nil {
		values["profile_picture_url"] = *data.ProfilePictureURL
	}

	if len(values) == 0 {
		return repository.ErrInvalidUserFound
	}

	pipe.HSet(ctx, idKey, values)
	pipe.Expire(ctx, idKey, time.Duration(r.cfg.Expires)*time.Minute)

	if usernameKey != "" {
		pipe.Set(ctx, usernameKey, data.ID, time.Duration(r.cfg.Expires)*time.Minute)
	}

	_, err := pipe.Exec(ctx)

	if err != nil {
		return fmt.Errorf("cache: set user for id %d failed: pipe exec: %w", data.ID, err)
	}

	return nil
}

func (r *UserRepoCacheRedis) GetUserByID(ctx context.Context, id int64) (*dto.GetUserDTOOutput, error) {
	out := &dto.GetUserDTOOutput{ID: id}
	values, err := r.client.HGetAll(ctx, r.idKey(id)).Result()

	if err != nil {
		return nil, fmt.Errorf("cache: get user by id for %d failed: HGetAll: %w", id, err)
	}

	if len(values) == 0 {
		return nil, repository.ErrUserNotFound
	}

	if username, ok := values["username"]; ok && username != "" {
		out.Username = username
	} else {
		return nil, repository.ErrInvalidUserFound
	}

	if createdAt, ok := values["created_at"]; ok && createdAt != "" {
		t, err := time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("cache: get user by id for %d failed: time.Parse createdAt for %s: %w", id, createdAt, err)
		}
		out.CreatedAt = t
	} else {
		return nil, repository.ErrInvalidUserFound
	}

	profilePictureURL, ok := values["profile_picture_url"]
	if !ok {
		out.ProfilePictureURL = ""
	} else {
		out.ProfilePictureURL = profilePictureURL
	}

	return out, nil
}

func (r *UserRepoCacheRedis) GetUserIDByUsername(ctx context.Context, username string) (int64, error) {
	id, err := r.client.Get(ctx, r.usernameKey(username)).Result()

	if err != nil {
		return 0, fmt.Errorf("cache: get user id by username for %s failed: Get: %w", username, err)
	}

	parsedID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("cache: get user id by username for %s failed: strconv.Parse id for %s: %w", username, id, err)
	}

	return parsedID, nil
}

func (r *UserRepoCacheRedis) DeleteUserByID(ctx context.Context, id int64) error {
	err := r.client.Del(ctx, r.idKey(id)).Err()

	if err != nil {
		return fmt.Errorf("cache delete user by id for %d failed: Del: %w")
	}

	return nil
}
