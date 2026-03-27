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

	cfg := &handlers.APIConfig{DB: dbQueries}
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.Handle("GET /admin/metrics", cfg.HandlerMetrics())
	mux.HandleFunc("POST /admin/reset", cfg.HandlerReset)
	mux.Handle("GET /api/healthz", handlers.HandlerHealth())
	mux.HandleFunc("POST /api/validate_chirp", handlers.HandlerValidateChirpy)
	mux.HandleFunc("POST /api/users", cfg.HandlerCreateUser)

	server.ListenAndServe()
}
