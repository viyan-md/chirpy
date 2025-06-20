package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/viyan-md/chirpy/internal/respond"
)

func HandleValidateChirpy(w http.ResponseWriter, r *http.Request) {
	const maxLen = 140
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	var params parameters

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > maxLen {
		respond.RespondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	clean := cleanInput(params.Body)

	respond.RespondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: clean,
	})
}

func cleanInput(text string) string {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}

	tokens := strings.Split(text, " ")

	for i, tok := range tokens {
		for _, prof := range profanities {
			if strings.EqualFold(tok, prof) {
				tokens[i] = "****"
				break
			}
		}
	}

	return strings.Join(tokens, " ")
}
