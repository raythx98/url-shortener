package controller

import (
	"encoding/json"
	"net/http"

	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/service"

	"github.com/raythx98/gohelpme/tool/httphelper"
	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"

	"github.com/go-playground/validator/v10"
)

type IRedirects interface {
	Redirect(w http.ResponseWriter, r *http.Request)
}

type Redirects struct {
	RedirectsService service.IRedirects
	Validator        *validator.Validate
	Log              logger.ILogger
}

func NewRedirects(service service.IRedirects, validate *validator.Validate, log logger.ILogger) *Redirects {
	return &Redirects{
		RedirectsService: service,
		Validator:        validate,
		Log:              log,
	}
}

func (c *Redirects) Redirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		reqctx.GetValue(ctx).SetError(err)
	}()
	
	req, err := httphelper.GetRequestBodyAndValidate[dto.RedirectRequest](ctx, r, c.Validator)
	if err != nil {
		return
	}

	resp, err := c.RedirectsService.Redirect(ctx, r.PathValue("shortLink"), req)
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

// Shorten example
//
//	@Summary		Shorten URL
//	@Description	Given a URL and custom settings, shorten it with a unique alias
//	@ID				shorten-url
//	@Accept			json
//	@Produce		json
//	@Param			req		body		dto.ShortenUrlRequest	true	"Shorten URL Request"
//	@Success		200		{string}	string					"Ok"
//	@Failure		422		{object}	dto.ErrorResponse		"Validation Error"
//	@Failure		500		{object}	dto.ErrorResponse		"Internal Server Error"
//	@Router			/url [post]

// Redirect example
//
//	@Summary		Redirects to full URL
//	@Description	Given an alias, redirects request to the full URL
//	@ID				redirect-alias
//	@Accept			json
//	@Produce		json
//	@Param			alias	path		string				true	"Alias"
//	@Success		303		{string}	interface{}			"Redirected"
//	@Failure		422		{object}	dto.ErrorResponse	"Validation Error"
//	@Failure		500		{object}	dto.ErrorResponse	"Internal Server Error"
//	@Router			/url/redirect/{alias} [post]
