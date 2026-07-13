package main

import (
	"fmt"
	"log"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	if _, err := fmt.Fprintf(w, "OK %s", r.URL.Path); err != nil {
		log.Println("response write error:", err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", health)

	log.Println("server listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("error:", err)
		return
	}
}
