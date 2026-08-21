package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api2/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"api":  "API 2 - Containerd",
			"info": "Funcionando correctamente",
		})
	})
	http.ListenAndServe(":8082", mux)
}