package dto

type GetUrlsResponse struct {
	Urls        []Url `json:"urls"`
	TotalClicks int64 `json:"total_clicks"`
}
