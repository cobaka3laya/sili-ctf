package service

import (
	"context"
	"time"
	"unicode/utf8"

	"github.com/cobaka3laya/sili-ctf/services/users/config"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/db/txcontext"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/domain"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/dto"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/repository"
	"github.com/cobaka3laya/sili-ctf/services/users/internal/service/utils"
)

type AuthService struct {
	userRepo              repository.UserRepository
	userCacheRepo         repository.UserCacheRepository
	userAuthDataRepo      repository.UserAuthDataRepository
	userAuthDataCacheRepo repository.UserAuthDataCacheRepository
	refreshTokenRepo      repository.RefreshTokenRepository
	refreshTokenCacheRepo repository.RefreshTokenCacheRepository
	cfg                   *config.Config
	tx                    txcontext.TxManager
}

type AuthServiceDeps struct {
	userRepo              repository.UserRepository
	userCacheRepo         repository.UserCacheRepository
	userAuthDataRepo      repository.UserAuthDataRepository
	userAuthDataCacheRepo repository.UserAuthDataCacheRepository
	refreshTokenRepo      repository.RefreshTokenRepository
	refreshTokenCacheRepo repository.RefreshTokenCacheRepository
	cfg                   *config.Config
	tx                    txcontext.TxManager
}

func NewAuthService(deps AuthServiceDeps) *AuthService {
	return &AuthService{
		userRepo:              deps.userRepo,
		userCacheRepo:         deps.userCacheRepo,
		userAuthDataRepo:      deps.userAuthDataRepo,
		userAuthDataCacheRepo: deps.userAuthDataCacheRepo,
		refreshTokenRepo:      deps.refreshTokenRepo,
		refreshTokenCacheRepo: deps.refreshTokenCacheRepo,
		cfg:                   deps.cfg,
		tx:                    deps.tx,
	}
}

func (s *AuthService) Login(ctx context.Context, data dto.LoginDTOInput) (*dto.LoginDTOOutput, error)

func (s *AuthService) Register(ctx context.Context, data dto.RegisterDTOInput) (*dto.RegisterDTOOutput, error) {
	if s.cfg.Service.Auth.MinPasswordLength > utf8.RuneCountInString(data.Password) {
		return nil, ErrPasswordTooShort
	}

	out := &dto.RegisterDTOOutput{}

	var outUser *domain.User
	var outUserAuthData *domain.UserAuthData
	var outRefreshToken *domain.RefreshToken

	hashedPassword, err := utils.HashPassword(data.Password)

	if err != nil {
		return nil, err
	}

	out.RefreshToken, err = utils.GenerateRefreshToken()

	if err != nil {
		return nil, err
	}

	hashedRefreshToken := utils.HashRefreshToken(out.RefreshToken)

	err = s.tx.WithinTx(
		ctx,
		func(ctx context.Context) error {
			user, err := s.userRepo.CreateUser(ctx, dto.CreateUserDTOInput{Username: data.Username})

			if err != nil {
				return err
			}

			outUser, err = domain.NewUser(user.Username, user.CreatedAt, &s.cfg.Data.User)

			if err != nil {
				return err
			}

			_, err = s.userAuthDataRepo.CreateUserAuthData(
				ctx,
				dto.CreateUserAuthDataDTOInput{
					UserID:         outUser.ID,
					HashedPassword: hashedPassword,
				},
			)

			if err != nil {
				return err
			}

			outUserAuthData, err = domain.NewUserAuthData(outUser.ID, hashedPassword)

			if err != nil {
				return err
			}

			refreshToken, err := s.refreshTokenRepo.CreateRefreshToken(
				ctx,
				dto.CreateRefreshTokenDTOInput{
					OwnerID:   outUser.ID,
					TokenHash: hashedRefreshToken,
					ExpiresAt: time.Now().Add(time.Duration(s.cfg.Data.Token.Refresh.TTL) * time.Minute),
				},
			)

			if err != nil {
				return err
			}

			outRefreshToken, err = domain.NewRefreshToken(outUser.ID, hashedRefreshToken, refreshToken.ExpiresAt, &s.cfg.Data.Token.Refresh)

			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	if outUser == nil || outUserAuthData == nil || outRefreshToken == nil {
		return nil, ErrAuthStateNil
	}

	out.AccessToken, err = utils.GenerateJWT(outUser.ID, time.Duration(s.cfg.Data.Token.Access.TTL)*time.Minute, &s.cfg.Secrets)

	if err != nil {
		return nil, err
	}

	err = s.userCacheRepo.SetUser(
		ctx,
		dto.SetUserDTOInput{
			ID:        outUser.ID,
			Username:  &outUser.Username,
			CreatedAt: &outUser.CreatedAt,
		},
	)

	out.CacheErrors.User = err

	err = s.userAuthDataCacheRepo.SetUserAuthData(
		ctx,
		dto.SetUserAuthDataByUserIDDTOInput{
			ID:             &outUserAuthData.ID,
			UserID:         outUserAuthData.UserID,
			HashedPassword: &outUserAuthData.HashedPassword,
		},
	)

	out.CacheErrors.UserAuthData = err

	err = s.refreshTokenCacheRepo.SetRefreshToken(
		ctx,
		dto.SetRefreshTokenDTOInput{
			ID:        outRefreshToken.ID,
			OwnerID:   outRefreshToken.OwnerID,
			TokenHash: outRefreshToken.TokenHash,
			ExpiresAt: outRefreshToken.ExpiresAt,
		},
	)

	out.CacheErrors.RefreshToken = err

	return out, nil
}

func (s *AuthService) Refresh(ctx context.Context, data dto.RefreshDTOInput) (*dto.RefreshDTOOutput, error)

func (s *AuthService) Logout(ctx context.Context, data dto.LogoutDTOInput) error
