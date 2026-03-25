package handlers

import (
	"encoding/json"
	"net/http"
)

func HandlerValidateChirpy(w http.ResponseWriter, r *http.Request) {

	type Chirps struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	var chirp Chirps
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(chirp.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Chirp body cannot be empty")
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	type ChirpValidResponse struct {
		Valid bool `json:"valid"`
	}

	respBody := ChirpValidResponse{Valid: true}

	respondWithJSON(w, http.StatusOK, respBody)

}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	respBody := ErrorResponse{Error: msg}

	respondWithJSON(w, http.StatusBadRequest, respBody)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	respBodyJSON, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respBodyJSON)
}
