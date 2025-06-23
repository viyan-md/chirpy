package api

import (
	"net/http"

	"github.com/viyan-md/chirpy/internal/respond"
)

func (cfg *APIConfig) HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirps []Chirp `json:"chirps"`
	}

	dbList, err := cfg.DB.GetChirps(r.Context())
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "couldn't load chirps", err)
		return
	}

	respond.RespondWithJSON(w, http.StatusOK, response{Chirps: mapChirps(dbList)}.Chirps)
}
