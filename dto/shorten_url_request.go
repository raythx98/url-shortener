package dto

import (
	"time"

	"github.com/raythx98/url-shortener/sqlc/url_mappings"
)

// ShortenUrlRequest request body for shorten url
type ShortenUrlRequest struct {
	// TODO: Add custom validation for URL Required
	Url          string         `json:"url" validate:"required"`
	ShortenedUrl *string        `json:"shortened_url"`
	CustomExpiry *time.Duration `json:"custom_expiry"`
	IsNoExpiry   *bool          `json:"is_no_expiry"`
}

func (s *ShortenUrlRequest) BindTo(p url_mappings.CreateUrlMappingParams) url_mappings.CreateUrlMappingParams {
	p.Url = s.Url
	return p
}
