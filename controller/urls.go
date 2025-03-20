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

// GetUrl
// @summary 	GetUrl
// @description Get Url details
// @tags 		Urls
// @param       id		path 		string 						true	"Url Id"
// @success     200		{object}	dto.GetUrlResponse    		 		"ok"
// @failure     400    	{object}  	errorhelper.ErrorResponse    		"bad request"
// @failure     401  	{object}  	errorhelper.ErrorResponse    		"unauthorized"
// @failure     422    	{object}  	errorhelper.ErrorResponse    		"validation error"
// @response    500   	{object}  	errorhelper.ErrorResponse    		"server error"
// @security 	BearerAuth
// @router 		/urls/v1/{id} [get]
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

// GetUrls
// @summary 	GetUrls
// @description Get User's Urls
// @tags 		Urls
// @param 		Authorization	header 		string 						true	"JWT token use `Bearer <token>`"
// @success     200          	{object}	dto.GetUrlsResponse    		 		"ok"
// @failure     400          	{object}  	errorhelper.ErrorResponse    		"bad request"
// @failure     401           	{object}  	errorhelper.ErrorResponse    		"unauthorized"
// @failure     422           	{object}  	errorhelper.ErrorResponse    		"validation error"
// @response    500          	{object}  	errorhelper.ErrorResponse    		"server error"
// @security 	BearerAuth
// @router 		/urls/v1 [get]
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

// CreateUrl
// @summary 	CreateUrl
// @description Create Shortened Urls
// @tags 		Urls
// @param       request 	body 		dto.CreateUrlRequest 		true	"Create Url Request"
// @success     200  		{object}  	dto.CreateUrlResponse        		"ok"
// @failure     400     	{object}  	errorhelper.ErrorResponse    		"bad request"
// @failure     401        	{object}  	errorhelper.ErrorResponse    		"unauthorized"
// @failure     422  		{object}  	errorhelper.ErrorResponse    		"validation error"
// @response    500        	{object}  	errorhelper.ErrorResponse    		"server error"
// @security 	BasicAuth
// @security 	BearerAuth
// @router 		/urls/v1 [post]
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

// DeleteUrl
// @summary 	DeleteUrl
// @description Delete Url
// @tags 		Urls
// @param       id		path 		string 						true	"Url Id"
// @success     200									    		 		"ok"
// @failure     400    	{object}  	errorhelper.ErrorResponse    		"bad request"
// @failure     401  	{object}  	errorhelper.ErrorResponse    		"unauthorized"
// @failure     422    	{object}  	errorhelper.ErrorResponse    		"validation error"
// @response    500   	{object}  	errorhelper.ErrorResponse    		"server error"
// @security 	BearerAuth
// @router 		/urls/v1/{id} [delete]
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
