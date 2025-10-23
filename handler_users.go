package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/deltron-fr/chirpy/internal/auth"
	"github.com/deltron-fr/chirpy/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	type users struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "error decoding json")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error hashing password")
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error creating user")
		return
	}

	userJson := users{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusOK, userJson)
}


func (cfg *apiConfig) handlerLoginUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	type users struct {
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	const ExpiresIn = 3600 * time.Second

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "error decoding json")
		return
	}

	user, err := cfg.db.GetUserEmail(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error fetching user")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, ExpiresIn)
	if err != nil {
		log.Printf("Error creating token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error creating token")
	}

	ok, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("error checking password hash: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error checking password hash")
		return
	}

	if !ok {
		log.Println("password does not match")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token")
	}

	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(24 * 60 * time.Hour),
	})
	if err != nil {
		log.Printf("error creating refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "error creating refresh token")
		return
	}


	userJson := users{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, userJson)
}


func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "error getting token")
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
		respondWithError(w, http.StatusInternalServerError, "error getting token")
		return
	}

	err = cfg.db.UpdateRevokedAt(req.Context(), token)
	if err != nil {
		log.Printf("Error updating revoked_at, err: %v", err)
		respondWithError(w, http.StatusUnauthorized, "cannot perform this operation")
		return
	}
	w.WriteHeader(http.StatusNoContent)

}