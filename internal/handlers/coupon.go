package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "log"
    "meli-coupon/internal/domain"
    "meli-coupon/internal/services"
    "meli-coupon/internal/monitoring"
)

type CouponHandler struct {
    service services.CouponService
}

func NewCouponHandler(service services.CouponService) *CouponHandler {
    return &CouponHandler{service: service}
}

func (h *CouponHandler) HandleCoupon(w http.ResponseWriter, r *http.Request) {
    startTime := time.Now()
    monitoring.RecordRequest()

    var req domain.CouponRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        monitoring.RecordError()
        http.Error(w, "Error al decodificar petición", http.StatusBadRequest)
        return
    }

    log.Printf("Petición recibida: %+v", req)

    response, err := h.service.ProcessCoupon(req)
    if err != nil {
        monitoring.RecordError()
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)

    monitoring.RecordResponseTime(time.Since(startTime))
}

func (h *CouponHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
    startTime := time.Now()
    monitoring.RecordRequest()

    stats := h.service.GetTopFavorites()
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)

    monitoring.RecordResponseTime(time.Since(startTime))
}