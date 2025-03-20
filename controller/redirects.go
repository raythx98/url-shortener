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

type IRedirects interface {
	Redirect(w http.ResponseWriter, r *http.Request)
}

type Redirects struct {
	RedirectsService service.IRedirects
	Validator        validator.IValidator
	Log              logger.ILogger
}

func NewRedirects(service service.IRedirects, validate validator.IValidator, log logger.ILogger) *Redirects {
	return &Redirects{
		RedirectsService: service,
		Validator:        validate,
		Log:              log,
	}
}

// Redirect
// @summary 	Redirect
// @description Redirect a short link to the full URL
// @tags 		Redirects
// @param       shortLink	path 		string 						true	"Short link"
// @param       request		body 		dto.RedirectRequest 		true	"Redirect Request"
// @success     200			{object}  	dto.RedirectResponse    	 		"ok"
// @failure     400			{object}  	errorhelper.ErrorResponse    		"bad request"
// @failure     401			{object}  	errorhelper.ErrorResponse    		"unauthorized"
// @failure     422			{object}  	errorhelper.ErrorResponse    		"validation error"
// @response    500         {object}  	errorhelper.ErrorResponse    		"server error"
// @security 	BasicAuth
// @router 		/urls/v1/redirect/{shortLink} [post]
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
