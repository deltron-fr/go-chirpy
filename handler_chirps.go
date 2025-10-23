package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/deltron-fr/chirpy/internal/auth"
	"github.com/deltron-fr/chirpy/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, req *http.Request) {
	
	type parameters struct {
		Body string `json:"body"`
	}

	type returnValid struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		CleanedBody string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error getting token")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		log.Printf("Error validating token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error validating token")
		return
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "error decoding json")
		return
	}

	trimmedChirp := strings.TrimSpace(params.Body)
	if len(trimmedChirp) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedMsg := checkBadWords(params.Body)

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: cleanedMsg,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error creating chirp")
		return
	}

	dataJson := returnValid{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		CleanedBody: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, dataJson)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	
	type Chirp struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		CleanedBody string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		log.Printf("error getting chirps: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error getting chirps")
		return
	}

	newChirps := make([]Chirp, len(chirps))	

	for i, chirp := range chirps {
		newChirps[i] = Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			CleanedBody: chirp.Body,
			UserID: chirp.UserID,
			}
	}

	respondWithJSON(w, http.StatusOK, newChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	type Chirp struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		CleanedBody string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("error converting uuid: %s", err)
		respondWithError(w, http.StatusBadRequest, "error converting id")
		return
	}

	chirp, err := cfg.db.GetChirp(context.Background(), chirpID)
	if err != nil {
		log.Printf("error getting chirp: %s", err)
		respondWithError(w, http.StatusNotFound, "chirp does not exist")
		return
	}

	chirpJson := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		CleanedBody: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirpJson)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {

	jwtToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error getting token")
		return
	}

	_, err = auth.ValidateJWT(jwtToken, cfg.secretKey)
	if err != nil {
		log.Printf("Error validating token: %s", err)
		respondWithError(w, http.StatusForbidden, "error validating token")
		return
	}

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		log.Printf("error converting uuid: %s", err)
		respondWithError(w, http.StatusBadRequest, "error converting id")
		return
	}

	rows, err := cfg.db.DeleteChirp(req.Context(), chirpID)
	if err != nil {
		log.Print("Error deleting chirp")
		respondWithError(w, http.StatusInternalServerError, "error deleting chirp")
		return
	}

	if rows == 0 {
		log.Printf("Error deleting chirp: %s", err)
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func checkBadWords(message string) string {
	m := make(map[string]bool)

	m["sharbert"] = true
	m["kerfuffle"] = true
	m["fornax"] = true

	words := strings.Split(message, " ")
	for i, word := range words {
		_, ok := m[strings.ToLower(word)]
		if !ok {
			continue
		} else {
			words[i] = "****"
		}
	}
	
	cleanedMessage := strings.Join(words, " ")
	return cleanedMessage

}