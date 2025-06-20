package main

import (
	"log"
	"net/http"

	"github.com/viyan-md/chirpy/internal/api"
)

func main() {
	const root = "assets"
	const port = "8080"

	apiCfg := api.APIConfig{}

	fsHandler := apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir(root)))

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", fsHandler))
	mux.HandleFunc("GET /api/healthz", api.HandleReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrix)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandleMetrixReset)
	mux.HandleFunc("POST /api/validate_chirp", api.HandleReadiness)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port %s\n", root, port)
	log.Fatalf("Error: %s", srv.ListenAndServe())
}
