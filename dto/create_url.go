package dto

type CreateUrlRequest struct {
	Title     string `json:"title" validate:"required"`
	FullUrl   string `json:"full_url" validate:"required"`
	CustomUrl string `json:"custom_url" validate:"omitempty,min=4,max=20"`
	Qr        string `json:"qr" validate:"required"`
}

type CreateUrlResponse struct {
	Id       int64  `json:"id"`
	ShortUrl string `json:"short_url"`
}
