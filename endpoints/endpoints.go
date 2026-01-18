package endpoints

import (
	"net/http"

	"github.com/raythx98/go-zap/resources"
	"github.com/raythx98/go-zap/tools/mw_helper"

	"github.com/raythx98/gohelpme/middleware"
)

func Register(mux *http.ServeMux, ctrls resources.Controllers, tools resources.Tools) {
	middlewares := mw_helper.GetMiddlewares(tools)

	mux.HandleFunc("OPTIONS /api/", middleware.Chain(ctrls.Auth.Register, middleware.CORS))

	mux.HandleFunc("POST /api/auth/v1/register", middleware.Chain(ctrls.Auth.Register, middlewares.DefaultBasicToken...))

	mux.HandleFunc("POST /api/auth/v1/login", middleware.Chain(ctrls.Auth.Login, middlewares.DefaultBasicToken...))
	mux.HandleFunc("POST /api/auth/v1/refresh", middleware.Chain(ctrls.Auth.Refresh, middlewares.DefaultRefreshToken...))
	mux.HandleFunc("POST /api/auth/v1/logout", middleware.Chain(ctrls.Auth.Logout, middlewares.DefaultAccessToken...))

	mux.HandleFunc("GET /api/urls/v1/{id}", middleware.Chain(ctrls.Urls.GetUrl, middlewares.DefaultAccessToken...))
	mux.HandleFunc("DELETE /api/urls/v1/{id}", middleware.Chain(ctrls.Urls.DeleteUrl, middlewares.DefaultAccessToken...))
	mux.HandleFunc("POST /api/urls/v1", middleware.Chain(ctrls.Urls.CreateUrl, middlewares.DefaultJwtOrBasicToken...))
	mux.HandleFunc("GET /api/urls/v1", middleware.Chain(ctrls.Urls.GetUrls, middlewares.DefaultAccessToken...))

	mux.HandleFunc("POST /api/urls/v1/redirect/{shortLink}", middleware.Chain(ctrls.Redirects.Redirect, middlewares.DefaultBasicToken...))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "API not found", http.StatusNotFound)
	})
}
