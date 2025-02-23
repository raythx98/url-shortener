package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/raythx98/gohelpme/errorhelper"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/sqlc/db"

	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"

	"github.com/jackc/pgx/v5/pgtype"
)

type IUrls interface {
	GetUrl(ctx context.Context, urlId string) (dto.GetUrlResponse, error)
	GetUrls(ctx context.Context) (dto.GetUrlsResponse, error)
	CreateUrl(ctx context.Context, req dto.CreateUrlRequest) (dto.CreateUrlResponse, error)
	DeleteUrl(ctx context.Context, urlId string) error
}

type Urls struct {
	Repo *db.Queries
	Log  logger.ILogger
}

func NewUrls(repo *db.Queries, log logger.ILogger) *Urls {
	return &Urls{
		Repo: repo,
		Log:  log,
	}
}

func (s *Urls) GetUrl(ctx context.Context, urlId string) (dto.GetUrlResponse, error) {
	parseInt, err := strconv.ParseInt(urlId, 10, 64)
	if err != nil {
		return dto.GetUrlResponse{}, err
	}

	getUrl, err := s.Repo.GetUrl(ctx, parseInt)
	if err != nil {
		return dto.GetUrlResponse{}, err
	}

	redirects, err := s.Repo.GetRedirectsByUrlId(ctx, pgtype.Int8{Int64: parseInt, Valid: true})
	if err != nil {
		return dto.GetUrlResponse{}, err
	}

	deviceMap := make(map[string]int)
	countryMap := make(map[string]int)
	for _, redirect := range redirects {
		if redirect.Device != "" {
			deviceMap[redirect.Device]++
		}
		if redirect.Country != "" {
			countryMap[redirect.Country]++
		}
	}

	unsortedDevices := make([]dto.Device, 0)
	unsortedCountries := make([]dto.Country, 0)
	for key, value := range deviceMap {
		unsortedDevices = append(unsortedDevices, dto.Device{Device: key, Count: value})
	}

	for key, value := range countryMap {
		unsortedCountries = append(unsortedCountries, dto.Country{Country: key, Count: value})
	}

	slices.SortFunc(unsortedDevices, func(a, b dto.Device) int {
		ac, ab := a.Count, b.Count
		if ac > ab {
			return -1
		}

		if ac == ab {
			return 0
		}

		return 1
	})

	slices.SortFunc(unsortedCountries, func(a, b dto.Country) int {
		ac, ab := a.Count, b.Count
		if ac > ab {
			return -1
		}

		if ac == ab {
			return 0
		}

		return 1
	})

	return dto.GetUrlResponse{
		Url: dto.Url{
			Id:        getUrl.ID,
			Title:     getUrl.Title,
			ShortUrl:  getUrl.ShortUrl,
			FullUrl:   getUrl.FullUrl,
			Qr:        getUrl.Qr,
			CreatedAt: getUrl.CreatedAt.Time,
		},
		Devices:     unsortedDevices[:min(5, len(unsortedDevices))],
		Countries:   unsortedCountries[:min(5, len(unsortedDevices))],
		TotalClicks: len(redirects),
	}, nil
}

func (s *Urls) GetUrls(ctx context.Context) (dto.GetUrlsResponse, error) {
	reqCtx := reqctx.GetValue(ctx)
	if reqCtx.UserId == nil {
		return dto.GetUrlsResponse{}, fmt.Errorf("user id not found")
	}

	urls, err := s.Repo.GetUrlsByUserId(ctx, pgtype.Int8{Int64: *reqCtx.UserId, Valid: true})
	if err != nil {
		return dto.GetUrlsResponse{}, err
	}

	totalClicks, err := s.Repo.GetUserTotalClicks(ctx, pgtype.Int8{Int64: *reqCtx.UserId, Valid: true})

	response := dto.GetUrlsResponse{
		Urls:        []dto.Url{},
		TotalClicks: totalClicks,
	}

	for _, each := range urls {
		response.Urls = append(response.Urls, dto.Url{
			Id:        each.ID,
			Title:     each.Title,
			ShortUrl:  each.ShortUrl,
			FullUrl:   each.FullUrl,
			Qr:        each.Qr,
			CreatedAt: each.CreatedAt.Time,
		})
	}
	return response, nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func generateAlphaNumeric(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *Urls) CreateUrl(ctx context.Context, req dto.CreateUrlRequest) (dto.CreateUrlResponse, error) {
	reqCtx := reqctx.GetValue(ctx)

	createUrlParams := db.CreateUrlParams{
		Title:   req.Title,
		FullUrl: req.FullUrl,
		Qr:      req.Qr,
	}
	if reqCtx.UserId != nil {
		createUrlParams.UserID = pgtype.Int8{Int64: *reqCtx.UserId, Valid: true}
	}

	if req.CustomUrl != "" {
		createUrlParams.ShortUrl = req.CustomUrl
	} else {
		createUrlParams.ShortUrl = generateAlphaNumeric(8)
	}

	_, err := s.Repo.GetUrlByShortUrl(ctx, createUrlParams.ShortUrl)
	if err == nil {
		return dto.CreateUrlResponse{}, &errorhelper.AppError{
			Code:    5,
			Message: "Url already taken",
		}
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return dto.CreateUrlResponse{}, err
	}

	if strings.EqualFold(createUrlParams.ShortUrl, "api") {
		return dto.CreateUrlResponse{}, fmt.Errorf("short url cannot be 'api'")
	}

	createdUrl, err := s.Repo.CreateUrl(ctx, createUrlParams)
	if err != nil {
		return dto.CreateUrlResponse{}, err
	}

	return dto.CreateUrlResponse{
		Id:       createdUrl.ID,
		ShortUrl: createdUrl.ShortUrl,
	}, nil
}

func (s *Urls) DeleteUrl(ctx context.Context, urlId string) error {
	reqCtx := reqctx.GetValue(ctx)
	if reqCtx.UserId == nil {
		return fmt.Errorf("user id not found")
	}

	parseInt, err := strconv.ParseInt(urlId, 10, 64)
	if err != nil {
		return err
	}

	return s.Repo.DeleteUrl(ctx, parseInt)
}
