package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/raythx98/url-shortener/configs"
	"github.com/raythx98/url-shortener/controller"
	_ "github.com/raythx98/url-shortener/docs"
	"github.com/raythx98/url-shortener/service"
	"github.com/raythx98/url-shortener/sqlc/url_mappings"
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

	defaultMiddlewares := []func(next http.Handler) http.Handler{
		middleware.CORS,
		middleware.JsonResponse,
		middleware.AddRequestId,
		middleware.Log(log),
		middleware.Recoverer(log),
		middleware.ErrorHandler,
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/url/redirect/{alias}", middleware.Chain(urlShortener.Redirect, defaultMiddlewares...))
	mux.Handle("/api/v1/url", middleware.Chain(urlShortener.Shorten, defaultMiddlewares...))
	mux.Handle("/api/v1/test", middleware.Chain(urlShortener.Test, defaultMiddlewares...))
	mux.Handle("/swagger/*", httpSwagger.Handler(httpSwagger.URL(
		fmt.Sprintf("http://localhost:%d/swagger/doc.json", config.ServerPort))))

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.ServerPort), mux)
	if err != nil {
		log.Fatal(ctx, "failed to listen and serve", logger.WithError(err))
	}

	log.Info(ctx, "Server stopped")
}

func registerRepos(pool *pgxpool.Pool) *url_mappings.Queries {
	return url_mappings.New(pool)
}

func registerControllers(urlShortenerSvc *service.UrlShortener, v *validator.Validate) *controller.UrlShortener {
	urlShortener := controller.New(urlShortenerSvc, v)
	func(controller.IUrlShortener) {}(urlShortener)
	return urlShortener
}

func registerServices(urlMappingRepo *url_mappings.Queries) *service.UrlShortener {
	urlShortenerSvc := service.New(urlMappingRepo)
	func(service.IUrlShortener) {}(urlShortenerSvc)
	return urlShortenerSvc
}

func createTools() *validator.Validate {
	return validator.New()
}
