package main

import (
	"net/http"

	"github.com/neuroshepherd/learn-http-servers/handlers"
)

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	cfg := &handlers.APIConfig{}
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.Handle("GET /api/metrics", cfg.HandlerMetrics())
	mux.Handle("POST /api/reset", cfg.HandlerMetricsReset())
	mux.HandleFunc("GET /api/healthz", handlers.HandlerHealth())

	server.ListenAndServe()
}
