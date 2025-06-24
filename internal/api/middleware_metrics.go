package api

import (
	"context"
	"log"
	"net/http"

	"github.com/viyan-md/chirpy/internal/auth"
	"github.com/viyan-md/chirpy/internal/respond"
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

func (cfg *APIConfig) MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respond.RespondWithError(w, http.StatusBadRequest, "couldn't parse token", err)
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
		if err != nil {
			respond.RespondWithError(w, http.StatusUnauthorized, "invalid token", err)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
