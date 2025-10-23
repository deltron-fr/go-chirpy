package main

import (
	"log"
	"net/http"
	"time"

	"github.com/deltron-fr/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "missing or invalid token")
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(req.Context(), token)
	if err != nil {
		log.Printf("Error fetching token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "token does not exist or is expired")
		return
	}

	JWTToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secretKey, 60 * time.Minute)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating token")
		return
	}

	type refresh struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, refresh{Token: JWTToken})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "missing or invalid token")
		return
	}

	err = cfg.db.UpdateRevokedAt(req.Context(), token)
	if err != nil {
		log.Printf("Error updating revoked_at, err: %v", err)
		respondWithError(w, http.StatusNotFound, "token not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)

}