package main

import (
	"github.com/raythx98/url-shortener/controller"
	"github.com/raythx98/url-shortener/service"
	"log"
	"net/http"

	"github.com/raythx98/gohelpme/middleware"
)

func main() {
	// converting our handler function to handler
	// type to make use of our middleware
	mux := http.NewServeMux()

	urlShortenerSvc := registerServices(&service.UrlShortener{})
	urlShortener := registerControllers(&controller.UrlShortener{UrlShortenerService: urlShortenerSvc})

	mux.Handle("/api/v1/url/redirect",
		middleware.JsonResponse(middleware.AddRequestId(middleware.Log(http.HandlerFunc(urlShortener.Redirect)))))
	mux.Handle("/api/v1/url",
		middleware.JsonResponse(middleware.AddRequestId(middleware.Log(http.HandlerFunc(urlShortener.Shorten)))))

	err := http.ListenAndServe(":5051", mux)
	log.Fatal(err)
}

func registerServices(urlShortener service.IUrlShortener) service.IUrlShortener {
	return urlShortener
}

func registerControllers(urlShortener controller.IUrlShortener) controller.IUrlShortener {
	return urlShortener
}
