package auth

import (
	"net/http"
	"strings"
)


func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	key := strings.Split(apiKey," ")[1]

	return key, nil
}