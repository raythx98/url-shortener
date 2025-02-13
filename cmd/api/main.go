package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/raythx98/url-shortener/configs"
	"github.com/raythx98/url-shortener/controller"
	_ "github.com/raythx98/url-shortener/docs"
	"github.com/raythx98/url-shortener/service"
	"github.com/raythx98/url-shortener/sqlc/db"
	"github.com/raythx98/url-shortener/tools/mysql"
	"github.com/raythx98/url-shortener/tools/zerologger"

	"github.com/raythx98/gohelpme/middleware"
	"github.com/raythx98/gohelpme/tool/logger"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title           URL Shortener Server
// @version         1.0

// @contact.name   Ray Toh
// @contact.url    https://www.raythx.com
// @contact.email  raythx98@gmail.com

// @host      localhost:5051
// @BasePath  /api/v1
func main() {
	validate := createTools()

	config := configs.Load()

	log := zerologger.New(true)

	ctx := context.Background()

	log.Debug(ctx, "configs loaded", logger.WithField("config", config))

	dbPool := mysql.NewPool(ctx, config, log)
	defer dbPool.Close()

	urlMappingRepo := registerRepos(dbPool)
	urlShortenerSvc := registerServices(urlMappingRepo)
	urlShortener := registerControllers(urlShortenerSvc, validate)

	defaultMiddlewares := []func(next http.HandlerFunc) http.HandlerFunc{
		middleware.CORS,
		middleware.JsonResponse,
		middleware.AddRequestId,
		middleware.Log(log),
		middleware.Recoverer(log),
		middleware.ErrorHandler,
	}

	defaultMiddlewaresAccess := []func(next http.HandlerFunc) http.HandlerFunc{
		middleware.CORS,
		middleware.JsonResponse,
		middleware.AddRequestId,
		middleware.Log(log),
		middleware.Recoverer(log),
		middleware.ErrorHandler,
		middleware.Auth(middleware.Access),
	}

	defaultMiddlewaresRefresh := []func(next http.HandlerFunc) http.HandlerFunc{
		middleware.CORS,
		middleware.JsonResponse,
		middleware.AddRequestId,
		middleware.Log(log),
		middleware.Recoverer(log),
		middleware.ErrorHandler,
		middleware.Auth(middleware.Refresh),
	}

	defaultMiddlewaresBasic := []func(next http.HandlerFunc) http.HandlerFunc{
		middleware.CORS,
		middleware.JsonResponse,
		middleware.AddRequestId,
		middleware.Log(log),
		middleware.Recoverer(log),
		middleware.ErrorHandler,
		middleware.Auth(middleware.Basic),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("OPTIONS /api/", middleware.Chain(urlShortener.Register, middleware.CORS))

	mux.HandleFunc("POST /api/auth/v1/register", middleware.Chain(urlShortener.Register, defaultMiddlewaresBasic...))
	mux.HandleFunc("POST /api/auth/v1/login", middleware.Chain(urlShortener.Login, defaultMiddlewaresBasic...))
	mux.HandleFunc("POST /api/auth/v1/refresh", middleware.Chain(urlShortener.Refresh, defaultMiddlewaresRefresh...))
	mux.HandleFunc("POST /api/auth/v1/logout", middleware.Chain(urlShortener.Logout, defaultMiddlewaresAccess...))

	mux.HandleFunc("GET /api/users/v1", middleware.Chain(urlShortener.Profile, defaultMiddlewaresAccess...))

	mux.HandleFunc("GET /api/urls/v1/{id}", middleware.Chain(urlShortener.GetUrl, defaultMiddlewaresAccess...))
	mux.HandleFunc("DELETE /api/urls/v1/{id}", middleware.Chain(urlShortener.DeleteUrl, defaultMiddlewaresAccess...))
	mux.HandleFunc("POST /api/urls/v1", middleware.Chain(urlShortener.CreateUrl, defaultMiddlewares...))
	mux.HandleFunc("GET /api/urls/v1", middleware.Chain(urlShortener.GetUrls, defaultMiddlewares...))

	mux.HandleFunc("POST /api/urls/v1/redirect/{alias}", middleware.Chain(urlShortener.Redirect, defaultMiddlewares...))

	//mux.Handle("/api/v1/url/redirect/{alias}", middleware.Chain(urlShortener.RedirectV2, defaultMiddlewaresAccess...))
	//mux.Handle("/api/v1/url", middleware.Chain(urlShortener.Shorten, defaultMiddlewaresAccess...))
	mux.HandleFunc("/swagger/*", httpSwagger.Handler(httpSwagger.URL(
		fmt.Sprintf("http://localhost:%d/swagger/doc.json", config.ServerPort))))

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.ServerPort), mux)
	if err != nil {
		log.Fatal(ctx, "failed to listen and serve", logger.WithError(err))
	}

	log.Info(ctx, "Server stopped")
}

func registerRepos(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
}

func registerControllers(urlShortenerSvc *service.UrlShortener, v *validator.Validate) *controller.UrlShortener {
	urlShortener := controller.New(urlShortenerSvc, v)
	func(controller.IUrlShortener) {}(urlShortener)
	return urlShortener
}

func registerServices(urlMappingRepo *db.Queries) *service.UrlShortener {
	urlShortenerSvc := service.New(urlMappingRepo)
	func(service.IUrlShortener) {}(urlShortenerSvc)
	return urlShortenerSvc
}

func createTools() *validator.Validate {
	return validator.New()
}
