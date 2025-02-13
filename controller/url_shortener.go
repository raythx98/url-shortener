package controller

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/raythx98/gohelpme/errorhelper"
	"github.com/raythx98/gohelpme/tool/jwt"
	"github.com/raythx98/gohelpme/tool/reqctx"
	"github.com/raythx98/url-shortener/sqlc/db"
	"github.com/rs/zerolog/log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/raythx98/gohelpme/tool/httphelper"
	"github.com/raythx98/url-shortener/service"

	"github.com/go-playground/validator/v10"
)

type IUrlShortener interface {
	//Shorten(w http.ResponseWriter, r *http.Request)
	//RedirectV2(w http.ResponseWriter, r *http.Request)
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

	request, err := httphelper.GetRequestBodyAndValidate[req](ctx, r, validator.New())
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	err = c.UrlShortenerService.GetDb().CreateUser(ctx, db.CreateUserParams{
		Email:    request.Email,
		Password: request.Password,
	})
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

	request, err := httphelper.GetRequestBodyAndValidate[req](ctx, r, validator.New())
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	user, err := c.UrlShortenerService.GetDb().GetUserByEmail(ctx, request.Email)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	if user.Password != request.Password {
		reqCtx.SetError(fmt.Errorf("Wrong Password!"))
		return
	}

	accessToken, _ := jwt.NewAccessToken(strconv.FormatInt(user.ID, 10))
	refreshToken, _ := jwt.NewRefreshToken(strconv.FormatInt(user.ID, 10))

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

	if reqCtx.UserId == nil {
		reqCtx.SetError(fmt.Errorf("user id not found"))
		return
	}

	accessToken, _ := jwt.NewAccessToken(strconv.FormatInt(*reqCtx.UserId, 10))
	refreshToken, _ := jwt.NewRefreshToken(strconv.FormatInt(*reqCtx.UserId, 10))

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
		Id    int64  `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	if reqCtx.UserId == nil {
		reqCtx.SetError(fmt.Errorf("user id not found"))
		return
	}

	user, err := c.UrlShortenerService.GetDb().GetUser(ctx, *reqCtx.UserId)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp{
		Id:    user.ID,
		Email: user.Email,
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

	token, err := request.BearerExtractor{}.ExtractToken(r)
	if err == nil {
		jwtToken, err := jwt.GetValidAccessToken(token)
		if err == nil {
			subject, err := jwtToken.Claims.GetSubject()
			if err == nil {
				parseInt, err := strconv.ParseInt(subject, 10, 64)
				if err == nil {
					reqCtx.SetUserId(parseInt)
				}
			}
		}
	}

	reqq, err := httphelper.GetRequestBodyAndValidate[req](ctx, r, validator.New())
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	createUrlParams := db.CreateUrlParams{
		Title:   reqq.Title,
		FullUrl: reqq.FullUrl,
		Qr:      reqq.Qr,
	}

	if reqCtx.UserId != nil {
		createUrlParams.UserID = pgtype.Int8{Int64: *reqCtx.UserId, Valid: true}
	}

	if reqq.CustomUrl != "" {
		createUrlParams.ShortUrl = reqq.CustomUrl
	} else {
		createUrlParams.ShortUrl = uuid.New().String()
	}

	createUrl, err := c.UrlShortenerService.GetDb().CreateUrl(ctx, createUrlParams)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp{Id: createUrl.ID, ShortUrl: createUrl.ShortUrl})
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	_, err = w.Write(marshal)
	reqCtx.SetError(err)
}

type url struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	ShortUrl  string    `json:"short_url"`
	FullUrl   string    `json:"full_url"`
	Qr        string    `json:"qr"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *UrlShortener) GetUrl(w http.ResponseWriter, r *http.Request) {
	type devices struct {
		Device string `json:"device"`
		Count  int    `json:"count"`
	}

	type countries struct {
		Country string `json:"country"`
		Count   int    `json:"count"`
	}

	type resp struct {
		Url         url         `json:"url"`
		TotalClicks int         `json:"total_clicks"`
		Devices     []devices   `json:"devices"`
		Countries   []countries `json:"countries"`
	}

	urlId := r.PathValue("id")
	fmt.Println("url_id", urlId)

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	parseInt, err := strconv.ParseInt(urlId, 10, 64)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	getUrl, err := c.UrlShortenerService.GetDb().GetUrl(ctx, parseInt)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	redirects, err := c.UrlShortenerService.GetDb().GetRedirectsByUrlId(ctx, pgtype.Int8{Int64: parseInt, Valid: true})
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	deviceMap := make(map[string]int)
	countryMap := make(map[string]int)
	for _, redirect := range redirects {
		if redirect.Device != "" {
			deviceMap[redirect.Device]++
		}
		if redirect.Country != "" {
			countryMap[redirect.Country]++
		}
	}

	unsortedDevices := make([]devices, 0)
	unsortedCountries := make([]countries, 0)
	for key, value := range deviceMap {
		unsortedDevices = append(unsortedDevices, devices{Device: key, Count: value})
	}

	for key, value := range countryMap {
		unsortedCountries = append(unsortedCountries, countries{Country: key, Count: value})
	}

	slices.SortFunc(unsortedDevices, func(a, b devices) int {
		ac, ab := a.Count, b.Count
		if ac > ab {
			return -1
		}

		if ac == ab {
			return 0
		}

		return 1
	})

	slices.SortFunc(unsortedCountries, func(a, b countries) int {
		ac, ab := a.Count, b.Count
		if ac > ab {
			return -1
		}

		if ac == ab {
			return 0
		}

		return 1
	})

	// single
	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp{
		Url: url{
			Id:        getUrl.ID,
			Title:     getUrl.Title,
			ShortUrl:  getUrl.ShortUrl,
			FullUrl:   getUrl.FullUrl,
			Qr:        getUrl.Qr,
			CreatedAt: getUrl.CreatedAt.Time,
		},
		Devices:     unsortedDevices[:min(5, len(unsortedDevices))],
		Countries:   unsortedCountries[:min(5, len(unsortedDevices))],
		TotalClicks: len(redirects),
	})
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
		TotalClicks int64 `json:"total_clicks"`
	}

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	token, err := request.BearerExtractor{}.ExtractToken(r)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	jwtToken, err := jwt.GetValidAccessToken(token)
	if err != nil {
		reqCtx.SetError(&errorhelper.AuthError{Err: fmt.Errorf("invalid access token")})
		return
	}

	subject, err := jwtToken.Claims.GetSubject()
	if err != nil {
		reqCtx.SetError(&errorhelper.AuthError{Err: fmt.Errorf("invalid subject")})
		return
	}

	parseInt, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		reqCtx.SetError(&errorhelper.AuthError{Err: fmt.Errorf("failed to parse subject")})
		return
	}

	reqCtx.SetUserId(parseInt)

	if reqCtx.UserId == nil {
		reqCtx.SetError(fmt.Errorf("user id not found"))
		return
	}

	urls, err := c.UrlShortenerService.GetDb().GetUrlsByUserId(ctx, pgtype.Int8{Int64: *reqCtx.UserId, Valid: true})
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	totalClicks, err := c.UrlShortenerService.GetDb().GetUserTotalClicks(ctx, pgtype.Int8{Int64: *reqCtx.UserId, Valid: true})

	response := resp{
		Urls:        []url{},
		TotalClicks: totalClicks,
	}

	for _, each := range urls {
		response.Urls = append(response.Urls, url{
			Id:        each.ID,
			Title:     each.Title,
			ShortUrl:  each.ShortUrl,
			FullUrl:   each.FullUrl,
			Qr:        each.Qr,
			CreatedAt: each.CreatedAt.Time,
		})
	}

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
	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)
	urlId := r.PathValue("id")

	parseInt, err := strconv.ParseInt(urlId, 10, 64)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	err = c.UrlShortenerService.GetDb().DeleteUrl(ctx, parseInt)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *UrlShortener) Redirect(w http.ResponseWriter, r *http.Request) {
	type req struct {
		City    string `json:"city"`
		Country string `json:"country"`
		Device  string `json:"device"`
	}

	type resp struct {
		FullUrl string `json:"full_url"`
	}

	ctx := r.Context()
	reqCtx := reqctx.GetValue(ctx)

	alias := r.PathValue("alias")
	fmt.Println("alias", alias)

	request, err := httphelper.GetRequestBodyAndValidate[req](ctx, r, validator.New())
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	getUrl, err := c.UrlShortenerService.GetDb().GetUrlByShortUrl(ctx, alias)
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	err = c.UrlShortenerService.GetDb().CreateRedirect(ctx, db.CreateRedirectParams{
		UrlID:   pgtype.Int8{Int64: getUrl.ID, Valid: true},
		Device:  request.Device,
		Country: request.Country,
		City:    request.City,
	})
	if err != nil {
		log.Error().Err(err).Msg("CreateRedirect")
	}

	w.WriteHeader(http.StatusOK)
	marshal, err := json.Marshal(resp{FullUrl: getUrl.FullUrl})
	if err != nil {
		reqCtx.SetError(err)
		return
	}

	_, err = w.Write(marshal)
	reqCtx.SetError(err)
}

//// Shorten example
////
////	@Summary		Shorten URL
////	@Description	Given a URL and custom settings, shorten it with a unique alias
////	@ID				shorten-url
////	@Accept			json
////	@Produce		json
////	@Param			req		body		dto.ShortenUrlRequest	true	"Shorten URL Request"
////	@Success		200		{string}	string					"Ok"
////	@Failure		422		{object}	dto.ErrorResponse		"Validation Error"
////	@Failure		500		{object}	dto.ErrorResponse		"Internal Server Error"
////	@Router			/url [post]
//func (c *UrlShortener) Shorten(w http.ResponseWriter, r *http.Request) {
//	ctx := r.Context()
//	reqCtx := reqctx.GetValue(ctx)
//
//	req, err := httphelper.GetRequestBodyAndValidate[dto.ShortenUrlRequest](ctx, r, validator.New())
//	if err != nil {
//		reqCtx.SetError(err)
//		return
//	}
//
//	url, err := c.UrlShortenerService.ShortenUrl(ctx, req)
//	if err != nil {
//		reqCtx.SetError(err)
//		return
//	}
//
//	w.WriteHeader(http.StatusCreated)
//	marshal, err := json.Marshal(url)
//	if err != nil {
//		reqCtx.SetError(err)
//		return
//	}
//
//	_, err = w.Write(marshal)
//	reqCtx.SetError(err)
//}

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
//func (c *UrlShortener) RedirectV2(w http.ResponseWriter, r *http.Request) {
//	ctx := r.Context()
//	reqCtx := reqctx.GetValue(ctx)
//
//	greetings := r.PathValue("alias")
//	url, err := c.UrlShortenerService.GetUrlWithShortened(ctx, greetings)
//	if url == "" || err != nil {
//		reqCtx.SetError(err)
//		return
//	}
//
//	if !strings.Contains(url, "://") { // TODO: Fix protocol not being added to full url
//		url = "http://" + url
//	}
//
//	http.Redirect(w, r, url, http.StatusSeeOther)
//}
