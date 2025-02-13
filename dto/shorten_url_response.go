package dto

//import (
//	"github.com/raythx98/url-shortener/sqlc"
//	"time"
//)
//
//// ShortenUrlResponse response body for shorten url
//type ShortenUrlResponse struct {
//	Url              string     `json:"url"`
//	ShortenedUrl     string     `json:"shortened_url"`
//	InactiveExpireAt *time.Time `json:"inactive_expire_at,omitempty"`
//	MustExpireAt     *time.Time `json:"must_expire_at,omitempty"`
//}
//
//func (s *ShortenUrlResponse) Bind(entity sqlc.UrlMapping) *ShortenUrlResponse {
//	var inactiveExpireAt, mustExpireAt *time.Time
//	if entity.InactiveExpireAt.Valid {
//		inactiveExpireAt = &entity.InactiveExpireAt.Time
//	}
//	if entity.MustExpireAt.Valid {
//		mustExpireAt = &entity.MustExpireAt.Time
//	}
//
//	s.Url = entity.Url
//	s.ShortenedUrl = entity.ShortenedUrl
//	s.InactiveExpireAt = inactiveExpireAt
//	s.MustExpireAt = mustExpireAt
//	return s
//}
