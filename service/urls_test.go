package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/mocks/github.com/raythx98/url-shortener/repositories"
	"github.com/raythx98/url-shortener/mocks/github.com/raythx98/url-shortener/service"
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
		log  *logger.MockILogger
		rand *random.MockIRandom
	}
	generateFields := func() fields {
		return fields{
			cfg:  service.NewMockConfigProvider(t),
			repo: repositories.NewMockIRepository(t),
			log:  logger.NewMockILogger(t),
			rand: random.NewMockIRandom(t),
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
						ID: 1, Title: "title", ShortUrl: "short.url", FullUrl: "full.url",
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
					Id: 1, Title: "title", ShortUrl: "short.url", FullUrl: "full.url",
					CreatedAt: time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC),
				},
				TotalClicks: 6,
				Devices:     []dto.Device{{Device: "mobile", Count: 3}, {Device: "desktop", Count: 2}, {Device: "pc", Count: 1}},
				Countries:   []dto.Country{{Country: "singapore", Count: 3}, {Country: "malaysia", Count: 2}, {Country: "australia", Count: 1}},
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
						ID: 1, Title: "title", ShortUrl: "short.url", FullUrl: "full.url",
						CreatedAt: pghelper.Time(timehelper.TimePtr(time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC))),
					}, nil)

				fields.repo.On("GetRedirectsByUrlId", args.ctx, inthelper.Int64Ptr(1)).Return(
					nil, nil)
			},
			wantResp: dto.GetUrlResponse{
				Url: dto.Url{
					Id: 1, Title: "title", ShortUrl: "short.url", FullUrl: "full.url",
					CreatedAt: time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC),
				},
				TotalClicks: 0,
				Devices:     []dto.Device{},
				Countries:   []dto.Country{},
			},
			wantErr: nil,
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

			r := NewUrls(tt.fields.cfg, tt.fields.repo, tt.fields.log, tt.fields.rand)

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
		log  *logger.MockILogger
		rand *random.MockIRandom
	}
	generateFields := func() fields {
		return fields{
			cfg:  service.NewMockConfigProvider(t),
			repo: repositories.NewMockIRepository(t),
			log:  logger.NewMockILogger(t),
			rand: random.NewMockIRandom(t),
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
							ID: 1, Title: "title", ShortUrl: "short.url", FullUrl: "long.url",
							CreatedAt: pghelper.Time(timehelper.TimePtr(time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC))),
						},
						{
							ID: 2, Title: "title2", ShortUrl: "short2.url", FullUrl: "long2.url",
							CreatedAt: pghelper.Time(timehelper.TimePtr(time.Date(2, 2, 3, 4, 5, 6, 7, time.UTC))),
						},
					}, nil)

				fields.repo.On("GetUserTotalClicks", args.ctx, inthelper.Int64Ptr(1)).Return(int64(10), nil)
			},
			wantResp: dto.GetUrlsResponse{
				Urls: []dto.Url{
					{
						Id: 1, Title: "title", ShortUrl: "short.url", FullUrl: "long.url",
						CreatedAt: time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC),
					},
					{
						Id: 2, Title: "title2", ShortUrl: "short2.url", FullUrl: "long2.url",
						CreatedAt: time.Date(2, 2, 3, 4, 5, 6, 7, time.UTC),
					},
				},
				TotalClicks: 10,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(&tt.args, &tt.fields)

			r := NewUrls(tt.fields.cfg, tt.fields.repo, tt.fields.log, tt.fields.rand)

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
		log  *logger.MockILogger
		rand *random.MockIRandom
	}
	generateFields := func() fields {
		return fields{
			cfg:  service.NewMockConfigProvider(t),
			repo: repositories.NewMockIRepository(t),
			log:  logger.NewMockILogger(t),
			rand: random.NewMockIRandom(t),
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

				fields.repo.On("CreateUrl", args.ctx, db.CreateUrlParams{
					UserID:   pghelper.Int8(inthelper.Int64Ptr(1)),
					Title:    "title",
					ShortUrl: "sh0rtieur1",
					FullUrl:  "full.url",
				}).Return(db.Url{
					ID:       1,
					ShortUrl: "sh0rtieur1",
				}, nil)
			},
			wantResp: dto.CreateUrlResponse{
				Id:       1,
				ShortUrl: "sh0rtieur1",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(&tt.args, &tt.fields)

			r := NewUrls(tt.fields.cfg, tt.fields.repo, tt.fields.log, tt.fields.rand)

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
		log  *logger.MockILogger
		rand *random.MockIRandom
	}
	generateFields := func() fields {
		return fields{
			cfg:  service.NewMockConfigProvider(t),
			repo: repositories.NewMockIRepository(t),
			log:  logger.NewMockILogger(t),
			rand: random.NewMockIRandom(t),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(&tt.args, &tt.fields)

			r := NewUrls(tt.fields.cfg, tt.fields.repo, tt.fields.log, tt.fields.rand)

			gotErr := r.DeleteUrl(tt.args.ctx, tt.args.urlId)

			if !errorhelper.IsEqual(gotErr, tt.wantErr) {
				t.Errorf("Urls.DeleteUrl() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}