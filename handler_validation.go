package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/deltron-fr/chirpy/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, req *http.Request) {
	
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type returnValid struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		CleanedBody string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
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
	user, err := cfg.db.GetUser(req.Context(), params.UserID)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error getting user")
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: cleanedMsg,
		UserID: user.ID,
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