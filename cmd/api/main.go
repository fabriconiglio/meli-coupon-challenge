package main

import (
    "log"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "meli-coupon/internal/handlers"
    "meli-coupon/internal/repository"
    "meli-coupon/internal/services"
)

func main() {
    // Configurar logging
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    // Inicializar dependencias
    meliRepo := repository.NewMeliRepository()
    couponService := services.NewCouponService(meliRepo)
    couponHandler := handlers.NewCouponHandler(couponService)

    // Configurar router
    router := mux.NewRouter()
    
    // Registrar rutas
    router.HandleFunc("/coupon/", couponHandler.HandleCoupon).Methods("POST")
    router.HandleFunc("/coupon/stats", couponHandler.HandleStats).Methods("GET")
    
    // Configurar servidor
    srv := &http.Server{
        Handler:      router,
        Addr:         ":8080",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    
    log.Printf("Server starting on port 8080...")
    log.Fatal(srv.ListenAndServe())
}