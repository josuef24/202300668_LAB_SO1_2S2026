package main

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Runtime string `json:"runtime"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := Response{
		Message: "¡API funcionando correctamente!",
		Status:  "success",
		Runtime: "Go Standard Library",
	}
	
	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)

	// Escucha en el puerto 8080
	http.ListenAndServe(":8080", mux)
}