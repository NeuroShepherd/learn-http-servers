package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/neuroshepherd/learn-http-servers/internal/auth"
	"github.com/neuroshepherd/learn-http-servers/internal/database"
)

func (cfg *APIConfig) HandlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type ReceivedChirpRequest struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	var receivedChirpReq ReceivedChirpRequest
	err := decoder.Decode(&receivedChirpReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(receivedChirpReq.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Chirp body cannot be empty")
		return
	}

	if len(receivedChirpReq.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// validate the jwt token from request header and make sure user ID from token matches
	// the user ID in the request body. If not, return 401 unauthorized
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid Authorization header")
		return
	}

	userIDFromToken, err := auth.ValidateJWT(tokenString, cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	chirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   receivedChirpReq.Body,
		UserID: userIDFromToken,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	type CreateChirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	respBody := CreateChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, respBody)

}

func (cfg *APIConfig) HandlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps")
		return
	}

	type GetAllChirpsResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	respBody := make([]GetAllChirpsResponse, len(chirps))
	for i, chirp := range chirps {
		respBody[i] = GetAllChirpsResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, respBody)
}

func (cfg *APIConfig) HandlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirpByID(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	type GetChirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	respBody := GetChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, respBody)

}
