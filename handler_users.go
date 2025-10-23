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
		IsChirpyRed bool `json:"is_chirpy_red"`
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
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, userJson)
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
		IsChirpyRed bool `json:"is_chirpy_red"`
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
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, userJson)
}


func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
	jwtToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("Error getting token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error getting token")
		return
	}

	userID, err := auth.ValidateJWT(jwtToken, cfg.secretKey)
	if err != nil {
		log.Printf("Error validating token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "error validating token")
		return
	}

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
		IsChirpyRed bool `json:"is_chirpy_red"`
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

	user, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: hashedPassword,
		ID: userID,
	})
	if err != nil {
		log.Printf("error updating user information, err: %v", err)
		respondWithError(w, http.StatusInternalServerError, "error updating user information")
		return
	}

	userJson := users{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: jwtToken,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, userJson)
}