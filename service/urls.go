package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/raythx98/url-shortener/tools/supabase"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/sqlc/db"
	"github.com/raythx98/url-shortener/tools/random"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	qrcode "github.com/skip2/go-qrcode"
)

type IUrls interface {
	GetUrl(ctx context.Context, urlId string) (dto.GetUrlResponse, error)
	GetUrls(ctx context.Context) (dto.GetUrlsResponse, error)
	CreateUrl(ctx context.Context, req dto.CreateUrlRequest, origin string) (dto.CreateUrlResponse, error)
	DeleteUrl(ctx context.Context, urlId string) error
}

type Urls struct {
	Repo   *db.Queries
	Log    logger.ILogger
	Random random.IRandom
}

func NewUrls(repo *db.Queries, log logger.ILogger, random random.IRandom) *Urls {
	return &Urls{
		Repo:   repo,
		Log:    log,
		Random: random,
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

func (s *Urls) CreateUrl(ctx context.Context, req dto.CreateUrlRequest, origin string) (dto.CreateUrlResponse, error) {
	reqCtx := reqctx.GetValue(ctx)

	createUrlParams := db.CreateUrlParams{
		Title:   req.Title,
		FullUrl: req.FullUrl,
	}
	if reqCtx.UserId != nil {
		createUrlParams.UserID = pgtype.Int8{Int64: *reqCtx.UserId, Valid: true}
	}

	if req.CustomUrl != "" {
		createUrlParams.ShortUrl = req.CustomUrl
	} else {
		createUrlParams.ShortUrl = s.Random.GenerateAlphaNum(8)
	}

	_, err := s.Repo.GetUrlByShortUrl(ctx, createUrlParams.ShortUrl)
	if err == nil {
		return dto.CreateUrlResponse{}, errorhelper.NewAppError(5, "Short url already taken", err)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return dto.CreateUrlResponse{}, err
	}

	if strings.EqualFold(createUrlParams.ShortUrl, "api") {
		return dto.CreateUrlResponse{}, fmt.Errorf("short url cannot be 'api'")
	}

	//upload file
	// get origin without 'http://' or 'https://' or `.` or `/` or `:` or `?`
	fileNamePrefix := strings.ReplaceAll(origin, "http://", "")
	fileNamePrefix = strings.ReplaceAll(fileNamePrefix, "https://", "")
	fileNamePrefix = strings.ReplaceAll(fileNamePrefix, ".", "")
	fileNamePrefix = strings.ReplaceAll(fileNamePrefix, "/", "")
	fileNamePrefix = strings.ReplaceAll(fileNamePrefix, ":", "")
	fileNamePrefix = strings.ReplaceAll(fileNamePrefix, "?", "")

	f, err := os.CreateTemp("", "qr-*.png")
	if err != nil {
		return dto.CreateUrlResponse{}, err
	}
	defer func(name string) {
		_ = os.Remove(name)
	}(f.Name()) // clean up

	fileName := fmt.Sprintf("%s-%s", fileNamePrefix, createUrlParams.ShortUrl)
	fmt.Println("filename: " + fileName)

	err = qrcode.WriteFile(fmt.Sprintf("%s/%s", origin, createUrlParams.ShortUrl),
		qrcode.Medium, 256, f.Name())
	if err != nil {
		return dto.CreateUrlResponse{}, err
	}

	if err := supabase.UpdateFile(fileName, f); err != nil {
		return dto.CreateUrlResponse{}, err
	}

	createUrlParams.Qr = "https://caqzitwuslrszkfwbmve.supabase.co/storage/v1/object/public/qrs//" + fileName

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
