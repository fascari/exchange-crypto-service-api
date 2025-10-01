package handler

import (
	"net/http"
	"time"

	httpjson "exchange-crypto-service-api/pkg/http"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	response := HealthResponse{
		Status:    "UP",
		Service:   "exchange-crypto-service-api",
		Version:   "1.0.0",
		Timestamp: time.Now(),
	}

	httpjson.WriteJSON(w, http.StatusOK, response)
}
