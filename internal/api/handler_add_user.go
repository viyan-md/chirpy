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

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *APIConfig) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respond.RespondWithError(w, http.StatusBadRequest, "invalid JSON payload", err)
		return
	}

	if params.Email == "" {
		respond.RespondWithError(w, http.StatusBadRequest, "email is required", nil)
		return
	}

	if params.Password == "" {
		respond.RespondWithError(w, http.StatusBadRequest, "password is required", nil)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}

	createUserParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	}

	dbuser, err := cfg.DB.CreateUser(r.Context(), createUserParams)
	if err != nil {
		respond.RespondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
		return
	}

	user := UserResponse{
		ID:        dbuser.ID,
		CreatedAt: dbuser.CreatedAt,
		UpdatedAt: dbuser.UpdatedAt,
		Email:     dbuser.Email,
	}
	respond.RespondWithJSON(w, http.StatusCreated, user)
}
