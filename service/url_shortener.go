package service

import (
	"github.com/raythx98/url-shortener/sqlc/db"
)

type IUrlShortener interface {
	//ShortenUrl(ctx context.Context, request dto.ShortenUrlRequest) (*dto.ShortenUrlResponse, error)
	//GetUrlWithShortened(ctx context.Context, shortenedUrl string) (string, error)
	GetDb() *db.Queries
}

type UrlShortener struct {
	UrlMappingRepo *db.Queries
}

func New(urlMappingRepo *db.Queries) *UrlShortener {
	return &UrlShortener{
		UrlMappingRepo: urlMappingRepo,
	}
}

func (s *UrlShortener) GetDb() *db.Queries {
	return s.UrlMappingRepo
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
//	mapping, err := s.UrlMappingRepo.CreateUrlMapping(ctx, params)
//	if err != nil {
//		return nil, err
//	}
//
//	return (&dto.ShortenUrlResponse{}).Bind(mapping), nil
//}
//
//func (s *UrlShortener) GetUrlWithShortened(ctx context.Context, shortenedUrl string) (string, error) {
//	mapping, err := s.UrlMappingRepo.GetUrlMapping(ctx, shortenedUrl)
//	if err != nil {
//		return "", err
//	}
//
//	return mapping.Url, nil
//}
