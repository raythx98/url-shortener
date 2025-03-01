package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/sqlc/db"
	"github.com/raythx98/url-shortener/tools/crypto"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/tool/jwthelper"
	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"

	"github.com/jackc/pgx/v5"
)

type IAuth interface {
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
	Refresh(ctx context.Context) (dto.LoginResponse, error)
}

type Auth struct {
	Repo   *db.Queries
	Jwt    jwthelper.IJwt
	Log    logger.ILogger
	Crypto crypto.ICrypto
}

func NewAuth(repo *db.Queries, log logger.ILogger, jwt jwthelper.IJwt, crypto crypto.ICrypto) *Auth {
	return &Auth{
		Repo:   repo,
		Log:    log,
		Jwt:    jwt,
		Crypto: crypto,
	}
}

func (s *Auth) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.Repo.GetUserByEmail(ctx, request.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.LoginResponse{}, errorhelper.NewAppError(3, "Email is not registered", err)
	}
	if err != nil {
		return dto.LoginResponse{}, err
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
	reqCtx := reqctx.GetValue(ctx)

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
