package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/raythx98/go-zap/dto"
	"github.com/raythx98/go-zap/mocks/github.com/raythx98/go-zap/repositories"
	"github.com/raythx98/go-zap/sqlc/db"
	"github.com/raythx98/go-zap/tools/pghelper"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/mocks/github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/inthelper"

	"github.com/stretchr/testify/mock"
)

func TestRedirect(t *testing.T) {
	type fields struct {
		repo *repositories.MockIRepository
		log  *logger.MockILogger
	}
	generateFields := func() fields {
		return fields{
			repo: repositories.NewMockIRepository(t),
			log:  logger.NewMockILogger(t),
		}
	}

	type args struct {
		ctx      context.Context
		shortUrl string
		req      dto.RedirectRequest
	}

	tests := []struct {
		name     string
		args     args
		fields   fields
		mocks    func(args *args, fields *fields, goroutineCalled chan struct{})
		wantResp dto.RedirectResponse
		wantErr  error
	}{
		{
			name: "success",
			args: args{
				ctx:      context.Background(),
				shortUrl: "shorturl",
				req: dto.RedirectRequest{
					City:    "Singapore",
					Country: "Singapore",
					Device:  "Mobile",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields, goroutineCalled chan struct{}) {
				fields.repo.On("GetUrlByShortUrl", args.ctx, "shorturl").Return(
					&db.Url{
						ID:      1,
						FullUrl: "full.url",
					}, nil)

				fields.repo.On("CreateRedirect", mock.Anything, db.CreateRedirectParams{
					UrlID:   pghelper.Int8(inthelper.Int64Ptr(1)),
					City:    "Singapore",
					Country: "Singapore",
					Device:  "Mobile",
				}).Return(nil).Run(func(args mock.Arguments) { goroutineCalled <- struct{}{} })
			},
			wantResp: dto.RedirectResponse{
				FullUrl: "full.url",
			},
			wantErr: nil,
		},
		{
			name: "db error",
			args: args{
				ctx:      context.Background(),
				shortUrl: "shorturl",
				req: dto.RedirectRequest{
					City:    "Singapore",
					Country: "Singapore",
					Device:  "Mobile",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields, goroutineCalled chan struct{}) {
				fields.repo.On("GetUrlByShortUrl", args.ctx, "shorturl").
					Return(nil, errors.New("GetUrlByShortUrl db error"))
			},
			wantResp: dto.RedirectResponse{},
			wantErr:  errors.New("GetUrlByShortUrl db error"),
		},
		{
			name: "short url not found",
			args: args{
				ctx:      context.Background(),
				shortUrl: "shorturl",
				req: dto.RedirectRequest{
					City:    "Singapore",
					Country: "Singapore",
					Device:  "Mobile",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields, goroutineCalled chan struct{}) {
				fields.repo.On("GetUrlByShortUrl", args.ctx, "shorturl").Return(nil, nil)
			},
			wantResp: dto.RedirectResponse{},
			wantErr:  errorhelper.NewAppError(4, "Invalid short url, please create a new one", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goroutineCalled := make(chan struct{})
			tt.mocks(&tt.args, &tt.fields, goroutineCalled)

			r := NewRedirects(tt.fields.repo, tt.fields.log)

			got, gotErr := r.Redirect(tt.args.ctx, tt.args.shortUrl, tt.args.req)

			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("Redirects.Redirect() = %v, want %v", got, tt.wantResp)
			}

			if !errorhelper.IsEqual(gotErr, tt.wantErr) {
				t.Errorf("Redirects.Redirect() error = %v, wantErr %v", gotErr, tt.wantErr)
			}

			select {
			case <-goroutineCalled:
				// Goroutine completed
			case <-time.After(5 * time.Second):
				// Wait for possible goroutine to complete
			}
		})
	}
}
