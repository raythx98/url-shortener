package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/raythx98/gohelpme/tool/httphelper"
	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/service"
	"net/http"
	"strings"
)

type IUrlShortener interface {
	Shorten(w http.ResponseWriter, r *http.Request)
	Redirect(w http.ResponseWriter, r *http.Request)
}

type UrlShortener struct {
	UrlShortenerService service.IUrlShortener
	Validator           *validator.Validate
}

func New(service service.IUrlShortener) UrlShortener {
	return UrlShortener{
		UrlShortenerService: service,
	}
}

func (c *UrlShortener) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req, err := httphelper.GetRequestBodyAndValidate[dto.ShortenUrlRequest](ctx, r, validator.New())
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("{\"response\": \"ERROR VALIDATION\"}"))
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Printf("req after unmarshal: %+v\n", req)

	url, err := c.UrlShortenerService.ShortenUrl(ctx, req.Url)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("{\"response\": \"ERROR\"}"))
		fmt.Printf("%+v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	marshal, err := json.Marshal(url)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("{\"response\": \"ERROR MARSHALLING\"}"))
		fmt.Printf("%+v\n", err)
		return
	}
	_, _ = w.Write(marshal)
}

func (c *UrlShortener) Redirect(w http.ResponseWriter, r *http.Request) {
	greetings := r.PathValue("alias")
	url, err := c.UrlShortenerService.GetUrlWithShortened(r.Context(), greetings)
	if url == "" || err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("{\"response\": \"ERROR\"}"))
		fmt.Printf("%+v\n", err)
		return
	}

	if !strings.Contains(url, "://") { // TODO: Fix protocol not being added to full url
		url = "http://" + url
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}
