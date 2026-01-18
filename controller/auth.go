package controller

import (
	"encoding/json"
	"net/http"

	"github.com/raythx98/go-zap/dto"
	"github.com/raythx98/go-zap/service"

	"github.com/raythx98/gohelpme/tool/httphelper"
	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"
	"github.com/raythx98/gohelpme/tool/validator"
)

type IAuth interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type Auth struct {
	AuthService service.IAuth
	Validator   validator.IValidator
	Log         logger.ILogger
}

func NewAuth(service service.IAuth, validate validator.IValidator, log logger.ILogger) *Auth {
	return &Auth{
		AuthService: service,
		Validator:   validate,
		Log:         log,
	}
}

// Register
// @summary 	Register
// @description Register a new user with account information
// @tags 		Auth
// @param       request		body		dto.RegisterRequest 		true	"Register Request"
// @Success     200   		{object}	dto.LoginResponse    				"ok"
// @failure     400        	{object}  	errorhelper.ErrorResponse    		"bad request"
// @failure     401       	{object}  	errorhelper.ErrorResponse    		"unauthorized"
// @failure     422       	{object}  	errorhelper.ErrorResponse    		"validation error"
// @response    500       	{object}  	errorhelper.ErrorResponse    		"server error"
// @security	BasicAuth
// @router 		/auth/v1/register [post]
func (c *Auth) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		reqctx.GetValue(ctx).SetError(err)
	}()

	req, err := httphelper.GetRequestBodyAndValidate[dto.RegisterRequest](ctx, r, c.Validator)
	if err != nil {
		return
	}

	resp, err := c.AuthService.Register(ctx, req)
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

// Login
// @summary 	Login
// @description Login with email and password
// @tags 		Auth
// @param		request		body 		dto.LoginRequest 			true	"LoginRequest Request"
// @success     200       	{object}  	dto.LoginResponse    				"ok"
// @failure     400       	{object}  	errorhelper.ErrorResponse    		"bad request"
// @failure     401        	{object}  	errorhelper.ErrorResponse    		"unauthorized"
// @failure     422      	{object}  	errorhelper.ErrorResponse    		"validation error"
// @response    500    		{object}  	errorhelper.ErrorResponse    		"server error"
// @security 	BasicAuth
// @router 		/auth/v1/login [post]
func (c *Auth) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		reqctx.GetValue(ctx).SetError(err)
	}()

	req, err := httphelper.GetRequestBodyAndValidate[dto.LoginRequest](ctx, r, c.Validator)
	if err != nil {
		return
	}

	resp, err := c.AuthService.Login(ctx, req)
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

// Refresh
// @summary 	Refresh
// @description Refresh session with refresh token
// @tags 		Auth
// @param 		Authorization	header 		string 						true	"JWT token use `Bearer <token>`"
// @success     200             {object}	dto.LoginResponse    				"ok"
// @failure     400             {object}  	errorhelper.ErrorResponse			"bad request"
// @failure     401             {object}  	errorhelper.ErrorResponse			"unauthorized"
// @failure     422             {object}  	errorhelper.ErrorResponse			"validation error"
// @response    500          	{object}  	errorhelper.ErrorResponse			"server error"
// @security 	BearerAuth
// @router 		/auth/v1/refresh [post]
func (c *Auth) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		reqctx.GetValue(ctx).SetError(err)
	}()

	resp, err := c.AuthService.Refresh(ctx)
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

// Logout
// @summary 	Logout
// @description Logout session
// @tags 		Auth
// @param 		Authorization	header 		string 						true	"JWT token use `Bearer <token>`"
// @success     200              										 		"ok"
// @failure     400         	{object}	errorhelper.ErrorResponse    		"bad request"
// @failure     401        		{object}	errorhelper.ErrorResponse    		"unauthorized"
// @failure     422          	{object}  	errorhelper.ErrorResponse    		"validation error"
// @response    500          	{object}  	errorhelper.ErrorResponse    		"server error"
// @security 	BearerAuth
// @router 		/auth/v1/logout [post]
func (c *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
