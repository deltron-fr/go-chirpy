package main

import (
	"log"
	"net/http"
	"sync/atomic"
)


func main() {
	const port = "8080"

	mux := http.NewServeMux()

	serverHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(serverHandler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerServerHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetHits)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChrips)

	server := &http.Server{
		Addr:    ":" + port,
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

