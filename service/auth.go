package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/raythx98/gohelpme/errorhelper"
	"strconv"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/sqlc/db"

	"github.com/raythx98/gohelpme/tool/jwt"
	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

type IAuth interface {
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
	Refresh(ctx context.Context) (dto.LoginResponse, error)
}

type Auth struct {
	Repo *db.Queries
	Log  logger.ILogger
}

func NewAuth(repo *db.Queries, log logger.ILogger) *Auth {
	return &Auth{
		Repo: repo,
		Log:  log,
	}
}

func (s *Auth) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.Repo.GetUserByEmail(ctx, request.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return dto.LoginResponse{}, &errorhelper.AppError{
			Code:    3,
			Message: "Email is not registered",
		}
	}
	if err != nil {
		return dto.LoginResponse{}, err
	}

	if user.Password != request.Password {
		return dto.LoginResponse{}, &errorhelper.AppError{
			Code:    2,
			Message: "Incorrect Password",
		}
	}

	accessToken, err := jwt.NewAccessToken(strconv.FormatInt(user.ID, 10))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	refreshToken, err := jwt.NewRefreshToken(strconv.FormatInt(user.ID, 10))
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

	accessToken, err := jwt.NewAccessToken(strconv.FormatInt(*reqCtx.UserId, 10))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	refreshToken, err := jwt.NewRefreshToken(strconv.FormatInt(*reqCtx.UserId, 10))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
