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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT secret must be set")
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

	polkaAPIKey := os.Getenv("POLKA_KEY")
	if polkaAPIKey == "" {
		log.Fatal("API key must be set")
	}

	const root = "assets"
	const port = "8080"

	apiCfg := api.APIConfig{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
		JWTSecret:      jwtSecret,
		API_Key:        polkaAPIKey,
	}

	fsHandler := apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir(root)))

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", fsHandler))
	mux.HandleFunc("GET /api/healthz", api.HandleReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrix)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandleMetrixReset)
	mux.HandleFunc("POST /api/users", apiCfg.HandleAddUser)
	mux.Handle("POST /api/chirps", apiCfg.MiddlewareAuth(http.HandlerFunc(apiCfg.HandleAddChirp)))
	mux.HandleFunc("GET /api/chirps", apiCfg.HandleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandleGetChirp)
	mux.HandleFunc("POST /api/login", apiCfg.HandleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.HandleRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.HandleRevokeRefreshToken)
	mux.HandleFunc("PUT /api/users", apiCfg.HandleChangeCredentials)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.HandleDeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.HandleUpgradeUserToRed)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", root, port)
	log.Fatalf("Error: %s", srv.ListenAndServe())
}
