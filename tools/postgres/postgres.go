package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/raythx98/url-shortener/configs"
	"github.com/raythx98/url-shortener/tools/db_tracer"

	"github.com/raythx98/gohelpme/tool/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg configs.IConfig, log logger.ILogger) *pgxpool.Pool {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable pool_max_conns=10",
		cfg.GetDbUsername(), cfg.GetDbPassword(), cfg.GetDbHost(), cfg.GetDbPort(), cfg.GetDbDefaultName())
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	config.ConnConfig.Tracer = &db_tracer.MultiQueryTracer{
		Tracers: []pgx.QueryTracer{
			// TODO: add tracer

			// logger
			&db_tracer.MyQueryTracer{
				Log: log,
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
