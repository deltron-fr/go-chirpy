package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	}) 

}

func main() {
	const port = "8080"

	
	mux := http.NewServeMux()

	serverHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(serverHandler))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerServerHits)
	mux.HandleFunc("/reset", apiCfg.handlerResetHits)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	log.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerServerHits(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		hits := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
		w.Write([]byte(hits))
}

func (cfg *apiConfig) handlerResetHits(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)

		cfg.fileserverHits.Swap(0)
}