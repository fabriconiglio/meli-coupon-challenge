package handlers

import (
    "encoding/json"
    "net/http"
    "meli-coupon/internal/monitoring"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
    health := map[string]string{
        "status": "ok",
        "version": "1.0.0",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
    metrics := monitoring.GetMetrics()
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(metrics)
}