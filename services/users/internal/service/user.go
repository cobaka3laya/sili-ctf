package service

import "github.com/cobaka3laya/sili-ctf/services/users/internal/repository"

type UserService struct {
	repo  repository.UserRepository
	cache repository.UserCacheRepository
}

func NewUserService(repo repository.UserRepository, cache repository.UserCacheRepository) *UserService {
	return &UserService{
		repo:  repo,
		cache: cache,
	}
}
