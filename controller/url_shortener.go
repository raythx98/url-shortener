package controller

import (
	"fmt"
	"github.com/raythx98/url-shortener/service"
	"net/http"
)

type IUrlShortener interface {
	Shorten(w http.ResponseWriter, r *http.Request)
	Redirect(w http.ResponseWriter, r *http.Request)
}

type UrlShortener struct {
	UrlShortenerService service.IUrlShortener
}

func New(service service.IUrlShortener) UrlShortener {
	return UrlShortener{
		UrlShortenerService: service,
	}
}

func (c *UrlShortener) Shorten(w http.ResponseWriter, r *http.Request) {
	url, err := c.UrlShortenerService.ShortenUrl(r.Context(), "www.google.com")
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("{\"response\": \"ERROR\"}"))
		fmt.Printf("%+v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(url))
}

func (c *UrlShortener) Redirect(w http.ResponseWriter, r *http.Request) {
	url, err := c.UrlShortenerService.GetUrlWithShortened(r.Context(), "test")
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("{\"response\": \"ERROR\"}"))
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}
