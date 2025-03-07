package controller

import (
	"encoding/json"
	"net/http"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/service"

	"github.com/raythx98/gohelpme/tool/httphelper"
	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"
	"github.com/raythx98/gohelpme/tool/validator"
)

type IUrls interface {
	GetUrl(w http.ResponseWriter, r *http.Request)
	GetUrls(w http.ResponseWriter, r *http.Request)
	CreateUrl(w http.ResponseWriter, r *http.Request)
	DeleteUrl(w http.ResponseWriter, r *http.Request)
}

type Urls struct {
	UrlsService service.IUrls
	Validator   validator.IValidator
	Log         logger.ILogger
}

func NewUrls(service service.IUrls, validate validator.IValidator, log logger.ILogger) *Urls {
	return &Urls{
		UrlsService: service,
		Validator:   validate,
		Log:         log,
	}
}

func (c *Urls) GetUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		reqctx.GetValue(ctx).SetError(err)
	}()

	resp, err := c.UrlsService.GetUrl(ctx, r.PathValue("id"))
	if err != nil {
		return
	}

	// single
	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp)
	if err != nil {
		return
	}

	_, err = w.Write(marshal)
}

func (c *Urls) GetUrls(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		reqctx.GetValue(ctx).SetError(err)
	}()

	resp, err := c.UrlsService.GetUrls(ctx)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp)
	if err != nil {
		return
	}

	_, err = w.Write(marshal)
}

func (c *Urls) CreateUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		reqctx.GetValue(ctx).SetError(err)
	}()

	req, err := httphelper.GetRequestBodyAndValidate[dto.CreateUrlRequest](ctx, r, validator.New())
	if err != nil {
		return
	}

	resp, err := c.UrlsService.CreateUrl(ctx, req, r.Header.Get("Origin"))
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp)
	if err != nil {
		return
	}

	_, err = w.Write(marshal)
}

func (c *Urls) DeleteUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		reqctx.GetValue(ctx).SetError(err)
	}()

	err = c.UrlsService.DeleteUrl(ctx, r.PathValue("id"))
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
}
