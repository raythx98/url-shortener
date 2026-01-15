package dto

import "time"

type Url struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	ShortUrl  string    `json:"short_url"`
	FullUrl   string    `json:"full_url"`
	CreatedAt time.Time `json:"created_at"`
}
