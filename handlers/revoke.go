package handlers

import (
	"net/http"

	"github.com/neuroshepherd/learn-http-servers/internal/auth"
)

func (cfg *APIConfig) HandlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing token")
		return
	}

	_, err = cfg.DB.RevokeRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
