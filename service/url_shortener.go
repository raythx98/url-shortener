package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/raythx98/url-shortener/sqlc/url_mappings"
	"github.com/raythx98/url-shortener/tools/sql_tool"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

type IUrlShortener interface {
	ShortenUrl(ctx context.Context, url string) (string, error)
	GetUrlWithShortened(ctx context.Context, shortenedUrl string) (string, error)
}

type UrlShortener struct {
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

func (s *UrlShortener) ShortenUrl(ctx context.Context, url string) (string, error) {
	// TODO: connect on startup
	config, err := pgxpool.ParseConfig("user=postgres password=password host=localhost port=5432 dbname=url_shortener sslmode=disable pool_max_conns=10")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	config.ConnConfig.Tracer = &myQueryTracer{
		log: log.Default(),
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if err != nil {
		return "", err
	}
	defer conn.Close()

	queries := url_mappings.New(conn)

	expireAt := time.Now().UTC().AddDate(1, 0, 0)
	params := url_mappings.CreateUrlMappingParams{
		ShortenedUrl:     uuid.NewString(),
		Url:              url,
		InactiveExpireAt: sql_tool.NewTime(nil),
		MustExpireAt:     sql_tool.NewTime(&expireAt),
	}
	mapping, err := queries.CreateUrlMapping(ctx, params)
	if err != nil {
		return "", err
	}

	return mapping.ShortenedUrl, nil
}

func (s *UrlShortener) GetUrlWithShortened(ctx context.Context, shortenedUrl string) (string, error) {
	conn, err := pgxpool.New(ctx, "user=postgres password=password host=localhost port=5432 dbname=url_shortener sslmode=disable") // TODO: connect on startup
	if err != nil {
		return "", err
	}
	defer conn.Close()

	queries := url_mappings.New(conn)

	mapping, err := queries.GetUrlMapping(ctx, shortenedUrl)
	if err != nil {
		return "", err
	}

	return mapping.Url, nil
}
