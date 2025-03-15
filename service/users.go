package service

import (
	"context"
	"fmt"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/repositories"
	"github.com/raythx98/url-shortener/tools/crypto"
	"github.com/raythx98/url-shortener/tools/reqctx"

	"github.com/raythx98/gohelpme/tool/logger"
)

type IUsers interface {
	GetProfile(ctx context.Context) (dto.ProfileResponse, error)
}

type Users struct {
	Repo   repositories.IRepository
	Log    logger.ILogger
	Crypto crypto.ICrypto
	ReqCtx reqctx.IReqCtx
}

func NewUsers(repo repositories.IRepository, log logger.ILogger, crypto crypto.ICrypto, reqCtx reqctx.IReqCtx) *Users {
	return &Users{
		Repo:   repo,
		Log:    log,
		Crypto: crypto,
		ReqCtx: reqCtx,
	}
}

func (s *Users) GetProfile(ctx context.Context) (dto.ProfileResponse, error) {
	reqCtx := s.ReqCtx.GetValue(ctx)
	if reqCtx.UserId == nil {
		return dto.ProfileResponse{}, fmt.Errorf("user id not found")
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
