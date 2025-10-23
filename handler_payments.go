package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/deltron-fr/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		log.Printf("error getting api key, err: %v", err)
		respondWithError(w, http.StatusUnauthorized, "missing or invalid apiKey")
		return
	}

	if apiKey != cfg.polkaKey {
		log.Print("error invalid api key")
		respondWithError(w, http.StatusUnauthorized, "invalid api key")
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

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	rows, err := cfg.db.UpgradeUser(req.Context(), params.Data.UserID)
	if err != nil {
		log.Printf("error upgrading user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "error upgrading user")
		return
	}

	if rows == 0 {
		log.Print("user not found")
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
