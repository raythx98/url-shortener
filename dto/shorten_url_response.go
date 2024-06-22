package dto

import "time"

// ShortenUrlResponse response body for shorten url
type ShortenUrlResponse struct {
	Url              string     `json:"url"`
	ShortenedUrl     string     `json:"shortened_url"`
	InactiveExpireAt *time.Time `json:"inactive_expire_at,omitempty"`
	MustExpireAt     *time.Time `json:"must_expire_at,omitempty"`
}
