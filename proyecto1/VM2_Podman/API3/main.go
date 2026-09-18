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

type CallResponse struct {
	ApiName    string `json:"apiname"`
	Message    string `json:"message"`
	Connection bool   `json:"connection"`
	Carnet     string `json:"carnet"`
}

func checkAPI(targetURL, targetAPI, targetVM string) CallResponse {
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(targetURL)
	if err != nil {
		return CallResponse{ApiName: targetAPI, Message: "ERROR: The " + targetAPI + " located on the " + targetVM + " is not working", Connection: false, Carnet: "202300668"}
	}
	defer resp.Body.Close()

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err == nil && health.Status == "UP" {
		return CallResponse{ApiName: targetAPI, Message: "The " + targetAPI + " located on the " + targetVM + " is working", Connection: true, Carnet: "202300668"}
	}
	return CallResponse{ApiName: targetAPI, Message: "ERROR: The " + targetAPI + " located on the " + targetVM + " is not working", Connection: false, Carnet: "202300668"}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{
			Status: "UP", Message: "API3 is Ready", Timestamp: time.Now().Format(time.RFC3339), VM: "VM2", Carnet: "202300668",
		})
	})

	mux.HandleFunc("/api3/202300668/call-api1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(checkAPI("http://192.168.122.115:8081/health", "API1", "VM1"))
	})

	mux.HandleFunc("/api3/202300668/call-api2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(checkAPI("http://192.168.122.115:8082/health", "API2", "VM1"))
	})

	http.ListenAndServe(":8080", mux)
}