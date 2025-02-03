package controller

import (
	"encoding/json"
	"github.com/raythx98/gohelpme/tool/reqctx"
	"net/http"
	"strings"

	"github.com/raythx98/gohelpme/tool/httphelper"
	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/service"

	"github.com/go-playground/validator/v10"
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

type testReq struct {
	Request int `json:"request"`
}

type testResp struct {
	Response int `json:"response"`
}

func (c *UrlShortener) Test(w http.ResponseWriter, r *http.Request) {
	//execute := func() error {
	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	_, err := httphelper.GetRequestBodyAndValidate[testReq](ctx, r, validator.New())
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	//reqCtx.SetError(fmt.Errorf("test error"))
	//return
	//return errorf.New("test error")
	//panic("implement me") // TODO: Figure out how to recover panic and continue logging

	//url, err := c.UrlShortenerService.ShortenUrl(ctx, req)
	//if err != nil {
	//	return err
	//}

	w.WriteHeader(http.StatusCreated)
	marshal, err := json.Marshal(testResp{Response: 1})
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	_, err = w.Write(marshal)
	reqCtx.SetError(err)
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
func (c *UrlShortener) Shorten(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	req, err := httphelper.GetRequestBodyAndValidate[dto.ShortenUrlRequest](ctx, r, validator.New())
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	url, err := c.UrlShortenerService.ShortenUrl(ctx, req)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	marshal, err := json.Marshal(url)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	_, err = w.Write(marshal)
	reqCtx.SetError(err)
}

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
func (c *UrlShortener) Redirect(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	greetings := r.PathValue("alias")
	url, err := c.UrlShortenerService.GetUrlWithShortened(ctx, greetings)
	if url == "" || err != nil {
		reqCtx.SetError(err)
		return
	}

	if !strings.Contains(url, "://") { // TODO: Fix protocol not being added to full url
		url = "http://" + url
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}
