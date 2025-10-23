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

	userJson := users{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	}

	respondWithJSON(w, http.StatusOK, userJson)
	

}