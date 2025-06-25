package api

import (
	"encoding/json"
	"net/http"

	"github.com/viyan-md/chirpy/internal/auth"
	"github.com/viyan-md/chirpy/internal/database"
	"github.com/viyan-md/chirpy/internal/respond"
)

func (cfg *APIConfig) HandleChangeCredentials(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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

	var params parameters
	if err = json.NewDecoder(r.Body).Decode(&params); err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "couldn't decode json parameters", err)
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "password hashing failed", err)
		return
	}

	updatedUser, err := cfg.DB.UpdateCredentials(r.Context(), database.UpdateCredentialsParams{Email: params.Email, HashedPassword: hashed, ID: userID})
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "failed to update user", err)
		return
	}

	respond.RespondWithJSON(w, 200, UserResponse{ID: updatedUser.ID, CreatedAt: updatedUser.CreatedAt, UpdatedAt: updatedUser.UpdatedAt, Email: updatedUser.Email, IsChirpyRed: updatedUser.IsChirpyRed})
}
