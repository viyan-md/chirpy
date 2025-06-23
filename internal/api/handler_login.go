package api

import (
	"encoding/json"
	"net/http"

	"github.com/viyan-md/chirpy/internal/auth"
	"github.com/viyan-md/chirpy/internal/respond"
)

func (cfg *APIConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		UserResponse
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "couldn't parse params", err)
		return
	}

	dbuser, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "couldn't load user", err)
		return
	}

	if err = auth.CheckPasswordHash(params.Password, dbuser.HashedPassword); err != nil {
		respond.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	user := UserResponse{
		ID:        dbuser.ID,
		CreatedAt: dbuser.CreatedAt,
		UpdatedAt: dbuser.UpdatedAt,
		Email:     dbuser.Email,
	}

	respond.RespondWithJSON(w, http.StatusOK, response{user})
}
