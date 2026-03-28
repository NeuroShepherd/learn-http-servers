package handlers

import (
	"net/http"
	"time"

	"github.com/neuroshepherd/learn-http-servers/internal/auth"
)

func (cfg *APIConfig) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	// look for a refresh token in the request header. If not found, return 401 unauthorized
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing refresh token")
		return
	}

	// validate the refresh token by checking if it exists in the database and is not expired
	user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	// if the refresh token is valid, generate a new JWT token and return it in the response
	newJWT, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Hour*1)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate new JWT token")
		return
	}

	type RefreshTokenResponse struct {
		JWT string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, RefreshTokenResponse{
		JWT: newJWT,
	})
}
