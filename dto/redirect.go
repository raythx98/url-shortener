package dto

type RedirectRequest struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Device  string `json:"device"`
}

type RedirectResponse struct {
	FullUrl string `json:"full_url"`
}
