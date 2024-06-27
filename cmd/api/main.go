package main

import (
	"context"
	"fmt"
	"github.com/raythx98/url-shortener/configs"
	"net/http"
	"os"

	"github.com/raythx98/url-shortener/controller"
	"github.com/raythx98/url-shortener/service"
	"github.com/raythx98/url-shortener/sqlc/url_mappings"
	"github.com/raythx98/url-shortener/tools/mysql"
	"github.com/raythx98/url-shortener/tools/zerologger"

	"github.com/raythx98/gohelpme/middleware"
	"github.com/raythx98/gohelpme/tool/logger"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	validate := createTools()

	config := configs.Load()

	log := zerologger.New(true)
	middleware.RegisterLogger(log)

	ctx := context.Background()
	printConfigs(ctx, log, config)

	dbPool := mysql.NewPool(ctx, config, log)
	defer dbPool.Close()

	urlMappingRepo := registerRepos(dbPool)
	urlShortenerSvc := registerServices(urlMappingRepo)
	urlShortener := registerControllers(urlShortenerSvc, validate)

	defaultMiddlewares := []func(next http.Handler) http.Handler{
		middleware.JsonResponse,
		middleware.AddRequestId,
		middleware.Log,
	}

	mux := http.NewServeMux()

	mux.Handle("/api/v1/url/redirect/{alias}", middleware.Chain(urlShortener.Redirect, defaultMiddlewares...))
	mux.Handle("/api/v1/url", middleware.Chain(urlShortener.Shorten, defaultMiddlewares...))

	_ = http.ListenAndServe(":5051", mux)
	os.Exit(1)
	//log.Fatal(err)
}

func printConfigs(ctx context.Context, log logger.ILogger, s *configs.Specification) {
	format := "Debug: %v\nServerPort: %d\nDBUser: %s\nDBPassword: %s\nDbHost: %s\nDbPort: %s\nDbDefaultDb: %s\n"
	fmt.Printf(format, s.Debug, s.ServerPort, s.DbUsername, s.DbPassword, s.DbHost, s.DbPort, s.DbDefaultName)
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
