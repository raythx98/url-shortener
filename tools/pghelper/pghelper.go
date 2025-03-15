package pghelper

import "github.com/jackc/pgx/v5/pgtype"

func Int8(userId *int64) pgtype.Int8 {
	if userId == nil {
		return pgtype.Int8{Valid: false}
	}

	return pgtype.Int8{Int64: *userId, Valid: true}
}
