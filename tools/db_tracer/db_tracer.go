package db_tracer

import (
	"context"

	"github.com/raythx98/gohelpme/tool/logger"

	"github.com/jackc/pgx/v5"
)

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

type MyQueryTracer struct {
	Log logger.ILogger
}

func (tracer *MyQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	tracer.Log.Info(ctx, "[begin-sql]",
		logger.WithField("sql", data.SQL),
		logger.WithField("args", data.Args))

	return ctx
}

func (tracer *MyQueryTracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	tracer.Log.Info(ctx, "[end-sql]",
		logger.WithField("sql error", data.Err),
		logger.WithField("command tag", data.CommandTag))
}
