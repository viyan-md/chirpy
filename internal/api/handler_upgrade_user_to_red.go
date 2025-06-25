package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/viyan-md/chirpy/internal/auth"
	"github.com/viyan-md/chirpy/internal/respond"
)

func (cfg *APIConfig) HandleUpgradeUserToRed(w http.ResponseWriter, r *http.Request) {
	polkaKey, err := auth.GetAPIKey(r.Header)
	if err != nil || polkaKey != cfg.API_Key {
		respond.RespondWithError(w, http.StatusUnauthorized, "missing or malformed Authorization header", nil)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "failed to decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respond.RespondWithError(w, http.StatusBadRequest, "user id is not valid", err)
		return
	}

	dbUser, err := cfg.DB.GetUserById(r.Context(), userID)
	if err != nil {
		respond.RespondWithError(w, http.StatusNotFound, "couldn't find user", err)
		return
	}

	if err = cfg.DB.UpgradeUserToRed(r.Context(), dbUser.ID); err != nil {
		log.Printf("Error: %v", err)
		respond.RespondWithError(w, http.StatusInternalServerError, "failed to upgrade user", err)
		return
	}

	w.WriteHeader(204)
}
