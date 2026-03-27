package main

import (
	"net/http"

	_ "github.com/lib/pq"
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
	mux.Handle("GET /admin/metrics", cfg.HandlerMetrics())
	mux.Handle("POST /admin/reset", cfg.HandlerMetricsReset())
	mux.Handle("GET /api/healthz", handlers.HandlerHealth())
	mux.HandleFunc("POST /api/validate_chirp", handlers.HandlerValidateChirpy)

	server.ListenAndServe()
}
