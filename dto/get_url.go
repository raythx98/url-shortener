package dto

type Device struct {
	Device string `json:"device"`
	Count  int    `json:"count"`
}

type Country struct {
	Country string `json:"country"`
	Count   int    `json:"count"`
}

type GetUrlResponse struct {
	Url         Url       `json:"url"`
	TotalClicks int       `json:"total_clicks"`
	Devices     []Device  `json:"devices"`
	Countries   []Country `json:"countries"`
}
