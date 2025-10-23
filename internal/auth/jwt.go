package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: uuid.UUID.String(userID),
	})

	mySigningKey := []byte(tokenSecret)

	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Printf("error: %s", err)
		return "", err
	}

	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	mySigningKey := []byte(tokenSecret)

	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func (token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return mySigningKey, nil
	})

	if err != nil {
		log.Printf("error: %s", err)
		return uuid.Nil, err
	}

	uuidValue, err:= token.Claims.GetSubject()
	if err != nil {
		log.Printf("error: %s", err)
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(uuidValue)
	if err != nil {
		log.Printf("error: %s", err)
		return uuid.Nil, err
	}

	return userID, nil
}


func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Values("Authorization")
	if tokenString == nil {
		return "", fmt.Errorf("no authorization header")
	}
	
	token := strings.Split(tokenString[0], " ")[1]
	return token, nil

}