package service

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/sqlc/url_mappings"
	"github.com/raythx98/url-shortener/tools/sql_tool"
	"time"

	"github.com/google/uuid"
)

type IUrlShortener interface {
	ShortenUrl(ctx context.Context, url string) (*dto.ShortenUrlResponse, error)
	GetUrlWithShortened(ctx context.Context, shortenedUrl string) (string, error)
}

type UrlShortener struct {
	DbPool *pgxpool.Pool
}

func (s *UrlShortener) ShortenUrl(ctx context.Context, url string) (*dto.ShortenUrlResponse, error) {
	expireAt := time.Now().UTC().AddDate(1, 0, 0)
	params := url_mappings.CreateUrlMappingParams{
		ShortenedUrl:     uuid.NewString()[:8],
		Url:              url,
		InactiveExpireAt: sql_tool.NewTime(nil),
		MustExpireAt:     sql_tool.NewTime(&expireAt),
	}
	mapping, err := url_mappings.New(s.DbPool).CreateUrlMapping(ctx, params)
	if err != nil {
		return nil, err
	}

	var inactiveExpireAt, mustExpireAt *time.Time
	if mapping.InactiveExpireAt.Valid {
		inactiveExpireAt = &mapping.InactiveExpireAt.Time
	}
	if mapping.MustExpireAt.Valid {
		mustExpireAt = &mapping.MustExpireAt.Time
	}

	return &dto.ShortenUrlResponse{
		Url:              mapping.Url,
		ShortenedUrl:     mapping.ShortenedUrl,
		InactiveExpireAt: inactiveExpireAt,
		MustExpireAt:     mustExpireAt,
	}, nil
}

func (s *UrlShortener) GetUrlWithShortened(ctx context.Context, shortenedUrl string) (string, error) {
	mapping, err := url_mappings.New(s.DbPool).GetUrlMapping(ctx, shortenedUrl)
	if err != nil {
		return "", err
	}
	return mapping.Url, nil
}
