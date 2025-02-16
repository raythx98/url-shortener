package dto

import "time"

type Url struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	ShortUrl  string    `json:"short_url"`
	FullUrl   string    `json:"full_url"`
	Qr        string    `json:"qr"`
	CreatedAt time.Time `json:"created_at"`
}
