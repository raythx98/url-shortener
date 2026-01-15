package mw_helper

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/raythx98/url-shortener/resources"

	"github.com/raythx98/gohelpme/middleware"
	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/reqctx"
	"gopkg.in/yaml.v3"
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
	rlConfig := loadRateLimitConfig(tools.Log)
	rateLimiter := middleware.NewRateLimiter(
		rlConfig,
		tools.Log,
		func(r *http.Request) (string, string) {
			identifier := middleware.ExtractIP(r)
			if reqCtx := reqctx.GetValue(r.Context()); reqCtx != nil && reqCtx.UserId != nil {
				identifier = fmt.Sprintf("user:%d", *reqCtx.UserId)
			}
			return identifier, fmt.Sprintf("%s:%s", r.Method, r.URL.Path)
		},
	)

	defaultMiddlewares := Middlewares{
		middleware.CORS,
		middleware.JsonResponse,
		middleware.AddRequestId,
		middleware.ReqCtx,
		middleware.JwtSubject(tools.Jwt),
		middleware.Log(tools.Log),
		rateLimiter.RateLimit,
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

func loadRateLimitConfig(log logger.ILogger) middleware.Config {
	ctx := context.Background()
	var cfg middleware.Config
	data, err := os.ReadFile("ratelimit.yaml")
	if err != nil {
		log.Warn(ctx, "failed to read ratelimit.yaml, using defaults", logger.WithError(err))
		return middleware.Config{
			Default: middleware.RateConfig{Rate: 1, Burst: 1},
		}
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Warn(ctx, "failed to unmarshal ratelimit.yaml, using defaults", logger.WithError(err))
		return middleware.Config{
			Default: middleware.RateConfig{Rate: 1, Burst: 1},
		}
	}

	return cfg
}
