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
	mux.Handle("/metrics", cfg.HandlerMetrics())
	mux.Handle("/reset", cfg.HandlerMetricsReset())
	mux.HandleFunc("/healthz", handlers.HandlerHealth())

	server.ListenAndServe()
}
