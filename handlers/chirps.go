package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
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

	sortAscDesc := r.URL.Query().Get("sort")
	if sortAscDesc != "" && sortAscDesc != "asc" && sortAscDesc != "desc" {
		respondWithError(w, http.StatusBadRequest, "Invalid sort parameter")
		return
	}

	if sortAscDesc == "" {
		sortAscDesc = "asc"
	}

	// this is very poorly designed, but in order to keep the code in line with the expected
	// response format for the bootdev assignment, i am not renaming or reworking the SQLC
	// queries and response structs. In a real world application, I would likely have separate
	// queries and response structs for getting all chirps vs getting chirps by author ID, rather
	// than overloading the GetAllChirps query to handle both cases based on the presence of a query
	// parameter OR I would at least rename the GetAllChirps query to something more generic like
	// GetChirps and then have the function take in an optional author ID parameter.
	authorString := r.URL.Query().Get("author_id")
	if authorString != "" {
		authorID, err := uuid.Parse(authorString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id")
			return
		}

		chirpsByAuthor, err := cfg.DB.GetChirpsByAuthorID(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve chirps by author")
			return
		}

		type GetChirpsByAuthorResponse struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Body      string    `json:"body"`
			UserID    uuid.UUID `json:"user_id"`
		}

		respBody := make([]GetChirpsByAuthorResponse, len(chirpsByAuthor))
		for i, chirp := range chirpsByAuthor {
			respBody[i] = GetChirpsByAuthorResponse{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			}
		}

		// sort the chirps by created_at based on sortAscDesc
		switch sortAscDesc {
		case "asc":
			sort.Slice(respBody, func(i, j int) bool {
				return respBody[i].CreatedAt.Before(respBody[j].CreatedAt)
			})
		case "desc":
			sort.Slice(respBody, func(i, j int) bool {
				return respBody[i].CreatedAt.After(respBody[j].CreatedAt)
			})
		}

		respondWithJSON(w, http.StatusOK, respBody)
		return
	}

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

	switch sortAscDesc {
	case "asc":
		sort.Slice(respBody, func(i, j int) bool {
			return respBody[i].CreatedAt.Before(respBody[j].CreatedAt)
		})
	case "desc":
		sort.Slice(respBody, func(i, j int) bool {
			return respBody[i].CreatedAt.After(respBody[j].CreatedAt)
		})
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

func (cfg *APIConfig) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	// validate the jwt token from request header and make sure user ID from token matches
	// the user ID of the chirp. If not, return 401 unauthorized
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

	chirp, err := cfg.DB.GetChirpByID(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	if chirp.UserID != userIDFromToken {
		respondWithError(w, http.StatusForbidden, "Forbidden action")
		return
	}

	err = cfg.DB.DeleteChirpByID(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *APIConfig) HandlerUpdateChirpyRedStatus(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid API key")
		return
	}

	if apiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}
	type PolkaWebhookRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	var webhookReq PolkaWebhookRequest
	err = decoder.Decode(&webhookReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if webhookReq.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.DB.UpdateUserChirpyRedStatus(context.Background(), database.UpdateUserChirpyRedStatusParams{
		ID:          webhookReq.Data.UserID,
		IsChirpyRed: true,
	})
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user status")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
