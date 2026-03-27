package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
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

	chirp.Body = CleanChirpContent(chirp.Body)

	type ChirpValidResponse struct {
		Valid       bool   `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}

	respBody := ChirpValidResponse{Valid: true, CleanedBody: chirp.Body}

	respondWithJSON(w, http.StatusOK, respBody)

}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	respBody := ErrorResponse{Error: msg}

	respondWithJSON(w, code, respBody)
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

func CleanChirpContent(chirp string) string {

	bannedWords := []string{"kerfuffle", "sharbert", "fornax"}
	chirpWords := strings.Split(chirp, " ")
	updatedChirp := ""

	for _, word := range bannedWords {
		for i, chirpWord := range chirpWords {
			if strings.EqualFold(word, chirpWord) {
				chirpWords[i] = "****"
			} else {
				chirpWords[i] = chirpWord
			}
		}
	}
	updatedChirp = strings.Join(chirpWords, " ")
	return updatedChirp

}
