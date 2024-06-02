package service

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/raythx98/url-shortener/sqlc/url_mappings"
	"github.com/raythx98/url-shortener/tools/sql_tool"
	"time"

	"github.com/google/uuid"
)

type IUrlShortener interface {
	ShortenUrl(ctx context.Context, url string) (string, error)
	GetUrlWithShortened(ctx context.Context, shortenedUrl string) (string, error)
}

type UrlShortener struct {
}

func (s *UrlShortener) ShortenUrl(ctx context.Context, url string) (string, error) {
	conn, err := pgx.Connect(ctx, "user=postgres password=password host=localhost port=5432 dbname=url_shortener sslmode=disable") // TODO: connect on startup
	if err != nil {
		return "", err
	}
	defer conn.Close(ctx)

	queries := url_mappings.New(conn)

	expireAt := time.Now().UTC().AddDate(1, 0, 0)
	params := url_mappings.CreateUrlMappingParams{
		ShortenedUrl:     uuid.NewString(),
		Url:              url,
		InactiveExpireAt: sql_tool.NewTime(nil),
		MustExpireAt:     sql_tool.NewTime(&expireAt),
	}
	mapping, err := queries.CreateUrlMapping(ctx, params)
	if err != nil {
		return "", err
	}

	return mapping.ShortenedUrl, nil
}

func (s *UrlShortener) GetUrlWithShortened(ctx context.Context, shortenedUrl string) (string, error) {
	conn, err := pgx.Connect(ctx, "user=postgres password=password host=localhost port=5432 dbname=url_shortener sslmode=disable") // TODO: connect on startup
	if err != nil {
		return "", err
	}
	defer conn.Close(ctx)

	queries := url_mappings.New(conn)

	mapping, err := queries.GetUrlMapping(ctx, shortenedUrl)
	if err != nil {
		return "", err
	}

	return mapping.Url, nil
}
