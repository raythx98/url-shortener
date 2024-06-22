package controller

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/raythx98/gohelpme/tool/httphelper"
	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/service"
	"github.com/raythx98/url-shortener/tools/error_tool"
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
		error_tool.Handle(w, err)
		return
	}

	//panic("implement me") TODO: Figure out how to recover panic and continue logging

	url, err := c.UrlShortenerService.ShortenUrl(ctx, req)
	if err != nil {
		error_tool.Handle(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	marshal, err := json.Marshal(url)
	if err != nil {
		error_tool.Handle(w, err)
		return
	}
	_, _ = w.Write(marshal)
}

func (c *UrlShortener) Redirect(w http.ResponseWriter, r *http.Request) {
	greetings := r.PathValue("alias")
	url, err := c.UrlShortenerService.GetUrlWithShortened(r.Context(), greetings)
	if url == "" || err != nil {
		error_tool.Handle(w, err)
		return
	}

	if !strings.Contains(url, "://") { // TODO: Fix protocol not being added to full url
		url = "http://" + url
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}
