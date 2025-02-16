package service

import (
	"context"
	"fmt"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/sqlc/db"

	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

type IUsers interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	GetProfile(ctx context.Context) (dto.ProfileResponse, error)
}

type Users struct {
	Repo *db.Queries
	Log  logger.ILogger
}

func NewUsers(repo *db.Queries, log logger.ILogger) *Users {
	return &Users{
		Repo: repo,
		Log:  log,
	}
}

func (s *Users) Register(ctx context.Context, req dto.RegisterRequest) error {
	return s.Repo.CreateUser(ctx, db.CreateUserParams{
		Email:    req.Email,
		Password: req.Password,
	})
}

func (s *Users) GetProfile(ctx context.Context) (dto.ProfileResponse, error) {
	reqCtx := reqctx.GetValue(ctx)
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
