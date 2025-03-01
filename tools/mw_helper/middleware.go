package mw_helper

import (
	"net/http"

	"github.com/raythx98/url-shortener/resources"

	"github.com/raythx98/gohelpme/middleware"
)

type Middlewares []func(next http.HandlerFunc) http.HandlerFunc

type Template struct {
	Default                Middlewares
	DefaultAccessToken     Middlewares
	DefaultRefreshToken    Middlewares
	DefaultBasicToken      Middlewares
	DefaultJwtOrBasicToken Middlewares
}

func GetMiddlewares(tools resources.Tools) Template {
	defaultMiddlewares := Middlewares{
		middleware.CORS,
		middleware.JsonResponse,
		middleware.AddRequestId,
		middleware.ReqCtx,
		middleware.JwtSubject(tools.Jwt),
		middleware.Log(tools.Log),
		middleware.Recoverer(),
		middleware.ErrorHandler,
	}
	return Template{
		Default:                defaultMiddlewares,
		DefaultAccessToken:     append(defaultMiddlewares, middleware.JwtAuth(tools.Jwt, middleware.AccessToken)),
		DefaultRefreshToken:    append(defaultMiddlewares, middleware.JwtAuth(tools.Jwt, middleware.RefreshToken)),
		DefaultBasicToken:      append(defaultMiddlewares, middleware.BasicAuth(tools.BasicAuth)),
		DefaultJwtOrBasicToken: append(defaultMiddlewares, middleware.JwtOrBasicAuth(tools.BasicAuth, tools.Jwt, middleware.AccessToken)),
	}
}
