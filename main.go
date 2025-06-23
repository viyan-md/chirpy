package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/viyan-md/chirpy/internal/api"
	"github.com/viyan-md/chirpy/internal/database"
)

func main() {
	godotenv.Load()

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	const root = "assets"
	const port = "8080"

	apiCfg := api.APIConfig{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
	}

	fsHandler := apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir(root)))

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", fsHandler))
	mux.HandleFunc("GET /api/healthz", api.HandleReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrix)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandleMetrixReset)
	mux.HandleFunc("POST /api/users", apiCfg.HandleAddUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandleAddChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandleGetChirp)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", root, port)
	log.Fatalf("Error: %s", srv.ListenAndServe())
}
