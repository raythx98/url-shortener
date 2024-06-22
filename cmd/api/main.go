package main

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raythx98/gohelpme/middleware"
	"github.com/raythx98/url-shortener/controller"
	"github.com/raythx98/url-shortener/service"
	"log"
	"net/http"
	"os"
)

func main() {
	validate := createTools()

	dbPool := createDb()
	defer dbPool.Close()

	urlShortenerSvc := registerServices(dbPool)
	urlShortener := registerControllers(urlShortenerSvc, validate)

	defaultMiddlewares := []func(next http.Handler) http.Handler{
		middleware.JsonResponse,
		middleware.AddRequestId,
		middleware.Log,
	}

	mux := http.NewServeMux()

	mux.Handle("/api/v1/url/redirect/{alias}", middleware.Chain(urlShortener.Redirect, defaultMiddlewares...))
	mux.Handle("/api/v1/url", middleware.Chain(urlShortener.Shorten, defaultMiddlewares...))

	err := http.ListenAndServe(":5051", mux)
	log.Fatal(err)
}

func registerControllers(urlShortenerSvc *service.UrlShortener, v *validator.Validate) *controller.UrlShortener {
	urlShortener := &controller.UrlShortener{UrlShortenerService: urlShortenerSvc, Validator: v}
	func(controller.IUrlShortener) {}(urlShortener)
	return urlShortener
}

func registerServices(pool *pgxpool.Pool) *service.UrlShortener {
	urlShortenerSvc := &service.UrlShortener{DbPool: pool}
	func(service.IUrlShortener) {}(urlShortenerSvc)
	return urlShortenerSvc
}

func createDb() *pgxpool.Pool {
	config, err := pgxpool.ParseConfig("user=postgres password=password host=localhost port=5432 dbname=url_shortener sslmode=disable pool_max_conns=10")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	config.ConnConfig.Tracer = &myQueryTracer{
		log: log.Default(),
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return pool
}

type myQueryTracer struct {
	log *log.Logger
}

func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	tracer.log.Println("Executing command", "sql", data.SQL, "args", data.Args)

	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func createTools() *validator.Validate {
	return validator.New()
}
