package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raythx98/gohelpme/middleware"
	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/url-shortener/controller"
	"github.com/raythx98/url-shortener/service"
	"github.com/raythx98/url-shortener/sqlc/url_mappings"
	"github.com/raythx98/url-shortener/tools/zerologger"
	"net/http"
	"os"
)

func main() {
	validate := createTools()

	log := zerologger.New(true)
	middleware.RegisterLogger(log)

	ctx := context.Background()
	dbPool := createDb(ctx, log)
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

func createDb(ctx context.Context, log logger.ILogger) *pgxpool.Pool {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable pool_max_conns=10",
		"postgres", "password", "localhost", "5432", "url_shortener")
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	config.ConnConfig.Tracer = &MultiQueryTracer{
		Tracers: []pgx.QueryTracer{
			// TODO: tracer

			// logger
			&myQueryTracer{
				log: log,
			},
		},
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return pool
}

type MultiQueryTracer struct {
	Tracers []pgx.QueryTracer
}

func (m *MultiQueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, t := range m.Tracers {
		ctx = t.TraceQueryStart(ctx, conn, data)
	}

	return ctx
}

func (m *MultiQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, t := range m.Tracers {
		t.TraceQueryEnd(ctx, conn, data)
	}
}

type myQueryTracer struct {
	log logger.ILogger
}

func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	tracer.log.Info(ctx, "[begin-sql]",
		logger.WithField("sql", data.SQL),
		logger.WithField("args", data.Args))

	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	tracer.log.Info(ctx, "[end-sql]",
		logger.WithField("sql error", data.Err),
		logger.WithField("command tag", data.CommandTag))
}

func createTools() *validator.Validate {
	return validator.New()
}
