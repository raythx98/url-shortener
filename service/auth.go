package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/repositories"
	"github.com/raythx98/url-shortener/sqlc/db"
	"github.com/raythx98/url-shortener/tools/crypto"
	"github.com/raythx98/url-shortener/tools/reqctx"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/tool/jwthelper"
	"github.com/raythx98/gohelpme/tool/logger"
)

type IAuth interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.LoginResponse, error)
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
	Refresh(ctx context.Context) (dto.LoginResponse, error)
}

type Auth struct {
	Repo   repositories.IRepository
	Jwt    jwthelper.IJwt
	Log    logger.ILogger
	Crypto crypto.ICrypto
	ReqCtx reqctx.IReqCtx
}

func NewAuth(repo repositories.IRepository, log logger.ILogger, jwt jwthelper.IJwt, crypto crypto.ICrypto,
	reqCtx reqctx.IReqCtx) *Auth {
	return &Auth{
		Repo:   repo,
		Log:    log,
		Jwt:    jwt,
		Crypto: crypto,
		ReqCtx: reqCtx,
	}
}

func (s *Auth) Register(ctx context.Context, req dto.RegisterRequest) (dto.LoginResponse, error) {
	existingUser, err := s.Repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if existingUser != nil {
		return dto.LoginResponse{}, errorhelper.NewAppError(1, "Email has already been registered", err)
	}

	encodedHashedPassword, err := s.Crypto.GenerateFromPassword(req.Password)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	user, err := s.Repo.CreateUser(ctx, db.CreateUserParams{
		Email:    req.Email,
		Password: encodedHashedPassword,
	})
	if err != nil {
		return dto.LoginResponse{}, err
	}

	accessToken, err := s.Jwt.NewAccessToken(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	refreshToken, err := s.Jwt.NewRefreshToken(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *Auth) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.Repo.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if user == nil {
		return dto.LoginResponse{}, errorhelper.NewAppError(3, "Email is not registered", err)
	}

	authenticated, err := s.Crypto.ComparePasswordAndHash(request.Password, user.Password)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	if !authenticated {
		return dto.LoginResponse{}, errorhelper.NewAppError(2, "Incorrect Password", err)
	}

	accessToken, err := s.Jwt.NewAccessToken(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	refreshToken, err := s.Jwt.NewRefreshToken(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *Auth) Refresh(ctx context.Context) (dto.LoginResponse, error) {
	reqCtx := s.ReqCtx.GetValue(ctx)

	if reqCtx.UserId == nil {
		return dto.LoginResponse{}, fmt.Errorf("user id not found")
	}

	accessToken, err := s.Jwt.NewAccessToken(strconv.FormatInt(*reqCtx.UserId, 10))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	refreshToken, err := s.Jwt.NewRefreshToken(strconv.FormatInt(*reqCtx.UserId, 10))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
