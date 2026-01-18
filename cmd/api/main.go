package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raythx98/go-zap/docs"
	"github.com/raythx98/go-zap/endpoints"
	"github.com/raythx98/go-zap/resources"
	"github.com/raythx98/go-zap/tools/config"

	"github.com/raythx98/gohelpme/tool/logger"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title           Zap Server
// @version         1.0

// @contact.name   Ray Toh
// @contact.url    https://raythx98.github.io/me/
// @contact.email  raythx98@gmail.com

// @host      raythx.com
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format: Bearer {token}
func main() {
	ctx := context.Background()

	cfg := config.Load()
	fmt.Printf("configs loaded: %+v\n", cfg)

	tool := resources.CreateTools(ctx, cfg)
	defer tool.Db.Pool().Close()

	repo := resources.RegisterRepos(ctx, tool)
	clients, err := resources.RegisterClients(ctx, cfg)
	if err != nil {
		tool.Log.Fatal(ctx, "failed to create clients", logger.WithError(err))
	}
	svcs := resources.RegisterServices(ctx, cfg, repo, clients, tool)
	ctrls := resources.RegisterControllers(ctx, svcs, tool)

	mux := http.NewServeMux()

	endpoints.Register(mux, ctrls, tool)

	var swaggerUrl string
	if cfg.IsDevelopment() {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.ServerPort)
		swaggerUrl = fmt.Sprintf("http://localhost:%d/swagger/doc.json", cfg.ServerPort)
	} else {
		docs.SwaggerInfo.Host = "129.150.49.141.sslip.io"
		swaggerUrl = "https://129.150.49.141.sslip.io/swagger/doc.json"
	}
	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		httpSwagger.Handler(httpSwagger.URL(swaggerUrl))(w, r)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: mux,
	}

	go func() {
		tool.Log.Info(ctx, "Server starting", logger.WithField("port", cfg.ServerPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			tool.Log.Fatal(ctx, "failed to listen and serve", logger.WithError(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	tool.Log.Info(ctx, "Server is shutting down...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		tool.Log.Fatal(ctx, "Server forced to shutdown", logger.WithError(err))
	}

	tool.Log.Info(ctx, "Server stopped")
}
