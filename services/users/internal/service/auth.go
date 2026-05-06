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

func (s *AuthService) Login(ctx context.Context, data dto.LoginDTOInput) (*dto.LoginDTOOutput, error) {
	if s.cfg.Service.Auth.MinPasswordLength > utf8.RuneCountInString(data.Password) {
		return nil, ErrPasswordTooShort
	}

	var out *dto.LoginDTOOutput
	var err error

	var user *dto.GetUserDTOOutput

	out.UserID, err = s.userCacheRepo.GetUserIDByUsername(ctx, data.Username)

	if err == nil {
		user, err = s.userCacheRepo.GetUserByID(ctx, out.UserID)
	}

	if err != nil {
		user, err = s.userRepo.GetUserByUsername(ctx, data.Username)
		if err != nil {
			return nil, loginError(data.Username, err)
		}
	}

	if user == nil {
		return nil, loginError(data.Username, ErrUserStateNil)
	}

	out.UserID = user.ID

	userAuthData, err := s.userAuthDataCacheRepo.GetUserAuthDataByUserID(ctx, out.UserID)

	if err != nil {
		userAuthData, err = s.userAuthDataRepo.GetUserAuthDataByUserID(ctx, out.UserID)
		if err != nil {
			return nil, loginError(data.Username, err)
		}
	}

	if userAuthData == nil {
		return nil, loginError(data.Username, ErrAuthDataStateNil)
	}

	hashedPassword, err := utils.HashPassword(data.Password)

	if hashedPassword != userAuthData.HashedPassword {
		return nil, loginError(data.Username, ErrPasswordIncorrect)
	}

	out.RefreshToken, err = utils.GenerateRefreshToken()

	if err != nil {
		return nil, loginError(data.Username, err)
	}

	hashedRefreshToken := utils.HashRefreshToken(out.RefreshToken)

	var outRefreshToken *domain.RefreshToken

	err = s.tx.WithinTx(
		ctx,
		func(ctx context.Context) error {
			refreshToken, err := s.refreshTokenRepo.CreateRefreshToken(
				ctx,
				dto.CreateRefreshTokenDTOInput{
					OwnerID:   out.UserID,
					TokenHash: hashedRefreshToken,
					ExpiresAt: time.Now().Add(time.Duration(s.cfg.Data.Token.Refresh.TTL) * time.Minute),
				},
			)

			if err != nil {
				return err
			}

			outRefreshToken, err = domain.NewRefreshToken(out.UserID, hashedRefreshToken, refreshToken.ExpiresAt, &s.cfg.Data.Token.Refresh)

			return err
		},
	)

	if err != nil {
		return nil, loginError(data.Username, err)
	}

	if outRefreshToken == nil {
		return nil, loginError(data.Username, ErrRefreshTokenStatenil)
	}

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

	out.AccessToken, err = utils.GenerateJWT(out.UserID, time.Duration(s.cfg.Data.Token.Access.TTL)*time.Minute, &s.cfg.Secrets)

	if err != nil {
		return nil, loginError(data.Username, err)
	}

	return out, nil
}

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
		return nil, registerError(data.Username, err)
	}

	out.RefreshToken, err = utils.GenerateRefreshToken()

	if err != nil {
		return nil, registerError(data.Username, err)
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
		return nil, registerError(data.Username, err)
	}

	if outUser == nil {
		return nil, registerError(data.Username, ErrUserStateNil)
	}

	if outUserAuthData == nil {
		return nil, registerError(data.Username, ErrAuthDataStateNil)
	}

	if outRefreshToken == nil {
		return nil, registerError(data.Username, ErrRefreshTokenStatenil)
	}

	out.AccessToken, err = utils.GenerateJWT(outUser.ID, time.Duration(s.cfg.Data.Token.Access.TTL)*time.Minute, &s.cfg.Secrets)

	if err != nil {
		return nil, registerError(data.Username, err)
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
