package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/sqlc/url_mappings"
	"github.com/raythx98/url-shortener/tools/sql_tool"
	"time"
)

type IUrlShortener interface {
	ShortenUrl(ctx context.Context, request dto.ShortenUrlRequest) (*dto.ShortenUrlResponse, error)
	GetUrlWithShortened(ctx context.Context, shortenedUrl string) (string, error)
}

type UrlShortener struct {
	UrlMappingRepo *url_mappings.Queries
}

func New(urlMappingRepo *url_mappings.Queries) *UrlShortener {
	return &UrlShortener{
		UrlMappingRepo: urlMappingRepo,
	}
}

func (s *UrlShortener) ShortenUrl(ctx context.Context, req dto.ShortenUrlRequest) (*dto.ShortenUrlResponse, error) {
	expireAt := time.Now().UTC().AddDate(1, 0, 0)

	params := req.BindTo(url_mappings.CreateUrlMappingParams{
		ShortenedUrl:     uuid.NewString()[:8],
		InactiveExpireAt: sql_tool.NewTime(nil),
		MustExpireAt:     sql_tool.NewTime(&expireAt),
	})

	mapping, err := s.UrlMappingRepo.CreateUrlMapping(ctx, params)
	if err != nil {
		return nil, err
	}

	return (&dto.ShortenUrlResponse{}).Bind(mapping), nil
}

func (s *UrlShortener) GetUrlWithShortened(ctx context.Context, shortenedUrl string) (string, error) {
	mapping, err := s.UrlMappingRepo.GetUrlMapping(ctx, shortenedUrl)
	if err != nil {
		return "", err
	}

	return mapping.Url, nil
}
