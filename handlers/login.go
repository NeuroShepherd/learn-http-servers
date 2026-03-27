package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neuroshepherd/learn-http-servers/internal/auth"
)

func (cfg *APIConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	var loginReq LoginRequest
	err := decoder.Decode(&loginReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(loginReq.Email) == 0 || strings.ReplaceAll(loginReq.Email, " ", "") == "" {
		respondWithError(w, http.StatusBadRequest, "Email cannot be empty")
		return
	}

	// need a sql query to get users my email
	user, err := cfg.DB.GetUserByEmail(context.Background(), loginReq.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// need to compare the password with the hashed password
	validPassBool, err := auth.CheckPasswordHash(loginReq.Password, user.HashedPassword)
	if err != nil || !validPassBool {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	type LoginResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
