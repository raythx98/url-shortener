package service

import (
	"context"
	"time"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/repositories"
	"github.com/raythx98/url-shortener/sqlc/db"
	"github.com/raythx98/url-shortener/tools/pghelper"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/tool/logger"
)

type IRedirects interface {
	Redirect(ctx context.Context, shortLink string, req dto.RedirectRequest) (dto.RedirectResponse, error)
}

type Redirects struct {
	Repo repositories.IRepository
	Log  logger.ILogger
}

func NewRedirects(repo repositories.IRepository, log logger.ILogger) *Redirects {
	return &Redirects{
		Repo: repo,
		Log:  log,
	}
}

func (s *Redirects) Redirect(ctx context.Context, shortLink string, req dto.RedirectRequest) (dto.RedirectResponse, error) {
	url, err := s.Repo.GetUrlByShortUrl(ctx, shortLink)
	if err != nil {
		return dto.RedirectResponse{}, err
	}
	if url == nil {
		return dto.RedirectResponse{}, errorhelper.NewAppError(4, "Invalid short url, please create a new one", err)
	}

	go func(ctx context.Context, cancel context.CancelFunc) {
		defer cancel()
		err = s.Repo.CreateRedirect(ctx, db.CreateRedirectParams{
			UrlID:   pghelper.Int8(&url.ID),
			Device:  req.Device,
			Country: req.Country,
			City:    req.City,
		})
		if err != nil {
			s.Log.Error(ctx, "create redirect", logger.WithError(err))
		}
	}(context.WithTimeout(context.Background(), 30*time.Second))

	return dto.RedirectResponse{
		FullUrl: url.FullUrl,
	}, nil
}
