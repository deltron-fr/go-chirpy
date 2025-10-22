package main


import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)


func handlerValidateChrips(w http.ResponseWriter, req *http.Request) {
	
	type parameters struct {
		Body string `json:"body"`
	}

	type returnValid struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "error decoding json")
		return
	}

	chirp := strings.TrimSpace(params.Body)

	if len(chirp) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedMsg := checkBadWords(params.Body)
	resBody := returnValid{
		CleanedBody: cleanedMsg,
	}

	respondWithJSON(w, http.StatusOK, resBody)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
    respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
		return
	}

	w.Write(data)
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