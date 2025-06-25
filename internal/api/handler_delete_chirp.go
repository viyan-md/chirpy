package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/viyan-md/chirpy/internal/auth"
	"github.com/viyan-md/chirpy/internal/database"
	"github.com/viyan-md/chirpy/internal/respond"
)

func (cfg *APIConfig) HandleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	id_param := r.PathValue("chirpID")
	chirpID, _ := uuid.Parse(id_param)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respond.RespondWithError(w, http.StatusUnauthorized, "missing or malformed Authorization header", nil)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
	if err != nil {
		respond.RespondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	c, err := cfg.DB.GetChirp(r.Context(), chirpID)
	if err != nil {
		respond.RespondWithError(w, http.StatusNotFound, "chirp not found", nil)
		return
	}

	if c.UserID != userID {
		respond.RespondWithError(w, http.StatusForbidden, "access denied", nil)
		return
	}

	if err = cfg.DB.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpID,
		UserID: userID,
	}); err != nil {
		respond.RespondWithError(w, http.StatusNotFound, "chirp not found", nil)
		return
	}

	w.WriteHeader(204)
}
