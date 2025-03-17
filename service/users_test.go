package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/mocks/github.com/raythx98/url-shortener/repositories"
	"github.com/raythx98/url-shortener/sqlc/db"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/mocks/github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

func TestGetProfile(t *testing.T) {
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
		ctx context.Context
	}

	tests := []struct {
		name     string
		args     args
		fields   fields
		mocks    func(args *args, fields *fields)
		wantResp dto.ProfileResponse
		wantErr  error
	}{
		{
			name: "success",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUser", args.ctx, int64(1)).Return(
					db.User{
						ID:    1,
						Email: "test@gmail.com",
					}, nil)
			},
			wantResp: dto.ProfileResponse{
				Id:    1,
				Email: "test@gmail.com",
				Role:  "authenticated",
			},
			wantErr: nil,
		},
		{
			name: "user id not in context",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("")),
			},
			fields:   generateFields(),
			mocks:    func(args *args, fields *fields) {},
			wantResp: dto.ProfileResponse{},
			wantErr:  errors.New("user id cannot be determined from reqctx"),
		},
		{
			name: "user id not in db",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUser", args.ctx, int64(1)).
					Return(db.User{}, errors.New("user id not found in db"))
			},
			wantResp: dto.ProfileResponse{},
			wantErr:  errors.New("user id not found in db"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(&tt.args, &tt.fields)

			r := NewUsers(tt.fields.repo, tt.fields.log)

			got, gotErr := r.GetProfile(tt.args.ctx)

			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("Users.GetProfile() = %v, want %v", got, tt.wantResp)
			}

			if !errorhelper.IsEqual(gotErr, tt.wantErr) {
				t.Errorf("Users.GetProfile() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}
