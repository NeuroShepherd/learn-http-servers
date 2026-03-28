package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/neuroshepherd/learn-http-servers/handlers"
	"github.com/neuroshepherd/learn-http-servers/internal/database"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	cfg := &handlers.APIConfig{DB: dbQueries, Platform: platform, JWTSecret: jwtSecret, PolkaKey: polkaKey}
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.Handle("GET /admin/metrics", cfg.HandlerMetrics())
	mux.HandleFunc("POST /admin/reset", cfg.HandlerReset)
	mux.Handle("GET /api/healthz", handlers.HandlerHealth())
	// mux.HandleFunc("POST /api/validate_chirp", handlers.HandlerValidateChirpy)
	mux.HandleFunc("POST /api/users", cfg.HandlerCreateUser)
	mux.HandleFunc("POST /api/chirps", cfg.HandlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.HandlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.HandlerGetChirpByID)
	mux.HandleFunc("POST /api/login", cfg.HandlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.HandlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.HandlerRevokeToken)
	mux.HandleFunc("PUT /api/users/", cfg.HandlerUpdateUser)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.HandlerDeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.HandlerUpdateChirpyRedStatus)

	server.ListenAndServe()
}
