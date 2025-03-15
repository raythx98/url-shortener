package main

import (
	"context"
	"fmt"
	"net/http"

	_ "github.com/raythx98/url-shortener/docs"
	"github.com/raythx98/url-shortener/endpoints"
	"github.com/raythx98/url-shortener/resources"
	"github.com/raythx98/url-shortener/tools/config"

	"github.com/raythx98/gohelpme/tool/logger"

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

	mux.HandleFunc("/swagger/*", httpSwagger.Handler(httpSwagger.URL(
		fmt.Sprintf("http://localhost:%d/swagger/doc.json", cfg.ServerPort))))

	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), mux)
	if err != nil {
		tool.Log.Fatal(ctx, "failed to listen and serve", logger.WithError(err))
	}

	tool.Log.Info(ctx, "Server stopped")
}
