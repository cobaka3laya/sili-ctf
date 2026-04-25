package service

import "github.com/cobaka3laya/sili-ctf/services/users/internal/repository"

type AuthService struct {
	repo  repository.UserAuthDataRepository
	cache repository.UserAuthDataCacheRepository
}

func NewAuthService(repo repository.UserAuthDataRepository, cache repository.UserAuthDataCacheRepository) *AuthService {
	return &AuthService{
		repo:  repo,
		cache: cache,
	}
}
