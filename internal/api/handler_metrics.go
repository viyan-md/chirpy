package api

import (
	"fmt"
	"net/http"
)

func (cfg *APIConfig) HandleMetrix(w http.ResponseWriter, r *http.Request) {
	html := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.FileserverHits.Load())
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(html))
}
