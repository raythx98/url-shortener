package dto

type CreateUrlRequest struct {
	Title     string `json:"title" validate:"required,max=255"`
	FullUrl   string `json:"full_url" validate:"required,max=2048"`
	CustomUrl string `json:"custom_url" validate:"omitempty,min=4,max=255,alphanum"`
}

type CreateUrlResponse struct {
	Id       int64  `json:"id"`
	ShortUrl string `json:"short_url"`
}
