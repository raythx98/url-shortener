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

func New(service service.IUrlShortener, validate *validator.Validate) *UrlShortener {
	return &UrlShortener{
		UrlShortenerService: service,
		Validator:           validate,
	}
}

func (c *UrlShortener) Shorten(w http.ResponseWriter, r *http.Request) {
	execute := func() error {
		ctx := r.Context()

		req, err := httphelper.GetRequestBodyAndValidate[dto.ShortenUrlRequest](ctx, r, validator.New())
		if err != nil {
			return err
		}

		//panic("implement me") TODO: Figure out how to recover panic and continue logging

		url, err := c.UrlShortenerService.ShortenUrl(ctx, req)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)
		marshal, err := json.Marshal(url)
		if err != nil {
			return err
		}

		_, err = w.Write(marshal)
		return err
	}

	if err := execute(); err != nil {
		error_tool.Handle(w, err) // TODO: Can we write a middleware to handle error in a more elegant way?
	}
}

func (c *UrlShortener) Redirect(w http.ResponseWriter, r *http.Request) {
	execute := func() error {
		greetings := r.PathValue("alias")
		url, err := c.UrlShortenerService.GetUrlWithShortened(r.Context(), greetings)
		if url == "" || err != nil {
			return err
		}

		if !strings.Contains(url, "://") { // TODO: Fix protocol not being added to full url
			url = "http://" + url
		}

		http.Redirect(w, r, url, http.StatusSeeOther)
		return nil
	}

	if err := execute(); err != nil {
		error_tool.Handle(w, err) // TODO: Can we write a middleware to handle error in a more elegant way?
	}
}
