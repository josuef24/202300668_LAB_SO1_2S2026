package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	VM        string `json:"VM"`
	Carnet    string `json:"carnet"`
}

func main() {
	mux := http.NewServeMux()

	// Endpoint /health obligatorio para API2
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := HealthResponse{
			Status:    "UP",
			Message:   "API2 is Ready",
			Timestamp: time.Now().Format(time.RFC3339),
			VM:        "VM1",
			Carnet:    "202300668",
		}
		json.NewEncoder(w).Encode(response)
	})

	// Aquí irán tus endpoints de llamada call-api1 y call-api3

	http.ListenAndServe(":8082", mux)
}