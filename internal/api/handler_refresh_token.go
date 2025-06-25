package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/viyan-md/chirpy/internal/auth"
	"github.com/viyan-md/chirpy/internal/respond"
)

func (cfg *APIConfig) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respond.RespondWithError(w, http.StatusUnauthorized, "missing or malformed Authorization header", nil)
		return
	}

	user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refresh)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.RespondWithError(w, http.StatusUnauthorized, "invalid or expired refresh token", nil)
		} else {
			respond.RespondWithError(w, http.StatusInternalServerError, "database lookup failed", err)
		}
		return
	}

	access, err := auth.MakeJWT(user.ID, cfg.JWTSecret)
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "could not create access token", err)
		return
	}

	respond.RespondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{Token: access})
}
