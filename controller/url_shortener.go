package controller

import (
	"encoding/json"
	"fmt"
	"github.com/raythx98/gohelpme/tool/jwt"
	"github.com/raythx98/gohelpme/tool/reqctx"
	"net/http"
	"strings"
	"time"

	"github.com/raythx98/gohelpme/tool/httphelper"
	"github.com/raythx98/url-shortener/dto"
	"github.com/raythx98/url-shortener/service"

	"github.com/go-playground/validator/v10"
)

type IUrlShortener interface {
	Shorten(w http.ResponseWriter, r *http.Request)
	RedirectV2(w http.ResponseWriter, r *http.Request)
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

func (c *UrlShortener) Register(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	_, err := httphelper.GetRequestBodyAndValidate[req](ctx, r, validator.New())
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *UrlShortener) Login(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	type resp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	_, err := httphelper.GetRequestBodyAndValidate[req](ctx, r, validator.New())
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	accessToken, _ := jwt.NewAccessToken()
	refreshToken, _ := jwt.NewRefreshToken()

	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp{AccessToken: accessToken, RefreshToken: refreshToken})
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	_, err = w.Write(marshal)
	reqCtx.SetError(err)
}

func (c *UrlShortener) Refresh(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	//_, err := httphelper.GetRequestBodyAndValidate[req](ctx, r, validator.New())
	//if err != nil {
	//	reqCtx.SetError(err)
	//	return
	//}

	accessToken, _ := jwt.NewAccessToken()
	refreshToken, _ := jwt.NewRefreshToken()

	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp{AccessToken: accessToken, RefreshToken: refreshToken})
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	_, err = w.Write(marshal)
	reqCtx.SetError(err)
}

func (c *UrlShortener) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (c *UrlShortener) Profile(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Id    string `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	//_, err := httphelper.GetRequestBodyAndValidate[req](ctx, r, validator.New())
	//if err != nil {
	//	reqCtx.SetError(err)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp{
		Id:    "400d32ba-842e-41dd-9f15-68a8b1ecdd21",
		Email: "raythx98@gmail.com",
		Role:  "authenticated",
	})
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	_, err = w.Write(marshal)
	reqCtx.SetError(err)
}

func (c *UrlShortener) Url(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		c.DeleteUrl(w, r)
		return
	default:
		c.GetUrl(w, r)
		return
	}
}

func (c *UrlShortener) Urls(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.CreateUrl(w, r)
		return
	default:
		c.GetUrls(w, r)
		return
	}
}

func (c *UrlShortener) CreateUrl(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Title     string `json:"title" validate:"required"`
		FullUrl   string `json:"full_url" validate:"required"`
		CustomUrl string `json:"custom_url"`
		Qr        string `json:"qr" validate:"required"`
	}

	type resp struct {
		Id       int64  `json:"id"`
		ShortUrl string `json:"short_url"`
	}

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	_, err := httphelper.GetRequestBodyAndValidate[req](ctx, r, validator.New())
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	//marshal, err := json.Marshal(resp{Id: 1001, ShortUrl: "ABC123"})
	marshal, err := json.Marshal(resp{Id: 1001, ShortUrl: "ABC123ABC123ABC123ABC123"})
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	_, err = w.Write(marshal)
	reqCtx.SetError(err)
}

type url struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	ShortUrl    string    `json:"short_url"`
	OriginalUrl string    `json:"original_url"`
	Qr          string    `json:"qr"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *UrlShortener) GetUrl(w http.ResponseWriter, r *http.Request) {
	type clicks struct {
		Device string `json:"device"`
		City   string `json:"city"`
	}

	type resp struct {
		Url    url      `json:"url"`
		Clicks []clicks `json:"clicks"`
	}

	urlId := r.PathValue("id")
	fmt.Println("url_id", urlId)

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	// single
	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp{
		Url: url{
			Id:          1004,
			Title:       "Example 4",
			ShortUrl:    "example4",
			OriginalUrl: "http://example4.com",
			Qr:          "https://ykygqfljwketdjfxanoh.supabase.co/storage/v1/object/public/qrs/qr-mb90b5",
			CreatedAt:   time.Time{},
		},
		Clicks: []clicks{
			{
				Device: "Mobile",
				City:   "Singapore",
			},
			{
				Device: "Desktop",
				City:   "London",
			},
			{
				Device: "Tablet",
				City:   "Sydney",
			},
		},
	},
	)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	_, err = w.Write(marshal)
	reqCtx.SetError(err)
}

func (c *UrlShortener) GetUrls(w http.ResponseWriter, r *http.Request) {

	type resp struct {
		Urls        []url `json:"urls"`
		TotalClicks int   `json:"total_clicks"`
	}

	response := resp{
		Urls:        []url{},
		TotalClicks: 1231,
	}

	response.Urls = append(response.Urls,
		url{
			Id:          1001,
			Title:       "Example",
			ShortUrl:    "example",
			OriginalUrl: "http://example.com",
			Qr:          "https://ykygqfljwketdjfxanoh.supabase.co/storage/v1/object/public/qrs/qr-d69833",
			CreatedAt:   time.Time{},
		},
		url{
			Id:          1002,
			Title:       "Example 2",
			ShortUrl:    "example2",
			OriginalUrl: "http://example2.com",
			Qr:          "https://ykygqfljwketdjfxanoh.supabase.co/storage/v1/object/public/qrs/qr-mb90b5",
			CreatedAt:   time.Time{},
		},
		url{
			Id:          1003,
			Title:       "Example 3",
			ShortUrl:    "example3",
			OriginalUrl: "http://example3.com",
			Qr:          "https://ykygqfljwketdjfxanoh.supabase.co/storage/v1/object/public/qrs/qr-d69833",
			CreatedAt:   time.Time{},
		},
		url{
			Id:          1004,
			Title:       "Example 4",
			ShortUrl:    "example4",
			OriginalUrl: "http://example4.com",
			Qr:          "https://ykygqfljwketdjfxanoh.supabase.co/storage/v1/object/public/qrs/qr-mb90b5",
			CreatedAt:   time.Time{},
		})

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	//_, err := httphelper.GetRequestBodyAndValidate[req](ctx, r, validator.New())
	//if err != nil {
	//	reqCtx.SetError(err)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(response)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	_, err = w.Write(marshal)
	reqCtx.SetError(err)
}

func (c *UrlShortener) DeleteUrl(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	//reqCtx := reqctx.GetValue(ctx)
	urlId := r.PathValue("id")

	fmt.Println("urlId", urlId)

	w.WriteHeader(http.StatusOK)
}

func (c *UrlShortener) Redirect(w http.ResponseWriter, r *http.Request) {

	type resp struct {
		OriginalUrl string `json:"original_url"`
	}

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	greetings := r.PathValue("alias")
	fmt.Println("alias", greetings)

	//_, err := httphelper.GetRequestBodyAndValidate[req](ctx, r, validator.New())
	//if err != nil {
	//	reqCtx.SetError(err)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp{OriginalUrl: "https://facebook.com"})
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
func (c *UrlShortener) RedirectV2(w http.ResponseWriter, r *http.Request) {
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
