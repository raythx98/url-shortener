package repositories

import (
	"context"
	"errors"

	"github.com/raythx98/go-zap/sqlc/db"
	"github.com/raythx98/go-zap/tools/pghelper"
	"github.com/raythx98/go-zap/tools/postgres"

	"github.com/jackc/pgx/v5"
)

type IRepository interface {
	CreateRedirect(ctx context.Context, arg db.CreateRedirectParams) error
	CreateUrl(ctx context.Context, arg db.CreateUrlParams) (db.Url, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	DeleteUrl(ctx context.Context, id int64) error
	GetRedirectsByUrlId(ctx context.Context, urlId *int64) ([]db.Redirect, error)
	GetUrl(ctx context.Context, id int64) (db.Url, error)
	GetUrlByShortUrl(ctx context.Context, shortUrl string) (*db.Url, error)
	GetUrlsByUserId(ctx context.Context, userId *int64) ([]db.Url, error)
	GetUserByEmail(ctx context.Context, email string) (*db.User, error)
	GetUserTotalClicks(ctx context.Context, userId *int64) (int64, error)
}

type Repository struct {
	db *db.Queries
}

func NewRepository(pg postgres.IPostgres) *Repository {
	return &Repository{
		db: db.New(pg.Pool()),
	}
}

func (r *Repository) CreateRedirect(ctx context.Context, arg db.CreateRedirectParams) error {
	return r.db.CreateRedirect(ctx, arg)
}

func (r *Repository) CreateUrl(ctx context.Context, arg db.CreateUrlParams) (db.Url, error) {
	return r.db.CreateUrl(ctx, arg)
}

func (r *Repository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return r.db.CreateUser(ctx, arg)
}

func (r *Repository) DeleteUrl(ctx context.Context, id int64) error {
	return r.db.DeleteUrl(ctx, id)
}

func (r *Repository) GetRedirectsByUrlId(ctx context.Context, urlId *int64) ([]db.Redirect, error) {
	return r.db.GetRedirectsByUrlId(ctx, pghelper.Int8(urlId))
}

func (r *Repository) GetUrl(ctx context.Context, id int64) (db.Url, error) {
	return r.db.GetUrl(ctx, id)
}

func (r *Repository) GetUrlByShortUrl(ctx context.Context, shortUrl string) (*db.Url, error) {
	url, err := r.db.GetUrlByShortUrl(ctx, shortUrl)
	if err == nil {
		return &url, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) GetUrlsByUserId(ctx context.Context, userId *int64) ([]db.Url, error) {
	return r.db.GetUrlsByUserId(ctx, pghelper.Int8(userId))
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := r.db.GetUserByEmail(ctx, email)
	if err == nil {
		return &user, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (r *Repository) GetUserTotalClicks(ctx context.Context, userId *int64) (int64, error) {
	return r.db.GetUserTotalClicks(ctx, pghelper.Int8(userId))
}
