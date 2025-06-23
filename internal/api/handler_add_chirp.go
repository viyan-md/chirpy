package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/viyan-md/chirpy/internal/database"
	"github.com/viyan-md/chirpy/internal/respond"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *APIConfig) HandleAddChirp(w http.ResponseWriter, r *http.Request) {
	const maxLen = 140

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > maxLen {
		respond.RespondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirpParams := database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserID,
	}

	chirp, err := cfg.DB.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	chirpResponse := response{
		Chirp: mapChirp(chirp),
	}

	respond.RespondWithJSON(w, http.StatusCreated, chirpResponse)

}

func mapChirp(src database.Chirp) Chirp {
	return Chirp{
		ID:        src.ID,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		Body:      src.Body,
		UserID:    src.UserID,
	}
}

func mapChirps(src []database.Chirp) []Chirp {
	dst := make([]Chirp, len(src))
	for i, c := range src {
		dst[i] = mapChirp(c)
	}
	return dst
}
