package api

import (
	"net/http"
)

func (cfg *APIConfig) HandleMetrixReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
}
