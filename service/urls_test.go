package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/mocks/github.com/raythx98/url-shortener/repositories"
	"github.com/raythx98/url-shortener/mocks/github.com/raythx98/url-shortener/service"
	"github.com/raythx98/url-shortener/mocks/github.com/raythx98/url-shortener/tools/aws"
	"github.com/raythx98/url-shortener/mocks/github.com/raythx98/url-shortener/tools/qrcode"
	"github.com/raythx98/url-shortener/mocks/github.com/raythx98/url-shortener/tools/random"
	"github.com/raythx98/url-shortener/sqlc/db"
	"github.com/raythx98/url-shortener/tools/pghelper"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/mocks/github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/inthelper"
	"github.com/raythx98/gohelpme/tool/reqctx"
	"github.com/raythx98/gohelpme/tool/timehelper"
)

func TestGetUrl(t *testing.T) {
	type fields struct {
		cfg  *service.MockConfigProvider
		repo *repositories.MockIRepository
		s3   *aws.MockIS3
		log  *logger.MockILogger
		rand *random.MockIRandom
		qr   *qrcode.MockIQrCode
	}
	generateFields := func() fields {
		return fields{
			cfg:  service.NewMockConfigProvider(t),
			repo: repositories.NewMockIRepository(t),
			s3:   aws.NewMockIS3(t),
			log:  logger.NewMockILogger(t),
			rand: random.NewMockIRandom(t),
			qr:   qrcode.NewMockIQrCode(t),
		}
	}

	type args struct {
		ctx   context.Context
		urlId string
	}

	tests := []struct {
		name     string
		args     args
		fields   fields
		mocks    func(args *args, fields *fields)
		wantResp dto.GetUrlResponse
		wantErr  error
	}{
		{
			name: "success, device & country sorted by count descending",
			args: args{
				ctx:   context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				urlId: "1",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrl", args.ctx, int64(1)).Return(
					db.Url{
						ID: 1, Title: "title", ShortUrl: "short.url", FullUrl: "full.url", Qr: "qrcode",
						CreatedAt: pghelper.Time(timehelper.TimePtr(time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC))),
					}, nil)

				fields.repo.On("GetRedirectsByUrlId", args.ctx, inthelper.Int64Ptr(1)).Return(
					[]db.Redirect{
						{Device: "mobile", Country: "malaysia"},
						{Device: "desktop", Country: "malaysia"},
						{Device: "pc", Country: "singapore"},
						{Device: "desktop", Country: "singapore"},
						{Device: "mobile", Country: "australia"},
						{Device: "mobile", Country: "singapore"},
					}, nil)
			},
			wantResp: dto.GetUrlResponse{
				Url: dto.Url{
					Id: 1, Title: "title", ShortUrl: "short.url", FullUrl: "full.url", Qr: "qrcode",
					CreatedAt: time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC),
				},
				TotalClicks: 6,
				Devices:     []dto.Device{{Device: "mobile", Count: 3}, {Device: "desktop", Count: 2}, {Device: "pc", Count: 1}},
				Countries:   []dto.Country{{Country: "singapore", Count: 3}, {Country: "malaysia", Count: 2}, {Country: "australia", Count: 1}},
			},
			wantErr: nil,
		},
		{
			name: "success, truncate to top 5 device & country",
			args: args{
				ctx:   context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				urlId: "1",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrl", args.ctx, int64(1)).Return(
					db.Url{
						ID: 1, Title: "title", ShortUrl: "short.url", FullUrl: "full.url", Qr: "qrcode",
						CreatedAt: pghelper.Time(timehelper.TimePtr(time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC))),
					}, nil)

				fields.repo.On("GetRedirectsByUrlId", args.ctx, inthelper.Int64Ptr(1)).Return(
					[]db.Redirect{
						{Device: "fifth", Country: "fifth"}, {Device: "fifth", Country: "fifth"},

						{Device: "truncated", Country: "truncated"},

						{Device: "second", Country: "second"}, {Device: "second", Country: "second"},
						{Device: "second", Country: "second"}, {Device: "second", Country: "second"},
						{Device: "second", Country: "second"},

						{Device: "first", Country: "first"}, {Device: "first", Country: "first"},
						{Device: "first", Country: "first"}, {Device: "first", Country: "first"},
						{Device: "first", Country: "first"}, {Device: "first", Country: "first"},

						{Device: "fourth", Country: "fourth"}, {Device: "fourth", Country: "fourth"},
						{Device: "fourth", Country: "fourth"},

						{Device: "third", Country: "third"}, {Device: "third", Country: "third"},
						{Device: "third", Country: "third"}, {Device: "third", Country: "third"},

						{Device: "truncated too", Country: "truncated too"},
					}, nil)
			},
			wantResp: dto.GetUrlResponse{
				Url: dto.Url{
					Id: 1, Title: "title", ShortUrl: "short.url", FullUrl: "full.url", Qr: "qrcode",
					CreatedAt: time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC),
				},
				TotalClicks: 22,
				Devices: []dto.Device{{Device: "first", Count: 6}, {Device: "second", Count: 5},
					{Device: "third", Count: 4}, {Device: "fourth", Count: 3}, {Device: "fifth", Count: 2}},
				Countries: []dto.Country{{Country: "first", Count: 6}, {Country: "second", Count: 5},
					{Country: "third", Count: 4}, {Country: "fourth", Count: 3}, {Country: "fifth", Count: 2}},
			},
			wantErr: nil,
		},
		{
			name: "no redirects",
			args: args{
				ctx:   context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				urlId: "1",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrl", args.ctx, int64(1)).Return(
					db.Url{
						ID: 1, Title: "title", ShortUrl: "short.url", FullUrl: "full.url", Qr: "qrcode",
						CreatedAt: pghelper.Time(timehelper.TimePtr(time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC))),
					}, nil)

				fields.repo.On("GetRedirectsByUrlId", args.ctx, inthelper.Int64Ptr(1)).Return(
					nil, nil)
			},
			wantResp: dto.GetUrlResponse{
				Url: dto.Url{
					Id: 1, Title: "title", ShortUrl: "short.url", FullUrl: "full.url", Qr: "qrcode",
					CreatedAt: time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC),
				},
				TotalClicks: 0,
				Devices:     []dto.Device{},
				Countries:   []dto.Country{},
			},
			wantErr: nil,
		},
		{
			name: "get redirects from db fail",
			args: args{
				ctx:   context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				urlId: "1",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrl", args.ctx, int64(1)).Return(
					db.Url{
						ID: 1, Title: "title", ShortUrl: "short.url", FullUrl: "full.url", Qr: "qrcode",
						CreatedAt: pghelper.Time(timehelper.TimePtr(time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC))),
					}, nil)

				fields.repo.On("GetRedirectsByUrlId", args.ctx, inthelper.Int64Ptr(1)).Return(
					nil, errors.New("get redirects from db fail"))
			},
			wantResp: dto.GetUrlResponse{},
			wantErr:  errors.New("get redirects from db fail"),
		},
		{
			name: "get url from db fail",
			args: args{
				ctx:   context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				urlId: "1",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrl", args.ctx, int64(1)).Return(db.Url{}, errors.New("get url from db fail"))
			},
			wantResp: dto.GetUrlResponse{},
			wantErr:  errors.New("get url from db fail"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(&tt.args, &tt.fields)

			r := NewUrls(tt.fields.cfg, tt.fields.repo, tt.fields.s3, tt.fields.log, tt.fields.rand, tt.fields.qr)

			got, gotErr := r.GetUrl(tt.args.ctx, tt.args.urlId)

			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("Urls.GetUrl() = %v, want %v", got, tt.wantResp)
			}

			if !errorhelper.IsEqual(gotErr, tt.wantErr) {
				t.Errorf("Urls.GetUrl() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetUrls(t *testing.T) {
	type fields struct {
		cfg  *service.MockConfigProvider
		repo *repositories.MockIRepository
		s3   *aws.MockIS3
		log  *logger.MockILogger
		rand *random.MockIRandom
		qr   *qrcode.MockIQrCode
	}
	generateFields := func() fields {
		return fields{
			cfg:  service.NewMockConfigProvider(t),
			repo: repositories.NewMockIRepository(t),
			s3:   aws.NewMockIS3(t),
			log:  logger.NewMockILogger(t),
			rand: random.NewMockIRandom(t),
			qr:   qrcode.NewMockIQrCode(t),
		}
	}

	type args struct {
		ctx context.Context
	}

	var tests = []struct {
		name     string
		args     args
		fields   fields
		mocks    func(args *args, fields *fields)
		wantResp dto.GetUrlsResponse
		wantErr  error
	}{
		{
			name: "success",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrlsByUserId", args.ctx, inthelper.Int64Ptr(1)).Return(
					[]db.Url{
						{
							ID: 1, Title: "title", ShortUrl: "short.url", FullUrl: "long.url", Qr: "qrcode",
							CreatedAt: pghelper.Time(timehelper.TimePtr(time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC))),
						},
						{
							ID: 2, Title: "title2", ShortUrl: "short2.url", FullUrl: "long2.url", Qr: "qrcode2",
							CreatedAt: pghelper.Time(timehelper.TimePtr(time.Date(2, 2, 3, 4, 5, 6, 7, time.UTC))),
						},
					}, nil)

				fields.repo.On("GetUserTotalClicks", args.ctx, inthelper.Int64Ptr(1)).Return(int64(10), nil)
			},
			wantResp: dto.GetUrlsResponse{
				Urls: []dto.Url{
					{
						Id: 1, Title: "title", ShortUrl: "short.url", FullUrl: "long.url", Qr: "qrcode",
						CreatedAt: time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC),
					},
					{
						Id: 2, Title: "title2", ShortUrl: "short2.url", FullUrl: "long2.url", Qr: "qrcode2",
						CreatedAt: time.Date(2, 2, 3, 4, 5, 6, 7, time.UTC),
					},
				},
				TotalClicks: 10,
			},
			wantErr: nil,
		},
		{
			name: "no urls",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrlsByUserId", args.ctx, inthelper.Int64Ptr(1)).Return(nil, nil)

				fields.repo.On("GetUserTotalClicks", args.ctx, inthelper.Int64Ptr(1)).Return(int64(0), nil)
			},
			wantResp: dto.GetUrlsResponse{
				Urls:        []dto.Url{},
				TotalClicks: 0,
			},
			wantErr: nil,
		},
		{
			name: "get total clicks from db fail",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrlsByUserId", args.ctx, inthelper.Int64Ptr(1)).Return(nil, nil)

				fields.repo.On("GetUserTotalClicks", args.ctx, inthelper.Int64Ptr(1)).
					Return(int64(0), errors.New("get clicks db error"))
			},
			wantResp: dto.GetUrlsResponse{},
			wantErr:  errors.New("get clicks db error"),
		},
		{
			name: "get urls from db fail",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrlsByUserId", args.ctx, inthelper.Int64Ptr(1)).
					Return(nil, errors.New("get urls db error"))
			},
			wantResp: dto.GetUrlsResponse{},
			wantErr:  errors.New("get urls db error"),
		},
		{
			name: "user ID not in context",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("")),
			},
			fields:   generateFields(),
			mocks:    func(args *args, fields *fields) {},
			wantResp: dto.GetUrlsResponse{},
			wantErr:  errors.New("user id cannot be determined from reqctx"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(&tt.args, &tt.fields)

			r := NewUrls(tt.fields.cfg, tt.fields.repo, tt.fields.s3, tt.fields.log, tt.fields.rand, tt.fields.qr)

			got, gotErr := r.GetUrls(tt.args.ctx)

			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("Urls.GetUrls() = %v, want %v", got, tt.wantResp)
			}

			if !errorhelper.IsEqual(gotErr, tt.wantErr) {
				t.Errorf("Urls.GetUrls() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestCreateUrl(t *testing.T) {
	type fields struct {
		cfg  *service.MockConfigProvider
		repo *repositories.MockIRepository
		s3   *aws.MockIS3
		log  *logger.MockILogger
		rand *random.MockIRandom
		qr   *qrcode.MockIQrCode
	}
	generateFields := func() fields {
		return fields{
			cfg:  service.NewMockConfigProvider(t),
			repo: repositories.NewMockIRepository(t),
			s3:   aws.NewMockIS3(t),
			log:  logger.NewMockILogger(t),
			rand: random.NewMockIRandom(t),
			qr:   qrcode.NewMockIQrCode(t),
		}
	}

	type args struct {
		ctx    context.Context
		req    dto.CreateUrlRequest
		origin string
	}

	tests := []struct {
		name     string
		args     args
		fields   fields
		mocks    func(args *args, fields *fields)
		wantResp dto.CreateUrlResponse
		wantErr  error
	}{
		{
			name: "success, logged in user with random short url",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				req: dto.CreateUrlRequest{
					Title:     "title",
					FullUrl:   "full.url",
					CustomUrl: "",
				},
				origin: "localhost",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.rand.On("GenerateAlphaNum", 10).Return("sh0rtieur1")

				fields.repo.On("GetUrlByShortUrl", args.ctx, "sh0rtieur1").Return(nil, nil)

				fields.qr.On("Encode", fmt.Sprintf("%s/%s", args.origin, "sh0rtieur1")).
					Return([]byte("png"), nil)

				fields.cfg.On("GetAwsS3Bucket").Return("bucket")

				fields.s3.On("Upload", args.ctx, "bucket", "sh0rtieur1.png", []byte("png"), "image/png").Return(nil)

				fields.cfg.On("GetAwsRegion").Return("region")

				fields.repo.On("CreateUrl", args.ctx, db.CreateUrlParams{
					UserID:   pghelper.Int8(inthelper.Int64Ptr(1)),
					Title:    "title",
					ShortUrl: "sh0rtieur1",
					FullUrl:  "full.url",
					Qr:       "https://bucket.s3.region.amazonaws.com/sh0rtieur1.png",
				}).Return(db.Url{
					ID:       1,
					ShortUrl: "sh0rtieur1",
					Qr:       "https://bucket.s3.region.amazonaws.com/sh0rtieur1.png",
				}, nil)
			},
			wantResp: dto.CreateUrlResponse{
				Id:       1,
				ShortUrl: "sh0rtieur1",
				Qr:       "https://bucket.s3.region.amazonaws.com/sh0rtieur1.png",
			},
			wantErr: nil,
		},
		{
			name: "success, non logged in user with custom url",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("")),
				req: dto.CreateUrlRequest{
					Title:     "title",
					FullUrl:   "full.url",
					CustomUrl: "custom.url",
				},
				origin: "localhost",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrlByShortUrl", args.ctx, "custom.url").Return(nil, nil)

				fields.qr.On("Encode", fmt.Sprintf("%s/%s", args.origin, "custom.url")).
					Return([]byte("png"), nil)

				fields.cfg.On("GetAwsS3Bucket").Return("bucket")

				fields.s3.On("Upload", args.ctx, "bucket", "custom.url.png", []byte("png"), "image/png").Return(nil)

				fields.cfg.On("GetAwsRegion").Return("region")

				fields.repo.On("CreateUrl", args.ctx, db.CreateUrlParams{
					UserID:   pghelper.Int8(nil),
					Title:    "title",
					ShortUrl: "custom.url",
					FullUrl:  "full.url",
					Qr:       "https://bucket.s3.region.amazonaws.com/custom.url.png",
				}).Return(db.Url{
					ID:       1,
					ShortUrl: "custom.url",
					Qr:       "https://bucket.s3.region.amazonaws.com/custom.url.png",
				}, nil)
			},
			wantResp: dto.CreateUrlResponse{
				Id:       1,
				ShortUrl: "custom.url",
				Qr:       "https://bucket.s3.region.amazonaws.com/custom.url.png",
			},
			wantErr: nil,
		},
		{
			name: "create url db error",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				req: dto.CreateUrlRequest{
					Title:     "title",
					FullUrl:   "full.url",
					CustomUrl: "custom.url",
				},
				origin: "localhost",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrlByShortUrl", args.ctx, "custom.url").Return(nil, nil)

				fields.qr.On("Encode", fmt.Sprintf("%s/%s", args.origin, "custom.url")).
					Return([]byte("png"), nil)

				fields.cfg.On("GetAwsS3Bucket").Return("bucket")

				fields.s3.On("Upload", args.ctx, "bucket", "custom.url.png", []byte("png"), "image/png").Return(nil)

				fields.cfg.On("GetAwsRegion").Return("region")

				fields.repo.On("CreateUrl", args.ctx, db.CreateUrlParams{
					UserID:   pghelper.Int8(inthelper.Int64Ptr(1)),
					Title:    "title",
					ShortUrl: "custom.url",
					FullUrl:  "full.url",
					Qr:       "https://bucket.s3.region.amazonaws.com/custom.url.png",
				}).Return(db.Url{}, errors.New("create url db error"))
			},
			wantResp: dto.CreateUrlResponse{},
			wantErr:  errors.New("create url db error"),
		},
		{
			name: "upload to s3 error",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				req: dto.CreateUrlRequest{
					Title:     "title",
					FullUrl:   "full.url",
					CustomUrl: "custom.url",
				},
				origin: "localhost",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrlByShortUrl", args.ctx, "custom.url").Return(nil, nil)

				fields.qr.On("Encode", fmt.Sprintf("%s/%s", args.origin, "custom.url")).
					Return([]byte("png"), nil)

				fields.cfg.On("GetAwsS3Bucket").Return("bucket")

				fields.s3.On("Upload", args.ctx, "bucket", "custom.url.png", []byte("png"), "image/png").
					Return(errors.New("upload error"))
			},
			wantResp: dto.CreateUrlResponse{},
			wantErr:  errors.New("upload error"),
		},
		{
			name: "create qr code error",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				req: dto.CreateUrlRequest{
					Title:     "title",
					FullUrl:   "full.url",
					CustomUrl: "custom.url",
				},
				origin: "localhost",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrlByShortUrl", args.ctx, "custom.url").Return(nil, nil)

				fields.qr.On("Encode", fmt.Sprintf("%s/%s", args.origin, "custom.url")).
					Return(nil, errors.New("qr code error"))
			},
			wantResp: dto.CreateUrlResponse{},
			wantErr:  errors.New("qr code error"),
		},
		{
			name: "short url taken",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				req: dto.CreateUrlRequest{
					Title:     "title",
					FullUrl:   "full.url",
					CustomUrl: "custom.url",
				},
				origin: "localhost",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrlByShortUrl", args.ctx, "custom.url").Return(&db.Url{ID: 1}, nil)
			},
			wantResp: dto.CreateUrlResponse{},
			wantErr:  errorhelper.NewAppError(5, "Short url already taken", nil),
		},
		{
			name: "check short url available db error",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				req: dto.CreateUrlRequest{
					Title:     "title",
					FullUrl:   "full.url",
					CustomUrl: "custom.url",
				},
				origin: "localhost",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUrlByShortUrl", args.ctx, "custom.url").
					Return(nil, errors.New("check short url available db error"))
			},
			wantResp: dto.CreateUrlResponse{},
			wantErr:  errors.New("check short url available db error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(&tt.args, &tt.fields)

			r := NewUrls(tt.fields.cfg, tt.fields.repo, tt.fields.s3, tt.fields.log, tt.fields.rand, tt.fields.qr)

			got, gotErr := r.CreateUrl(tt.args.ctx, tt.args.req, tt.args.origin)

			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("Urls.CreateUrl() = %v, want %v", got, tt.wantResp)
			}

			if !errorhelper.IsEqual(gotErr, tt.wantErr) {
				t.Errorf("Urls.CreateUrl() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestDeleteUrl(t *testing.T) {
	type fields struct {
		cfg  *service.MockConfigProvider
		repo *repositories.MockIRepository
		s3   *aws.MockIS3
		log  *logger.MockILogger
		rand *random.MockIRandom
		qr   *qrcode.MockIQrCode
	}
	generateFields := func() fields {
		return fields{
			cfg:  service.NewMockConfigProvider(t),
			repo: repositories.NewMockIRepository(t),
			s3:   aws.NewMockIS3(t),
			log:  logger.NewMockILogger(t),
			rand: random.NewMockIRandom(t),
			qr:   qrcode.NewMockIQrCode(t),
		}
	}

	type args struct {
		ctx   context.Context
		urlId string
	}

	tests := []struct {
		name    string
		args    args
		fields  fields
		mocks   func(args *args, fields *fields)
		wantErr error
	}{
		{
			name: "success",
			args: args{
				ctx:   context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				urlId: "1",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("DeleteUrl", args.ctx, int64(1)).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "delete db error",
			args: args{
				ctx:   context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
				urlId: "1",
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("DeleteUrl", args.ctx, int64(1)).Return(errors.New("delete db error"))
			},
			wantErr: errors.New("delete db error"),
		},
		{
			name: "user ID not in context",
			args: args{
				ctx:   context.WithValue(context.Background(), reqctx.Key, reqctx.New("")),
				urlId: "1",
			},
			fields:  generateFields(),
			mocks:   func(args *args, fields *fields) {},
			wantErr: errors.New("user id cannot be determined from reqctx"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(&tt.args, &tt.fields)

			r := NewUrls(tt.fields.cfg, tt.fields.repo, tt.fields.s3, tt.fields.log, tt.fields.rand, tt.fields.qr)

			gotErr := r.DeleteUrl(tt.args.ctx, tt.args.urlId)

			if !errorhelper.IsEqual(gotErr, tt.wantErr) {
				t.Errorf("Urls.DeleteUrl() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}
