// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package url_mappings

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type UrlMapping struct {
	ID               int64
	ShortenedUrl     string
	Url              string
	CreatedAt        pgtype.Timestamp
	InactiveExpireAt pgtype.Timestamp
	MustExpireAt     pgtype.Timestamp
}
