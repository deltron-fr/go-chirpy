package main

import (
	"net/http"
	"log"
)

func main() {
	const port = "8080"

	serverHandler := http.FileServer(http.Dir("."))

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	})

	mux.Handle("/app/", http.StripPrefix("/app/", serverHandler))
	
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	log.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}