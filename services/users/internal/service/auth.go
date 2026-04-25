package service

import (
	"context"

	"github.com/cobaka3laya/sili-ctf/services/users/internal/dto"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/repository"
)

type AuthService struct {
	userRepo              repository.UserRepository
	userCacheRepo         repository.UserCacheRepository
	userAuthDataRepo      repository.UserAuthDataRepository
	userAuthDataCacheRepo repository.UserAuthDataCacheRepository
	refreshTokenRepo      repository.RefreshTokenRepository
	refreshTokenCacheRepo repository.RefreshTokenCacheRepository
}

type AuthServiceDeps struct {
	repos struct {
		userRepo              repository.UserRepository
		userCacheRepo         repository.UserCacheRepository
		userAuthDataRepo      repository.UserAuthDataRepository
		userAuthDataCacheRepo repository.UserAuthDataCacheRepository
		refreshTokenRepo      repository.RefreshTokenRepository
		refreshTokenCacheRepo repository.RefreshTokenCacheRepository
	}
}

func NewAuthService(deps AuthServiceDeps) *AuthService {
	return &AuthService{
		userRepo:              deps.repos.userRepo,
		userCacheRepo:         deps.repos.userCacheRepo,
		userAuthDataRepo:      deps.repos.userAuthDataRepo,
		userAuthDataCacheRepo: deps.repos.userAuthDataCacheRepo,
		refreshTokenRepo:      deps.repos.refreshTokenRepo,
		refreshTokenCacheRepo: deps.repos.refreshTokenCacheRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, data dto.LoginDTOInput) (*dto.LoginDTOOutput, error)

func (s *AuthService) Register(ctx context.Context, data dto.RegisterDTOInput) (*dto.RegisterDTOOutput, error)

func (s *AuthService) Refresh(ctx context.Context, data dto.RefreshDTOInput) (*dto.RefreshDTOOutput, error)

func (s *AuthService) Logout(ctx context.Context, data dto.LogoutDTOInput) error
