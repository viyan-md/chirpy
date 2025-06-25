package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/viyan-md/chirpy/internal/auth"
	"github.com/viyan-md/chirpy/internal/respond"
)

func (cfg *APIConfig) HandleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respond.RespondWithError(w, http.StatusUnauthorized, "missing or malformed Authorization header", nil)
		return
	}

	if err := cfg.DB.RevokeToken(r.Context(), refresh); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.RespondWithError(w, http.StatusUnauthorized, "invalid or expired refresh token", nil)
		} else {
			respond.RespondWithError(w, http.StatusInternalServerError, "failed to revoke token", err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
