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

type IUsers interface {
	Register(w http.ResponseWriter, r *http.Request)
	GetProfile(w http.ResponseWriter, r *http.Request)
}

type Users struct {
	UsersService service.IUsers
	Validator    *validator.Validate
	Log          logger.ILogger
}

func NewUsers(service service.IUsers, validate *validator.Validate, log logger.ILogger) *Users {
	return &Users{
		UsersService: service,
		Validator:    validate,
		Log:          log,
	}
}

func (c *Users) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		reqctx.GetValue(ctx).SetError(err)
	}()

	req, err := httphelper.GetRequestBodyAndValidate[dto.RegisterRequest](ctx, r, c.Validator)
	if err != nil {
		return
	}

	if err = c.UsersService.Register(ctx, req); err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Users) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error
	defer func() {
		reqctx.GetValue(ctx).SetError(err)
	}()

	resp, err := c.UsersService.GetProfile(ctx)
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
