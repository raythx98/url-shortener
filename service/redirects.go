package service

import (
	"context"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/sqlc/db"
	
	"github.com/raythx98/gohelpme/tool/logger"

	"github.com/jackc/pgx/v5/pgtype"
)

type IRedirects interface {
	Redirect(ctx context.Context, shortLink string, req dto.RedirectRequest) (dto.RedirectResponse, error)
}

type Redirects struct {
	Repo *db.Queries
	Log  logger.ILogger
}

func NewRedirects(repo *db.Queries, log logger.ILogger) *Redirects {
	return &Redirects{
		Repo: repo,
		Log:  log,
	}
}

func (s *Redirects) Redirect(ctx context.Context, shortLink string, req dto.RedirectRequest) (dto.RedirectResponse, error) {
	getUrl, err := s.Repo.GetUrlByShortUrl(ctx, shortLink)
	if err != nil {
		return dto.RedirectResponse{}, err
	}

	err = s.Repo.CreateRedirect(ctx, db.CreateRedirectParams{
		UrlID:   pgtype.Int8{Int64: getUrl.ID, Valid: true},
		Device:  req.Device,
		Country: req.Country,
		City:    req.City,
	})
	if err != nil {
		s.Log.Error(ctx, "create redirect", logger.WithError(err))
	}

	return dto.RedirectResponse{
		FullUrl: getUrl.FullUrl,
	}, nil
}

//func (s *UrlShortener) ShortenUrl(ctx context.Context, req dto.ShortenUrlRequest) (*dto.ShortenUrlResponse, error) {
//	expireAt := time.Now().UTC().AddDate(1, 0, 0)
//
//	params := req.BindTo(db.CreateUrlMappingParams{
//		ShortenedUrl:     uuid.NewString()[:8],
//		InactiveExpireAt: pg_helper.NewTime(nil),
//		MustExpireAt:     pg_helper.NewTime(&expireAt),
//	})
//
//	mapping, err := s.Repo.CreateUrlMapping(ctx, params)
//	if err != nil {
//		return nil, err
//	}
//
//	return (&dto.ShortenUrlResponse{}).Bind(mapping), nil
//}
//
//func (s *UrlShortener) GetUrlWithShortened(ctx context.Context, shortenedUrl string) (string, error) {
//	mapping, err := s.Repo.GetUrlMapping(ctx, shortenedUrl)
//	if err != nil {
//		return "", err
//	}
//
//	return mapping.Url, nil
//}
