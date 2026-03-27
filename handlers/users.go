package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *APIConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type CreateUserRequest struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	var createUserReq CreateUserRequest
	err := decoder.Decode(&createUserReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(createUserReq.Email) == 0 || strings.ReplaceAll(createUserReq.Email, " ", "") == "" {
		respondWithError(w, http.StatusBadRequest, "Email cannot be empty")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), createUserReq.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	type CreateUserResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	respBody := CreateUserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusCreated, respBody)

}
