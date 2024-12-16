package main

import (
    "log"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "meli-coupon/internal/handlers"
    "meli-coupon/internal/repository"
    "meli-coupon/internal/services"
    "meli-coupon/internal/cache"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    cache := cache.NewMemoryCache()
    meliRepo := repository.NewMeliRepository(cache)
    couponService := services.NewCouponService(meliRepo)
    couponHandler := handlers.NewCouponHandler(couponService)

    router := mux.NewRouter()
    
    // Core endpoints
    router.HandleFunc("/coupon/", couponHandler.HandleCoupon).Methods("POST")
    router.HandleFunc("/coupon/stats", couponHandler.HandleStats).Methods("GET")
    
    // Monitoring endpoints
    router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
    router.HandleFunc("/metrics", handlers.MetricsHandler).Methods("GET")
    
    srv := &http.Server{
        Handler:      router,
        Addr:         ":8080",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }
    
    log.Printf("Server starting on port 8080...")
    log.Fatal(srv.ListenAndServe())
}