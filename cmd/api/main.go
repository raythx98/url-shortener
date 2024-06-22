package main

import (
	"github.com/raythx98/gohelpme/middleware"
	"github.com/raythx98/url-shortener/controller"
	"github.com/raythx98/url-shortener/service"
	"log"
	"net/http"
)

func main() {
	// converting our handler function to handler
	// type to make use of our middleware
	mux := http.NewServeMux()

	urlShortenerSvc := registerServices(&service.UrlShortener{})
	urlShortener := registerControllers(&controller.UrlShortener{UrlShortenerService: urlShortenerSvc})

	defaultMiddlewares := []func(next http.Handler) http.Handler{
		middleware.JsonResponse,
		middleware.AddRequestId,
		middleware.Log,
	}

	mux.Handle("/api/v1/url/redirect/{alias}", middleware.Chain(urlShortener.Redirect, defaultMiddlewares...))
	mux.Handle("/api/v1/url", middleware.Chain(urlShortener.Shorten, defaultMiddlewares...))

	err := http.ListenAndServe(":5051", mux)
	log.Fatal(err)
}

func registerServices(urlShortener service.IUrlShortener) service.IUrlShortener {
	return urlShortener
}

func registerControllers(urlShortener controller.IUrlShortener) controller.IUrlShortener {
	return urlShortener
}
