// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Redirect struct {
	ID        int64
	UrlID     pgtype.Int8
	Device    string
	Country   string
	City      string
	CreatedAt pgtype.Timestamp
}

type Url struct {
	ID        int64
	UserID    pgtype.Int8
	Title     string
	ShortUrl  string
	FullUrl   string
	Qr        string
	CreatedAt pgtype.Timestamp
	IsDeleted bool
}

type User struct {
	ID        int64
	Email     string
	Password  string
	CreatedAt pgtype.Timestamp
	IsActive  bool
}
