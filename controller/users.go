package controller

import (
	"encoding/json"
	"net/http"

	"github.com/raythx98/url-shortener/service"

	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"
	"github.com/raythx98/gohelpme/tool/validator"
)

type IUsers interface {
	GetProfile(w http.ResponseWriter, r *http.Request)
}

type Users struct {
	UsersService service.IUsers
	Validator    validator.IValidator
	Log          logger.ILogger
}

func NewUsers(service service.IUsers, validate validator.IValidator, log logger.ILogger) *Users {
	return &Users{
		UsersService: service,
		Validator:    validate,
		Log:          log,
	}
}

// GetProfile
// @summary 	GetProfile
// @description Get User's Profile
// @tags 		Users
// @param 		Authorization	header 		string 						true	"JWT token use `Bearer <token>`"
// @success     200             {object}	dto.ProfileResponse    		 		"ok"
// @failure     400             {object}  	errorhelper.ErrorResponse    		"bad request"
// @failure     401             {object}  	errorhelper.ErrorResponse    		"unauthorized"
// @failure     422             {object}  	errorhelper.ErrorResponse    		"validation error"
// @response    500          	{object}  	errorhelper.ErrorResponse    		"server error"
// @security 	BearerAuth
// @router 		/users/v1 [get]
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
