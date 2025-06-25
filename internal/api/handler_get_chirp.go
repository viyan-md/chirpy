package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/viyan-md/chirpy/internal/respond"
)

func (cfg *APIConfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	id_param := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(id_param)

	log.Printf("Chirp parsed: %v", chirpID)

	if err != nil {
		respond.RespondWithError(w, http.StatusBadRequest, "invalid chirp id", err)
		return
	}

	type response struct {
		Chirp
	}

	chirp, err := cfg.DB.GetChirp(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.RespondWithError(w, http.StatusNotFound, "chirp not found", nil)
		} else {
			respond.RespondWithError(w, http.StatusInternalServerError, "couldn't load chirp", err)
		}
		return
	}

	respond.RespondWithJSON(w, http.StatusOK, response{Chirp: mapChirp(chirp)}.Chirp)
}
