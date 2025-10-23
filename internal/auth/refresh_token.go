package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func MakeRefreshToken() (string, error) {
	s := make([]byte, 32)
	_, err := rand.Read(s)
	if err != nil {
		log.Printf("error: %v", err)
		return "", err	
	}

	token := hex.EncodeToString(s)

	return token, nil
}