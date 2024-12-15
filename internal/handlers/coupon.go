package handlers

import (
    "encoding/json"
    "net/http"
    "log"
    "meli-coupon/internal/domain"
    "meli-coupon/internal/services"
)

// CouponHandler maneja las peticiones HTTP relacionadas con cupones
type CouponHandler struct {
    service services.CouponService
}

// NewCouponHandler crea una nueva instancia del handler
func NewCouponHandler(service services.CouponService) *CouponHandler {
    return &CouponHandler{service: service}
}

// HandleCoupon procesa la petición POST del cupón
func (h *CouponHandler) HandleCoupon(w http.ResponseWriter, r *http.Request) {
    var req domain.CouponRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Error al decodificar petición", http.StatusBadRequest)
        return
    }

    log.Printf("Petición recibida: %+v", req)

    response, err := h.service.ProcessCoupon(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// HandleStats procesa la petición GET de estadísticas
func (h *CouponHandler) HandleStats(w http.ResponseWriter, r *http.Request) {
    stats := h.service.GetTopFavorites()
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(stats)
}