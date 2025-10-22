package main

import (
	"net/http"
	"fmt"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) handlerServerHits(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	hits := fmt.Sprintf("<html>\n	<body>\n		<h1>Welcome, Chirpy Admin</h1>\n		<p>Chirpy has been visited %d times!</p>\n	</body>\n</html>",
		cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}

func (cfg *apiConfig) handlerResetHits(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Swap(0)
}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}