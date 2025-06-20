package api

import (
	"log"
	"net/http"
)

func (cfg *APIConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		cfg.FileserverHits.Add(1)

		log.Printf("➤ fileserver hit: %s %s; total=%d",
			r.Method, r.URL.Path, cfg.FileserverHits.Load(),
		)
	})
}
