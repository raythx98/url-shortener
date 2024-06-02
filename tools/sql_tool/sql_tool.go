package sql_tool

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

func NewTime(newTime *time.Time) pgtype.Timestamp {
	if newTime == nil || newTime.IsZero() {
		return pgtype.Timestamp{Valid: false}
	}
	return pgtype.Timestamp{
		Time:  *newTime,
		Valid: true,
	}
}
