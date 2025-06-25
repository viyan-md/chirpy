package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/viyan-md/chirpy/internal/auth"
	"github.com/viyan-md/chirpy/internal/database"
	"github.com/viyan-md/chirpy/internal/respond"
)

type TokenResponse struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *APIConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		TokenResponse
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

	token, err := auth.MakeJWT(dbuser.ID, cfg.JWTSecret)
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "couldn't create token", err)
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "could not create refresh token", err)
		return
	}

	_, err = cfg.DB.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{Token: refreshToken, UserID: dbuser.ID})
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "could not save refresh token", err)
		return
	}

	respond.RespondWithJSON(w, http.StatusOK, response{TokenResponse{Id: dbuser.ID, CreatedAt: dbuser.CreatedAt, UpdatedAt: dbuser.UpdatedAt, Email: dbuser.Email, Token: token, RefreshToken: refreshToken, IsChirpyRed: dbuser.IsChirpyRed}})
}
