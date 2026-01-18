package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/raythx98/go-zap/dto"
	"github.com/raythx98/go-zap/mocks/github.com/raythx98/go-zap/repositories"
	"github.com/raythx98/go-zap/mocks/github.com/raythx98/go-zap/tools/crypto"
	"github.com/raythx98/go-zap/sqlc/db"

	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/mocks/github.com/raythx98/gohelpme/tool/jwthelper"
	"github.com/raythx98/gohelpme/mocks/github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"
)

func TestRegister(t *testing.T) {
	type fields struct {
		repo   *repositories.MockIRepository
		jwt    *jwthelper.MockIJwt
		log    *logger.MockILogger
		crypto *crypto.MockICrypto
	}
	generateFields := func() fields {
		return fields{
			repo:   repositories.NewMockIRepository(t),
			jwt:    jwthelper.NewMockIJwt(t),
			log:    logger.NewMockILogger(t),
			crypto: crypto.NewMockICrypto(t),
		}
	}

	type args struct {
		ctx context.Context
		req dto.RegisterRequest
	}

	tests := []struct {
		name     string
		args     args
		fields   fields
		mocks    func(args *args, fields *fields)
		wantResp dto.LoginResponse
		wantErr  error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: dto.RegisterRequest{
					Email:    "register@gmail.com",
					Password: "register_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "register@gmail.com").
					Return(nil, nil)

				fields.crypto.On("GenerateFromPassword", "register_password").
					Return("encoded_password", nil)

				fields.repo.On("CreateUser", args.ctx, db.CreateUserParams{
					Email:    "register@gmail.com",
					Password: "encoded_password",
				}).Return(db.User{
					ID: 1,
				}, nil)

				fields.jwt.On("NewAccessToken", "1").Return("access_token", nil)

				fields.jwt.On("NewRefreshToken", "1").Return("refresh_token", nil)
			},
			wantResp: dto.LoginResponse{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
			},
			wantErr: nil,
		},
		{
			name: "fail to generate refresh token",
			args: args{
				ctx: context.Background(),
				req: dto.RegisterRequest{
					Email:    "register@gmail.com",
					Password: "register_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "register@gmail.com").
					Return(nil, nil)

				fields.crypto.On("GenerateFromPassword", "register_password").
					Return("encoded_password", nil)

				fields.repo.On("CreateUser", args.ctx, db.CreateUserParams{
					Email:    "register@gmail.com",
					Password: "encoded_password",
				}).Return(db.User{
					ID: 1,
				}, nil)

				fields.jwt.On("NewAccessToken", "1").Return("access_token", nil)

				fields.jwt.On("NewRefreshToken", "1").Return("", errors.New("refresh error"))
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("refresh error"),
		},
		{
			name: "fail to generate access token",
			args: args{
				ctx: context.Background(),
				req: dto.RegisterRequest{
					Email:    "register@gmail.com",
					Password: "register_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "register@gmail.com").
					Return(nil, nil)

				fields.crypto.On("GenerateFromPassword", "register_password").
					Return("encoded_password", nil)

				fields.repo.On("CreateUser", args.ctx, db.CreateUserParams{
					Email:    "register@gmail.com",
					Password: "encoded_password",
				}).Return(db.User{
					ID: 1,
				}, nil)

				fields.jwt.On("NewAccessToken", "1").Return("", errors.New("access token error"))
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("access token error"),
		},
		{
			name: "create user db error",
			args: args{
				ctx: context.Background(),
				req: dto.RegisterRequest{
					Email:    "register@gmail.com",
					Password: "register_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "register@gmail.com").
					Return(nil, nil)

				fields.crypto.On("GenerateFromPassword", "register_password").
					Return("encoded_password", nil)

				fields.repo.On("CreateUser", args.ctx, db.CreateUserParams{
					Email:    "register@gmail.com",
					Password: "encoded_password",
				}).Return(db.User{}, errors.New("create user db error"))

			},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("create user db error"),
		},
		{
			name: "encode password error",
			args: args{
				ctx: context.Background(),
				req: dto.RegisterRequest{
					Email:    "register@gmail.com",
					Password: "register_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "register@gmail.com").
					Return(nil, nil)

				fields.crypto.On("GenerateFromPassword", "register_password").
					Return("encoded_password", errors.New("encode password error"))
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("encode password error"),
		},
		{
			name: "get existing user db error",
			args: args{
				ctx: context.Background(),
				req: dto.RegisterRequest{
					Email:    "register@gmail.com",
					Password: "register_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "register@gmail.com").
					Return(nil, errors.New("get existing user db error"))
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("get existing user db error"),
		},
		{
			name: "email already exists",
			args: args{
				ctx: context.Background(),
				req: dto.RegisterRequest{
					Email:    "register@gmail.com",
					Password: "register_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "register@gmail.com").
					Return(&db.User{ID: 2}, nil)
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errorhelper.NewAppError(1, "Email has already been registered", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(&tt.args, &tt.fields)

			r := NewAuth(tt.fields.repo, tt.fields.log, tt.fields.jwt, tt.fields.crypto)

			got, gotErr := r.Register(tt.args.ctx, tt.args.req)

			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("Auth.Register() = %v, want %v", got, tt.wantResp)
			}

			if !errorhelper.IsEqual(gotErr, tt.wantErr) {
				t.Errorf("Auth.Register() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	type fields struct {
		repo   *repositories.MockIRepository
		jwt    *jwthelper.MockIJwt
		log    *logger.MockILogger
		crypto *crypto.MockICrypto
	}
	generateFields := func() fields {
		return fields{
			repo:   repositories.NewMockIRepository(t),
			jwt:    jwthelper.NewMockIJwt(t),
			log:    logger.NewMockILogger(t),
			crypto: crypto.NewMockICrypto(t),
		}
	}

	type args struct {
		ctx context.Context
		req dto.LoginRequest
	}

	tests := []struct {
		name     string
		args     args
		fields   fields
		mocks    func(args *args, fields *fields)
		wantResp dto.LoginResponse
		wantErr  error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: dto.LoginRequest{
					Email:    "login@gmail.com",
					Password: "login_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "login@gmail.com").Return(
					&db.User{
						ID:       1,
						Password: "hash_password",
					}, nil)

				fields.crypto.On("ComparePasswordAndHash", "login_password", "hash_password").
					Return(true, nil)

				fields.jwt.On("NewAccessToken", "1").Return("access_token", nil)

				fields.jwt.On("NewRefreshToken", "1").Return("refresh_token", nil)
			},
			wantResp: dto.LoginResponse{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
			},
			wantErr: nil,
		},
		{
			name: "refresh token error",
			args: args{
				ctx: context.Background(),
				req: dto.LoginRequest{
					Email:    "login@gmail.com",
					Password: "login_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "login@gmail.com").Return(
					&db.User{
						ID:       1,
						Password: "hash_password",
					}, nil)

				fields.crypto.On("ComparePasswordAndHash", "login_password", "hash_password").
					Return(true, nil)

				fields.jwt.On("NewAccessToken", "1").Return("access_token", nil)

				fields.jwt.On("NewRefreshToken", "1").Return("", errors.New("refresh error"))
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("refresh error"),
		},
		{
			name: "create access token error",
			args: args{
				ctx: context.Background(),
				req: dto.LoginRequest{
					Email:    "login@gmail.com",
					Password: "login_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "login@gmail.com").Return(
					&db.User{
						ID:       1,
						Password: "hash_password",
					}, nil)

				fields.crypto.On("ComparePasswordAndHash", "login_password", "hash_password").
					Return(true, nil)

				fields.jwt.On("NewAccessToken", "1").Return("", errors.New("access token error"))
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("access token error"),
		},
		{
			name: "wrong password",
			args: args{
				ctx: context.Background(),
				req: dto.LoginRequest{
					Email:    "login@gmail.com",
					Password: "login_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "login@gmail.com").Return(
					&db.User{
						ID:       1,
						Password: "hash_password",
					}, nil)

				fields.crypto.On("ComparePasswordAndHash", "login_password", "hash_password").
					Return(false, nil)
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errorhelper.NewAppError(2, "Incorrect Password", nil),
		},
		{
			name: "compare password error",
			args: args{
				ctx: context.Background(),
				req: dto.LoginRequest{
					Email:    "login@gmail.com",
					Password: "login_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "login@gmail.com").Return(
					&db.User{
						ID:       1,
						Password: "hash_password",
					}, nil)

				fields.crypto.On("ComparePasswordAndHash", "login_password", "hash_password").
					Return(false, errors.New("compare password error"))
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("compare password error"),
		},
		{
			name: "not registered",
			args: args{
				ctx: context.Background(),
				req: dto.LoginRequest{
					Email:    "login@gmail.com",
					Password: "login_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "login@gmail.com").Return(
					nil, nil)
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errorhelper.NewAppError(3, "Email is not registered", nil),
		},
		{
			name: "get user db error",
			args: args{
				ctx: context.Background(),
				req: dto.LoginRequest{
					Email:    "login@gmail.com",
					Password: "login_password",
				},
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.repo.On("GetUserByEmail", args.ctx, "login@gmail.com").Return(
					nil, errors.New("get user db error"))
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("get user db error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(&tt.args, &tt.fields)

			r := NewAuth(tt.fields.repo, tt.fields.log, tt.fields.jwt, tt.fields.crypto)

			got, gotErr := r.Login(tt.args.ctx, tt.args.req)

			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("Auth.Login() = %v, want %v", got, tt.wantResp)
			}

			if !errorhelper.IsEqual(gotErr, tt.wantErr) {
				t.Errorf("Auth.Login() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestRefresh(t *testing.T) {
	type fields struct {
		repo   *repositories.MockIRepository
		jwt    *jwthelper.MockIJwt
		log    *logger.MockILogger
		crypto *crypto.MockICrypto
	}
	generateFields := func() fields {
		return fields{
			repo:   repositories.NewMockIRepository(t),
			jwt:    jwthelper.NewMockIJwt(t),
			log:    logger.NewMockILogger(t),
			crypto: crypto.NewMockICrypto(t),
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
		wantResp dto.LoginResponse
		wantErr  error
	}{
		{
			name: "success",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.jwt.On("NewAccessToken", "1").Return("access_token", nil)

				fields.jwt.On("NewRefreshToken", "1").Return("refresh_token", nil)
			},
			wantResp: dto.LoginResponse{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
			},
			wantErr: nil,
		},
		{
			name: "refresh token fail",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.jwt.On("NewAccessToken", "1").Return("access_token", nil)

				fields.jwt.On("NewRefreshToken", "1").Return("", errors.New("refresh fail"))
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("refresh fail"),
		},
		{
			name: "create access token fail",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("").SetUserId(1)),
			},
			fields: generateFields(),
			mocks: func(args *args, fields *fields) {
				fields.jwt.On("NewAccessToken", "1").Return("", errors.New("access token fail"))
			},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("access token fail"),
		},
		{
			name: "user id not in context",
			args: args{
				ctx: context.WithValue(context.Background(), reqctx.Key, reqctx.New("")),
			},
			fields:   generateFields(),
			mocks:    func(args *args, fields *fields) {},
			wantResp: dto.LoginResponse{},
			wantErr:  errors.New("user id cannot be determined from reqctx"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mocks(&tt.args, &tt.fields)

			r := NewAuth(tt.fields.repo, tt.fields.log, tt.fields.jwt, tt.fields.crypto)

			got, gotErr := r.Refresh(tt.args.ctx)

			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("Auth.Refresh() = %v, want %v", got, tt.wantResp)
			}

			if !errorhelper.IsEqual(gotErr, tt.wantErr) {
				t.Errorf("Auth.Refresh() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}
