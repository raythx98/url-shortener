package service

import (
	"context"
	"fmt"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/repositories"

	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

type IUsers interface {
	GetProfile(ctx context.Context) (dto.ProfileResponse, error)
}

type Users struct {
	Repo repositories.IRepository
	Log  logger.ILogger
}

func NewUsers(repo repositories.IRepository, log logger.ILogger) *Users {
	return &Users{
		Repo: repo,
		Log:  log,
	}
}

func (s *Users) GetProfile(ctx context.Context) (dto.ProfileResponse, error) {
	reqCtx := reqctx.GetValue(ctx)
	if reqCtx == nil || reqCtx.UserId == nil {
		return dto.ProfileResponse{}, fmt.Errorf("user id not found, reqCtx: %+v", reqCtx)
	}

	user, err := s.Repo.GetUser(ctx, *reqCtx.UserId)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	return dto.ProfileResponse{
		Id:    user.ID,
		Email: user.Email,
		Role:  "authenticated",
	}, nil
}
